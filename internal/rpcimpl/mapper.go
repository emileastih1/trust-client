package rpcimpl

import (
	"bulletin-board-api/internal/service"
	"fmt"

	pb "bulletin-board-api/gen/go"

	"github.com/google/uuid"
)

func MapGetAddressOwnershipResponseToGAOProto(resp *service.GetAddressOwnershipResponse) (res *pb.GetAddressOwnershipResponse, err error) {
	if resp == nil {
		return nil, nil
	}
	res = &pb.GetAddressOwnershipResponse{
		Vasp: resp.Vasp.ToProtobuf(),
	}
	res.AddressOwnershipProof, err = resp.AddressOwnershipProof.ToProtobuf()
	if err != nil {
		return res, fmt.Errorf("mapper error %w", err)
	}

	return res, nil
}

func MapCAOPRequestProtoToUpdateAddressOwnershipProofRequest(request *pb.CreateAddressOwnershipProofRequest) (res *service.UpdateAddressOwnershipProofRequest, err error) {
	sig := request.GetSignature()
	prefix := request.GetPrefix()
	proofType := request.GetProofType()
	rID := uuid.MustParse(request.GetRegistrationId())

	auxProofData := make([]service.AuxProofData, 0, len(request.GetAuxProofData()))
	for _, item := range request.GetAuxProofData() {
		auxProofData = append(auxProofData, service.AuxProofData{
			Type: item.Type,
			Data: item.Data,
		})
	}

	return &service.UpdateAddressOwnershipProofRequest{
		Chain:          request.GetChain(),
		IOU:            request.GetIou(),
		Address:        request.GetAddress(),
		Signature:      &sig,
		Prefix:         &prefix,
		RegistrationID: rID,
		ProofType:      &proofType,
		AuxProofData:   auxProofData,
	}, nil
}

func MapANSProtoToCreateAddressOwnershipRequest(request *pb.CreateAddressOwnershipRequest) *service.CreateAddressOwnershipRequest {
	return &service.CreateAddressOwnershipRequest{
		Chain:   request.GetChain(),
		Address: request.GetAddress(),
	}
}
