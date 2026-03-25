package rpcimpl_test

import (
	"bulletin-board-api/internal/rpcimpl"
	"testing"

	servicemock "bulletin-board-api/internal/service/mocks"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type ServiceTester struct {
	mockAddressService *servicemock.AddressService
	mockVaspService    *servicemock.VaspService
	mockUtilityService *servicemock.UtilityService
	underTesting       *rpcimpl.Server
}

func (s *ServiceTester) AssertAll(t *testing.T) {
	t.Helper()
	s.mockUtilityService.AssertExpectations(t)
	s.mockAddressService.AssertExpectations(t)
	s.mockVaspService.AssertExpectations(t)
}

func NewServiceTester() *ServiceTester {
	mockAddressService := &servicemock.AddressService{}
	mockVaspService := &servicemock.VaspService{}
	mockUtilityService := &servicemock.UtilityService{}

	underTesting := rpcimpl.NewRPCI(
		mockUtilityService,
		mockVaspService,
		mockAddressService,
		zap.NewNop(),
		grpc.NewServer(),
	)
	return &ServiceTester{
		mockAddressService: mockAddressService,
		mockVaspService:    mockVaspService,
		mockUtilityService: mockUtilityService,
		underTesting:       underTesting,
	}
}
