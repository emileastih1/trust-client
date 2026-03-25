package service_test

import (
	"bulletin-board-api/internal/metrics"
	"bulletin-board-api/internal/service"
	"testing"

	addressRepoMock "bulletin-board-api/internal/repository/mocks"
	repoMocks "bulletin-board-api/internal/repository/mocks"

	encrypterMock "bulletin-board-api/internal/encryption/mocks"
	vaspServiceMock "bulletin-board-api/internal/service/mocks"

	"go.uber.org/zap"
)

const (
	testAddressLocation = "test_address_location"
	testErrorMessage    = "test error message"
	testDomain          = "test.com"
)

// AddressServiceTester ...
type AddressServiceTester struct {
	mockVaspService *vaspServiceMock.VaspService
	mockAddressRepo *addressRepoMock.AddressRepository
	mockEncrypter   *encrypterMock.MockEncrypter

	underTesting service.AddressService
}

func NewAddressServiceTester() *AddressServiceTester {
	mockVaspService := &vaspServiceMock.VaspService{}
	mockAddressRepo := &addressRepoMock.AddressRepository{}
	mockEncrypter := &encrypterMock.MockEncrypter{}

	underTesting := service.NewAddressService(mockVaspService, mockAddressRepo, zap.NewNop(), mockEncrypter, metrics.NewTest())

	return &AddressServiceTester{
		mockVaspService: mockVaspService,
		mockAddressRepo: mockAddressRepo,
		mockEncrypter:   mockEncrypter,
		underTesting:    underTesting,
	}
}

func (te *AddressServiceTester) AssertAll(t *testing.T) {
	t.Helper()
	te.mockVaspService.AssertExpectations(t)
	te.mockAddressRepo.AssertExpectations(t)
}

// VaspServiceTester ...
type VaspServiceTester struct {
	mockVaspRepo *repoMocks.VaspRepository

	underTesting service.VaspService
}

func NewVaspServiceTester() *VaspServiceTester {
	mockVaspRepo := &repoMocks.VaspRepository{}

	underTesting := service.NewVaspService(mockVaspRepo, zap.NewNop())

	return &VaspServiceTester{
		mockVaspRepo: mockVaspRepo,
		underTesting: underTesting,
	}
}

func (te *VaspServiceTester) AssertAll(t *testing.T) {
	t.Helper()
	te.mockVaspRepo.AssertExpectations(t)
}
