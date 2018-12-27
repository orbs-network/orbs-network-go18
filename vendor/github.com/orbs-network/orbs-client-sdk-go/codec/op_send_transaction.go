package codec

import (
	"github.com/orbs-network/orbs-client-sdk-go/crypto/digest"
	"github.com/orbs-network/orbs-client-sdk-go/crypto/keys"
	"github.com/orbs-network/orbs-client-sdk-go/crypto/signature"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/orbs-network/orbs-spec/types/go/protocol/client"
	"github.com/pkg/errors"
	"time"
)

type SendTransactionRequest struct {
	ProtocolVersion uint32
	VirtualChainId  uint32
	Timestamp       time.Time
	NetworkType     NetworkType
	PublicKey       []byte
	ContractName    string
	MethodName      string
	InputArguments  []interface{}
}

type SendTransactionResponse struct {
	RequestStatus     RequestStatus
	TxHash            []byte
	ExecutionResult   ExecutionResult
	OutputArguments   []interface{}
	OutputEvents      []*Event
	TransactionStatus TransactionStatus
	BlockHeight       uint64
	BlockTimestamp    time.Time
}

func EncodeSendTransactionRequest(req *SendTransactionRequest, privateKey []byte) ([]byte, []byte, error) {
	// validate
	if req.ProtocolVersion != 1 {
		return nil, nil, errors.Errorf("expected ProtocolVersion 1, %d given", req.ProtocolVersion)
	}
	if len(req.PublicKey) != keys.ED25519_PUBLIC_KEY_SIZE_BYTES {
		return nil, nil, errors.Errorf("expected PublicKey length %d, %d given", keys.ED25519_PUBLIC_KEY_SIZE_BYTES, len(req.PublicKey))
	}
	if len(privateKey) != keys.ED25519_PRIVATE_KEY_SIZE_BYTES {
		return nil, nil, errors.Errorf("expected PrivateKey length %d, %d given", keys.ED25519_PRIVATE_KEY_SIZE_BYTES, len(privateKey))
	}

	// encode method arguments
	inputArgumentArray, err := PackedMethodArgumentsEncode(req.InputArguments)
	if err != nil {
		return nil, nil, err
	}

	// encode network type
	networkType, err := networkTypeEncode(req.NetworkType)
	if err != nil {
		return nil, nil, err
	}

	// encode request
	res := (&client.SendTransactionRequestBuilder{
		SignedTransaction: &protocol.SignedTransactionBuilder{
			Transaction: &protocol.TransactionBuilder{
				ProtocolVersion: primitives.ProtocolVersion(req.ProtocolVersion),
				VirtualChainId:  primitives.VirtualChainId(req.VirtualChainId),
				Timestamp:       primitives.TimestampNano(req.Timestamp.UnixNano()),
				Signer: &protocol.SignerBuilder{
					Scheme: protocol.SIGNER_SCHEME_EDDSA,
					Eddsa: &protocol.EdDSA01SignerBuilder{
						NetworkType:     networkType,
						SignerPublicKey: primitives.Ed25519PublicKey(req.PublicKey),
					},
				},
				ContractName:       primitives.ContractName(req.ContractName),
				MethodName:         primitives.MethodName(req.MethodName),
				InputArgumentArray: inputArgumentArray,
			},
			Signature: make([]byte, signature.ED25519_SIGNATURE_SIZE_BYTES),
		},
	}).Build()

	// sign
	txHash := digest.CalcTxHash(res.SignedTransaction().Transaction())
	sig, err := signature.SignEd25519(primitives.Ed25519PrivateKey(privateKey), txHash)
	if err != nil {
		return nil, nil, err
	}
	res.SignedTransaction().MutateSignature(sig)

	// return
	return res.Raw(), digest.GenerateTxId(txHash, res.SignedTransaction().Transaction().Timestamp()), nil
}

func DecodeSendTransactionResponse(buf []byte) (*SendTransactionResponse, error) {
	// decode response
	res := client.SendTransactionResponseReader(buf)
	if !res.IsValid() {
		return nil, errors.New("response is corrupt and cannot be decoded")
	}

	// decode request status
	requestStatus, err := requestStatusDecode(res.RequestStatus())
	if err != nil {
		return nil, err
	}

	// decode execution result
	executionResult, err := executionResultDecode(res.TransactionReceipt().ExecutionResult())
	if err != nil {
		return nil, err
	}

	// decode method arguments
	outputArgumentArray, err := PackedMethodArgumentsDecode(res.TransactionReceipt().RawOutputArgumentArrayWithHeader())
	if err != nil {
		return nil, err
	}

	// decode events
	outputEventArray, err := PackedEventsDecode(res.TransactionReceipt().RawOutputEventsArrayWithHeader())
	if err != nil {
		return nil, err
	}

	// decode transaction status
	transactionStatus, err := transactionStatusDecode(res.TransactionStatus())
	if err != nil {
		return nil, err
	}

	// return
	return &SendTransactionResponse{
		RequestStatus:     requestStatus,
		TxHash:            res.TransactionReceipt().Txhash(),
		ExecutionResult:   executionResult,
		OutputArguments:   outputArgumentArray,
		OutputEvents:      outputEventArray,
		TransactionStatus: transactionStatus,
		BlockHeight:       uint64(res.BlockHeight()),
		BlockTimestamp:    time.Unix(0, int64(res.BlockTimestamp())),
	}, nil
}
