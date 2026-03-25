package interceptor

import (
	"bulletin-board-api/internal/metrics"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func GetRequestMetricsInterceptor(sd *metrics.StatsD) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		service, method := getServiceAndMethod(info.FullMethod)

		options := []string{
			"service", service,
			"method", method,
		}

		// record latency
		timer := metrics.NewTimer(sd, metrics.MetricBulletinBoardGRPCLatency, options...)
		defer timer.Emit()

		// execute handler
		resp, err := handler(ctx, req)

		// post processing - metrics
		timer.AppendTags([]string{"status", status.Code(err).String()}...)

		increaseMetrics(sd, service, method, err)

		return resp, err
	}
}

func increaseMetrics(
	sd *metrics.StatsD,
	service string,
	method string,
	err error,
) {
	options := []string{
		"service", service,
		"method", method,
		"status", status.Code(err).String(),
	}
	sd.Incr(
		metrics.MetricBulletinBoardGRPCResponse,
		options...,
	)
}
