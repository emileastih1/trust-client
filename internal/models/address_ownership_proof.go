package models

import (
	"bulletin-board-api/internal/lib"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	pb "bulletin-board-api/gen/go"

	"google.golang.org/protobuf/types/known/wrapperspb"

	"gorm.io/gorm"

	"github.com/google/uuid"
)

var errJSONMarshal = errors.New("failed to umarshal json value")

type AddressOwnership struct {
	ID                uuid.UUID `gorm:"type:uuid;primary_key"` // registration_id
	Address           string
	Chain             string
	EncryptedVaspUUID *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt
}

type AddressOwnershipProof struct {
	Ownership          *AddressOwnership `gorm:"embedded"`
	Signature          *string
	Prefix             *string
	IOU                *bool
	ProofSubmittedTime *time.Time
	ProofType          *string
	AuxProofData       JSON `sql:"type:JSON NOT NULL DEFAULT '[]'::JSON"`
}

type JSON json.RawMessage

func (j *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errJSONMarshal
	}

	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
	if err != nil {
		return fmt.Errorf("unmarshal error %w", err)
	}
	*j = JSON(result)
	return nil
}

func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	driver, err := json.RawMessage(j).MarshalJSON()
	if err != nil {
		return driver, fmt.Errorf("gorm EncodedSecret error %w", err)
	}
	return driver, nil
}

func (a *AddressOwnership) TableName() string {
	return "address_ownership_proofs"
}

func (a *AddressOwnershipProof) TableName() string {
	return "address_ownership_proofs"
}

func (a *AddressOwnership) BeforeCreate(_ *gorm.DB) error {
	a.ID = uuid.New()
	return nil
}

func (a *AddressOwnershipProof) ToProtobuf() (*pb.AddressOwnershipProof, error) {
	if a == nil {
		return nil, nil
	}
	v := &pb.AddressOwnershipProof{}
	if a.Ownership != nil {
		v.Id = a.Ownership.ID.String()
		v.Address = a.Ownership.Address
		v.Chain = a.Ownership.Chain
	}
	if a.Signature != nil {
		v.Signature = &wrapperspb.StringValue{Value: lib.StringPtrToString(a.Signature)}
	}
	if a.Prefix != nil {
		v.Prefix = &wrapperspb.StringValue{Value: lib.StringPtrToString(a.Prefix)}
	}
	if a.IOU != nil {
		v.Iou = &wrapperspb.BoolValue{Value: lib.BoolPtrToBool(a.IOU)}
	}
	if a.ProofType != nil {
		v.ProofType = *a.ProofType
	}
	if a.AuxProofData != nil {
		var data []*pb.AuxProofData
		err := json.Unmarshal(a.AuxProofData, &data)
		if err != nil {
			return v, fmt.Errorf("error umarshal aux proof data: %w", err)
		}
		v.AuxProofData = data
	}

	return v, nil
}
