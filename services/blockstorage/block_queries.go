package blockstorage

import (
	"context"
	"fmt"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/orbs-network/orbs-spec/types/go/services"
	"github.com/pkg/errors"
	"time"
)

func (s *service) GetLastCommittedBlockHeight(ctx context.Context, input *services.GetLastCommittedBlockHeightInput) (*services.GetLastCommittedBlockHeightOutput, error) {
	b, err := s.persistence.GetLastBlock()
	if err != nil {
		return nil, err
	}
	return &services.GetLastCommittedBlockHeightOutput{
		LastCommittedBlockHeight:    getBlockHeight(b),
		LastCommittedBlockTimestamp: getBlockTimestamp(b),
	}, nil
}

func (s *service) loadTransactionsBlockHeader(height primitives.BlockHeight) (*services.GetTransactionsBlockHeaderOutput, error) {
	txBlock, err := s.persistence.GetTransactionsBlock(height)
	if err != nil {
		return nil, err
	}
	return &services.GetTransactionsBlockHeaderOutput{
		TransactionsBlockProof:    txBlock.BlockProof,
		TransactionsBlockHeader:   txBlock.Header,
		TransactionsBlockMetadata: txBlock.Metadata,
	}, nil
}

func (s *service) GetTransactionsBlockHeader(ctx context.Context, input *services.GetTransactionsBlockHeaderInput) (result *services.GetTransactionsBlockHeaderOutput, err error) {
	err = s.persistence.GetBlockTracker().WaitForBlock(ctx, input.BlockHeight)
	if err == nil {
		return s.loadTransactionsBlockHeader(input.BlockHeight)
	}
	return nil, err
}

func (s *service) loadResultsBlockHeader(height primitives.BlockHeight) (*services.GetResultsBlockHeaderOutput, error) {
	txBlock, err := s.persistence.GetResultsBlock(height)
	if err != nil {
		return nil, err
	}
	return &services.GetResultsBlockHeaderOutput{
		ResultsBlockProof:  txBlock.BlockProof,
		ResultsBlockHeader: txBlock.Header,
	}, nil
}

func (s *service) GetResultsBlockHeader(ctx context.Context, input *services.GetResultsBlockHeaderInput) (result *services.GetResultsBlockHeaderOutput, err error) {
	err = s.persistence.GetBlockTracker().WaitForBlock(ctx, input.BlockHeight)
	if err == nil {
		return s.loadResultsBlockHeader(input.BlockHeight)
	}
	return nil, err
}

func (s *service) GetTransactionReceipt(ctx context.Context, input *services.GetTransactionReceiptInput) (*services.GetTransactionReceiptOutput, error) {
	endGraceNano := s.config.BlockTransactionReceiptQueryGraceEnd().Nanoseconds()
	startGraceNano := s.config.BlockTransactionReceiptQueryGraceStart().Nanoseconds()
	txExpireNano := s.config.BlockTransactionReceiptQueryExpirationWindow().Nanoseconds()

	start := input.TransactionTimestamp - primitives.TimestampNano(startGraceNano)
	end := input.TransactionTimestamp + primitives.TimestampNano(endGraceNano+txExpireNano)

	// TODO(v1): sanity check, this is really useless here right now, but we were going to refactor this, and when we were going to, this was here to remind us to have a sanity check on this query
	if end < start || end-start > primitives.TimestampNano(time.Hour.Nanoseconds()) {
		receipt, err := s.createEmptyTransactionReceiptResult(ctx)
		if err != nil {
			return nil, err
		}
		// TODO((https://github.com/orbs-network/orbs-network-go/issues/448): probably don't fail here
		return receipt, errors.Errorf("failed to search for blocks on tx timestamp of %d, hash %s", input.TransactionTimestamp, input.Txhash)
	}

	blockPair, txIdx, err := s.persistence.GetBlockByTx(input.Txhash, start, end)
	if err != nil {
		return nil, err
	}
	if blockPair == nil {
		receipt, err := s.createEmptyTransactionReceiptResult(ctx)
		if err != nil {
			return nil, err
		}
		return receipt, nil
	}

	return &services.GetTransactionReceiptOutput{
		TransactionReceipt: blockPair.ResultsBlock.TransactionReceipts[txIdx],
		BlockHeight:        blockPair.ResultsBlock.Header.BlockHeight(),
		BlockTimestamp:     blockPair.ResultsBlock.Header.Timestamp(),
	}, nil
}

// Returns a slice of blocks containing first and last
// TODO kill this method signature or use a larger page size without returning too many blocks
func (s *service) GetBlockSlice(first primitives.BlockHeight, last primitives.BlockHeight) ([]*protocol.BlockPairContainer, primitives.BlockHeight, primitives.BlockHeight, error) {
	blocks := make([]*protocol.BlockPairContainer, 0, last-first+1)
	err := s.persistence.ScanBlocks(first, 1, func(firstInPage primitives.BlockHeight, page []*protocol.BlockPairContainer) bool {
		blocks = append(blocks, page...)
		return firstInPage < last
	})
	if err != nil {
		return nil, 0, 0, errors.Wrap(err, "failed getting block slice")
	}
	if len(blocks) == 0 {
		return nil, 0, 0, fmt.Errorf("could not find blocks in height range %d-%d", first, last)
	}
	return blocks, first, first + primitives.BlockHeight(len(blocks)) - 1, nil
}

func (s *service) createEmptyTransactionReceiptResult(ctx context.Context) (*services.GetTransactionReceiptOutput, error) {
	out, err := s.GetLastCommittedBlockHeight(ctx, &services.GetLastCommittedBlockHeightInput{})
	if err != nil {
		return nil, err
	}
	return &services.GetTransactionReceiptOutput{
		TransactionReceipt: nil,
		BlockHeight:        out.LastCommittedBlockHeight,
		BlockTimestamp:     out.LastCommittedBlockTimestamp,
	}, nil
}
