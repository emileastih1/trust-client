package rpcimpl_test

import (
	"bulletin-board-api/internal/models"
	"context"
	"errors"
	"testing"

	pb "bulletin-board-api/gen/go"
	serviceErrors "bulletin-board-api/internal/errors"

	vaspServiceMock "bulletin-board-api/internal/service/mocks"

	"github.com/stretchr/testify/assert"

	"github.com/google/uuid"
)

type vaspServiceCallSetup struct {
	want        bool
	expectation func(mock *vaspServiceMock.VaspService)
}

func TestServer_GetVasp(t *testing.T) {
	testVASPID := uuid.New()
	testVASP := &models.VASP{
		ID:   testVASPID,
		Name: "test vasp name",
	}
	tests := map[string]struct {
		in           *pb.GetVaspRequest
		assert       func(t *testing.T, res *pb.GetVaspResponse, err error)
		vaspRepoCall vaspServiceCallSetup
	}{
		"success": {
			in: &pb.GetVaspRequest{VaspId: testVASPID.String()},
			assert: func(t *testing.T, res *pb.GetVaspResponse, err error) {
				t.Helper()
				assert.NoError(t, err)
				assert.Equal(t, res.Vasp.Id, testVASPID.String())
			},
			vaspRepoCall: vaspServiceCallSetup{
				want: true,
				expectation: func(mock *vaspServiceMock.VaspService) {
					mock.On("GetVasp", context.Background(), testVASPID).
						Return(testVASP, nil)
				},
			},
		},
		"invalid vasp uuid": {
			in: &pb.GetVaspRequest{VaspId: "invalid"},
			assert: func(t *testing.T, res *pb.GetVaspResponse, err error) {
				t.Helper()
				assert.True(t, serviceErrors.EqualErrorCode(err, serviceErrors.InvalidArgument))
				assert.Nil(t, res)
			},
		},
		"vasp service error": {
			in: &pb.GetVaspRequest{VaspId: testVASPID.String()},
			assert: func(t *testing.T, res *pb.GetVaspResponse, err error) {
				t.Helper()
				assert.True(t, serviceErrors.EqualErrorCode(errors.Unwrap(err), serviceErrors.Internal))
				assert.Nil(t, res)
			},
			vaspRepoCall: vaspServiceCallSetup{
				want: true,
				expectation: func(mock *vaspServiceMock.VaspService) {
					mock.On("GetVasp", context.Background(), testVASPID).
						Return(nil, serviceErrors.New(testErrorMessage, serviceErrors.Internal))
				},
			},
		},
	}

	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			// given
			te := NewServiceTester()
			if tt.vaspRepoCall.want {
				tt.vaspRepoCall.expectation(te.mockVaspService)
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
