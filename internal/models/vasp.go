package models

import (
	pb "bulletin-board-api/gen/go"

	uuid "github.com/google/uuid"
)

type VASP struct {
	ID                             uuid.UUID
	Name                           string `mapstructure:"name"`
	Domain                         string `mapstructure:"domain"`
	PIIEndpoint                    string `mapstructure:"payload_endpoint"`
	PIIRequestEndpoint             string `mapstructure:"pii_request_endpoint"`
	PublicKey                      string `mapstructure:"public_key"`
	PublicKeyVersion               int32  `mapstructure:"public_key_version"`
	LEI                            string `mapstructure:"lei"`
	ReturnAddressEndpoint          string `mapstructure:"return_address_endpoint"`
	ReturnFundConfirmationEndpoint string `mapstructure:"return_fund_confirmation_endpoint"`
}

func (v *VASP) ToProtobuf() *pb.VASP {
	if v == nil {
		return nil
	}

	return &pb.VASP{
		Id:                 v.ID.String(),
		Name:               v.Name,
		Domain:             v.Domain,
		PiiEndpoint:        v.PIIEndpoint,
		PiiRequestEndpoint: v.PIIRequestEndpoint,
		PublicKey:          v.PublicKey,
		PublicKeyVersion:   v.PublicKeyVersion,
		Lei:                v.LEI,
	}
}
