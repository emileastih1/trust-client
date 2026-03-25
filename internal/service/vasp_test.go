package service_test

import (
	"bulletin-board-api/internal/models"
	"context"
	"testing"

	serviceErrors "bulletin-board-api/internal/errors"

	repoMocks "bulletin-board-api/internal/repository/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type vaspRepoCallSetup struct {
	want        bool
	expectation func(t *testing.T, mock *repoMocks.VaspRepository)
}

func TestVaspService_GetVasp(t *testing.T) {
	// setup
	testVASPUUID := uuid.New()

	tests := map[string]struct {
		in           uuid.UUID
		assert       func(t *testing.T, res *models.VASP, err error)
		vaspRepoCall vaspRepoCallSetup
	}{
		"success": {
			in: testVASPUUID,
			assert: func(t *testing.T, res *models.VASP, err error) {
				t.Helper()
				assert.NoError(t, err)
				assert.NotEmpty(t, res)
				assert.Equal(t, testVASPUUID, res.ID)
				assert.Equal(t, testDomain, res.Domain)
				assert.Equal(t, "pii_endpoint", res.PIIEndpoint)
				assert.Equal(t, "pii_request_endpoint", res.PIIRequestEndpoint)
				assert.Equal(t, "public_key", res.PublicKey)
				assert.Equal(t, "lei", res.LEI)
			},
			vaspRepoCall: vaspRepoCallSetup{
				want: true,
				expectation: func(t *testing.T, mock *repoMocks.VaspRepository) {
					t.Helper()
					mock.On("Get", context.Background(), testVASPUUID).
						Return(&models.VASP{
							ID:                 testVASPUUID,
							Domain:             testDomain,
							PIIEndpoint:        "pii_endpoint",
							PIIRequestEndpoint: "pii_request_endpoint",
							PublicKey:          "public_key",
							LEI:                "lei",
						}, nil)
				},
			},
		},
		"failure, vasp not found": {
			in: testVASPUUID,
			assert: func(t *testing.T, res *models.VASP, err error) {
				t.Helper()
				assert.Empty(t, res)
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.NotFound))
			},
			vaspRepoCall: vaspRepoCallSetup{
				want: true,
				expectation: func(t *testing.T, mock *repoMocks.VaspRepository) {
					t.Helper()
					mock.On("Get", context.Background(), testVASPUUID).
						Return(nil, nil)
				},
			},
		},
	}

	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			// given
			te := NewVaspServiceTester()
			if tt.vaspRepoCall.want {
				tt.vaspRepoCall.expectation(t, te.mockVaspRepo)
			}
			// when
			res, err := te.underTesting.GetVasp(context.Background(), tt.in)

			// then
			if tt.assert != nil {
				tt.assert(t, res, err)
			}
			te.AssertAll(t)
		})
	}
}

func TestVaspService_GetVasps(t *testing.T) {
	// setup
	testVASPUUID := uuid.New()
	testVasp := &models.VASP{
		ID:     testVASPUUID,
		Domain: testDomain,
	}
	tests := map[string]struct {
		assert       func(t *testing.T, res map[string]*models.VASP, err error)
		vaspRepoCall vaspRepoCallSetup
	}{
		"success": {
			assert: func(t *testing.T, res map[string]*models.VASP, err error) {
				t.Helper()
				assert.NoError(t, err)
				assert.NotEmpty(t, res)
				assert.Equal(t, 1, len(res))
			},
			vaspRepoCall: vaspRepoCallSetup{
				want: true,
				expectation: func(t *testing.T, mock *repoMocks.VaspRepository) {
					t.Helper()
					mock.On("GetAll", context.Background()).
						Return(map[string]*models.VASP{testVASPUUID.String(): testVasp}, nil)
				},
			},
		},
	}

	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			// given
			te := NewVaspServiceTester()
			if tt.vaspRepoCall.want {
				tt.vaspRepoCall.expectation(t, te.mockVaspRepo)
			}
			// when
			res, err := te.underTesting.GetVasps(context.Background())

			// then
			if tt.assert != nil {
				tt.assert(t, res, err)
			}
			te.AssertAll(t)
		})
	}
}

func TestVaspService_SearchVasp(t *testing.T) {
	// setup
	testVASPUUID := uuid.New()

	tests := map[string]struct {
		in           string
		assert       func(t *testing.T, res *models.VASP, err error)
		vaspRepoCall vaspRepoCallSetup
	}{
		"success": {
			in: testDomain,
			assert: func(t *testing.T, res *models.VASP, err error) {
				t.Helper()
				assert.NoError(t, err)
				assert.NotEmpty(t, res)
				assert.Equal(t, testVASPUUID, res.ID)
				assert.Equal(t, testDomain, res.Domain)
			},
			vaspRepoCall: vaspRepoCallSetup{
				want: true,
				expectation: func(t *testing.T, mock *repoMocks.VaspRepository) {
					t.Helper()
					mock.On("GetByDomain", context.Background(), testDomain).
						Return(&models.VASP{
							ID:     testVASPUUID,
							Domain: testDomain,
						})
				},
			},
		},
		"vasp not found with given domain": {
			in: testDomain,
			assert: func(t *testing.T, res *models.VASP, err error) {
				t.Helper()
				assert.Empty(t, res)
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.NotFound))
			},
			vaspRepoCall: vaspRepoCallSetup{
				want: true,
				expectation: func(t *testing.T, mock *repoMocks.VaspRepository) {
					t.Helper()
					mock.On("GetByDomain", context.Background(), testDomain).
						Return(nil)
				},
			},
		},
	}

	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			// given
			te := NewVaspServiceTester()
			if tt.vaspRepoCall.want {
				tt.vaspRepoCall.expectation(t, te.mockVaspRepo)
			}
			// when
			res, err := te.underTesting.SearchVasp(context.Background(), tt.in)

			// then
			if tt.assert != nil {
				tt.assert(t, res, err)
			}
			te.AssertAll(t)
		})
	}
}
