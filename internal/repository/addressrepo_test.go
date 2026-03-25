package repository_test

import (
	"bulletin-board-api/internal/constants"
	"bulletin-board-api/internal/lib"
	"bulletin-board-api/internal/models"
	"bulletin-board-api/internal/repository"
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var errTest = errors.New(testErrorMessage)

type addressRepoInputData struct {
	proof *models.AddressOwnershipProof
}

func TestAddressRepository_GetOne(t *testing.T) {
	testUUID := uuid.New()
	testRegistrationID := uuid.New()
	getOneSQL := `SELECT * FROM "address_ownership_proofs" WHERE 
					address = $1 AND chain = $2 AND "address_ownership_proofs"."deleted_at" IS NULL 
					ORDER BY "address_ownership_proofs"."id" LIMIT 1`
	tests := map[string]struct {
		in       *repository.AddressRepoGetRequest
		assert   func(t *testing.T, res *models.AddressOwnershipProof, err error)
		gormCall gormCallSetup
	}{
		"get one success": {
			in: &repository.AddressRepoGetRequest{
				Address: testAddressLocation,
				Chain:   constants.ChainEthereum,
			},
			assert: func(t *testing.T, res *models.AddressOwnershipProof, err error) {
				t.Helper()
				assert.Equal(t, res.Ownership.ID, testUUID)
				assert.NoError(t, err)
			},
			gormCall: gormCallSetup{
				want: true,
				expectation: func(mock sqlmock.Sqlmock) {
					mock.ExpectQuery(getOneSQL).
						WithArgs(testAddressLocation, constants.ChainEthereum).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).
							AddRow(testUUID))
				},
			},
		},
		"get one success with id": {
			in: &repository.AddressRepoGetRequest{
				Address: testAddressLocation,
				Chain:   constants.ChainEthereum,
				ID:      &testRegistrationID,
			},
			assert: func(t *testing.T, res *models.AddressOwnershipProof, err error) {
				t.Helper()
				assert.Equal(t, res.Ownership.ID, testRegistrationID)
				assert.NoError(t, err)
			},
			gormCall: gormCallSetup{
				want: true,
				expectation: func(mock sqlmock.Sqlmock) {
					sqlQuery := `SELECT * FROM "address_ownership_proofs" WHERE 
						address = $1 AND chain = $2 AND id = $3 AND "address_ownership_proofs"."deleted_at" IS NULL 
						ORDER BY "address_ownership_proofs"."id" LIMIT 1`
					mock.ExpectQuery(sqlQuery).
						WithArgs(testAddressLocation, constants.ChainEthereum, testRegistrationID.String()).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).
							AddRow(testRegistrationID))
				},
			},
		},
		"success, record not found, should return empty with no error": {
			in: &repository.AddressRepoGetRequest{
				Address: testAddressLocation,
				Chain:   constants.ChainEthereum,
			},
			assert: func(t *testing.T, res *models.AddressOwnershipProof, err error) {
				t.Helper()
				assert.Nil(t, res)
				assert.NoError(t, err)
			},
			gormCall: gormCallSetup{
				want: true,
				expectation: func(mock sqlmock.Sqlmock) {
					mock.ExpectQuery(getOneSQL).
						WithArgs(testAddressLocation, constants.ChainEthereum).
						WillReturnError(gorm.ErrRecordNotFound)
				},
			},
		},
		"failure, get one gorm failure": {
			in: &repository.AddressRepoGetRequest{
				Address: testAddressLocation,
				Chain:   constants.ChainEthereum,
			},
			assert: func(t *testing.T, res *models.AddressOwnershipProof, err error) {
				t.Helper()
				assert.Empty(t, res)
				assert.True(t, strings.Contains(err.Error(), testErrorMessage))
			},
			gormCall: gormCallSetup{
				want: true,
				expectation: func(mock sqlmock.Sqlmock) {
					mock.ExpectQuery(getOneSQL).
						WithArgs(testAddressLocation, constants.ChainEthereum).
						WillReturnError(fmt.Errorf("%w", errTest))
				},
			},
		},
	}

	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			// given
			te := NewAddressRepoTester(t)
			if tt.gormCall.want {
				tt.gormCall.expectation(te.mockGorm)
			}

			// when
			res, err := te.underTesting.GetOne(context.Background(), tt.in)

			// then
			if tt.assert != nil {
				tt.assert(t, res, err)
			}
			te.AssertAll(t)
		})
	}
}

func TestAddressRepository_CreateOne(t *testing.T) {
	testProof := &models.AddressOwnership{
		Address: testAddressLocation,
		Chain:   constants.ChainEthereum,
	}

	tests := map[string]struct {
		in struct {
			proof *models.AddressOwnership
		}
		assert   func(t *testing.T, err error)
		gormCall gormCallSetup
	}{
		"insert one success": {
			in: struct {
				proof *models.AddressOwnership
			}{
				proof: testProof,
			},
			assert: func(t *testing.T, err error) {
				t.Helper()
				assert.NoError(t, err)
			},
			gormCall: gormCallSetup{
				want: true,
				expectation: func(mock sqlmock.Sqlmock) {
					insertSQL := `INSERT INTO "address_ownership_proofs"
									("id","address","chain","encrypted_vasp_uuid","created_at","updated_at","deleted_at")
									VALUES ($1,$2,$3,$4,$5,$6,$7)`
					mock.ExpectBegin()
					mock.ExpectExec(insertSQL).
						WithArgs(AnyUUID{}, testAddressLocation, constants.ChainEthereum, nil, AnyTime{}, AnyTime{}, nil).
						WillReturnResult(sqlmock.NewResult(1, 1))
					mock.ExpectCommit()
				},
			},
		},
		"insert one gorm failure, should error": {
			in: struct {
				proof *models.AddressOwnership
			}{
				proof: testProof,
			},
			assert: func(t *testing.T, err error) {
				t.Helper()
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), testErrorMessage))
			},
			gormCall: gormCallSetup{
				want: true,
				expectation: func(mock sqlmock.Sqlmock) {
					insertSQL := `INSERT INTO "address_ownership_proofs"
									("id","address","chain","encrypted_vasp_uuid","created_at","updated_at","deleted_at")
									VALUES ($1,$2,$3,$4,$5,$6,$7)`
					mock.ExpectBegin()
					mock.ExpectExec(insertSQL).
						WithArgs(AnyUUID{}, testAddressLocation, constants.ChainEthereum, nil, AnyTime{}, AnyTime{}, nil).
						WillReturnError(fmt.Errorf("%w", errTest))
					mock.ExpectRollback()
				},
			},
		},
	}

	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			// given
			te := NewAddressRepoTester(t)
			if tt.gormCall.want {
				tt.gormCall.expectation(te.mockGorm)
			}

			// when
			err := te.underTesting.CreateOneOwnership(context.Background(), tt.in.proof)

			// then
			if tt.assert != nil {
				tt.assert(t, err)
			}

			te.AssertAll(t)
		})
	}
}

func TestAddressRepository_UpdateOne(t *testing.T) {
	testRegistrationID := uuid.New()
	testProofType := "test_proof_type"
	testProof := &models.AddressOwnershipProof{
		Ownership: &models.AddressOwnership{
			ID:      testRegistrationID,
			Address: testAddressLocation,
			Chain:   constants.ChainEthereum,
		},
		Signature:          nil,
		Prefix:             nil,
		IOU:                lib.NewBoolPointer(true),
		ProofSubmittedTime: nil,
		ProofType:          &testProofType,
		AuxProofData:       models.JSON{},
	}

	tests := map[string]struct {
		in       addressRepoInputData
		assert   func(t *testing.T, err error)
		gormCall gormCallSetup
	}{
		"update one success": {
			in: addressRepoInputData{
				proof: testProof,
			},
			assert: func(t *testing.T, err error) {
				t.Helper()
				assert.NoError(t, err)
			},
			gormCall: gormCallSetup{
				want: true,
				expectation: func(mock sqlmock.Sqlmock) {
					sqlStatement := `UPDATE "address_ownership_proofs"
										SET "aux_proof_data"=$1,"iou"=$2,"prefix"=$3,"proof_submitted_time"=$4,"proof_type"=$5,"signature"=$6,"updated_at"=$7
										WHERE "address_ownership_proofs"."deleted_at" IS NULL AND "id" = $8`
					mock.ExpectBegin()
					mock.ExpectExec(sqlStatement).
						WithArgs(nil, true, nil, nil, testProofType, nil, AnyTime{}, testRegistrationID).
						WillReturnResult(sqlmock.NewResult(1, 1))
					mock.ExpectCommit()
				},
			},
		},
		"update one gorm failure, should error": {
			in: addressRepoInputData{
				proof: testProof,
			},
			assert: func(t *testing.T, err error) {
				t.Helper()
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), testErrorMessage))
			},
			gormCall: gormCallSetup{
				want: true,
				expectation: func(mock sqlmock.Sqlmock) {
					sqlStatement := `UPDATE "address_ownership_proofs" 
										SET "aux_proof_data"=$1,"iou"=$2,"prefix"=$3,"proof_submitted_time"=$4,"proof_type"=$5,"signature"=$6,"updated_at"=$7
										WHERE "address_ownership_proofs"."deleted_at" IS NULL AND "id" = $8`
					mock.ExpectBegin()
					mock.ExpectExec(sqlStatement).
						WithArgs(nil, true, nil, nil, testProofType, nil, AnyTime{}, testRegistrationID).
						WillReturnError(fmt.Errorf("%w", errTest))
					mock.ExpectRollback()
				},
			},
		},
	}

	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			// given
			te := NewAddressRepoTester(t)
			if tt.gormCall.want {
				tt.gormCall.expectation(te.mockGorm)
			}

			// when
			err := te.underTesting.UpdateOne(context.Background(), tt.in.proof)

			// then
			if tt.assert != nil {
				tt.assert(t, err)
			}
			te.AssertAll(t)
		})
	}
}

func TestAddressRepository_DeleteOne(t *testing.T) {
	testRegistrationID := uuid.New()
	testProof := &models.AddressOwnershipProof{
		Ownership: &models.AddressOwnership{
			ID:      testRegistrationID,
			Address: testAddressLocation,
			Chain:   constants.ChainEthereum,
		},
		Signature:          nil,
		Prefix:             nil,
		IOU:                lib.NewBoolPointer(true),
		ProofSubmittedTime: nil,
	}

	tests := map[string]struct {
		in       addressRepoInputData
		assert   func(t *testing.T, err error)
		gormCall gormCallSetup
	}{
		"delete one success": {
			in: addressRepoInputData{
				proof: testProof,
			},
			assert: func(t *testing.T, err error) {
				t.Helper()
				assert.NoError(t, err)
			},
			gormCall: gormCallSetup{
				want: true,
				expectation: func(mock sqlmock.Sqlmock) {
					sqlStatement := `UPDATE "address_ownership_proofs" SET "deleted_at"=$1 
										WHERE address = $2 AND chain = $3 AND "address_ownership_proofs"."id" = $4 
										AND "address_ownership_proofs"."deleted_at" IS NULL`
					mock.ExpectBegin()
					mock.ExpectExec(sqlStatement).
						WithArgs(AnyTime{}, testAddressLocation, constants.ChainEthereum, testRegistrationID).
						WillReturnResult(sqlmock.NewResult(1, 1))
					mock.ExpectCommit()
				},
			},
		},
		"delete one gorm failure, should error": {
			in: addressRepoInputData{
				proof: testProof,
			},
			assert: func(t *testing.T, err error) {
				t.Helper()
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), testErrorMessage))
			},
			gormCall: gormCallSetup{
				want: true,
				expectation: func(mock sqlmock.Sqlmock) {
					sqlStatement := `UPDATE "address_ownership_proofs" SET "deleted_at"=$1
										WHERE address = $2 AND chain = $3 AND "address_ownership_proofs"."id" = $4
										AND "address_ownership_proofs"."deleted_at" IS NULL`
					mock.ExpectBegin()
					mock.ExpectExec(sqlStatement).
						WithArgs(AnyTime{}, testAddressLocation, constants.ChainEthereum, testRegistrationID).
						WillReturnError(fmt.Errorf("%w", errTest))
					mock.ExpectRollback()
				},
			},
		},
	}

	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			// given
			te := NewAddressRepoTester(t)
			if tt.gormCall.want {
				tt.gormCall.expectation(te.mockGorm)
			}

			// when
			err := te.underTesting.DeleteOne(context.Background(), tt.in.proof)

			// then
			if tt.assert != nil {
				tt.assert(t, err)
			}
			te.AssertAll(t)
		})
	}
}
