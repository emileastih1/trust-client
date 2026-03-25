package service

import (
	"bulletin-board-api/internal/models"
	"bulletin-board-api/internal/repository"
	"context"

	serviceErrors "bulletin-board-api/internal/errors"

	"github.com/google/uuid"

	"go.uber.org/zap"
)

/*
*****************

	Interface

*****************.
*/
type VaspService interface {
	GetVasp(ctx context.Context, vaspID uuid.UUID) (*models.VASP, error)
	GetVasps(ctx context.Context) (map[string]*models.VASP, error)
	SearchVasp(ctx context.Context, domain string) (*models.VASP, error)
}

/*
*****************

	Implementation

*****************.
*/
type vaspService struct {
	vaspRepo repository.VaspRepository
	logger   *zap.Logger
}

func NewVaspService(
	vaspRepo repository.VaspRepository,
	logger *zap.Logger,
) VaspService {
	return &vaspService{
		vaspRepo: vaspRepo,
		logger:   logger,
	}
}

func (s *vaspService) GetVasp(ctx context.Context, vaspID uuid.UUID) (resp *models.VASP, err error) {
	if res := s.vaspRepo.Get(ctx, vaspID); res != nil {
		return res, nil
	}
	return nil, serviceErrors.New("vasp not found", serviceErrors.NotFound)
}

func (s *vaspService) GetVasps(ctx context.Context) (vasps map[string]*models.VASP, err error) {
	return s.vaspRepo.GetAll(ctx), nil
}

func (s *vaspService) SearchVasp(ctx context.Context, domain string) (*models.VASP, error) {
	if res := s.vaspRepo.GetByDomain(ctx, domain); res != nil {
		return res, nil
	}
	return nil, serviceErrors.New("vasp not found", serviceErrors.NotFound)
}
