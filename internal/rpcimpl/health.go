package rpcimpl

import (
	"context"
	"fmt"

	pb "bulletin-board-api/gen/go"
)

func (s *Server) Health(ctx context.Context, _ *pb.HealthRequest) (res *pb.HealthResponse, err error) {
	resp, err := s.utilityService.Health(ctx)
	if err != nil {
		return res, fmt.Errorf("health check failed: %w", err)
	}
	return &pb.HealthResponse{
		Message: resp,
	}, nil
}
