package rpcimpl

import (
	"context"
	"fmt"

	pb "bulletin-board-api/gen/go"
	serviceErrors "bulletin-board-api/internal/errors"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// CreateAddressOwnership.
func (s *Server) CreateAddressOwnership(ctx context.Context, request *pb.CreateAddressOwnershipRequest) (res *pb.CreateAddressOwnershipResponse, err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "CreateAddressOwnership", tracerType)
	defer func() {
		span.Finish(tracer.WithError(err))
	}()

	if err = validateCommonRequest(request); err != nil {
		return res, err
	}
	internalReq := MapANSProtoToCreateAddressOwnershipRequest(request)
	internalResp, err := s.addressService.CreateAddressOwnership(ctx, internalReq)
	if err != nil {
		return nil, fmt.Errorf("address service error: %w", err)
	}
	return &pb.CreateAddressOwnershipResponse{
		RegistrationId: internalResp.RegistrationID.String(),
	}, nil
}

// CreateAddressOwnershipProof.
func (s *Server) CreateAddressOwnershipProof(ctx context.Context, request *pb.CreateAddressOwnershipProofRequest) (res *pb.CreateAddressOwnershipProofResponse, err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "CreateAddressOwnershipProof", tracerType)
	defer func() {
		span.Finish(tracer.WithError(err))
	}()

	if err = validateCreateAddressOwnershipProofRequest(request); err != nil {
		return res, err
	}
	if err = validateCommonRequest(request); err != nil {
		return res, err
	}

	if !request.Iou {
		if err = validateChainSupportsProof(request.Chain); err != nil {
			return res, err
		}
	}

	internalRequest, err := MapCAOPRequestProtoToUpdateAddressOwnershipProofRequest(request)
	if err != nil {
		return res, serviceErrors.Wrap(err, "invalid request", serviceErrors.InvalidArgument)
	}
	internalResponse, err := s.addressService.UpdateAddressOwnershipProof(ctx, internalRequest)
	if err != nil {
		return nil, fmt.Errorf("address service error: %w", err)
	}

	proof, err := internalResponse.AddressOwnershipProof.ToProtobuf()
	if err != nil {
		return res, fmt.Errorf("address mapper error: %w", err)
	}
	return &pb.CreateAddressOwnershipProofResponse{RegistrationId: proof.Id}, nil
}

// GetAddressOwnership get ownership of an address.
func (s *Server) GetAddressOwnership(ctx context.Context, request *pb.GetAddressOwnershipRequest) (res *pb.GetAddressOwnershipResponse, err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "GetAddressOwnership", tracerType)
	defer func() {
		span.Finish(tracer.WithError(err))
	}()

	if err = validateCommonRequest(request); err != nil {
		return res, err
	}
	response, err := s.addressService.GetAddressOwnership(ctx, request.GetAddress(), request.GetChain())
	if err != nil {
		return nil, fmt.Errorf("address service error: %w", err)
	}

	return MapGetAddressOwnershipResponseToGAOProto(response)
}

// DeleteAddressOwnership deletes an address ownership and also deletes the address ownership proof.
func (s *Server) DeleteAddressOwnership(ctx context.Context, request *pb.DeleteAddressOwnershipRequest) (res *pb.DeleteAddressOwnershipResponse, err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "DeleteAddressOwnership", tracerType)
	defer func() {
		span.Finish(tracer.WithError(err))
	}()

	if err = validateCommonRequest(request); err != nil {
		return res, err
	}
	err = s.addressService.DeleteAddressOwnership(ctx, request.GetAddress(), request.GetChain())
	if err != nil {
		return &pb.DeleteAddressOwnershipResponse{Message: "failure"}, fmt.Errorf("address service error: %w", err)
	}
	return &pb.DeleteAddressOwnershipResponse{Message: "success"}, nil
}
