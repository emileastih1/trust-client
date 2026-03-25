package service

import (
	"bulletin-board-api/internal/repository"
)

func MapUpdateAddressOwnershipProofRequestToGetRequest(request *UpdateAddressOwnershipProofRequest) *repository.AddressRepoGetRequest {
	return &repository.AddressRepoGetRequest{
		Address: request.Address,
		Chain:   request.Chain,
		ID:      &request.RegistrationID,
	}
}
