package service_test

import (
	"bulletin-board-api/internal/constants"
	encrypterMock "bulletin-board-api/internal/encryption/mocks"
	"bulletin-board-api/internal/lib"
	"bulletin-board-api/internal/models"
	repo "bulletin-board-api/internal/repository"
	"bulletin-board-api/internal/service"
	"context"
	"errors"
	"fmt"
	"testing"

	serviceErrors "bulletin-board-api/internal/errors"

	addressRepoMock "bulletin-board-api/internal/repository/mocks"

	vaspServiceMock "bulletin-board-api/internal/service/mocks"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
)

var (
	errTest           = errors.New(testErrorMessage)
	testInvalidString = "<script>alert(1)</script>"
	testEncryptedUUID = "encryptedUUID"
)

type vaspCallSetup struct {
	want        bool
	expectation func(mock *vaspServiceMock.VaspService)
}

type addressRepoCallSetup struct {
	want        bool
	expectation func(mock *addressRepoMock.AddressRepository)
}

type encrypterCallSetup struct {
	want        bool
	expectation func(mock *encrypterMock.MockEncrypter)
}

func TestAddressService_GetAddressOwnership(t *testing.T) {
	// setup
	testOwningVASPUUID := uuid.New()
	testRequestingVASP := models.VASP{
		ID:     uuid.New(),
		Name:   "test vasp",
		Domain: "test.com",
	}
	testCtx := context.WithValue(context.Background(), constants.RequestVASPCtxKey{}, testRequestingVASP)
	testVASP := models.VASP{
		ID:   testOwningVASPUUID,
		Name: "test vasp",
	}
	testRegistrationID := uuid.New()
	testProof := models.AddressOwnershipProof{
		Ownership: &models.AddressOwnership{
			ID:                testRegistrationID,
			Address:           testAddressLocation,
			Chain:             constants.ChainEthereum,
			EncryptedVaspUUID: &testEncryptedUUID,
		},
	}
	tests := map[string]struct {
		in struct {
			Chain   string
			Address string
		}
		contextOverride context.Context
		assert          func(t *testing.T, res *service.GetAddressOwnershipResponse, err error)
		vaspCall        vaspCallSetup
		addressRepoCall addressRepoCallSetup
		encypterCall    encrypterCallSetup
	}{
		"success, proof not provided, should return valid VASP and nil proof": {
			in: struct {
				Chain   string
				Address string
			}{
				Chain:   constants.ChainEthereum,
				Address: testAddressLocation,
			},
			assert: func(t *testing.T, res *service.GetAddressOwnershipResponse, err error) {
				t.Helper()
				assert.NotEmpty(t, res)
				assert.NoError(t, err)

				// nil proof
				assert.Nil(t, res.AddressOwnershipProof)

				// valid VASP
				assert.NotEmpty(t, res.Vasp)
				assert.Equal(t, testOwningVASPUUID, res.Vasp.ID)
			},
			vaspCall: vaspCallSetup{
				want: true,
				expectation: func(mock *vaspServiceMock.VaspService) {
					mock.On("GetVasp", testifyMock.AnythingOfType("*context.valueCtx"), testOwningVASPUUID).
						Return(&testVASP, nil)
				},
			},
			encypterCall: encrypterCallSetup{
				want: true,
				expectation: func(mock *encrypterMock.MockEncrypter) {
					mock.On("Decrypt", testifyMock.Anything, testEncryptedUUID).Return(testOwningVASPUUID.String(), nil)
				},
			},
			addressRepoCall: addressRepoCallSetup{
				want: true,
				expectation: func(mock *addressRepoMock.AddressRepository) {
					mock.On("GetOne", testifyMock.AnythingOfType("*context.valueCtx"),
						testifyMock.Anything).
						Return(&testProof, nil)
				},
			},
		},
		"failure, vasp is not set in context": {
			in: struct {
				Chain   string
				Address string
			}{
				Chain:   constants.ChainEthereum,
				Address: testAddressLocation,
			},
			assert: func(t *testing.T, res *service.GetAddressOwnershipResponse, err error) {
				t.Helper()
				assert.Nil(t, res)
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.PermissionDenied))
			},
			contextOverride: context.Background(),
		},
		"success, proof provided, should return valid VASP and non-nil proof": {
			in: struct {
				Chain   string
				Address string
			}{
				Chain:   constants.ChainEthereum,
				Address: testAddressLocation,
			},
			assert: func(t *testing.T, res *service.GetAddressOwnershipResponse, err error) {
				t.Helper()
				assert.NotEmpty(t, res)
				assert.NoError(t, err)

				// valid proof
				assert.NotEmpty(t, res.AddressOwnershipProof)
				assert.Equal(t, testRegistrationID, res.AddressOwnershipProof.Ownership.ID)
				assert.True(t, *res.AddressOwnershipProof.IOU)

				// valid VASP
				assert.NotEmpty(t, res.Vasp)
				assert.Equal(t, testOwningVASPUUID, res.Vasp.ID)
			},
			vaspCall: vaspCallSetup{
				want: true,
				expectation: func(mock *vaspServiceMock.VaspService) {
					mock.On("GetVasp", testifyMock.AnythingOfType("*context.valueCtx"), testOwningVASPUUID).
						Return(&testVASP, nil)
				},
			},
			encypterCall: encrypterCallSetup{
				want: true,
				expectation: func(mock *encrypterMock.MockEncrypter) {
					mock.On("Decrypt", testifyMock.Anything, testEncryptedUUID).Return(testOwningVASPUUID.String(), nil)
				},
			},
			addressRepoCall: addressRepoCallSetup{
				want: true,
				expectation: func(mock *addressRepoMock.AddressRepository) {
					proof := &models.AddressOwnershipProof{}
					_ = copier.Copy(proof, testProof)
					proof.IOU = lib.NewBoolPointer(true)
					mock.On("GetOne", testifyMock.AnythingOfType("*context.valueCtx"),
						testifyMock.Anything).
						Return(proof, nil)
				},
			},
		},
		"failure, no owning VASP, should return not found error": {
			in: struct {
				Chain   string
				Address string
			}{
				Chain:   constants.ChainEthereum,
				Address: testAddressLocation,
			},
			assert: func(t *testing.T, res *service.GetAddressOwnershipResponse, err error) {
				t.Helper()
				assert.Nil(t, res)
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.NotFound))
			},
			addressRepoCall: addressRepoCallSetup{
				want: true,
				expectation: func(mock *addressRepoMock.AddressRepository) {
					mock.On("GetOne", testifyMock.AnythingOfType("*context.valueCtx"),
						testifyMock.Anything).
						Return(nil, nil)
				},
			},
		},
		"failure, get VASP failure, should return error": {
			in: struct {
				Chain   string
				Address string
			}{
				Chain:   constants.ChainEthereum,
				Address: testAddressLocation,
			},
			assert: func(t *testing.T, res *service.GetAddressOwnershipResponse, err error) {
				t.Helper()
				assert.Nil(t, res)
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.NotFound))
			},
			vaspCall: vaspCallSetup{
				want: true,
				expectation: func(mock *vaspServiceMock.VaspService) {
					mock.On("GetVasp", testifyMock.AnythingOfType("*context.valueCtx"), testOwningVASPUUID).
						Return(nil, fmt.Errorf("%w", errTest))
				},
			},
			encypterCall: encrypterCallSetup{
				want: true,
				expectation: func(mock *encrypterMock.MockEncrypter) {
					mock.On("Decrypt", testifyMock.Anything, testEncryptedUUID).Return(testOwningVASPUUID.String(), nil)
				},
			},
			addressRepoCall: addressRepoCallSetup{
				want: true,
				expectation: func(mock *addressRepoMock.AddressRepository) {
					mock.On("GetOne", testifyMock.AnythingOfType("*context.valueCtx"),
						testifyMock.Anything).
						Return(&testProof, nil)
				},
			},
		},
		"failure, address repo failure, should return error": {
			in: struct {
				Chain   string
				Address string
			}{
				Chain:   constants.ChainEthereum,
				Address: testAddressLocation,
			},
			assert: func(t *testing.T, res *service.GetAddressOwnershipResponse, err error) {
				t.Helper()
				assert.Nil(t, res)
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.Internal))
			},
			addressRepoCall: addressRepoCallSetup{
				want: true,
				expectation: func(mock *addressRepoMock.AddressRepository) {
					mock.On("GetOne", testifyMock.AnythingOfType("*context.valueCtx"),
						testifyMock.Anything).
						Return(nil, fmt.Errorf("%w", errTest))
				},
			},
		},
	}

	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			// given
			te := NewAddressServiceTester()
			requestCtx := testCtx
			if tt.vaspCall.want {
				tt.vaspCall.expectation(te.mockVaspService)
			}
			if tt.addressRepoCall.want {
				tt.addressRepoCall.expectation(te.mockAddressRepo)
			}
			if tt.encypterCall.want {
				tt.encypterCall.expectation(te.mockEncrypter)
			}
			if tt.contextOverride != nil {
				requestCtx = tt.contextOverride
			}

			// when
			res, err := te.underTesting.GetAddressOwnership(requestCtx, tt.in.Address, tt.in.Chain)

			// then
			if tt.assert != nil {
				tt.assert(t, res, err)
			}
			te.AssertAll(t)
		})
	}
}

func TestAddressService_CreateAddressOwnership(t *testing.T) {
	// setup
	testRequestingVASP := models.VASP{
		ID:     uuid.New(),
		Name:   "test vasp",
		Domain: "test.com",
	}
	testCtx := context.WithValue(context.Background(), constants.RequestVASPCtxKey{}, testRequestingVASP)
	testRegistrationID := uuid.New()
	testExistingProof := models.AddressOwnershipProof{
		Ownership: &models.AddressOwnership{
			ID:                testRegistrationID,
			Address:           testAddressLocation,
			Chain:             constants.ChainEthereum,
			EncryptedVaspUUID: &testEncryptedUUID,
		},
	}
	tests := map[string]struct {
		in              *service.CreateAddressOwnershipRequest
		contextOverride context.Context
		assert          func(t *testing.T, res *service.CreateAddressOwnershipResponse, err error)
		addressRepoCall addressRepoCallSetup
		encrypterCall   encrypterCallSetup
	}{
		"success, create address": {
			in: &service.CreateAddressOwnershipRequest{
				Chain:   constants.ChainEthereum,
				Address: testAddressLocation,
			},
			assert: func(t *testing.T, res *service.CreateAddressOwnershipResponse, err error) {
				t.Helper()
				assert.NoError(t, err)
				assert.NotEmpty(t, res)
				assert.Equal(t, testRegistrationID, res.RegistrationID)
				assert.Equal(t, testAddressLocation, res.Address)
				assert.Equal(t, constants.ChainEthereum, res.Chain)
			},
			addressRepoCall: addressRepoCallSetup{
				want: true,
				expectation: func(mock *addressRepoMock.AddressRepository) {
					mock.On("GetOne", testifyMock.AnythingOfType("*context.valueCtx"), testifyMock.Anything).
						Return(nil, nil)
					mock.On("CreateOneOwnership", testifyMock.AnythingOfType("*context.valueCtx"), testifyMock.Anything).Return(nil).Run(func(args testifyMock.Arguments) {
						arg, _ := args.Get(1).(*models.AddressOwnership)
						arg.ID = testRegistrationID
					})
				},
			},
			encrypterCall: encrypterCallSetup{
				want: true,
				expectation: func(mock *encrypterMock.MockEncrypter) {
					mock.On("Encrypt", testifyMock.Anything, testRequestingVASP.ID.String()).Return(testEncryptedUUID, nil)
				},
			},
		},
		"success, pre-claimed address owned by the same vasp": {
			in: &service.CreateAddressOwnershipRequest{
				Chain:   constants.ChainEthereum,
				Address: testAddressLocation,
			},
			assert: func(t *testing.T, res *service.CreateAddressOwnershipResponse, err error) {
				t.Helper()
				assert.NoError(t, err)
				assert.Equal(t, testRegistrationID, res.RegistrationID)
				assert.Equal(t, testAddressLocation, res.Address)
				assert.Equal(t, constants.ChainEthereum, res.Chain)
			},
			addressRepoCall: addressRepoCallSetup{
				want: true,
				expectation: func(mock *addressRepoMock.AddressRepository) {
					mock.On("GetOne", testifyMock.AnythingOfType("*context.valueCtx"), &repo.AddressRepoGetRequest{
						Address: testAddressLocation,
						Chain:   constants.ChainEthereum,
					}).Return(&models.AddressOwnershipProof{
						Ownership: &models.AddressOwnership{
							ID:                testRegistrationID,
							Address:           testAddressLocation,
							Chain:             constants.ChainEthereum,
							EncryptedVaspUUID: &testEncryptedUUID,
						},
					}, nil)
				},
			},
			encrypterCall: encrypterCallSetup{
				want: true,
				expectation: func(mock *encrypterMock.MockEncrypter) {
					mock.On("Decrypt", testifyMock.Anything, testEncryptedUUID).Return(testRequestingVASP.ID.String(), nil)
				},
			},
		},
		"failure, vasp is not set in context": {
			in: &service.CreateAddressOwnershipRequest{
				Chain:   constants.ChainEthereum,
				Address: testAddressLocation,
			},
			assert: func(t *testing.T, res *service.CreateAddressOwnershipResponse, err error) {
				t.Helper()
				assert.Nil(t, res)
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.PermissionDenied))
			},
			contextOverride: context.Background(),
		},
		"failure, pre-claimed address owned by a different vasp, should error": {
			in: &service.CreateAddressOwnershipRequest{
				Chain:   constants.ChainEthereum,
				Address: testAddressLocation,
			},
			addressRepoCall: addressRepoCallSetup{
				want: true,
				expectation: func(mock *addressRepoMock.AddressRepository) {
					mock.On("GetOne", testifyMock.AnythingOfType("*context.valueCtx"), &repo.AddressRepoGetRequest{
						Address: testAddressLocation,
						Chain:   constants.ChainEthereum,
					}).Return(&models.AddressOwnershipProof{
						Ownership: &models.AddressOwnership{
							ID:                testRegistrationID,
							Address:           testAddressLocation,
							Chain:             constants.ChainEthereum,
							EncryptedVaspUUID: &testEncryptedUUID,
						},
					}, nil)
				},
			},
			encrypterCall: encrypterCallSetup{
				want: true,
				expectation: func(mock *encrypterMock.MockEncrypter) {
					mock.On("Decrypt", testifyMock.Anything, testEncryptedUUID).Return(uuid.New().String(), nil)
				},
			},
			assert: func(t *testing.T, res *service.CreateAddressOwnershipResponse, err error) {
				t.Helper()
				assert.Nil(t, res)
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.AlreadyExists))
			},
		},
		"failure, pre-claimed address owned by the same vasp, repo failure, should return error": {
			in: &service.CreateAddressOwnershipRequest{
				Chain:   constants.ChainEthereum,
				Address: testAddressLocation,
			},
			assert: func(t *testing.T, res *service.CreateAddressOwnershipResponse, err error) {
				t.Helper()
				assert.Nil(t, res)
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.Internal))
			},
			addressRepoCall: addressRepoCallSetup{
				want: true,
				expectation: func(mock *addressRepoMock.AddressRepository) {
					mock.On("GetOne", testifyMock.AnythingOfType("*context.valueCtx"), &repo.AddressRepoGetRequest{
						Address: testAddressLocation,
						Chain:   constants.ChainEthereum,
					}).Return(nil, fmt.Errorf("%w", errTest))
				},
			},
		},
		"failure, address not owned, address ownership call error, should error": {
			in: &service.CreateAddressOwnershipRequest{
				Chain:   constants.ChainEthereum,
				Address: testAddressLocation,
			},
			assert: func(t *testing.T, res *service.CreateAddressOwnershipResponse, err error) {
				t.Helper()
				assert.Nil(t, res)
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.Internal))
			},
			addressRepoCall: addressRepoCallSetup{
				want: true,
				expectation: func(mock *addressRepoMock.AddressRepository) {
					mock.On("GetOne", testifyMock.AnythingOfType("*context.valueCtx"), testifyMock.Anything).Return(nil, nil)
					mock.On("CreateOneOwnership", testifyMock.Anything, testifyMock.Anything).Return(fmt.Errorf("%w", errTest))
				},
			},
			encrypterCall: encrypterCallSetup{
				want: true,
				expectation: func(mock *encrypterMock.MockEncrypter) {
					mock.On("Encrypt", testifyMock.Anything, testRequestingVASP.ID.String()).Return(testEncryptedUUID, nil)
				},
			},
		},
		"failure, proof already submitted, can not re-submit, should conflict error": {
			in: &service.CreateAddressOwnershipRequest{
				Chain:   constants.ChainEthereum,
				Address: testAddressLocation,
			},
			assert: func(t *testing.T, res *service.CreateAddressOwnershipResponse, err error) {
				t.Helper()
				assert.Nil(t, res)
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.AlreadyExists))
			},
			addressRepoCall: addressRepoCallSetup{
				want: true,
				expectation: func(mock *addressRepoMock.AddressRepository) {
					mock.On("GetOne", testifyMock.AnythingOfType("*context.valueCtx"), testifyMock.Anything).
						Return(&testExistingProof, nil)
				},
			},
			encrypterCall: encrypterCallSetup{
				want: true,
				expectation: func(mock *encrypterMock.MockEncrypter) {
					mock.On("Decrypt", testifyMock.Anything, testEncryptedUUID).Return(testEncryptedUUID, nil)
				},
			},
		},
	}

	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			// given
			te := NewAddressServiceTester()
			requestCtx := testCtx
			if tt.addressRepoCall.want {
				tt.addressRepoCall.expectation(te.mockAddressRepo)
			}

			if tt.encrypterCall.want {
				tt.encrypterCall.expectation(te.mockEncrypter)
			}

			if tt.contextOverride != nil {
				requestCtx = tt.contextOverride
			}

			// when
			res, err := te.underTesting.CreateAddressOwnership(requestCtx, tt.in)

			// then
			if tt.assert != nil {
				tt.assert(t, res, err)
			}
			te.AssertAll(t)
		})
	}
}

func TestAddressService_UpdateAddressOwnershipProof(t *testing.T) {
	// setup
	testSignature := "test signature"
	testPrefix := "test prefix"
	testProofType := constants.BitcoinP2PKH
	testInvalidProofType := "invalid"
	testAuxProofData := []service.AuxProofData{}
	testRequestingVASPUUID := uuid.New()
	testAlternativeVASPUUID := uuid.New()

	testRequestingVASP := models.VASP{
		ID:     testRequestingVASPUUID,
		Name:   "test vasp",
		Domain: "test.com",
	}
	testCtx := context.WithValue(context.Background(), constants.RequestVASPCtxKey{}, testRequestingVASP)
	testRegistrationID := uuid.New()
	testExistingProof := models.AddressOwnershipProof{
		Ownership: &models.AddressOwnership{
			ID:                testRegistrationID,
			Address:           testAddressLocation,
			Chain:             constants.ChainEthereum,
			EncryptedVaspUUID: &testEncryptedUUID,
		},
	}

	tests := map[string]struct {
		in              *service.UpdateAddressOwnershipProofRequest
		contextOverride context.Context
		assert          func(t *testing.T, res *service.UpdateAddressOwnershipProofResponse, err error)
		addressRepoCall addressRepoCallSetup
		encrypterCall   encrypterCallSetup
	}{
		"success, with signature and prefix": {
			in: &service.UpdateAddressOwnershipProofRequest{
				Chain:          constants.ChainEthereum,
				IOU:            false,
				Address:        testAddressLocation,
				Signature:      &testSignature,
				Prefix:         &testPrefix,
				RegistrationID: testRegistrationID,
				ProofType:      &testProofType,
				AuxProofData:   testAuxProofData,
			},
			assert: func(t *testing.T, res *service.UpdateAddressOwnershipProofResponse, err error) {
				t.Helper()
				assert.NoError(t, err)
				assert.NotEmpty(t, res)
				assert.False(t, *res.AddressOwnershipProof.IOU)
				assert.Equal(t, testRegistrationID, res.AddressOwnershipProof.Ownership.ID)
				assert.Equal(t, testSignature, *res.AddressOwnershipProof.Signature)
				assert.Equal(t, testPrefix, *res.AddressOwnershipProof.Prefix)
			},
			encrypterCall: encrypterCallSetup{
				want: true,
				expectation: func(mock *encrypterMock.MockEncrypter) {
					mock.On("Decrypt", testifyMock.Anything, testEncryptedUUID).Return(testRequestingVASPUUID.String(), nil)
				},
			},
			addressRepoCall: addressRepoCallSetup{
				want: true,
				expectation: func(mock *addressRepoMock.AddressRepository) {
					mock.On("GetOne", testifyMock.AnythingOfType("*context.valueCtx"), testifyMock.Anything).
						Return(&testExistingProof, nil)
					mock.On("UpdateOne", testifyMock.AnythingOfType("*context.valueCtx"),
						testifyMock.MatchedBy(func(proof *models.AddressOwnershipProof) bool {
							return proof.Ownership.ID == testRegistrationID
						})).Return(nil)
				},
			},
		},
		"success, with IOU": {
			in: &service.UpdateAddressOwnershipProofRequest{
				Chain:          constants.ChainEthereum,
				IOU:            true,
				Address:        testAddressLocation,
				RegistrationID: testRegistrationID,
			},
			assert: func(t *testing.T, res *service.UpdateAddressOwnershipProofResponse, err error) {
				t.Helper()
				assert.NoError(t, err)
				assert.NotEmpty(t, res)
				assert.True(t, *res.AddressOwnershipProof.IOU)
				assert.Equal(t, testRegistrationID, res.AddressOwnershipProof.Ownership.ID)
			},
			encrypterCall: encrypterCallSetup{
				want: true,
				expectation: func(mock *encrypterMock.MockEncrypter) {
					mock.On("Decrypt", testifyMock.Anything, testEncryptedUUID).Return(testRequestingVASPUUID.String(), nil)
				},
			},
			addressRepoCall: addressRepoCallSetup{
				want: true,
				expectation: func(mock *addressRepoMock.AddressRepository) {
					mock.On("GetOne", testifyMock.AnythingOfType("*context.valueCtx"), testifyMock.Anything).
						Return(&testExistingProof, nil)
					mock.On("UpdateOne", testifyMock.AnythingOfType("*context.valueCtx"),
						testifyMock.MatchedBy(func(proof *models.AddressOwnershipProof) bool {
							return proof.Ownership.ID == testRegistrationID
						})).Return(nil)
				},
			},
		},
		"failure, existing registration not exist, should return invalid arguments": {
			in: &service.UpdateAddressOwnershipProofRequest{
				Chain:          constants.ChainEthereum,
				IOU:            false,
				Address:        testAddressLocation,
				Signature:      &testSignature,
				Prefix:         &testPrefix,
				RegistrationID: testRegistrationID,
			},
			assert: func(t *testing.T, res *service.UpdateAddressOwnershipProofResponse, err error) {
				t.Helper()
				assert.Nil(t, res)
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.InvalidArgument))
			},
			addressRepoCall: addressRepoCallSetup{
				want: true,
				expectation: func(mock *addressRepoMock.AddressRepository) {
					mock.On("GetOne", testifyMock.AnythingOfType("*context.valueCtx"), testifyMock.Anything).
						Return(nil, nil)
				},
			},
		},
		"failure, address repo get one failure": {
			in: &service.UpdateAddressOwnershipProofRequest{
				Chain:          constants.ChainEthereum,
				IOU:            false,
				Address:        testAddressLocation,
				Signature:      &testSignature,
				Prefix:         &testPrefix,
				RegistrationID: testRegistrationID,
			},
			assert: func(t *testing.T, res *service.UpdateAddressOwnershipProofResponse, err error) {
				t.Helper()
				assert.Nil(t, res)
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.Internal))
			},
			addressRepoCall: addressRepoCallSetup{
				want: true,
				expectation: func(mock *addressRepoMock.AddressRepository) {
					mock.On("GetOne", testifyMock.AnythingOfType("*context.valueCtx"), testifyMock.Anything).
						Return(nil, fmt.Errorf("%w", errTest))
				},
			},
		},
		"failure, once non-iou is set, can not go back to iou mode": {
			in: &service.UpdateAddressOwnershipProofRequest{
				Chain:          constants.ChainEthereum,
				IOU:            true,
				Address:        testAddressLocation,
				RegistrationID: testRegistrationID,
			},
			assert: func(t *testing.T, res *service.UpdateAddressOwnershipProofResponse, err error) {
				t.Helper()
				assert.Nil(t, res)
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.AlreadyExists))
				assert.Equal(t, "can not change from non-iou to iou", err.Error())
			},
			encrypterCall: encrypterCallSetup{
				want: true,
				expectation: func(mock *encrypterMock.MockEncrypter) {
					mock.On("Decrypt", testifyMock.Anything, testEncryptedUUID).Return(testRequestingVASPUUID.String(), nil)
				},
			},
			addressRepoCall: addressRepoCallSetup{
				want: true,
				expectation: func(mock *addressRepoMock.AddressRepository) {
					proof := &models.AddressOwnershipProof{}
					_ = copier.Copy(proof, testExistingProof)

					proof.IOU = lib.NewBoolPointer(false)
					proof.Signature = &testSignature
					proof.Prefix = &testPrefix
					mock.On("GetOne", testifyMock.AnythingOfType("*context.valueCtx"), testifyMock.Anything).
						Return(proof, nil)
				},
			},
		},
		"failure, vasp is not set in context": {
			in: &service.UpdateAddressOwnershipProofRequest{
				Chain:          constants.ChainEthereum,
				IOU:            true,
				Address:        testAddressLocation,
				RegistrationID: testRegistrationID,
			},
			assert: func(t *testing.T, res *service.UpdateAddressOwnershipProofResponse, err error) {
				t.Helper()
				assert.Nil(t, res)
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.PermissionDenied))
				assert.Equal(t, "vasp is not set", err.Error())
			},
			contextOverride: context.Background(),
		},
		"failure, providing proof for an unclaimed address": {
			in: &service.UpdateAddressOwnershipProofRequest{
				Chain:          constants.ChainEthereum,
				IOU:            false,
				Address:        testAddressLocation,
				Signature:      &testSignature,
				Prefix:         &testPrefix,
				RegistrationID: testRegistrationID,
			},
			assert: func(t *testing.T, res *service.UpdateAddressOwnershipProofResponse, err error) {
				t.Helper()
				assert.Nil(t, res)
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.InvalidArgument))
				assert.Equal(t, "updating proof for an unclaimed address", err.Error())
			},
			addressRepoCall: addressRepoCallSetup{
				want: true,
				expectation: func(mock *addressRepoMock.AddressRepository) {
					mock.On("GetOne", testifyMock.AnythingOfType("*context.valueCtx"), testifyMock.Anything).Return(nil, nil)
				},
			},
		},
		"failure, encrypted_vaspUUID, requesting vasp and owning vasp does not match": {
			in: &service.UpdateAddressOwnershipProofRequest{
				Chain:          constants.ChainEthereum,
				IOU:            false,
				Address:        testAddressLocation,
				Signature:      &testSignature,
				Prefix:         &testPrefix,
				RegistrationID: testRegistrationID,
			},
			assert: func(t *testing.T, res *service.UpdateAddressOwnershipProofResponse, err error) {
				t.Helper()
				assert.Nil(t, res)
				assert.Equal(t, "requesting VASP and owner VASP does not match", err.Error())
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.PermissionDenied))
			},
			addressRepoCall: addressRepoCallSetup{
				want: true,
				expectation: func(mock *addressRepoMock.AddressRepository) {
					mock.On("GetOne", testifyMock.AnythingOfType("*context.valueCtx"), testifyMock.Anything).
						Return(&testExistingProof, nil)
				},
			},
			encrypterCall: encrypterCallSetup{
				want: true,
				expectation: func(mock *encrypterMock.MockEncrypter) {
					mock.On("Decrypt", testifyMock.Anything, testEncryptedUUID).Return(testAlternativeVASPUUID.String(), nil)
				},
			},
		},
		"failure, both iou and (signature, prefix), are set": {
			in: &service.UpdateAddressOwnershipProofRequest{
				Chain:          constants.ChainEthereum,
				IOU:            true,
				Address:        testAddressLocation,
				Signature:      &testSignature,
				Prefix:         &testPrefix,
				RegistrationID: testRegistrationID,
			},
			assert: func(t *testing.T, res *service.UpdateAddressOwnershipProofResponse, err error) {
				t.Helper()
				assert.Nil(t, res)
				assert.Equal(t, "can not set proof data when IOU is true", err.Error())

				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.InvalidArgument))
			},
			addressRepoCall: addressRepoCallSetup{
				want: true,
				expectation: func(mock *addressRepoMock.AddressRepository) {
					mock.On("GetOne", testifyMock.AnythingOfType("*context.valueCtx"), testifyMock.Anything).
						Return(&testExistingProof, nil)
				},
			},
			encrypterCall: encrypterCallSetup{
				want: true,
				expectation: func(mock *encrypterMock.MockEncrypter) {
					mock.On("Decrypt", testifyMock.Anything, testEncryptedUUID).Return(testRequestingVASPUUID.String(), nil)
				},
			},
		},
		"failure, both iou and (proofType, auxProofData), are set": {
			in: &service.UpdateAddressOwnershipProofRequest{
				Chain:          constants.ChainEthereum,
				IOU:            true,
				Address:        testAddressLocation,
				RegistrationID: testRegistrationID,
				ProofType:      &testProofType,
				AuxProofData:   testAuxProofData,
			},
			assert: func(t *testing.T, res *service.UpdateAddressOwnershipProofResponse, err error) {
				t.Helper()
				assert.Nil(t, res)
				assert.Equal(t, "can not set proof data when IOU is true", err.Error())
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.InvalidArgument))
			},
			addressRepoCall: addressRepoCallSetup{
				want: true,
				expectation: func(mock *addressRepoMock.AddressRepository) {
					mock.On("GetOne", testifyMock.AnythingOfType("*context.valueCtx"), testifyMock.Anything).
						Return(&testExistingProof, nil)
				},
			},
			encrypterCall: encrypterCallSetup{
				want: true,
				expectation: func(mock *encrypterMock.MockEncrypter) {
					mock.On("Decrypt", testifyMock.Anything, testEncryptedUUID).Return(testRequestingVASPUUID.String(), nil)
				},
			},
		},
		"failure, iou false, invalid (signature, prefix), prefix nil": {
			in: &service.UpdateAddressOwnershipProofRequest{
				Chain:          constants.ChainEthereum,
				IOU:            false,
				Address:        testAddressLocation,
				Signature:      &testSignature,
				Prefix:         nil,
				RegistrationID: testRegistrationID,
				ProofType:      &testProofType,
				AuxProofData:   testAuxProofData,
			},
			assert: func(t *testing.T, res *service.UpdateAddressOwnershipProofResponse, err error) {
				t.Helper()
				assert.Nil(t, res)
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.InvalidArgument))
				assert.Equal(t, "invalid proof, prefix is empty", err.Error())
			},
			addressRepoCall: addressRepoCallSetup{
				want: true,
				expectation: func(mock *addressRepoMock.AddressRepository) {
					mock.On("GetOne", testifyMock.AnythingOfType("*context.valueCtx"), testifyMock.Anything).
						Return(&testExistingProof, nil)
				},
			},
			encrypterCall: encrypterCallSetup{
				want: true,
				expectation: func(mock *encrypterMock.MockEncrypter) {
					mock.On("Decrypt", testifyMock.Anything, testEncryptedUUID).Return(testRequestingVASPUUID.String(), nil)
				},
			},
		},
		"failure, iou false, invalid (signature, prefix), signature nil": {
			in: &service.UpdateAddressOwnershipProofRequest{
				Chain:          constants.ChainEthereum,
				IOU:            false,
				Address:        testAddressLocation,
				Signature:      nil,
				Prefix:         &testPrefix,
				RegistrationID: testRegistrationID,
				ProofType:      &testProofType,
				AuxProofData:   testAuxProofData,
			},
			assert: func(t *testing.T, res *service.UpdateAddressOwnershipProofResponse, err error) {
				t.Helper()
				assert.Nil(t, res)
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.InvalidArgument))
				assert.Equal(t, "invalid proof, signature is empty", err.Error())
			},
			addressRepoCall: addressRepoCallSetup{
				want: true,
				expectation: func(mock *addressRepoMock.AddressRepository) {
					mock.On("GetOne", testifyMock.AnythingOfType("*context.valueCtx"), testifyMock.Anything).
						Return(&testExistingProof, nil)
				},
			},
			encrypterCall: encrypterCallSetup{
				want: true,
				expectation: func(mock *encrypterMock.MockEncrypter) {
					mock.On("Decrypt", testifyMock.Anything, testEncryptedUUID).Return(testRequestingVASPUUID.String(), nil)
				},
			},
		},
		"failure, update one repo failure": {
			in: &service.UpdateAddressOwnershipProofRequest{
				Chain:          constants.ChainEthereum,
				IOU:            false,
				Address:        testAddressLocation,
				Signature:      &testSignature,
				Prefix:         &testPrefix,
				RegistrationID: testRegistrationID,
				ProofType:      &testProofType,
				AuxProofData:   testAuxProofData,
			},
			assert: func(t *testing.T, res *service.UpdateAddressOwnershipProofResponse, err error) {
				t.Helper()
				assert.Nil(t, res)
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.Internal))
				assert.Equal(t, "update address ownership failed: test error message", err.Error())
			},
			addressRepoCall: addressRepoCallSetup{
				want: true,
				expectation: func(mock *addressRepoMock.AddressRepository) {
					mock.On("GetOne", testifyMock.AnythingOfType("*context.valueCtx"), testifyMock.Anything).
						Return(&testExistingProof, nil)
					mock.On("UpdateOne", testifyMock.AnythingOfType("*context.valueCtx"),
						testifyMock.MatchedBy(func(proof *models.AddressOwnershipProof) bool {
							return proof.Ownership.ID == testRegistrationID
						})).Return(fmt.Errorf("%w", errTest))
				},
			},
			encrypterCall: encrypterCallSetup{
				want: true,
				expectation: func(mock *encrypterMock.MockEncrypter) {
					mock.On("Decrypt", testifyMock.Anything, testEncryptedUUID).Return(testRequestingVASPUUID.String(), nil)
				},
			},
		},
		"failure, proof type, invalid proof type": {
			in: &service.UpdateAddressOwnershipProofRequest{
				Chain:          constants.ChainEthereum,
				IOU:            false,
				Address:        testAddressLocation,
				Signature:      &testSignature,
				Prefix:         &testPrefix,
				RegistrationID: testRegistrationID,
				ProofType:      &testInvalidProofType,
				AuxProofData:   testAuxProofData,
			},
			assert: func(t *testing.T, res *service.UpdateAddressOwnershipProofResponse, err error) {
				t.Helper()
				assert.Nil(t, res)
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.InvalidArgument))
				assert.Equal(t, "invalid proof type and data combination", err.Error())
			},
			addressRepoCall: addressRepoCallSetup{
				want: true,
				expectation: func(mock *addressRepoMock.AddressRepository) {
					mock.On("GetOne", testifyMock.AnythingOfType("*context.valueCtx"), testifyMock.Anything).
						Return(&testExistingProof, nil)
				},
			},
			encrypterCall: encrypterCallSetup{
				want: true,
				expectation: func(mock *encrypterMock.MockEncrypter) {
					mock.On("Decrypt", testifyMock.Anything, testEncryptedUUID).Return(testRequestingVASPUUID.String(), nil)
				},
			},
		},
	}

	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			// given
			te := NewAddressServiceTester()
			requestCtx := testCtx
			if tt.addressRepoCall.want {
				tt.addressRepoCall.expectation(te.mockAddressRepo)
			}

			if tt.encrypterCall.want {
				tt.encrypterCall.expectation(te.mockEncrypter)
			}

			if tt.contextOverride != nil {
				requestCtx = tt.contextOverride
			}

			// when
			res, err := te.underTesting.UpdateAddressOwnershipProof(requestCtx, tt.in)

			// then
			if tt.assert != nil {
				tt.assert(t, res, err)
			}
			te.AssertAll(t)
		})
	}
}

func TestAddressService_UpdateAddressOwnershipProof_ProofTypeFocused(t *testing.T) {
	// setup
	testSignature := "test signature"
	testPrefix := "test prefix"
	testInvalidProofType := "invalid"
	testValidAuxProofData := service.AuxProofData{
		Type: constants.RedeemScript,
		Data: map[string]string{
			constants.Script: "javascript",
		},
	}
	testAuxProofDataList := []service.AuxProofData{testValidAuxProofData}

	testRequestingVASPUUID := uuid.New()

	testRequestingVASP := models.VASP{
		ID:     testRequestingVASPUUID,
		Name:   "test vasp",
		Domain: "test.com",
	}
	testCtx := context.WithValue(context.Background(), constants.RequestVASPCtxKey{}, testRequestingVASP)
	testRegistrationID := uuid.New()
	testExistingProof := models.AddressOwnershipProof{
		Ownership: &models.AddressOwnership{
			ID:                testRegistrationID,
			Address:           testAddressLocation,
			Chain:             constants.ChainEthereum,
			EncryptedVaspUUID: &testEncryptedUUID,
		},
	}
	tests := map[string]struct {
		in          *service.UpdateAddressOwnershipProofRequest
		expectError bool
	}{
		"valid BITCOIN_P2PKH": {
			in: &service.UpdateAddressOwnershipProofRequest{
				ProofType:    lib.GetStringPtr(constants.BitcoinP2PKH),
				AuxProofData: []service.AuxProofData{},
			},
		},
		"invalid BITCOIN_P2PKH": {
			in: &service.UpdateAddressOwnershipProofRequest{
				ProofType:    lib.GetStringPtr(constants.BitcoinP2PKH),
				AuxProofData: testAuxProofDataList,
			},
			expectError: true,
		},
		"valid BITCOIN_P2WPKH": {
			in: &service.UpdateAddressOwnershipProofRequest{
				ProofType:    lib.GetStringPtr(constants.BitcoinP2WPKH),
				AuxProofData: []service.AuxProofData{},
			},
		},
		"invalid BITCOIN_P2WPKH": {
			in: &service.UpdateAddressOwnershipProofRequest{
				ProofType:    lib.GetStringPtr(constants.BitcoinP2WPKH),
				AuxProofData: testAuxProofDataList,
			},
			expectError: true,
		},
		"valid ETHEREUM_EOA": {
			in: &service.UpdateAddressOwnershipProofRequest{
				ProofType:    lib.GetStringPtr(constants.EthereumEOA),
				AuxProofData: []service.AuxProofData{},
			},
		},
		"invalid ETHEREUM_EOA": {
			in: &service.UpdateAddressOwnershipProofRequest{
				ProofType:    lib.GetStringPtr(constants.EthereumEOA),
				AuxProofData: testAuxProofDataList,
			},
			expectError: true,
		},
		"valid TRON_EOA": {
			in: &service.UpdateAddressOwnershipProofRequest{
				ProofType:    lib.GetStringPtr(constants.TronEOA),
				AuxProofData: []service.AuxProofData{},
			},
		},
		"invalid TRON_EOA": {
			in: &service.UpdateAddressOwnershipProofRequest{
				ProofType:    lib.GetStringPtr(constants.TronEOA),
				AuxProofData: testAuxProofDataList,
			},
			expectError: true,
		},
		"valid BITCOIN_P2SH": {
			in: &service.UpdateAddressOwnershipProofRequest{
				ProofType:    lib.GetStringPtr(constants.BitcoinP2SH),
				AuxProofData: testAuxProofDataList,
			},
		},
		"valid empty data BITCOIN_P2SH": {
			in: &service.UpdateAddressOwnershipProofRequest{
				ProofType:    lib.GetStringPtr(constants.BitcoinP2SH),
				AuxProofData: []service.AuxProofData{},
			},
			expectError: false,
		},
		"invalid 1 invalid data BITCOIN_P2SH": {
			in: &service.UpdateAddressOwnershipProofRequest{
				ProofType: lib.GetStringPtr(constants.BitcoinP2SH),
				AuxProofData: []service.AuxProofData{{
					Type: constants.EthCreate,
					Data: map[string]string{
						constants.Salt: "123123",
					},
				}},
			},
			expectError: true,
		},
		"invalid 2 invalid data BITCOIN_P2SH": {
			in: &service.UpdateAddressOwnershipProofRequest{
				ProofType: lib.GetStringPtr(constants.BitcoinP2SH),
				AuxProofData: []service.AuxProofData{{
					Type: constants.EthCreate,
					Data: map[string]string{
						constants.Nonce: "123123",
					},
				}},
			},
			expectError: true,
		},
		"invalid 3 multiple redeem script data BITCOIN_P2SH": {
			in: &service.UpdateAddressOwnershipProofRequest{
				ProofType: lib.GetStringPtr(constants.BitcoinP2SH),
				AuxProofData: []service.AuxProofData{{
					Type: constants.RedeemScript,
					Data: map[string]string{
						constants.Script: "123123",
					},
				}, {
					Type: constants.RedeemScript,
					Data: map[string]string{
						constants.Script: "456456",
					},
				}},
			},
			expectError: true,
		},
		"valid BITCOIN_P2L": {
			in: &service.UpdateAddressOwnershipProofRequest{
				ProofType:    lib.GetStringPtr(constants.BitcoinP2L),
				AuxProofData: []service.AuxProofData{},
			},
		},
		"invalid BITCOIN_P2L": {
			in: &service.UpdateAddressOwnershipProofRequest{
				ProofType:    lib.GetStringPtr(constants.BitcoinP2L),
				AuxProofData: testAuxProofDataList,
			},
			expectError: true,
		},
		"valid ethereum contract": {
			in: &service.UpdateAddressOwnershipProofRequest{
				ProofType: lib.GetStringPtr(constants.EthereumContract),
				AuxProofData: []service.AuxProofData{{
					Type: constants.EthCreate,
					Data: map[string]string{
						constants.Nonce: "123123",
					},
				}, {
					Type: constants.EthCreate2,
					Data: map[string]string{
						constants.Salt:         "0x14e26bf8489b5ecf8bed66029345a1fbfce9aba797dee64db4bd6403515916db",
						constants.InitCodeHash: "0x5c5cd361ee6815781aa97a4b22dcbb49a1aed5548b57e848e34b5f3e62e3154e",
					},
				}},
			},
		},
		"invalid ethereum contract invalid eth_create2 data": {
			in: &service.UpdateAddressOwnershipProofRequest{
				ProofType: lib.GetStringPtr(constants.EthereumContract),
				AuxProofData: []service.AuxProofData{{
					Type: constants.EthCreate2,
					Data: map[string]string{
						constants.Salt: "0x14e26bf8489b5ecf8bed66029345a1fbfce9aba797dee64db4bd6403515916db",
					},
				}},
			},
			expectError: true,
		},
		"invalid ethereum contract invalid eth_create2 data 2": {
			in: &service.UpdateAddressOwnershipProofRequest{
				ProofType: lib.GetStringPtr(constants.EthereumContract),
				AuxProofData: []service.AuxProofData{{
					Type: constants.EthCreate2,
					Data: map[string]string{
						constants.InitCodeHash: "0x5c5cd361ee6815781aa97a4b22dcbb49a1aed5548b57e848e34b5f3e62e3154e",
					},
				}},
			},
			expectError: true,
		},
		"invalid ethereum contract": {
			in: &service.UpdateAddressOwnershipProofRequest{
				ProofType: lib.GetStringPtr(constants.EthereumContract),
				AuxProofData: []service.AuxProofData{{
					Type: constants.RedeemScript,
					Data: map[string]string{
						constants.Script: "123123",
					},
				}, {
					Type: constants.EthCreate2,
					Data: map[string]string{
						constants.Salt:         "0x14e26bf8489b5ecf8bed66029345a1fbfce9aba797dee64db4bd6403515916db",
						constants.InitCodeHash: "0x5c5cd361ee6815781aa97a4b22dcbb49a1aed5548b57e848e34b5f3e62e3154e",
					},
				}},
			},
			expectError: true,
		},
		"invalid empty data ethereum contract": {
			in: &service.UpdateAddressOwnershipProofRequest{
				ProofType:    lib.GetStringPtr(constants.EthereumContract),
				AuxProofData: []service.AuxProofData{},
			},
			expectError: true,
		},

		"invalid proof type": {
			in: &service.UpdateAddressOwnershipProofRequest{
				ProofType:    &testInvalidProofType,
				AuxProofData: testAuxProofDataList,
			},
			expectError: true,
		},
		"invalid proof type 2": {
			in: &service.UpdateAddressOwnershipProofRequest{
				ProofType:    &testInvalidString,
				AuxProofData: []service.AuxProofData{},
			},
			expectError: true,
		},
		"invalid proof data key": {
			in: &service.UpdateAddressOwnershipProofRequest{
				ProofType: lib.GetStringPtr(constants.EthereumContract),
				AuxProofData: []service.AuxProofData{
					{
						Type: testInvalidString,
						Data: map[string]string{
							constants.Salt:         "0x14e26bf8489b5ecf8bed66029345a1fbfce9aba797dee64db4bd6403515916db",
							constants.InitCodeHash: "0x5c5cd361ee6815781aa97a4b22dcbb49a1aed5548b57e848e34b5f3e62e3154e",
						},
					},
				},
			},
			expectError: true,
		},
	}
	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			// given
			te := NewAddressServiceTester()
			requestCtx := testCtx

			te.mockAddressRepo.On("GetOne", testifyMock.AnythingOfType("*context.valueCtx"), testifyMock.Anything).
				Return(&testExistingProof, nil)
			te.mockEncrypter.On("Decrypt", testifyMock.Anything, testEncryptedUUID).Return(testRequestingVASPUUID.String(), nil)

			if !tt.expectError {
				te.mockAddressRepo.On("UpdateOne", testifyMock.AnythingOfType("*context.valueCtx"),
					testifyMock.MatchedBy(func(proof *models.AddressOwnershipProof) bool {
						return proof.Ownership.ID == testRegistrationID
					})).Return(nil)
			}

			// when
			res, err := te.underTesting.UpdateAddressOwnershipProof(requestCtx, &service.UpdateAddressOwnershipProofRequest{
				Chain:          constants.ChainEthereum,
				IOU:            false,
				Address:        testAddressLocation,
				Signature:      &testSignature,
				Prefix:         &testPrefix,
				RegistrationID: testRegistrationID,
				ProofType:      tt.in.ProofType,
				AuxProofData:   tt.in.AuxProofData,
			})

			// then
			if tt.expectError {
				assert.Nil(t, res)
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.InvalidArgument))
				assert.Equal(t, "invalid proof type and data combination", err.Error())
			} else {
				assert.NotNil(t, res)
				assert.NoError(t, err)
			}

			te.AssertAll(t)
		})
	}
}

func TestAddressService_DeleteAddressOwnership(t *testing.T) {
	// setup
	testRequestingVASPUUID := uuid.New()

	testRequestingVASP := models.VASP{
		ID:     testRequestingVASPUUID,
		Name:   "test vasp",
		Domain: "test.com",
	}
	testCtx := context.WithValue(context.Background(), constants.RequestVASPCtxKey{}, testRequestingVASP)
	testProof := models.AddressOwnershipProof{
		Ownership: &models.AddressOwnership{
			Address:           testAddressLocation,
			Chain:             constants.ChainEthereum,
			EncryptedVaspUUID: &testEncryptedUUID,
		},
	}

	tests := map[string]struct {
		in struct {
			Address string
			Chain   string
		}
		contextOverride context.Context
		assert          func(t *testing.T, err error)
		addressRepoCall addressRepoCallSetup
		encrypterCall   encrypterCallSetup
	}{
		"success, delete one": {
			in: struct {
				Address string
				Chain   string
			}{Address: testAddressLocation, Chain: constants.ChainEthereum},
			assert: func(t *testing.T, err error) {
				t.Helper()
				assert.NoError(t, err)
			},
			encrypterCall: encrypterCallSetup{
				want: true,
				expectation: func(mock *encrypterMock.MockEncrypter) {
					mock.On("Decrypt", testifyMock.Anything, testEncryptedUUID).Return(testRequestingVASPUUID.String(), nil)
				},
			},
			addressRepoCall: addressRepoCallSetup{
				want: true,
				expectation: func(mock *addressRepoMock.AddressRepository) {
					mock.On("GetOne", testifyMock.AnythingOfType("*context.valueCtx"), testifyMock.Anything).
						Return(&testProof, nil)
					mock.On("DeleteOne", testifyMock.AnythingOfType("*context.valueCtx"), testifyMock.Anything).
						Return(nil)
				},
			},
		},
		"failure, delete one, address does not exist, should return not found": {
			in: struct {
				Address string
				Chain   string
			}{Address: testAddressLocation, Chain: constants.ChainEthereum},
			assert: func(t *testing.T, err error) {
				t.Helper()
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.NotFound))
			},
			addressRepoCall: addressRepoCallSetup{
				want: true,
				expectation: func(mock *addressRepoMock.AddressRepository) {
					mock.On("GetOne", testifyMock.AnythingOfType("*context.valueCtx"), testifyMock.Anything).Return(nil, nil)
				},
			},
		},
		"failure, requesting vasp is not set": {
			in: struct {
				Address string
				Chain   string
			}{Address: testAddressLocation, Chain: constants.ChainEthereum},
			assert: func(t *testing.T, err error) {
				t.Helper()
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.PermissionDenied))
			},
			contextOverride: context.Background(),
		},
		"failure, requesting vasp and owning vasp does not match": {
			in: struct {
				Address string
				Chain   string
			}{Address: testAddressLocation, Chain: constants.ChainEthereum},
			encrypterCall: encrypterCallSetup{
				want: true,
				expectation: func(mock *encrypterMock.MockEncrypter) {
					mock.On("Decrypt", testifyMock.Anything, testEncryptedUUID).Return(uuid.New().String(), nil)
				},
			},
			addressRepoCall: addressRepoCallSetup{
				want: true,
				expectation: func(mock *addressRepoMock.AddressRepository) {
					mock.On("GetOne", testifyMock.AnythingOfType("*context.valueCtx"), testifyMock.Anything).
						Return(&testProof, nil)
				},
			},
			assert: func(t *testing.T, err error) {
				t.Helper()
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.PermissionDenied))
			},
		},
		"failure, repo delete one error": {
			in: struct {
				Address string
				Chain   string
			}{Address: testAddressLocation, Chain: constants.ChainEthereum},
			assert: func(t *testing.T, err error) {
				t.Helper()
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.Internal))
			},
			encrypterCall: encrypterCallSetup{
				want: true,
				expectation: func(mock *encrypterMock.MockEncrypter) {
					mock.On("Decrypt", testifyMock.Anything, testEncryptedUUID).Return(testRequestingVASPUUID.String(), nil)
				},
			},
			addressRepoCall: addressRepoCallSetup{
				want: true,
				expectation: func(mock *addressRepoMock.AddressRepository) {
					mock.On("GetOne", testifyMock.AnythingOfType("*context.valueCtx"), testifyMock.Anything).
						Return(&testProof, nil)
					mock.On("DeleteOne", testifyMock.AnythingOfType("*context.valueCtx"), testifyMock.Anything).
						Return(fmt.Errorf("%w", errTest))
				},
			},
		},
	}

	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			// given
			te := NewAddressServiceTester()
			requestCtx := testCtx
			if tt.addressRepoCall.want {
				tt.addressRepoCall.expectation(te.mockAddressRepo)
			}
			if tt.encrypterCall.want {
				tt.encrypterCall.expectation(te.mockEncrypter)
			}
			if tt.contextOverride != nil {
				requestCtx = tt.contextOverride
			}

			// when
			err := te.underTesting.DeleteAddressOwnership(requestCtx, tt.in.Address, tt.in.Chain)

			// then
			if tt.assert != nil {
				tt.assert(t, err)
			}
			te.AssertAll(t)
		})
	}
}
