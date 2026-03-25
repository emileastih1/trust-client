package interceptor

import (
	"bulletin-board-api/internal/constants"
	"context"

	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func LoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, _ := metadata.FromIncomingContext(ctx)

		// Start setting additional tags that will show up in unary logs
		tags := grpc_ctxtags.NewTags()

		// extract and set request_id
		requestID := md.Get(constants.RequestRequestIDHeaderKey)
		if len(requestID) > 0 {
			tags = tags.Set("request_id", requestID[0])
		}

		// extract vasp_logging_id and set it
		vaspLoggingID := md.Get(constants.RequestVaspLoggingIDHeaderKey)
		if len(vaspLoggingID) > 0 {
			tags = tags.Set("vasp_logging_id", vaspLoggingID[0])
		}

		// Done building tags. Writing it to the context.
		ctx = grpc_ctxtags.SetInContext(ctx, tags)

		// execute handler
		resp, err := handler(ctx, req)

		return resp, err
	}
}
