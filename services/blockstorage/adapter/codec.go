package adapter

import (
	"encoding/binary"
	"fmt"
	"github.com/orbs-network/membuffers/go"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/pkg/errors"
	"io"
)

// TODO V1 write test for violation of: blockPair == nil || blockPair.TransactionsBlock == nil || blockPair.ResultsBlock == nil
// TODO V1 write test for violation of: other nil values
type blockHeader struct {
	FixedSize    uint32
	ReceiptsSize uint32
	DiffsSize    uint32
	TxsSize      uint32
}

func (h *blockHeader) addFixed(m membuffers.Message) {
	h.FixedSize += 4 + uint32(len(m.Raw()))
}

func (h *blockHeader) addReceipt(receipt *protocol.TransactionReceipt) {
	h.ReceiptsSize += 4 + uint32(len(receipt.Raw()))
}

func (h *blockHeader) addDiff(diff *protocol.ContractStateDiff) {
	h.DiffsSize += 4 + uint32(len(diff.Raw()))
}

func (h *blockHeader) addTx(tx *protocol.SignedTransaction) {
	h.TxsSize += 4 + uint32(len(tx.Raw()))
}

func (h *blockHeader) totalSize() uint32 {
	return h.FixedSize + h.DiffsSize + h.ReceiptsSize + h.TxsSize
}

func (h *blockHeader) write(w io.Writer) error {
	err := binary.Write(w, binary.LittleEndian, h)
	if err != nil {
		return err
	}
	return nil
}

func (h *blockHeader) read(r io.Reader) error {
	err := binary.Read(r, binary.LittleEndian, h)
	if err != nil {
		return err
	}
	return nil
}

func writeMessage(writer io.Writer, message membuffers.Message) error {
	err := binary.Write(writer, binary.LittleEndian, uint32(len(message.Raw())))
	if err != nil {
		return err
	}
	err = binary.Write(writer, binary.LittleEndian, message.Raw())
	if err != nil {
		return err
	}
	return nil
}

func readChunk(reader io.Reader) ([]byte, error) {
	var chunkSize uint32
	err := binary.Read(reader, binary.LittleEndian, &chunkSize)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read chunk size from disk")
	}

	chunk := make([]byte, chunkSize)
	n, err := reader.Read(chunk)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read block chunk from disk")
	}
	if n != len(chunk) {
		return nil, fmt.Errorf("read %d bytes in block chuck while expecting %d", n, len(chunk))
	}
	return chunk, nil
}

func encode(block *protocol.BlockPairContainer, w io.Writer) error {
	tb := block.TransactionsBlock
	rb := block.ResultsBlock

	// write header
	serializationHeader := &blockHeader{}
	serializationHeader.addFixed(tb.Header)
	serializationHeader.addFixed(tb.Metadata)
	serializationHeader.addFixed(tb.BlockProof)
	serializationHeader.addFixed(rb.Header)
	serializationHeader.addFixed(rb.BlockProof)
	serializationHeader.addFixed(rb.TransactionsBloomFilter)
	for _, receipt := range rb.TransactionReceipts {
		serializationHeader.addReceipt(receipt)
	}
	for _, diff := range rb.ContractStateDiffs {
		serializationHeader.addDiff(diff)
	}
	for _, tx := range tb.SignedTransactions {
		serializationHeader.addTx(tx)
	}

	err := serializationHeader.write(w)

	if err != nil {
		return errors.Wrap(err, "failed to write block header")
	}

	// write buffers
	err = writeMessage(w, tb.Header)
	if err != nil {
		return errors.Wrap(err, "failed to write block tx header")
	}
	err = writeMessage(w, tb.Metadata)
	if err != nil {
		return errors.Wrap(err, "failed to write block tx metadata")
	}
	err = writeMessage(w, tb.BlockProof)
	if err != nil {
		return errors.Wrap(err, "failed to write block tx proof")
	}
	err = writeMessage(w, rb.Header)
	if err != nil {
		return errors.Wrap(err, "failed to write block results header")
	}
	err = writeMessage(w, rb.BlockProof)
	if err != nil {
		return errors.Wrap(err, "failed to write block results proof")
	}
	err = writeMessage(w, rb.TransactionsBloomFilter)
	if err != nil {
		return errors.Wrap(err, "failed to write block results tx bloom filter")
	}

	for _, receipt := range rb.TransactionReceipts {
		err = writeMessage(w, receipt)
		if err != nil {
			return errors.Wrap(err, "failed to write block tx receipts")
		}

	}
	for _, diff := range rb.ContractStateDiffs {
		err = writeMessage(w, diff)
		if err != nil {
			return errors.Wrap(err, "failed to write block contract diffs")
		}
	}
	for _, tx := range tb.SignedTransactions {
		err = writeMessage(w, tx)
		if err != nil {
			return errors.Wrap(err, "failed to write block signed transactions")
		}
	}
	return nil
}

func decode(r io.Reader) (*protocol.BlockPairContainer, error) {
	serializationHeader := &blockHeader{}
	err := serializationHeader.read(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read block header")
	}

	tbHeaderChunk, err := readChunk(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read chunk")
	}
	tbHeader := protocol.TransactionsBlockHeaderReader(tbHeaderChunk)

	tbMetadataChunk, err := readChunk(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read chunk")
	}
	tbMetadata := protocol.TransactionsBlockMetadataReader(tbMetadataChunk)

	tbBlockProofChunk, err := readChunk(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read chunk")
	}
	tbBlockProof := protocol.TransactionsBlockProofReader(tbBlockProofChunk)

	rbHeaderChunk, err := readChunk(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read chunk")
	}
	rbHeader := protocol.ResultsBlockHeaderReader(rbHeaderChunk)

	rbBlockProofChunk, err := readChunk(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read chunk")
	}
	rbBlockProof := protocol.ResultsBlockProofReader(rbBlockProofChunk)

	rbBloomChunk, err := readChunk(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read chunk")
	}
	rbBloomFilter := protocol.TransactionsBloomFilterReader(rbBloomChunk)

	// TODO V1 add validations : - 1) IsValid() on each membuff 2) check that num of bytes read match header 3) perform structural validations?
	receipts := make([]*protocol.TransactionReceipt, 0, rbHeader.NumTransactionReceipts())
	for i := 0; i < cap(receipts); i++ {
		chunk, err := readChunk(r)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read receipt chunk")
		}
		receipts = append(receipts, protocol.TransactionReceiptReader(chunk))
	}

	stateDiffs := make([]*protocol.ContractStateDiff, 0, rbHeader.NumContractStateDiffs())
	for i := 0; i < cap(stateDiffs); i++ {
		chunk, err := readChunk(r)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read contract state diff chunk")
		}
		stateDiffs = append(stateDiffs, protocol.ContractStateDiffReader(chunk))
	}

	txs := make([]*protocol.SignedTransaction, 0, tbHeader.NumSignedTransactions())
	for i := 0; i < cap(txs); i++ {
		chunk, err := readChunk(r)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read contract state diff chunk")
		}
		txs = append(txs, protocol.SignedTransactionReader(chunk))
	}

	blockPair := &protocol.BlockPairContainer{
		TransactionsBlock: &protocol.TransactionsBlockContainer{
			Header:             tbHeader,
			Metadata:           tbMetadata,
			SignedTransactions: txs,
			BlockProof:         tbBlockProof,
		},
		ResultsBlock: &protocol.ResultsBlockContainer{
			Header:                  rbHeader,
			TransactionsBloomFilter: rbBloomFilter,
			TransactionReceipts:     receipts,
			ContractStateDiffs:      stateDiffs,
			BlockProof:              rbBlockProof,
		},
	}

	return blockPair, nil
}
