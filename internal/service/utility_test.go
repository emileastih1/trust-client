package service_test

import (
	"bulletin-board-api/internal/service"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"
)

func TestUtilityService_Health(t *testing.T) {
	// given
	utilService := service.NewUtilityService(zap.NewNop())

	// when
	resp, err := utilService.Health(context.Background())

	// then
	assert.NoError(t, err)
	assert.Equal(t, resp, "healthy")
}
