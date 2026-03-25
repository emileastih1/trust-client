package service

import (
	"context"

	"go.uber.org/zap"
)

/*
*****************

	Interface

*****************.
*/
type UtilityService interface {
	// Health Check
	Health(ctx context.Context) (string, error)
}

/*
*****************

	Implementation

*****************.
*/
type utilityService struct {
	logger *zap.Logger
}

func NewUtilityService(
	logger *zap.Logger,
) UtilityService {
	return &utilityService{
		logger: logger,
	}
}

func (s *utilityService) Health(_ context.Context) (string, error) {
	return "healthy", nil
}
