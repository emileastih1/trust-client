//go:build wireinject
// +build wireinject

package cmd

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/google/wire"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"

	pb "bulletin-board-api/gen/go"
	"bulletin-board-api/internal/config"
	"bulletin-board-api/internal/constants"
	"bulletin-board-api/internal/encryption"
	serviceErrors "bulletin-board-api/internal/errors"
	"bulletin-board-api/internal/interceptor"
	"bulletin-board-api/internal/lib"
	"bulletin-board-api/internal/metrics"
	"bulletin-board-api/internal/repository"
	"bulletin-board-api/internal/rpcimpl"
	"bulletin-board-api/internal/secrets"
	bulletinboardservice "bulletin-board-api/internal/service"
	"bulletin-board-api/internal/tracer"
)

const (
	grpcPort    = 7001
	httpPort    = 7000
	ServiceName = "bulletin-board-api"
)

// **********************
// Service App ...
type App struct {
	listener      net.Listener
	gsrv          *grpc.Server
	gatewayServer *http.Server
	rpcimpl       *rpcimpl.Server
	logger        *zap.Logger
	ddTracer      func()
	addressRepo   repository.AddressRepository
}

func (a App) Start() error {

	err := profiler.Start(
		profiler.WithService(ServiceName),
		profiler.WithEnv(viper.GetString("stage")),
		profiler.WithProfileTypes(
			profiler.CPUProfile,
			profiler.HeapProfile,
			profiler.GoroutineProfile,
		),
	)

	if err != nil {
		a.logger.Sugar().Error("failed to start profiler", zap.Error(err))
		return err
	}

	defer profiler.Stop()
	// make sure we flush all call tracings
	defer a.ddTracer()

	a.logger.Sugar().Info("starting server")
	go func() {
		a.logger.Sugar().Info(fmt.Sprintf("GRPC server listening at http://0.0.0.0:%d", grpcPort))
		a.gsrv.Serve(a.listener)
	}()

	a.logger.Sugar().Info(fmt.Sprintf("HTTP server listening at http://0.0.0.0:%d", httpPort))
	err = a.gatewayServer.ListenAndServe()
	if err != nil {
		a.logger.Sugar().Fatal("grpc gateway start error", zap.Error(err))
		return err
	}

	return nil
}

func NewLogger() (*zap.Logger, error) {
	var cfg zap.Config
	if viper.GetString("stage") == constants.StageProduction {
		cfg = zap.NewProductionConfig()
		cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}
	cfg.EncoderConfig = zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	cfg.Encoding = "json"

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	defer logger.Sync()
	return logger, nil
}

func NewListener(logger *zap.Logger) (net.Listener, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", grpcPort))
	if err != nil {
		logger.Sugar().Fatal("start listener failed", zap.Error(err))
	}
	return lis, err
}

func NewGRPCServer(logger *zap.Logger, vaspService bulletinboardservice.VaspService, sd *metrics.StatsD) *grpc.Server {
	keepAliveParams := grpc.KeepaliveParams(
		keepalive.ServerParameters{
			Time:              10 * time.Second,
			Timeout:           5 * time.Second,
			MaxConnectionIdle: 55 * time.Second,
			MaxConnectionAge:  1800 * time.Second,
		})
	opts := []grpc_zap.Option{
		grpc_zap.WithDurationField(func(duration time.Duration) zapcore.Field {
			return zap.Int64("grpc.time_ms", duration.Milliseconds())
		}),
	}
	opt := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			interceptor.GetRequestMetricsInterceptor(sd),
			interceptor.LoggingInterceptor(),
			grpc_zap.UnaryServerInterceptor(logger, opts...),
			interceptor.GetRequestContextInterceptor(logger, vaspService),
		)),
		keepAliveParams,
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             5 * time.Second,
			PermitWithoutStream: true,
		}),
	}
	// grpc_zap.ReplaceGrpcLoggerV2(logger)
	server := grpc.NewServer(
		opt...,
	)
	return server
}

func setResponseHeader(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
	w.Header().Set("Cache-Control", "no-cache, no-store")
	return nil
}

func NewGrpcGateway(logger *zap.Logger) *http.Server {
	conn, err := grpc.DialContext(
		context.Background(),
		fmt.Sprintf("0.0.0.0:%d", grpcPort),
		grpc.WithInsecure(),
	)
	if err != nil {
		logger.Sugar().Error("Failed to dial server:", zap.Error(err))
	}

	// IMPORTANT: please note this is serving as a temporary hack to solve the XSS issue on the boolean field IOU
	// The root cause of that error is that grpc-gateway does not wrap the JSON unmarshal errors, and return
	// it directly to the client. The json unmarshal function echos back the user input, and we don't have control
	// over that lib, unless we provide a fix in grpc-gateway. A bug is already reported in github repo for grpc-gateway.
	customErrorHandler := func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, writer http.ResponseWriter, request *http.Request, err error) {
		if err != nil {
			s, ok := status.FromError(err)
			if ok && s.Code() == codes.InvalidArgument &&
				!strings.Contains(s.Message(), serviceErrors.ValidationErrorPrefix) {
				err = status.New(codes.InvalidArgument, "invalid request").Err()
			}
		}
		runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, writer, request, err)
	}

	gwmux := runtime.NewServeMux(
		runtime.WithForwardResponseOption(setResponseHeader),
		runtime.WithErrorHandler(customErrorHandler),
	)

	err = pb.RegisterBulletinBoardServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		logger.Sugar().Error("Failed to register gateway:", zap.Error(err))
	}

	gwServer := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", httpPort),
		Handler: gwmux,
	}

	return gwServer
}

func InitConfig() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("BBAPI")
	stage := viper.GetString("stage")
	if len(stage) == 0 {
		panic(errors.New("missing required variable"))
	}
	viper.SetConfigName(stage)
	viper.AddConfigPath("./configs")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("error loading config: %v", err))
	}

	viper.SetConfigName("common")
	viper.AddConfigPath("./configs")
	err = viper.MergeInConfig()
	if err != nil {
		panic(fmt.Sprintf("error loading config: %v", err))
	}

	v2 := viper.New()
	v2.AutomaticEnv()
	v2.SetEnvPrefix("BBAPI")
	v2.SetConfigName(stage)
	v2.AddConfigPath("./configs/vasp")
	err = v2.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("error loading config: %v", err))
	}

	err = viper.MergeConfigMap(v2.AllSettings())
	if err != nil {
		panic(fmt.Sprintf("error loading vasp: %v", err))
	}

	if err := config.ProcessCertificateBundle(); err != nil {
		panic(fmt.Sprintf("error processing certificate bundle: %v", err))
	}

	if stage == constants.StageLocal && (!viper.IsSet(constants.AddressSecret) || len(viper.GetString(constants.AddressSecret)) == 0) {
		panic(fmt.Errorf("address secret must be set"))
	}
}

func InitApp() (*App, error) {
	wire.Build(
		metrics.NewStatsD,
		repository.DBProvider,
		rpcimpl.NewRPCI,
		NewListener,
		NewGRPCServer,
		NewGrpcGateway,
		NewLogger,
		tracer.StartTracer,
		secrets.Provider,
		config.Provider,
		bulletinboardservice.Provider,
		encryption.Provider,
		wire.Struct(new(App), "*"),
	)

	return &App{}, nil
}

// **********************
// Bootstrap ...
type BootstrapApp struct {
	logger *zap.Logger
	svc    *secretsmanager.SecretsManager
}

func (a BootstrapApp) BootstrapSecrets() error {
	if viper.GetString("stage") == constants.StageLocal {
		return nil
	}

	generatedSalt, err := secrets.GenerateSecureLen256Salt()
	if err != nil {
		return fmt.Errorf("initialize secret failed %w", err)
	}

	description := "Address secret key"
	input := &secretsmanager.CreateSecretInput{
		Description:                 &description,
		ForceOverwriteReplicaSecret: lib.NewBoolPointer(false), // this should prevent us from overriding existing secrets
		KmsKeyId:                    nil,
		Name:                        &constants.AddressSecret,
		SecretString:                &generatedSalt,
	}
	res, err := a.svc.CreateSecret(input)
	if err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) {
			if awsErr.Code() == secretsmanager.ErrCodeResourceExistsException {
				return nil
			}
		} else {
			return fmt.Errorf("initialize secret failed %w", err)
		}
	}
	a.logger.Sugar().Infow("successfully created secret",
		zap.String("arn", *res.ARN),
		zap.String("name", *res.Name),
	)

	return nil
}

func InitBootstrapApp() (*BootstrapApp, error) {
	wire.Build(
		NewLogger,
		secrets.Provider,
		config.Provider,
		wire.Struct(new(BootstrapApp), "*"),
	)

	return &BootstrapApp{}, nil
}
