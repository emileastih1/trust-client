package repository

import (
	"bulletin-board-api/internal/models"
	"context"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

/*
*****************

	Interface

*****************.
*/
type VaspRepository interface {
	Get(ctx context.Context, vaspID uuid.UUID) *models.VASP
	GetAll(ctx context.Context) map[string]*models.VASP
	GetByDomain(ctx context.Context, domain string) *models.VASP
}

/*
*****************

	Implementation

*****************.
*/
type vaspRepository struct {
	vasps         map[string]*models.VASP
	domainToVasps map[string]*models.VASP
}

func NewVaspRepository() VaspRepository {
	vaspsValues := make(map[string]*models.VASP)
	if err := viper.UnmarshalKey("vasps", &vaspsValues); err != nil {
		panic(err)
	}

	if len(vaspsValues) == 0 {
		return &vaspRepository{
			vasps:         vaspsValues,
			domainToVasps: vaspsValues,
		}
	}

	for key := range vaspsValues {
		vaspUUID, err := uuid.Parse(key)
		if err != nil {
			panic(err)
		}
		vaspsValues[key].ID = vaspUUID
	}

	domainToVaspValues := make(map[string]*models.VASP)
	for _, value := range vaspsValues {
		domainToVaspValues[value.Domain] = value
	}

	return &vaspRepository{
		vasps:         vaspsValues,
		domainToVasps: domainToVaspValues,
	}
}

func (s *vaspRepository) GetByDomain(_ context.Context, domain string) (res *models.VASP) {
	if vasp, ok := s.domainToVasps[domain]; ok {
		return vasp
	}
	return nil
}

func (s *vaspRepository) Get(_ context.Context, vaspID uuid.UUID) *models.VASP {
	if vasp, ok := s.vasps[vaspID.String()]; ok {
		return vasp
	}
	return nil
}

func (s *vaspRepository) GetAll(_ context.Context) map[string]*models.VASP {
	return s.vasps
}
