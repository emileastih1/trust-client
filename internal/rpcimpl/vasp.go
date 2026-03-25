package rpcimpl

import (
	"context"
	"fmt"

	pb "bulletin-board-api/gen/go"
	serviceErrors "bulletin-board-api/internal/errors"

	"github.com/google/uuid"
)

func (s *Server) GetVasp(ctx context.Context, request *pb.GetVaspRequest) (res *pb.GetVaspResponse, err error) {
	vaspUUID, err := uuid.Parse(request.VaspId)
	if err != nil {
		return nil, serviceErrors.New("invalid VASP UUID", serviceErrors.InvalidArgument)
	}

	vasp, err := s.vaspService.GetVasp(ctx, vaspUUID)
	if err != nil {
		return nil, fmt.Errorf("rpc get vasp failed: %w", err)
	}

	return &pb.GetVaspResponse{
		Vasp: vasp.ToProtobuf(),
	}, nil
}

func (s *Server) GetVasps(ctx context.Context, _ *pb.GetVaspsRequest) (res *pb.GetVaspsResponse, err error) {
	vasps, err := s.vaspService.GetVasps(ctx)
	if err != nil {
		return nil, fmt.Errorf("rpc get vasps failed: %w", err)
	}
	vaspsPb := []*pb.VASP{}
	for _, v := range vasps {
		vaspsPb = append(vaspsPb, v.ToProtobuf())
	}
	return &pb.GetVaspsResponse{
		Vasps: vaspsPb,
	}, nil
}
