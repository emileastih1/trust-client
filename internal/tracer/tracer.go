package tracer

import (
	"bulletin-board-api/internal/constants"
	"fmt"
	"os"

	"go.uber.org/zap"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func StartTracer(logger *zap.Logger) func() {
	statsdHost := os.Getenv("DD_AGENT_HOST")
	if statsdHost == "" {
		statsdHost = "127.0.0.1"
	}
	logger.Sugar().Infow("Tracer StatsD initialization", "host", statsdHost)
	opts := []tracer.StartOption{
		tracer.WithAgentAddr(fmt.Sprintf("%s:8126", statsdHost)),
		tracer.WithAnalytics(true),
		tracer.WithDogstatsdAddress(fmt.Sprintf("%s:8125", statsdHost)),
		tracer.WithGlobalTag(ext.AnalyticsEvent, true),
		tracer.WithGlobalTag("version", os.Getenv("CODEFLOW_RELEASE_ID")),
		tracer.WithLogger((*ddZapLogger)(logger)),
		tracer.WithRuntimeMetrics(),
		tracer.WithService(constants.BulletinBoardServiceName),
	}
	tracer.Start(opts...)
	return tracer.Stop
}

type ddZapLogger zap.Logger

func (log *ddZapLogger) Log(msg string) {
	if log == nil {
		return
	}
	(*zap.Logger)(log).Info(msg)
}
