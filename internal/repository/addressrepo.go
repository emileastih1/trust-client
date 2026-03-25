package repository

import (
	"bulletin-board-api/internal/models"
	"context"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

const CreateBatchSize = 25

type ProofTypeStat struct {
	Chain     string
	ProofType string
	Count     int64
}

/*
*****************

	Interface

*****************.
*/
type AddressRepository interface {
	GetOne(ctx context.Context, request *AddressRepoGetRequest) (*models.AddressOwnershipProof, error)
	UpdateOne(ctx context.Context, proof *models.AddressOwnershipProof) error
	DeleteOne(ctx context.Context, proof *models.AddressOwnershipProof) error
	CreateOneOwnership(ctx context.Context, ownership *models.AddressOwnership) error

	FindByEmptyVaspUUID(ctx context.Context, limit int) ([]*models.AddressOwnershipProof, error)
	CountEmptyVaspUUID(ctx context.Context) (int64, error)
	UpdateOwnership(_ context.Context, proof *models.AddressOwnershipProof) error

	/*
		Batch APIs are not supported for now, will re-evaluate in the future and build if needed.
		UpdateMany(ctx context.Context, proofs []*models.AddressOwnershipProof) error
		GetMany(ctx context.Context, requests []*AddressRepoGetRequest) ([]*models.AddressOwnershipProof, error)
		DeleteMany(ctx context.Context, proofs []*models.AddressOwnershipProof) error
	*/
}

/*
*****************

	Request/Response Structures

*****************.
*/
type AddressRepoGetRequest struct {
	Address string
	Chain   string
	ID      *uuid.UUID
}

type AddressRepoCreateRequest struct {
	AddressOwnershipProof *models.AddressOwnershipProof
}

/*
*****************

	Implementation

*****************.
*/
type addressRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewAddressRepository(db *gorm.DB, logger *zap.Logger) AddressRepository {
	return &addressRepository{
		db:     db.Session(&gorm.Session{}),
		logger: logger,
	}
}

func (r *addressRepository) CreateOneOwnership(_ context.Context, ownership *models.AddressOwnership) error {
	return r.db.Create(&ownership).Error
}

func (r *addressRepository) DeleteOne(_ context.Context, proof *models.AddressOwnershipProof) error {
	tx := r.db.Begin()
	err := tx.
		Where("address = ?", proof.Ownership.Address).
		Where("chain = ?", proof.Ownership.Chain).
		Delete(&proof).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *addressRepository) GetOne(_ context.Context, request *AddressRepoGetRequest) (res *models.AddressOwnershipProof, err error) {
	tx := r.db.Clauses(dbresolver.Write).Where("address = ?", request.Address).
		Where("chain = ?", request.Chain)
	if request.ID != nil {
		tx.Where("id = ?", request.ID)
	}
	if err = tx.First(&res).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return res, nil
}

func (r *addressRepository) UpdateOwnership(_ context.Context, proof *models.AddressOwnershipProof) (err error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	err = tx.Model(&proof).
		Select("encrypted_vasp_uuid").
		Updates(map[string]interface{}{
			"encrypted_vasp_uuid": proof.Ownership.EncryptedVaspUUID,
		}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *addressRepository) UpdateOne(_ context.Context, proof *models.AddressOwnershipProof) (err error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	err = tx.Model(&proof).
		Select("prefix", "signature", "iou", "proof_submitted_time", "proof_type", "aux_proof_data").
		Updates(map[string]interface{}{
			"prefix":               proof.Prefix,
			"signature":            proof.Signature,
			"iou":                  proof.IOU,
			"proof_submitted_time": proof.ProofSubmittedTime,
			"proof_type":           proof.ProofType,
			"aux_proof_data":       proof.AuxProofData,
		}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *addressRepository) FindByEmptyVaspUUID(_ context.Context, limit int) (res []*models.AddressOwnershipProof, err error) {
	err = r.db.Clauses(dbresolver.Read).Where("encrypted_vasp_uuid IS NULL").Limit(limit).Find(&res).Error

	return res, err
}

func (r *addressRepository) CountEmptyVaspUUID(_ context.Context) (int64, error) {
	var count int64
	err := r.db.Clauses(dbresolver.Read).Model(&models.AddressOwnershipProof{}).Where("encrypted_vasp_uuid IS NULL").Count(&count).Error
	return count, err
}
