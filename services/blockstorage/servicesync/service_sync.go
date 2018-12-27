package servicesync

import (
	"context"
	"fmt"
	"github.com/orbs-network/orbs-network-go/instrumentation/log"
	"github.com/orbs-network/orbs-network-go/instrumentation/trace"
	"github.com/orbs-network/orbs-network-go/services/blockstorage/adapter"
	"github.com/orbs-network/orbs-network-go/synchronization"
	"github.com/orbs-network/orbs-network-go/synchronization/supervised"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
)

type BlockPairCommitter interface {
	commitBlockPair(ctx context.Context, committedBlockPair *protocol.BlockPairContainer) (next primitives.BlockHeight, err error)
	getServiceName() string
}

type blockSource interface {
	GetBlockTracker() *synchronization.BlockTracker
	ScanBlocks(from primitives.BlockHeight, pageSize uint8, f adapter.CursorFunc) error
	GetLastBlock() (*protocol.BlockPairContainer, error)
}

func syncToTopBlock(ctx context.Context, source blockSource, committer BlockPairCommitter, logger log.BasicLogger) (primitives.BlockHeight, error) {
	topBlock, err := source.GetLastBlock()
	if err != nil {
		return 0, err
	}

	// try to commit the top block
	requestedHeight := syncOneBlock(ctx, topBlock, committer, logger)

	// scan all available blocks starting the requested height
	committedHeight := requestedHeight - 1
	err = source.ScanBlocks(requestedHeight, 1, func(h primitives.BlockHeight, page []*protocol.BlockPairContainer) bool {
		requestedHeight = syncOneBlock(ctx, page[0], committer, logger)
		committedHeight = h
		return requestedHeight == h+1
	})
	if err != nil {
		return 0, err
	}

	return committedHeight, nil
}

func syncOneBlock(ctx context.Context, block *protocol.BlockPairContainer, committer BlockPairCommitter, logger log.BasicLogger) primitives.BlockHeight {
	h := block.ResultsBlock.Header.BlockHeight()
	// log transactions
	for _, tx := range block.ResultsBlock.TransactionReceipts {
		logger.Info("attempt service sync for block", log.BlockHeight(h), log.Transaction(tx.Txhash()))
	}
	// notify the receiving service of a new block
	requestedHeight, err := committer.commitBlockPair(ctx, block)
	if err != nil {
		logger.Error("failed committing block", log.Error(err), log.BlockHeight(h))
		panic(fmt.Sprintf("failed committing block at height %d", h))
	}
	// if receiving service keep requesting the current height we are stuck
	if h == requestedHeight {
		// TODO (https://github.com/orbs-network/orbs-network-go/issues/617)
		logger.Error("committer requested same block height in response to commit", log.BlockHeight(h))
	}
	return requestedHeight
}

func NewServiceBlockSync(ctx context.Context, logger log.BasicLogger, source blockSource, committer BlockPairCommitter) {
	ctx = trace.NewContext(ctx, committer.getServiceName())
	logger = logger.WithTags(trace.LogFieldFrom(ctx))
	supervised.GoForever(ctx, logger, func() {

		var height primitives.BlockHeight
		var err error
		for err == nil {
			err = source.GetBlockTracker().WaitForBlock(ctx, height+1)
			if err != nil {
				logger.Info("failed waiting for block", log.Error(err))
				return
			}
			height, err = syncToTopBlock(ctx, source, committer, logger)
		}
	})
}
