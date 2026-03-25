package interceptor

import (
	"bulletin-board-api/internal/constants"
	"bulletin-board-api/internal/service"
	"context"
	"fmt"
	"regexp"

	serviceErrors "bulletin-board-api/internal/errors"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	health                      = "Health"
	getAddressOwnership         = "GetAddressOwnership"
	createAddressOwnershipProof = "CreateAddressOwnershipProof"
)

var (
	methodRegex = regexp.MustCompile(`\/(.+)\/(.+)$`)

	// Presence in this map implies that we should not log errors for the method/code combination.
	expectedErrors = map[string]map[serviceErrors.Code]struct{}{
		getAddressOwnership:         {serviceErrors.NotFound: {}},
		createAddressOwnershipProof: {serviceErrors.AlreadyExists: {}},
	}
)

func getServiceAndMethod(fullMethod string) (string, string) {
	methodParts := methodRegex.FindStringSubmatch(fullMethod)
	if len(methodParts) == 0 {
		return fullMethod, ""
	}

	return methodParts[1], methodParts[2]
}

func GetRequestContextInterceptor(logger *zap.Logger, vaspService service.VaspService) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		_, method := getServiceAndMethod(info.FullMethod)

		// setting requesting vasp
		if method != health {
			// setting vasp to context
			md, ok := metadata.FromIncomingContext(ctx)
			headers := md.Get(constants.RequestVASPHeaderKey)
			if len(headers) == 0 {
				return nil, status.Error(codes.PermissionDenied, "missing VASP identity")
			} else if len(headers) > 1 {
				return nil, status.Error(codes.PermissionDenied, "multiple VASP identity")
			}

			vaspDomain := md.Get(constants.RequestVASPHeaderKey)[0]
			vasp, err := vaspService.SearchVasp(ctx, vaspDomain)
			if !ok || err != nil || vasp == nil {
				logger.Sugar().Errorw("invalid VASP from request", vaspDomain)
				return nil, status.Error(codes.PermissionDenied, fmt.Sprintf("invalid requesting VASP: %v", vaspDomain))
			}

			ctx = context.WithValue(ctx, constants.RequestVASPCtxKey{}, *vasp)
		}

		resp, err := handler(ctx, req)
		if err != nil && !isExpected(method, err) {
			logger.Error("request failed", zap.String("method", method), zap.Error(err))
		}
		return resp, serviceErrors.ToGrpcError(err)
	}
}

// Returns true if the method/error combination is an expected output that does not need to be logged as an error.
func isExpected(method string, err error) bool {
	methodErrors, expected := expectedErrors[method]
	if !expected {
		// Unexpected method.
		return false
	}

	e := serviceErrors.AsError(err)
	if e == nil {
		// Unexpected error type.
		return false
	}

	_, expected = methodErrors[e.Code]
	return expected
}
