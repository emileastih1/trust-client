package rpcimpl_test

import (
	"bulletin-board-api/internal/constants"
	"bulletin-board-api/internal/lib"
	"bulletin-board-api/internal/models"
	"bulletin-board-api/internal/service"
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"

	pb "bulletin-board-api/gen/go"

	serviceErrors "bulletin-board-api/internal/errors"

	servicemock "bulletin-board-api/internal/service/mocks"
)

const (
	testErrorMessage   = "test error message from rpcimpl"
	testValidAddress   = "9c7bdd2e46af171213aba0217343c518d2b8c3612e7001e768a9bcbb69eb82086031eff8d4c1363894615c143566405fb23be942c26d6fdd4944f9a52090f523"
	testValidSignature = "test_valid_signature"
	testValidPrefix    = "test_valid_prefix"
	testValidProofType = "BITCOIN_P2SH"
	testInvalidString  = "<script>alert(1)</script>"
)

var errTest = errors.New(testErrorMessage)

type addressServiceCallSetup struct {
	want        bool
	expectation func(mock *servicemock.AddressService)
}

func TestServer_CreateAddressOwnershipProof(t *testing.T) {
	// setup
	testValidNonce := "123123123"
	testRegistrationID := uuid.New()
	testProofType := "BITCOIN_P2PKH"
	testValidAuxProofData := pb.AuxProofData{
		Type: "ETH_CREATE",
		Data: map[string]string{
			constants.Nonce: testValidNonce,
		},
	}

	testAuxProofData := []service.AuxProofData{{
		Type: constants.EthCreate,
		Data: map[string]string{
			constants.Nonce: testValidNonce,
		},
	}}

	tests := map[string]struct {
		in                 *pb.CreateAddressOwnershipProofRequest
		assert             func(t *testing.T, res *pb.CreateAddressOwnershipProofResponse, err error)
		addressServiceCall addressServiceCallSetup
	}{
		"success": {
			in: &pb.CreateAddressOwnershipProofRequest{
				Address:        testValidAddress,
				RegistrationId: testRegistrationID.String(),
				Chain:          constants.ChainBitcoin,
				Iou:            false,
				ProofType:      testProofType,
				AuxProofData:   []*pb.AuxProofData{&testValidAuxProofData},
			},
			assert: func(t *testing.T, res *pb.CreateAddressOwnershipProofResponse, err error) {
				t.Helper()
				assert.NoError(t, err)
				assert.Equal(t, testRegistrationID.String(), res.GetRegistrationId())
			},
			addressServiceCall: addressServiceCallSetup{
				want: true,
				expectation: func(mock *servicemock.AddressService) {
					expectedResp := &service.UpdateAddressOwnershipProofResponse{
						AddressOwnershipProof: &models.AddressOwnershipProof{
							Ownership: &models.AddressOwnership{
								ID:      testRegistrationID,
								Address: testValidAddress,
								Chain:   constants.ChainBitcoin,
							},
							IOU:       lib.NewBoolPointer(false),
							ProofType: &testProofType,
						},
					}
					mock.On("UpdateAddressOwnershipProof", testifyMock.AnythingOfType("*context.valueCtx"), &service.UpdateAddressOwnershipProofRequest{
						Chain:          constants.ChainBitcoin,
						IOU:            false,
						Address:        testValidAddress,
						Signature:      lib.StringToStringPtr(""),
						Prefix:         lib.StringToStringPtr(""),
						RegistrationID: testRegistrationID,
						ProofType:      &testProofType,
						AuxProofData:   testAuxProofData,
					}).Return(expectedResp, nil)
				},
			},
		},
		"success, ripple": {
			in: &pb.CreateAddressOwnershipProofRequest{
				Address:        testValidAddress,
				RegistrationId: testRegistrationID.String(),
				Chain:          constants.ChainRipple,
				Iou:            false,
				ProofType:      constants.RippleClassic,
				AuxProofData: []*pb.AuxProofData{
					{
						Data: map[string]string{
							constants.EncryptedPublicKey: "M3Iw40JVyfgeAVJMtoYufDyz/ltf3OmMT4650RlKTXFL/AbW2ih3P1W+O47qOH9g4pZX0y8v/eCWQffkjTaRoCV5MjG6",
						},
					},
				},
			},
			assert: func(t *testing.T, res *pb.CreateAddressOwnershipProofResponse, err error) {
				t.Helper()
				assert.NoError(t, err)
				assert.Equal(t, testRegistrationID.String(), res.GetRegistrationId())
			},
			addressServiceCall: addressServiceCallSetup{
				want: true,
				expectation: func(mock *servicemock.AddressService) {
					proofType := constants.RippleClassic
					expectedResp := &service.UpdateAddressOwnershipProofResponse{
						AddressOwnershipProof: &models.AddressOwnershipProof{
							Ownership: &models.AddressOwnership{
								ID:      testRegistrationID,
								Address: testValidAddress,
								Chain:   constants.ChainRipple,
							},
							IOU:       lib.NewBoolPointer(false),
							ProofType: &proofType,
						},
					}
					mock.On("UpdateAddressOwnershipProof", testifyMock.AnythingOfType("*context.valueCtx"), testifyMock.Anything).Return(expectedResp, nil)
				},
			},
		},
		"failure, invalid address": {
			in: &pb.CreateAddressOwnershipProofRequest{
				Address:        "testValidAddress",
				RegistrationId: uuid.New().String(),
				Chain:          constants.ChainBitcoin,
				Iou:            true,
			},
			assert: func(t *testing.T, res *pb.CreateAddressOwnershipProofResponse, err error) {
				t.Helper()
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.InvalidArgument))
				assert.Equal(t, "invalid address", err.Error())
				assert.Nil(t, res)
			},
		},
		"failure, invalid chain": {
			in: &pb.CreateAddressOwnershipProofRequest{
				Address:        testValidAddress,
				RegistrationId: uuid.New().String(),
				Chain:          "invalid chain",
				Iou:            true,
			},
			assert: func(t *testing.T, res *pb.CreateAddressOwnershipProofResponse, err error) {
				t.Helper()
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.InvalidArgument))
				assert.Equal(t, "invalid chain: invalid chain", err.Error())
				assert.Nil(t, res)
			},
		},
		"failure, invalid signature": {
			in: &pb.CreateAddressOwnershipProofRequest{
				Address:        testValidAddress,
				RegistrationId: uuid.New().String(),
				Chain:          constants.ChainBitcoin,
				Iou:            true,
				Prefix:         testValidPrefix,
				Signature:      testInvalidString,
			},
			assert: func(t *testing.T, res *pb.CreateAddressOwnershipProofResponse, err error) {
				t.Helper()
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.InvalidArgument))
				assert.Equal(t, "invalid signature", err.Error())
				assert.Nil(t, res)
			},
		},
		"failure, invalid proof type": {
			in: &pb.CreateAddressOwnershipProofRequest{
				Address:        testValidAddress,
				RegistrationId: uuid.New().String(),
				Chain:          constants.ChainBitcoin,
				Iou:            true,
				Prefix:         testValidPrefix,
				Signature:      testValidSignature,
				ProofType:      testInvalidString,
			},
			assert: func(t *testing.T, res *pb.CreateAddressOwnershipProofResponse, err error) {
				t.Helper()
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.InvalidArgument))
				assert.Equal(t, "invalid proof type", err.Error())
				assert.Nil(t, res)
			},
		},
		"failure, invalid aux proof data type": {
			in: &pb.CreateAddressOwnershipProofRequest{
				Address:        testValidAddress,
				RegistrationId: uuid.New().String(),
				Chain:          constants.ChainBitcoin,
				Iou:            true,
				Prefix:         testValidPrefix,
				Signature:      testValidSignature,
				ProofType:      testValidProofType,
				AuxProofData: []*pb.AuxProofData{
					{
						Type: testInvalidString,
						Data: nil,
					},
				},
			},
			assert: func(t *testing.T, res *pb.CreateAddressOwnershipProofResponse, err error) {
				t.Helper()
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.InvalidArgument))
				assert.Equal(t, "invalid aux proof data", err.Error())
				assert.Nil(t, res)
			},
		},
		"failure, invalid aux proof data key": {
			in: &pb.CreateAddressOwnershipProofRequest{
				Address:        testValidAddress,
				RegistrationId: uuid.New().String(),
				Chain:          constants.ChainBitcoin,
				Iou:            true,
				Prefix:         testValidPrefix,
				Signature:      testValidSignature,
				ProofType:      testValidProofType,
				AuxProofData: []*pb.AuxProofData{
					{
						Type: testValidProofType,
						Data: map[string]string{
							testInvalidString: "value",
						},
					},
				},
			},
			assert: func(t *testing.T, res *pb.CreateAddressOwnershipProofResponse, err error) {
				t.Helper()
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.InvalidArgument))
				assert.Equal(t, "invalid aux proof data", err.Error())
				assert.Nil(t, res)
			},
		},
		"failure, invalid aux proof data value": {
			in: &pb.CreateAddressOwnershipProofRequest{
				Address:        testValidAddress,
				RegistrationId: uuid.New().String(),
				Chain:          constants.ChainBitcoin,
				Iou:            true,
				Prefix:         testValidPrefix,
				Signature:      testValidSignature,
				ProofType:      testValidProofType,
				AuxProofData: []*pb.AuxProofData{
					{
						Type: testValidProofType,
						Data: map[string]string{
							constants.Nonce: testInvalidString,
						},
					},
				},
			},
			assert: func(t *testing.T, res *pb.CreateAddressOwnershipProofResponse, err error) {
				t.Helper()
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.InvalidArgument))
				assert.Equal(t, "invalid aux proof data", err.Error())
				assert.Nil(t, res)
			},
		},
		"failure, invalid registration id": {
			in: &pb.CreateAddressOwnershipProofRequest{
				Address:        testValidAddress,
				RegistrationId: "invalid",
				Chain:          constants.ChainBitcoin,
				Iou:            true,
				Prefix:         testValidPrefix,
				Signature:      testValidSignature,
				ProofType:      testValidProofType,
				AuxProofData:   nil,
			},
			assert: func(t *testing.T, res *pb.CreateAddressOwnershipProofResponse, err error) {
				t.Helper()
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.InvalidArgument))
				assert.Equal(t, "invalid registration id", err.Error())
				assert.Nil(t, res)
			},
		},
		"failure, unexpected internal error, should map to grpc internal error": {
			in: &pb.CreateAddressOwnershipProofRequest{
				Address:        testValidAddress,
				RegistrationId: testRegistrationID.String(),
				Chain:          constants.ChainBitcoin,
				Iou:            false,
				ProofType:      testProofType,
				AuxProofData:   []*pb.AuxProofData{&testValidAuxProofData},
			},
			assert: func(t *testing.T, res *pb.CreateAddressOwnershipProofResponse, err error) {
				t.Helper()
				assert.EqualError(t, err, "address service error: "+testErrorMessage)
				assert.Nil(t, res)
			},
			addressServiceCall: addressServiceCallSetup{
				want: true,
				expectation: func(mock *servicemock.AddressService) {
					mock.On("UpdateAddressOwnershipProof", testifyMock.AnythingOfType("*context.valueCtx"), &service.UpdateAddressOwnershipProofRequest{
						Chain:          constants.ChainBitcoin,
						IOU:            false,
						Address:        testValidAddress,
						Signature:      lib.StringToStringPtr(""),
						Prefix:         lib.StringToStringPtr(""),
						RegistrationID: testRegistrationID,
						ProofType:      &testProofType,
						AuxProofData:   testAuxProofData,
					}).Return(nil, fmt.Errorf("%w", errTest))
				},
			},
		},
		"failure, not allowed error, should map to grpc not allowed error": {
			in: &pb.CreateAddressOwnershipProofRequest{
				Address:        testValidAddress,
				RegistrationId: testRegistrationID.String(),
				Chain:          constants.ChainBitcoin,
				Iou:            false,
				ProofType:      testProofType,
				AuxProofData:   []*pb.AuxProofData{&testValidAuxProofData},
			},
			assert: func(t *testing.T, res *pb.CreateAddressOwnershipProofResponse, err error) {
				t.Helper()
				assert.True(t, serviceErrors.EqualErrorCode(errors.Unwrap(err), serviceErrors.PermissionDenied))
				assert.Nil(t, res)
			},
			addressServiceCall: addressServiceCallSetup{
				want: true,
				expectation: func(mock *servicemock.AddressService) {
					mock.On("UpdateAddressOwnershipProof", testifyMock.AnythingOfType("*context.valueCtx"), &service.UpdateAddressOwnershipProofRequest{
						Chain:          constants.ChainBitcoin,
						IOU:            false,
						Address:        testValidAddress,
						Signature:      lib.StringToStringPtr(""),
						Prefix:         lib.StringToStringPtr(""),
						RegistrationID: testRegistrationID,
						ProofType:      &testProofType,
						AuxProofData:   testAuxProofData,
					}).Return(nil, serviceErrors.New(testErrorMessage, serviceErrors.PermissionDenied))
				},
			},
		},
	}

	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			// given
			te := NewServiceTester()
			if tt.addressServiceCall.want {
				tt.addressServiceCall.expectation(te.mockAddressService)
			}

			// when
			res, err := te.underTesting.CreateAddressOwnershipProof(context.Background(), tt.in)

			// then
			if tt.assert != nil {
				tt.assert(t, res, err)
			}
			te.AssertAll(t)
		})
	}
}

func TestServer_CreateAddressOwnership(t *testing.T) {
	// setup
	testRegistrationID := uuid.New()
	tests := map[string]struct {
		in                 *pb.CreateAddressOwnershipRequest
		assert             func(t *testing.T, res *pb.CreateAddressOwnershipResponse, err error)
		addressServiceCall addressServiceCallSetup
	}{
		"success": {
			in: &pb.CreateAddressOwnershipRequest{
				Address: testValidAddress,
				Chain:   constants.ChainBitcoin,
			},
			assert: func(t *testing.T, res *pb.CreateAddressOwnershipResponse, err error) {
				t.Helper()
				assert.Equal(t, testRegistrationID.String(), res.RegistrationId)
				assert.NoError(t, err)
			},
			addressServiceCall: addressServiceCallSetup{
				want: true,
				expectation: func(mock *servicemock.AddressService) {
					mock.On("CreateAddressOwnership", testifyMock.AnythingOfType("*context.valueCtx"), &service.CreateAddressOwnershipRequest{
						Chain:   constants.ChainBitcoin,
						Address: testValidAddress,
					}).Return(&service.CreateAddressOwnershipResponse{
						Address:        testValidAddress,
						Chain:          constants.ChainBitcoin,
						RegistrationID: testRegistrationID,
					}, nil)
				},
			},
		},
		"failure, invalid address": {
			in: &pb.CreateAddressOwnershipRequest{
				Address: testInvalidString,
				Chain:   constants.ChainBitcoin,
			},
			assert: func(t *testing.T, res *pb.CreateAddressOwnershipResponse, err error) {
				t.Helper()
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.InvalidArgument))
				assert.Equal(t, "invalid address", err.Error())
				assert.Nil(t, res)
			},
		},

		"failure, invalid chain": {
			in: &pb.CreateAddressOwnershipRequest{
				Address: testValidAddress,
				Chain:   testInvalidString,
			},
			assert: func(t *testing.T, res *pb.CreateAddressOwnershipResponse, err error) {
				t.Helper()
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.InvalidArgument))
				assert.Equal(t, "invalid chain: &lt;script&gt;alert(1)&lt;/script&gt;", err.Error())
				assert.Nil(t, res)
			},
		},
		"failure, unexpected internal error, should map to grpc internal error": {
			in: &pb.CreateAddressOwnershipRequest{
				Address: testValidAddress,
				Chain:   constants.ChainBitcoin,
			},
			assert: func(t *testing.T, res *pb.CreateAddressOwnershipResponse, err error) {
				t.Helper()
				assert.EqualError(t, err, "address service error: "+testErrorMessage)
				assert.Nil(t, res)
			},
			addressServiceCall: addressServiceCallSetup{
				want: true,
				expectation: func(mock *servicemock.AddressService) {
					mock.On("CreateAddressOwnership", testifyMock.AnythingOfType("*context.valueCtx"), &service.CreateAddressOwnershipRequest{
						Chain:   constants.ChainBitcoin,
						Address: testValidAddress,
					}).Return(nil, fmt.Errorf("%w", errTest))
				},
			},
		},
		"failure, conflict error, should map to grpc already exists error": {
			in: &pb.CreateAddressOwnershipRequest{
				Address: testValidAddress,
				Chain:   constants.ChainBitcoin,
			},
			assert: func(t *testing.T, res *pb.CreateAddressOwnershipResponse, err error) {
				t.Helper()
				assert.True(t, serviceErrors.EqualErrorCode(errors.Unwrap(err), serviceErrors.AlreadyExists))
				assert.Nil(t, res)
			},
			addressServiceCall: addressServiceCallSetup{
				want: true,
				expectation: func(mock *servicemock.AddressService) {
					mock.On("CreateAddressOwnership", testifyMock.AnythingOfType("*context.valueCtx"), &service.CreateAddressOwnershipRequest{
						Chain:   constants.ChainBitcoin,
						Address: testValidAddress,
					}).Return(nil, serviceErrors.New(testErrorMessage, serviceErrors.AlreadyExists))
				},
			},
		},
		"failure, not found error, should map to grpc not found error": {
			in: &pb.CreateAddressOwnershipRequest{
				Address: testValidAddress,
				Chain:   constants.ChainBitcoin,
			},
			assert: func(t *testing.T, res *pb.CreateAddressOwnershipResponse, err error) {
				t.Helper()
				assert.True(t, serviceErrors.EqualErrorCode(errors.Unwrap(err), serviceErrors.NotFound))
				assert.Nil(t, res)
			},
			addressServiceCall: addressServiceCallSetup{
				want: true,
				expectation: func(mock *servicemock.AddressService) {
					mock.On("CreateAddressOwnership", testifyMock.AnythingOfType("*context.valueCtx"), &service.CreateAddressOwnershipRequest{
						Chain:   constants.ChainBitcoin,
						Address: testValidAddress,
					}).Return(nil, serviceErrors.New(testErrorMessage, serviceErrors.NotFound))
				},
			},
		},
		"failure, not allowed error, should map to grpc not allowed error": {
			in: &pb.CreateAddressOwnershipRequest{
				Address: testValidAddress,
				Chain:   constants.ChainBitcoin,
			},
			assert: func(t *testing.T, res *pb.CreateAddressOwnershipResponse, err error) {
				t.Helper()
				assert.True(t, serviceErrors.EqualErrorCode(errors.Unwrap(err), serviceErrors.PermissionDenied))
				assert.Nil(t, res)
			},
			addressServiceCall: addressServiceCallSetup{
				want: true,
				expectation: func(mock *servicemock.AddressService) {
					mock.On("CreateAddressOwnership", testifyMock.AnythingOfType("*context.valueCtx"), &service.CreateAddressOwnershipRequest{
						Chain:   constants.ChainBitcoin,
						Address: testValidAddress,
					}).Return(nil, serviceErrors.New(testErrorMessage, serviceErrors.PermissionDenied))
				},
			},
		},
	}

	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			// given
			te := NewServiceTester()
			if tt.addressServiceCall.want {
				tt.addressServiceCall.expectation(te.mockAddressService)
			}

			// when
			res, err := te.underTesting.CreateAddressOwnership(context.Background(), tt.in)

			// then
			if tt.assert != nil {
				tt.assert(t, res, err)
			}
			te.AssertAll(t)
		})
	}
}

func TestServer_GetAddressOwnership(t *testing.T) {
	// setup
	testRegistrationID := uuid.New()
	testVASPUUID := uuid.New()
	tests := map[string]struct {
		in                 *pb.GetAddressOwnershipRequest
		assert             func(t *testing.T, res *pb.GetAddressOwnershipResponse, err error)
		addressServiceCall addressServiceCallSetup
	}{
		"success": {
			in: &pb.GetAddressOwnershipRequest{
				Address: testValidAddress,
				Chain:   constants.ChainBitcoin,
			},
			assert: func(t *testing.T, res *pb.GetAddressOwnershipResponse, err error) {
				t.Helper()
				assert.Equal(t, testRegistrationID.String(), res.AddressOwnershipProof.Id)
				assert.Equal(t, testValidAddress, res.AddressOwnershipProof.Address)
				assert.Equal(t, testVASPUUID.String(), res.Vasp.Id)
				assert.NoError(t, err)
			},
			addressServiceCall: addressServiceCallSetup{
				want: true,
				expectation: func(mock *servicemock.AddressService) {
					mock.On("GetAddressOwnership", testifyMock.AnythingOfType("*context.valueCtx"), testValidAddress, constants.ChainBitcoin).
						Return(&service.GetAddressOwnershipResponse{
							AddressOwnershipProof: &models.AddressOwnershipProof{
								Ownership: &models.AddressOwnership{
									ID:      testRegistrationID,
									Address: testValidAddress,
									Chain:   constants.ChainBitcoin,
								},
							},
							Vasp: &models.VASP{
								ID: testVASPUUID,
							},
						}, nil)
				},
			},
		},
		"failure, invalid address": {
			in: &pb.GetAddressOwnershipRequest{
				Address: testInvalidString,
				Chain:   constants.ChainBitcoin,
			},
			assert: func(t *testing.T, res *pb.GetAddressOwnershipResponse, err error) {
				t.Helper()
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.InvalidArgument))
				assert.Equal(t, "invalid address", err.Error())
				assert.Nil(t, res)
			},
		},

		"failure, invalid chain": {
			in: &pb.GetAddressOwnershipRequest{
				Address: testValidAddress,
				Chain:   testInvalidString,
			},
			assert: func(t *testing.T, res *pb.GetAddressOwnershipResponse, err error) {
				t.Helper()
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.InvalidArgument))
				assert.Equal(t, "invalid chain: &lt;script&gt;alert(1)&lt;/script&gt;", err.Error())
				assert.Nil(t, res)
			},
		},
		"failure, unexpected internal error, should map to grpc internal error": {
			in: &pb.GetAddressOwnershipRequest{
				Address: testValidAddress,
				Chain:   constants.ChainBitcoin,
			},
			assert: func(t *testing.T, res *pb.GetAddressOwnershipResponse, err error) {
				t.Helper()
				assert.EqualError(t, err, "address service error: "+testErrorMessage)
				assert.Nil(t, res)
			},
			addressServiceCall: addressServiceCallSetup{
				want: true,
				expectation: func(mock *servicemock.AddressService) {
					mock.On("GetAddressOwnership", testifyMock.AnythingOfType("*context.valueCtx"), testValidAddress, constants.ChainBitcoin).
						Return(nil, fmt.Errorf("%w", errTest))
				},
			},
		},
		"failure, conflict error, should map to grpc already exists error": {
			in: &pb.GetAddressOwnershipRequest{
				Address: testValidAddress,
				Chain:   constants.ChainBitcoin,
			},
			assert: func(t *testing.T, res *pb.GetAddressOwnershipResponse, err error) {
				t.Helper()
				assert.True(t, serviceErrors.EqualErrorCode(errors.Unwrap(err), serviceErrors.AlreadyExists))
				assert.Nil(t, res)
			},
			addressServiceCall: addressServiceCallSetup{
				want: true,
				expectation: func(mock *servicemock.AddressService) {
					mock.On("GetAddressOwnership", testifyMock.AnythingOfType("*context.valueCtx"), testValidAddress, constants.ChainBitcoin).
						Return(nil, serviceErrors.New(testErrorMessage, serviceErrors.AlreadyExists))
				},
			},
		},
		"failure, not found error, should map to grpc not found error": {
			in: &pb.GetAddressOwnershipRequest{
				Address: testValidAddress,
				Chain:   constants.ChainBitcoin,
			},
			assert: func(t *testing.T, res *pb.GetAddressOwnershipResponse, err error) {
				t.Helper()
				assert.True(t, serviceErrors.EqualErrorCode(errors.Unwrap(err), serviceErrors.NotFound))
				assert.Nil(t, res)
			},
			addressServiceCall: addressServiceCallSetup{
				want: true,
				expectation: func(mock *servicemock.AddressService) {
					mock.On("GetAddressOwnership", testifyMock.AnythingOfType("*context.valueCtx"), testValidAddress, constants.ChainBitcoin).
						Return(nil, serviceErrors.New(testErrorMessage, serviceErrors.NotFound))
				},
			},
		},
		"failure, not allowed error, should map to grpc not allowed error": {
			in: &pb.GetAddressOwnershipRequest{
				Address: testValidAddress,
				Chain:   constants.ChainBitcoin,
			},
			assert: func(t *testing.T, res *pb.GetAddressOwnershipResponse, err error) {
				t.Helper()
				assert.True(t, serviceErrors.EqualErrorCode(errors.Unwrap(err), serviceErrors.PermissionDenied))
				assert.Nil(t, res)
			},
			addressServiceCall: addressServiceCallSetup{
				want: true,
				expectation: func(mock *servicemock.AddressService) {
					mock.On("GetAddressOwnership", testifyMock.AnythingOfType("*context.valueCtx"), testValidAddress, constants.ChainBitcoin).
						Return(nil, serviceErrors.New(testErrorMessage, serviceErrors.PermissionDenied))
				},
			},
		},
	}

	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			// given
			te := NewServiceTester()
			if tt.addressServiceCall.want {
				tt.addressServiceCall.expectation(te.mockAddressService)
			}

			// when
			res, err := te.underTesting.GetAddressOwnership(context.Background(), tt.in)

			// then
			if tt.assert != nil {
				tt.assert(t, res, err)
			}
			te.AssertAll(t)
		})
	}
}

func TestServer_DeleteAddressOwnership(t *testing.T) {
	tests := map[string]struct {
		in                 *pb.DeleteAddressOwnershipRequest
		assert             func(t *testing.T, res *pb.DeleteAddressOwnershipResponse, err error)
		addressServiceCall addressServiceCallSetup
	}{
		"success": {
			in: &pb.DeleteAddressOwnershipRequest{
				Address: testValidAddress,
				Chain:   constants.ChainBitcoin,
			},
			assert: func(t *testing.T, res *pb.DeleteAddressOwnershipResponse, err error) {
				t.Helper()
				assert.Equal(t, "success", res.Message)
				assert.NoError(t, err)
			},
			addressServiceCall: addressServiceCallSetup{
				want: true,
				expectation: func(mock *servicemock.AddressService) {
					mock.On("DeleteAddressOwnership", testifyMock.AnythingOfType("*context.valueCtx"), testValidAddress, constants.ChainBitcoin).
						Return(nil)
				},
			},
		},
		"failure, invalid address": {
			in: &pb.DeleteAddressOwnershipRequest{
				Address: testInvalidString,
				Chain:   constants.ChainBitcoin,
			},
			assert: func(t *testing.T, res *pb.DeleteAddressOwnershipResponse, err error) {
				t.Helper()
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.InvalidArgument))
				assert.Equal(t, "invalid address", err.Error())
				assert.Nil(t, res)
			},
		},

		"failure, invalid chain": {
			in: &pb.DeleteAddressOwnershipRequest{
				Address: testValidAddress,
				Chain:   testInvalidString,
			},
			assert: func(t *testing.T, res *pb.DeleteAddressOwnershipResponse, err error) {
				t.Helper()
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.InvalidArgument))
				assert.Equal(t, "invalid chain: &lt;script&gt;alert(1)&lt;/script&gt;", err.Error())
				assert.Nil(t, res)
			},
		},
		"failure, unexpected internal error, should map to grpc internal error": {
			in: &pb.DeleteAddressOwnershipRequest{
				Address: testValidAddress,
				Chain:   constants.ChainBitcoin,
			},
			assert: func(t *testing.T, res *pb.DeleteAddressOwnershipResponse, err error) {
				t.Helper()
				assert.EqualError(t, err, "address service error: "+testErrorMessage)
				assert.Equal(t, "failure", res.Message)
			},
			addressServiceCall: addressServiceCallSetup{
				want: true,
				expectation: func(mock *servicemock.AddressService) {
					mock.On("DeleteAddressOwnership", testifyMock.AnythingOfType("*context.valueCtx"), testValidAddress, constants.ChainBitcoin).
						Return(fmt.Errorf("%w", errTest))
				},
			},
		},
	}

	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			// given
			te := NewServiceTester()
			if tt.addressServiceCall.want {
				tt.addressServiceCall.expectation(te.mockAddressService)
			}

			// when
			res, err := te.underTesting.DeleteAddressOwnership(context.Background(), tt.in)

			// then
			if tt.assert != nil {
				tt.assert(t, res, err)
			}
			te.AssertAll(t)
		})
	}
}
