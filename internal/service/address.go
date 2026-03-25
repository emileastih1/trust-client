//nolint:gofumpt
package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"bulletin-board-api/internal/constants"
	"bulletin-board-api/internal/encryption"
	"bulletin-board-api/internal/lib"
	"bulletin-board-api/internal/metrics"
	"bulletin-board-api/internal/models"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gorm.io/gorm"

	serviceErrors "bulletin-board-api/internal/errors"

	repo "bulletin-board-api/internal/repository"
)

const DuplicateErrorMessageStr = "duplicate key value violates unique constraint"

var tracerType = tracer.SpanType("function")

/*
*****************

	Interface

*****************.
*/
type AddressService interface {
	GetAddressOwnership(ctx context.Context, address string, chain string) (*GetAddressOwnershipResponse, error)
	CreateAddressOwnership(ctx context.Context, request *CreateAddressOwnershipRequest) (*CreateAddressOwnershipResponse, error)
	UpdateAddressOwnershipProof(ctx context.Context, request *UpdateAddressOwnershipProofRequest) (*UpdateAddressOwnershipProofResponse, error)
	DeleteAddressOwnership(ctx context.Context, address string, chain string) error
}

/*
*****************

	Request/Response Structures

*****************.
*/
type CreateAddressOwnershipRequest struct {
	Chain   string
	Address string
}

type UpdateAddressOwnershipProofRequest struct {
	Chain          string
	IOU            bool
	Address        string
	Signature      *string
	Prefix         *string
	RegistrationID uuid.UUID
	ProofType      *string
	AuxProofData   []AuxProofData
}

func (u *UpdateAddressOwnershipProofRequest) ValidateAndCreateProof() (*models.AddressOwnershipProof, error) {
	proof := &models.AddressOwnershipProof{
		Ownership: &models.AddressOwnership{
			Address: u.Address,
			Chain:   u.Chain,
			ID:      u.RegistrationID,
		},
	}

	if u.IOU { //nolint
		proof.IOU = lib.NewBoolPointer(true)
		if u.Prefix != nil && len(*u.Prefix) != 0 ||
			u.Signature != nil && len(*u.Signature) != 0 ||
			u.ProofType != nil && len(*u.ProofType) != 0 ||
			len(u.AuxProofData) != 0 {
			return nil, serviceErrors.New("can not set proof data when IOU is true", serviceErrors.InvalidArgument)
		}
	} else {
		proof.IOU = lib.NewBoolPointer(false)
		if u.Prefix == nil || len(*u.Prefix) == 0 {
			return nil, serviceErrors.New("invalid proof, prefix is empty", serviceErrors.InvalidArgument)
		}
		proof.Prefix = u.Prefix

		if u.Signature == nil || len(*u.Signature) == 0 {
			return nil, serviceErrors.New("invalid proof, signature is empty", serviceErrors.InvalidArgument)
		}
		proof.Signature = u.Signature
		proof.ProofSubmittedTime = lib.CurrentTimeUTCPtr()

		if u.ProofType == nil || len(*u.ProofType) == 0 {
			return nil, serviceErrors.New("proof type must not be null when IOU is false", serviceErrors.InvalidArgument)
		}
		proof.ProofType = u.ProofType
		if !u.ValidateProof() {
			return nil, serviceErrors.New("invalid proof type and data combination", serviceErrors.InvalidArgument)
		}
		if u.AuxProofData != nil {
			res, err := json.Marshal(u.AuxProofData)
			if err != nil {
				return nil, serviceErrors.New("invalid aux proof data", serviceErrors.InvalidArgument)
			}
			proof.AuxProofData = res
		}
	}
	return proof, nil
}

func (u *UpdateAddressOwnershipProofRequest) ValidateProof() bool {
	if u == nil || u.ProofType == nil {
		return false
	}
	switch *u.ProofType {
	case constants.BitcoinP2PKH,
		constants.BitcoinCashP2PKH,
		constants.LitecoinP2PKH,
		constants.LitecoinP2WPKH,
		constants.LitecoinP2L,
		constants.BitcoinP2WPKH,
		constants.BitcoinP2SHP2WPKH,
		constants.BitcoinP2L,
		constants.EthereumEOA,
		constants.SolanaEd25519,
		constants.StellarAccount,
		constants.AvalancheAccount,
		constants.EvmEOA,
		constants.TronEOA,
		constants.PolkadotECDSA,
		constants.PolkadotEd25519,
		constants.PolkadotSr25519,
		constants.AptosECDSA,
		constants.CosmosSecp256k1,
		constants.NearImplicit,
		constants.FilecoinECDSA,
		constants.DashP2PKH,
		constants.DogecoinP2PKH,
		constants.BitcoinP2TR,
		constants.DfinitySecp256k1,
		constants.EosioAll,
		constants.TezosSecp256k1,
		constants.EthereumClassicEOA,
		constants.AlgorandEd25519,
		constants.HederaEd25519,
		constants.SeiSecp256k1:
		return len(u.AuxProofData) == 0

	case constants.BitcoinP2SH, constants.BitcoinP2WSH, constants.BitcoinCashP2SH, constants.LitecoinP2SH, constants.LitecoinP2WSH:
		if len(u.AuxProofData) == 0 {
			// partial proof without script
			return true
		}
		return len(u.AuxProofData) == 1 &&
			u.AuxProofData[0].Type == constants.RedeemScript &&
			u.AuxProofData[0].Validate()

	case constants.BitcoinP2SHP2WSH, constants.LitecoinP2SHP2WSH:
		return len(u.AuxProofData) == 1 &&
			u.AuxProofData[0].Type == constants.WitnessScript &&
			u.AuxProofData[0].Validate()
	case constants.SolanaPda,
		constants.RippleClassic,
		constants.CardanoBase,
		constants.CardanoPointer,
		constants.CardanoEnterprise,
		constants.CardanoReward,
		constants.CardanoByron,
		constants.AptosEd25519,
		constants.AptosMultiEd25519,
		constants.AptosMinpkbls,
		constants.AptosMinSigbls,
		constants.CosmosSecp256r1,
		constants.NearNamed,
		constants.FilecoinBLS,
		constants.DashP2SH,
		constants.DogecoinP2SH,
		constants.DfinitySecp256r1,
		constants.DfinityEd25519,
		constants.TezosEd25519,
		constants.TezosSecp256r1,
		constants.TonEd25519,
		constants.SuiEd25519:
		return len(u.AuxProofData) == 1

	case constants.EthereumContract, constants.EvmContract, constants.TronContract:
		if len(u.AuxProofData) == 0 {
			return false
		}
		for _, item := range u.AuxProofData {
			if !lib.Contains([]string{constants.EthCreate, constants.EthCreate2, constants.EthAllow, constants.TronCreate, constants.TronCreate2, constants.TronAllow}, item.Type) {
				return false
			}
			if !item.Validate() {
				return false
			}
		}
		return true
	}
	return false
}

type AuxProofData struct {
	Type string            `json:"type"`
	Data map[string]string `json:"data"`
}

func (a *AuxProofData) Validate() bool {
	switch a.Type {
	case constants.RedeemScript, constants.WitnessScript:
		return containExactKeys(a.Data, constants.RedeemScriptValidDataType)
	case constants.EthCreate:
		return containExactKeys(a.Data, constants.EthCreateValidDataType)
	case constants.EthCreate2:
		return containExactKeys(a.Data, constants.EthCreate2ValidDataType)
	case constants.EthAllow:
		return containExactKeys(a.Data, constants.EthAllowValidDataType)
	case constants.TronCreate:
		return containExactKeys(a.Data, constants.TronCreateValidDataType)
	case constants.TronCreate2:
		return containExactKeys(a.Data, constants.TronCreate2ValidDataType)
	case constants.TronAllow:
		return containExactKeys(a.Data, constants.TronAllowValidDataType)
	}
	return false
}

func containExactKeys(m map[string]string, keys []string) bool {
	if len(m) != len(keys) {
		return false
	}
	for _, key := range keys {
		if _, ok := m[key]; !ok {
			return false
		}
	}
	return true
}

type GetAddressOwnershipProofResponse struct {
	VASP                  models.VASP
	AddressOwnershipProof models.AddressOwnershipProof
}

type CreateAddressOwnershipResponse struct {
	Address        string
	Chain          string
	RegistrationID uuid.UUID
}

type GetAddressOwnershipResponse struct {
	Vasp                  *models.VASP
	AddressOwnershipProof *models.AddressOwnershipProof
}

type UpdateAddressOwnershipProofResponse struct {
	AddressOwnershipProof *models.AddressOwnershipProof
}

/*
*****************
Implementation
*****************.
*/
type addressService struct {
	vaspService VaspService
	addressRepo repo.AddressRepository
	encrypter   encryption.Encrypter
	logger      *zap.Logger
	sd          *metrics.StatsD
}

func NewAddressService(
	vaspService VaspService,
	addressRepo repo.AddressRepository,
	logger *zap.Logger,
	encrypter encryption.Encrypter,
	sd *metrics.StatsD,
) AddressService {
	return &addressService{
		vaspService: vaspService,
		addressRepo: addressRepo,
		logger:      logger,
		encrypter:   encrypter,
		sd:          sd,
	}
}

func (s *addressService) getVaspFromContext(ctx context.Context) (res *models.VASP, err error) {
	vasp, ok := ctx.Value(constants.RequestVASPCtxKey{}).(models.VASP)
	if !ok || vasp == (models.VASP{}) || len(vasp.ID.String()) == 0 {
		return res, serviceErrors.New("vasp is not set", serviceErrors.PermissionDenied)
	}

	return &vasp, nil
}

func (s *addressService) GetAddressOwnership(ctx context.Context, address string, chain string) (res *GetAddressOwnershipResponse, err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "GetAddressOwnership", tracerType)

	defer func() {
		span.Finish(tracer.WithError(err))
	}()

	logger := ctxzap.Extract(ctx).With(
		zap.String("chain", chain),
	)

	// validate request vasp
	_, err = s.getVaspFromContext(ctx)
	if err != nil {
		return nil, err
	}

	proofResp, err := s.addressRepo.GetOne(ctx, &repo.AddressRepoGetRequest{
		Address: address,
		Chain:   chain,
	})

	if err != nil {
		logger.Error("load address ownership failed", zap.Error(err))
		return nil, serviceErrors.Wrap(err, "load address ownership failed", serviceErrors.Internal)
	}

	if proofResp == nil {
		// no owner of the address and chain; return empty response
		return nil, serviceErrors.New("owning VASP not found", serviceErrors.NotFound)
	}

	if proofResp.Ownership.EncryptedVaspUUID == nil {
		return nil, serviceErrors.Wrap(err, "lookup address owership failed", serviceErrors.Internal)
	}

	// Decrypt the encrypted vasp uuid
	decryptedVaspUUID, err := s.encrypter.Decrypt(ctx, *proofResp.Ownership.EncryptedVaspUUID)
	if err != nil {
		return nil, serviceErrors.Wrap(err, "decrypt vasp uuid failed", serviceErrors.Internal)
	}

	owningVASPUUID, err := uuid.Parse(decryptedVaspUUID)
	if err != nil {
		return nil, serviceErrors.Wrap(err, "parse vasp uuid failed", serviceErrors.Internal)
	}

	// Get VASP info
	vasp, err := s.vaspService.GetVasp(ctx, owningVASPUUID)
	if err != nil {
		// no owner of the address and chain; return empty response
		logger.Warn("owning VASP not found", zap.String("deprecated_vasp_uuid", decryptedVaspUUID))
		return nil, serviceErrors.New("owning VASP not found", serviceErrors.NotFound)
	}

	var proof *models.AddressOwnershipProof
	// IOU should be either true or false if a proof is updated.
	// If IOU is nil, it means the address is claimed but neither proof or iou was submitted. We return nil for proof in this case.
	if proofResp.IOU != nil {
		proof = proofResp
	}

	return &GetAddressOwnershipResponse{
		Vasp:                  vasp,
		AddressOwnershipProof: proof,
	}, nil
}

// Only used when an address is first claimed
// It creates empty AddressOwnershipProof record.
func (s *addressService) CreateAddressOwnership(ctx context.Context, request *CreateAddressOwnershipRequest) (resp *CreateAddressOwnershipResponse, err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "CreateAddressOwnership", tracerType)

	defer func() {
		span.Finish(tracer.WithError(err))
	}()

	// 1. get requesting vasp
	vasp, err := s.getVaspFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// New migration workflow
	// 1. get the existing proof
	proof, err := s.addressRepo.GetOne(ctx, &repo.AddressRepoGetRequest{
		Address: request.Address,
		Chain:   request.Chain,
	})

	if err != nil {
		return nil, serviceErrors.Wrap(err, "load address ownership failed", serviceErrors.Internal)
	}

	if proof != nil {
		// Decrypt the encrypted vasp uuid
		vaspUUID, err := s.encrypter.Decrypt(ctx, *proof.Ownership.EncryptedVaspUUID)
		if err != nil {
			s.logger.Error("decrypt vasp uuid failed", zap.Error(err))
			return nil, serviceErrors.Wrap(err, "decrypt vasp uuid failed", serviceErrors.Internal)
		}

		if vaspUUID != vasp.ID.String() {
			// Current owner is different from the requesting vasp
			// This is a potential abusive behavior we should monitor
			// Log details for handling and investigating the duplicate address claims
			s.logger.Error("duplicate address ownership proof details",
				zap.String("address", request.Address),
				zap.String("chain", request.Chain),
				zap.String("owning_vasp", vaspUUID),
			)

			return nil, serviceErrors.New("address is already claimed by another VASP", serviceErrors.AlreadyExists)
		}

		return &CreateAddressOwnershipResponse{
			Address:        request.Address,
			Chain:          request.Chain,
			RegistrationID: proof.Ownership.ID,
		}, nil
	}

	encryptedVaspUUID, err := s.encrypter.Encrypt(ctx, vasp.ID.String())
	if err != nil {
		s.logger.Error("encrypt vasp uuid failed", zap.Error(err))
		return nil, serviceErrors.Wrap(err, "encrypt vasp uuid failed", serviceErrors.Internal)
	}

	// insert new proof
	ownership := &models.AddressOwnership{
		ID:                uuid.New(),
		Address:           request.Address,
		Chain:             request.Chain,
		EncryptedVaspUUID: &encryptedVaspUUID,
	}

	err = s.addressRepo.CreateOneOwnership(ctx, ownership)

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), DuplicateErrorMessageStr) {
			s.logger.Info("address is already claimed")
			return nil, serviceErrors.New("address is already claimed", serviceErrors.AlreadyExists)
		}

		s.logger.Error("create address ownership failed", zap.Error(err))
		return nil, serviceErrors.Wrap(err, "create address ownership failed", serviceErrors.Internal)
	}

	return &CreateAddressOwnershipResponse{
		Address:        ownership.Address,
		Chain:          ownership.Chain,
		RegistrationID: ownership.ID,
	}, nil
}

func (s *addressService) validateAndCreateProof(ctx context.Context, request *UpdateAddressOwnershipProofRequest) (proof *models.AddressOwnershipProof, err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "validateAndCreateProof", tracerType)

	defer func() {
		span.Finish(tracer.WithError(err))
	}()

	// 1. check if the VASP is the owner of the address. If the requesting VASP is trying to update a proof which
	// is associated with an ownership claimed by a different VASP, then reject the call.
	requestingVASP, err := s.getVaspFromContext(ctx)
	if err != nil {
		return nil, err
	}
	// 2. get existing proofs
	existingProof, err := s.addressRepo.GetOne(ctx, MapUpdateAddressOwnershipProofRequestToGetRequest(request))
	if err != nil {
		return nil, serviceErrors.Wrap(err, "load address ownership failed", serviceErrors.Internal)
	}

	if existingProof == nil {
		return nil, serviceErrors.New("updating proof for an unclaimed address", serviceErrors.InvalidArgument)
	}

	if existingProof.IOU != nil && !*existingProof.IOU && request.IOU {
		return nil, serviceErrors.New("can not change from non-iou to iou", serviceErrors.AlreadyExists)
	}

	if existingProof.Ownership.EncryptedVaspUUID == nil {
		return nil, serviceErrors.New("submitting proof for an unclaimed address", serviceErrors.NotFound)
	}

	decryptedVaspUUID, err := s.encrypter.Decrypt(ctx, *existingProof.Ownership.EncryptedVaspUUID)
	if err != nil {
		return nil, serviceErrors.New("can't decrypt vasp UUID", serviceErrors.Internal)
	}

	vaspUUID, err := uuid.Parse(decryptedVaspUUID)

	if err != nil {
		return nil, serviceErrors.New("malformed vasp UUID", serviceErrors.Internal)
	}

	if requestingVASP.ID != vaspUUID {
		return nil, serviceErrors.New("requesting VASP and owner VASP does not match", serviceErrors.PermissionDenied)
	}

	// 3. the proof is valid or it's an IOU
	return request.ValidateAndCreateProof()
}

// UpdateAddressOwnershipProof is used to update existing AddressOwnershipProof.
func (s *addressService) UpdateAddressOwnershipProof(
	ctx context.Context,
	request *UpdateAddressOwnershipProofRequest,
) (resp *UpdateAddressOwnershipProofResponse, err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "UpdateAddressOwnershipProof", tracerType)
	defer func() {
		span.Finish(tracer.WithError(err))
	}()

	proof, err := s.validateAndCreateProof(ctx, request)
	if err != nil {
		return nil, err
	}

	err = s.addressRepo.UpdateOne(ctx, proof)
	if err != nil {
		return resp, serviceErrors.Wrap(err, "update address ownership failed", serviceErrors.Internal)
	}
	return &UpdateAddressOwnershipProofResponse{
		AddressOwnershipProof: proof,
	}, nil
}

func (s *addressService) DeleteAddressOwnership(ctx context.Context, address string, chain string) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "DeleteAddressOwnership", tracerType)
	defer func() {
		span.Finish(tracer.WithError(err))
	}()

	// 1. Get Vasp from context
	requestingVASP, err := s.getVaspFromContext(ctx)
	if err != nil {
		return err
	}
	// address does not exist
	proof, err := s.addressRepo.GetOne(ctx, &repo.AddressRepoGetRequest{
		Address: address,
		Chain:   chain,
	})
	if err != nil {
		return serviceErrors.Wrap(err, "address ownership lookup failure", serviceErrors.Internal)
	}

	if proof == nil {
		return serviceErrors.New("address not found", serviceErrors.NotFound)
	}

	if proof.Ownership.EncryptedVaspUUID == nil {
		return serviceErrors.New("address not claimed", serviceErrors.NotFound)
	}

	decryptedVaspUUID, err := s.encrypter.Decrypt(ctx, *proof.Ownership.EncryptedVaspUUID)
	if err != nil {
		return serviceErrors.Wrap(err, "failed to decrypt vasp uuid", serviceErrors.Internal)
	}

	if decryptedVaspUUID != requestingVASP.ID.String() {
		return serviceErrors.New("can not delete address ownership that is not owned by the current requesting VASP", serviceErrors.PermissionDenied)
	}

	// 2. delete address ownership data.
	proofToDelete := &models.AddressOwnershipProof{
		Ownership: &models.AddressOwnership{
			Address: address,
			Chain:   chain,
		},
	}
	err = s.addressRepo.DeleteOne(ctx, proofToDelete)
	if err != nil {
		return serviceErrors.Wrap(err, "delete address ownership proof failed", serviceErrors.Internal)
	}
	return nil
}
