package rpcimpl

import (
	pb "bulletin-board-api/gen/go"
	"bulletin-board-api/internal/service"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"go.uber.org/zap"

	"google.golang.org/grpc"
)

var tracerType = tracer.SpanType("function")

type Server struct {
	utilityService service.UtilityService
	vaspService    service.VaspService
	addressService service.AddressService
	logger         *zap.Logger
	pb.UnimplementedBulletinBoardServiceServer
}

func NewRPCI(
	utilityService service.UtilityService,
	vaspService service.VaspService,
	addressService service.AddressService,
	logger *zap.Logger,
	gsrv *grpc.Server,
) *Server {
	s := &Server{
		utilityService: utilityService,
		vaspService:    vaspService,
		addressService: addressService,
		logger:         logger,
	}
	pb.RegisterBulletinBoardServiceServer(gsrv, s)

	return s
}
