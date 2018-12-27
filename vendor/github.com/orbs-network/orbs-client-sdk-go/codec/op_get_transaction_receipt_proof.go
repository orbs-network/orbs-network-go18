package codec

import (
	"github.com/orbs-network/orbs-client-sdk-go/crypto/digest"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/orbs-spec/types/go/protocol/client"
	"github.com/pkg/errors"
	"time"
)

type GetTransactionReceiptProofRequest struct {
	ProtocolVersion uint32
	VirtualChainId  uint32
	TxId            []byte
}

type GetTransactionReceiptProofResponse struct {
	RequestStatus     RequestStatus
	PackedProof       []byte
	TransactionStatus TransactionStatus
	BlockHeight       uint64
	BlockTimestamp    time.Time
	PackedReceipt     []byte
}

func EncodeGetTransactionReceiptProofRequest(req *GetTransactionReceiptProofRequest) ([]byte, error) {
	// validate
	if req.ProtocolVersion != 1 {
		return nil, errors.Errorf("expected ProtocolVersion 1, %d given", req.ProtocolVersion)
	}
	if len(req.TxId) != digest.TX_ID_SIZE_BYTES {
		return nil, errors.Errorf("expected TxId length %d, %d given", digest.TX_ID_SIZE_BYTES, len(req.TxId))
	}

	// extract txid
	txHash, txTimestamp, err := digest.ExtractTxId(req.TxId)
	if err != nil {
		return nil, err
	}

	// encode request
	res := (&client.GetTransactionReceiptProofRequestBuilder{
		ProtocolVersion:      primitives.ProtocolVersion(req.ProtocolVersion),
		VirtualChainId:       primitives.VirtualChainId(req.VirtualChainId),
		TransactionTimestamp: txTimestamp,
		Txhash:               primitives.Sha256(txHash),
	}).Build()

	// return
	return res.Raw(), nil
}

func DecodeGetTransactionReceiptProofResponse(buf []byte) (*GetTransactionReceiptProofResponse, error) {
	// decode response
	res := client.GetTransactionReceiptProofResponseReader(buf)
	if !res.IsValid() {
		return nil, errors.New("response is corrupt and cannot be decoded")
	}

	// decode request status
	requestStatus, err := requestStatusDecode(res.RequestStatus())
	if err != nil {
		return nil, err
	}

	// decode transaction status
	transactionStatus, err := transactionStatusDecode(res.TransactionStatus())
	if err != nil {
		return nil, err
	}

	// return
	return &GetTransactionReceiptProofResponse{
		RequestStatus:     requestStatus,
		PackedProof:       res.PackedProof(),
		TransactionStatus: transactionStatus,
		BlockHeight:       uint64(res.BlockHeight()),
		BlockTimestamp:    time.Unix(0, int64(res.BlockTimestamp())),
		PackedReceipt:     res.PackedReceipt(),
	}, nil
}
