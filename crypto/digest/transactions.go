package digest

import (
	"github.com/orbs-network/membuffers/go"
	"github.com/orbs-network/orbs-network-go/crypto/hash"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/pkg/errors"
)

const (
	TX_ID_SIZE_BYTES = 8 + 32
)

func CalcTxHash(transaction *protocol.Transaction) primitives.Sha256 {
	return hash.CalcSha256(transaction.Raw())
}

func CalcSignedTxHashes(signedTransactions []*protocol.SignedTransaction) []primitives.Sha256 {
	res := make([]primitives.Sha256, len(signedTransactions))
	for i := 0; i < len(signedTransactions); i++ {
		res[i] = CalcTxHash(signedTransactions[i].Transaction())
	}
	return res
}

func CalcReceiptHash(receipt *protocol.TransactionReceipt) primitives.Sha256 {
	return hash.CalcSha256(receipt.Raw())
}

func CalcReceiptHashes(receipts []*protocol.TransactionReceipt) []primitives.Sha256 {
	res := make([]primitives.Sha256, len(receipts))
	for i := 0; i < len(receipts); i++ {
		res[i] = CalcReceiptHash(receipts[i])
	}
	return res
}

func CalcTxId(transaction *protocol.Transaction) []byte {
	return GenerateTxId(CalcTxHash(transaction), transaction.Timestamp())
}

func GenerateTxId(txHash primitives.Sha256, txTimestamp primitives.TimestampNano) []byte {
	res := make([]byte, TX_ID_SIZE_BYTES)
	membuffers.WriteUint64(res, uint64(txTimestamp))
	copy(res[8:], txHash)

	return res
}

func ExtractTxId(txId []byte) (txHash primitives.Sha256, txTimestamp primitives.TimestampNano, err error) {
	if len(txId) != TX_ID_SIZE_BYTES {
		err = errors.Errorf("txid has invalid length %d", len(txId))
		return
	}
	txTimestamp = primitives.TimestampNano(membuffers.GetUint64(txId))
	txHash = txId[8:]
	return
}
