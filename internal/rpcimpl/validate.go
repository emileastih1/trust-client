package rpcimpl

import (
	"bulletin-board-api/internal/constants"
	"bulletin-board-api/internal/lib"
	"errors"
	"fmt"
	"regexp"

	"github.com/google/uuid"

	pb "bulletin-board-api/gen/go"

	serviceErrors "bulletin-board-api/internal/errors"
)

var (
	errInvalidRequest = errors.New("invalid request")
	// common characters.
	expression = regexp.MustCompile(`^[A-Za-z0-9_].+$`)
	// SHA512(address) in HEX-encoded string has 128 characters (64 bytes).
	addressExpression = regexp.MustCompile("^[a-f0-9]{128}$")
	// int value.
	intExpression = regexp.MustCompile("^[0-9]{1,128}$")

	// 32 bytes hex string in Ethereum convention (with 0x prefix, lowercase).
	// salt is 32 bytes hex (0x) string (64 characters)
	// init_code_hash is keccak256 hash of the init_code in hex (0x) string (32 bytes)
	// example: 0x14e26bf8489b5ecf8bed66029345a1fbfce9aba797dee64db4bd6403515916db
	ethereumHex32Expression = regexp.MustCompile("^0x?[a-f0-9]{64}$")

	// EIP-55: Mixed-case checksum address encoding.
	ethereumChecksumHexExpression = regexp.MustCompile("^0x?[a-fA-F0-9]{40}$")

	// Hex string in Bitcoin convention (without 0x prefix, lowercase).
	// No length constraints.
	// P2SH the redeem script is capped at 520 bytes
	// P2WSH allows maximum script size of 10,000 bytes.
	bitcoinHexExpression = regexp.MustCompile("^[a-f0-9]+$")

	base64EncodedString = regexp.MustCompile("^[a-zA-Z0-9+/=]+$")

	chainExpression = regexp.MustCompile("^([A-Z_])+$")

	maxStringLen       = 4096
	maxParamValueLen   = 2048
	maxAuxProofDataLen = 8

	chainsWithProofSupport = []string{
		constants.ChainBitcoin,
		constants.ChainEthereum,
		constants.ChainSolana,
		constants.ChainRipple,
		constants.ChainLiteCoin,
		constants.ChainStellar,
		constants.ChainCronos,
		constants.ChainAvalancheCChain,
		constants.ChainAvalanchePChain,
		constants.ChainAvalancheXChain,
		constants.ChainBitcoinCash,
		constants.ChainTron,
		constants.ChainFilecoin,
		constants.ChainArbitrum,
		constants.ChainNear,
		constants.ChainCelo,
		constants.ChainAxelar,
		constants.ChainPolygonPOS,
		constants.ChainCosmos,
		constants.ChainAptos,
		constants.ChainPolkadot,
		constants.ChainCardano,
		constants.ChainBase,
		constants.ChainDogeCoin,
		constants.ChainDash,
		constants.ChainDfinity,
		constants.ChainAlgorand,
		constants.ChainEosio,
		constants.ChainTezos,
		constants.ChainTon,
		constants.ChainBsc,
		constants.ChainOptimism,
		constants.ChainEthereumClassic,
		constants.ChainSui,
		constants.ChainHedera,
		constants.ChainSei,
		constants.ChainHyperliquid,
		constants.ChainMonad,
		constants.ChainArc,
	}

	validProofTypes = []string{
		constants.BitcoinP2PKH,
		constants.BitcoinP2SH,
		constants.BitcoinP2SHP2WPKH,
		constants.BitcoinP2SHP2WSH,
		constants.BitcoinP2WPKH,
		constants.BitcoinP2WSH,
		constants.BitcoinP2L,
		constants.EthereumEOA,
		constants.EthereumContract,
		constants.SolanaEd25519,
		constants.SolanaPda,
		constants.RippleClassic,
		constants.LitecoinP2PKH,
		constants.LitecoinP2WPKH,
		constants.LitecoinP2SH,
		constants.LitecoinP2WSH,
		constants.LitecoinP2SHP2WPKH,
		constants.LitecoinP2SHP2WSH,
		constants.LitecoinP2L,
		constants.StellarAccount,
		constants.EvmEOA,
		constants.EvmContract,
		constants.AvalancheAccount,
		constants.BitcoinCashP2PKH,
		constants.BitcoinCashP2SH,
		constants.TronEOA,
		constants.FilecoinBLS,
		constants.FilecoinECDSA,
		constants.NearImplicit,
		constants.NearNamed,
		constants.CosmosSecp256k1,
		constants.CosmosSecp256r1,
		constants.AptosECDSA,
		constants.AptosEd25519,
		constants.AptosMultiEd25519,
		constants.AptosMinpkbls,
		constants.AptosMinSigbls,
		constants.PolkadotECDSA,
		constants.PolkadotEd25519,
		constants.PolkadotSr25519,
		constants.CardanoBase,
		constants.CardanoEnterprise,
		constants.CardanoPointer,
		constants.CardanoReward,
		constants.CardanoByron,
		constants.DogecoinP2PKH,
		constants.DogecoinP2SH,
		constants.DashP2PKH,
		constants.DashP2SH,
		constants.BitcoinP2TR,
		constants.DfinitySecp256k1,
		constants.DfinitySecp256r1,
		constants.DfinityEd25519,
		constants.AlgorandEd25519,
		constants.EosioAll,
		constants.TezosEd25519,
		constants.TezosSecp256k1,
		constants.TezosSecp256r1,
		constants.TonEd25519,
		constants.EthereumClassicEOA,
		constants.SuiEd25519,
		constants.HederaEd25519,
		constants.SeiSecp256k1,
		constants.TronContract,
	}

	validAuxProofDataType = []string{
		constants.EthCreate,
		constants.EthCreate2,
		constants.EthAllow,
		constants.RedeemScript,
		constants.WitnessScript,
		constants.TronCreate,
		constants.TronCreate2,
		constants.TronAllow,
	}

	validKnownAuxProofDataKeysAndValueRegex = map[string]*regexp.Regexp{
		constants.Nonce:              intExpression,
		constants.Salt:               ethereumHex32Expression,
		constants.InitCodeHash:       ethereumHex32Expression,
		constants.ClaimAddress:       ethereumChecksumHexExpression,
		constants.Script:             bitcoinHexExpression,
		constants.EncryptedProgramID: base64EncodedString,
		constants.EncryptedSeeds:     base64EncodedString,
		constants.EncryptedPublicKey: base64EncodedString,
	}
)

// CommonRequest ...
type CommonRequest interface {
	GetAddress() string
	GetChain() string
}

func validateCommonRequest(r CommonRequest) error {
	if r == nil {
		return serviceErrors.New("invalid request", serviceErrors.InvalidArgument)
	}
	if !isAddressValid(r.GetAddress()) {
		return serviceErrors.New("invalid address", serviceErrors.InvalidArgument)
	}
	if !chainExpression.MatchString(r.GetChain()) {
		return serviceErrors.New("invalid chain: "+lib.Sanitize(r.GetChain()), serviceErrors.InvalidArgument)
	}
	return nil
}

func validateChainSupportsProof(chain string) error {
	if !lib.Contains(chainsWithProofSupport, chain) {
		return serviceErrors.New("invalid chain: "+lib.Sanitize(chain), serviceErrors.InvalidArgument)
	}
	return nil
}

func isAuxProofDataItemKeyValid(key string) bool {
	_, found := validKnownAuxProofDataKeysAndValueRegex[key]
	return found && len(key) != 0 && isStringValid(key)
}

func isStringValid(data string) bool {
	return expression.MatchString(data) && len(data) <= maxStringLen
}

func isAuxProofDataItemValueValid(value string, key string) bool {
	regex, ok := validKnownAuxProofDataKeysAndValueRegex[key]
	// we do not allow unknown keys
	if !ok {
		return false
	}

	return regex.MatchString(value) && len(value) <= maxParamValueLen
}

func isNullableStringValid(data string) bool {
	if len(data) == 0 {
		return true
	}
	return isStringValid(data)
}

func isAddressValid(address string) bool {
	return len(address) != 0 && addressExpression.MatchString(address)
}

func isAuxProofDataTypeValid(data string) bool {
	return len(data) == 0 || lib.Contains(validAuxProofDataType, data)
}

func isAuxProofDataValid(auxProofData []*pb.AuxProofData) bool {
	// empty is valid in certain context.
	if len(auxProofData) == 0 {
		return true
	}
	if len(auxProofData) > maxAuxProofDataLen {
		return false
	}

	for i := 0; i < len(auxProofData); i++ {
		if !isAuxProofDataTypeValid(auxProofData[i].Type) {
			return false
		}

		// proof data, once provided, must not be empty.
		if len(auxProofData[i].Data) == 0 {
			return false
		}
		for k, v := range auxProofData[i].Data {
			if !isAuxProofDataItemKeyValid(k) || !isAuxProofDataItemValueValid(v, k) {
				return false
			}
		}
	}

	return true
}

func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func validateCreateAddressOwnershipProofRequest(r *pb.CreateAddressOwnershipProofRequest) error {
	if r == nil {
		return fmt.Errorf("%w: request is nil", errInvalidRequest)
	}

	if !isNullableStringValid(r.GetSignature()) {
		return serviceErrors.New("invalid signature", serviceErrors.InvalidArgument)
	}

	if !isNullableStringValid(r.GetProofType()) ||
		len(r.GetProofType()) != 0 && !lib.Contains(validProofTypes, r.GetProofType()) {
		return serviceErrors.New("invalid proof type", serviceErrors.InvalidArgument)
	}

	if !isAuxProofDataValid(r.GetAuxProofData()) {
		return serviceErrors.New("invalid aux proof data", serviceErrors.InvalidArgument)
	}

	if !isValidUUID(r.GetRegistrationId()) {
		return serviceErrors.New("invalid registration id", serviceErrors.InvalidArgument)
	}

	return nil
}
