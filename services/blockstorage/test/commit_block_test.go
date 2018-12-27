package test

import (
	"context"
	"github.com/orbs-network/orbs-network-go/crypto/digest"
	"github.com/orbs-network/orbs-network-go/test"
	"github.com/orbs-network/orbs-network-go/test/builders"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCommitBlockSavesToPersistentStorage(t *testing.T) {
	test.WithContext(func(ctx context.Context) {
		harness := newBlockStorageHarness().withSyncBroadcast(1).start(ctx)

		blockCreated := time.Now()
		blockHeight := primitives.BlockHeight(1)

		_, err := harness.commitBlock(ctx, builders.BlockPair().WithHeight(blockHeight).WithBlockCreated(blockCreated).Build())

		require.NoError(t, err)
		require.EqualValues(t, 1, harness.numOfWrittenBlocks())

		harness.verifyMocks(t, 1)

		lastCommittedBlockHeight := harness.getLastBlockHeight(ctx, t)

		require.EqualValues(t, blockHeight, lastCommittedBlockHeight.LastCommittedBlockHeight, "block height in storage should be the same")
		require.EqualValues(t, blockCreated.UnixNano(), lastCommittedBlockHeight.LastCommittedBlockTimestamp, "timestamp in storage should be the same")

	})
	// TODO(https://github.com/orbs-network/orbs-network-go/issues/569) Spec: If any of the intra block syncs (StateStorage, TransactionPool) is blocking and waiting, wake it up.
}

func TestCommitBlockDoesNotUpdateCommittedBlockHeightAndTimestampIfStorageFails(t *testing.T) {
	test.WithContext(func(ctx context.Context) {
		harness := newBlockStorageHarness().withSyncBroadcast(1).start(ctx)

		blockCreated := time.Now()
		blockHeight := primitives.BlockHeight(1)

		harness.commitBlock(ctx, builders.BlockPair().WithHeight(blockHeight).WithBlockCreated(blockCreated).Build())
		require.EqualValues(t, 1, harness.numOfWrittenBlocks())

		harness.failNextBlocks()

		_, err := harness.commitBlock(ctx, builders.BlockPair().WithHeight(blockHeight+1).Build())
		require.EqualError(t, err, "could not write a block", "error should be returned if storage fails")

		harness.verifyMocks(t, 1)

		lastCommittedBlockHeight := harness.getLastBlockHeight(ctx, t)

		require.EqualValues(t, blockHeight, lastCommittedBlockHeight.LastCommittedBlockHeight, "block height should not update as storage was unavailable")
		require.EqualValues(t, blockCreated.UnixNano(), lastCommittedBlockHeight.LastCommittedBlockTimestamp, "timestamp should not update as storage was unavailable")

	})
}

func TestCommitBlockReturnsErrorWhenProtocolVersionMismatches(t *testing.T) {
	test.WithContext(func(ctx context.Context) {
		harness := newBlockStorageHarness().withSyncBroadcast(1).start(ctx)

		_, err := harness.commitBlock(ctx, builders.BlockPair().WithProtocolVersion(99999).Build())

		require.EqualError(t, err, "protocol version mismatch in transactions block header")
	})
}

func TestCommitBlockDiscardsBlockIfAlreadyExists(t *testing.T) {
	test.WithContext(func(ctx context.Context) {
		harness := newBlockStorageHarness().withSyncBroadcast(1).start(ctx)
		blockPair := builders.BlockPair().Build()

		harness.commitBlock(ctx, blockPair)
		_, err := harness.commitBlock(ctx, blockPair)

		require.NoError(t, err)

		require.EqualValues(t, 1, harness.numOfWrittenBlocks(), "block should be written only once")
		harness.verifyMocks(t, 1)
	})
}

func TestCommitBlockReturnsErrorIfBlockExistsButHasDifferentTimestamp(t *testing.T) {
	test.WithContext(func(ctx context.Context) {
		harness := newBlockStorageHarness().withSyncBroadcast(1).start(ctx)

		blockPair := builders.BlockPair()
		harness.commitBlock(ctx, blockPair.Build())

		mutatedBlockPair := blockPair.WithBlockCreated(time.Now().Add(1 * time.Hour)).Build()
		_, err := harness.commitBlock(ctx, mutatedBlockPair)

		require.EqualError(t, err, "FORK!! block already in storage, timestamp mismatch", "same block, different timestamp should return an error")
		require.EqualValues(t, 1, harness.numOfWrittenBlocks(), "only one block should have been written")
		harness.verifyMocks(t, 1)
	})
}

func TestCommitBlockReturnsErrorIfBlockExistsButHasDifferentTxBlock(t *testing.T) {
	test.WithContext(func(ctx context.Context) {
		harness := newBlockStorageHarness().withSyncBroadcast(1).start(ctx)

		blockPair := builders.BlockPair()
		harness.commitBlock(ctx, blockPair.Build())

		mutatedBlock := blockPair.Build()
		mutatedBlock.TransactionsBlock.Header.MutateNumSignedTransactions(999)

		_, err := harness.commitBlock(ctx, mutatedBlock)

		require.EqualError(t, err, "FORK!! block already in storage, transaction block header mismatch", "same block, different timestamp should return an error")
		require.EqualValues(t, 1, harness.numOfWrittenBlocks(), "only one block should have been written")
		harness.verifyMocks(t, 1)
	})
}

func TestCommitBlockReturnsErrorIfBlockExistsButHasDifferentRxBlock(t *testing.T) {
	test.WithContext(func(ctx context.Context) {
		harness := newBlockStorageHarness().withSyncBroadcast(1).start(ctx)

		blockPair := builders.BlockPair()
		harness.commitBlock(ctx, blockPair.Build())

		mutatedBlock := blockPair.Build()
		mutatedBlock.ResultsBlock.Header.MutateNumTransactionReceipts(999)

		_, err := harness.commitBlock(ctx, mutatedBlock)

		require.EqualError(t, err, "FORK!! block already in storage, results block header mismatch", "same block, different timestamp should return an error")
		require.EqualValues(t, 1, harness.numOfWrittenBlocks(), "only one block should have been written")
		harness.verifyMocks(t, 1)
	})
}

func TestCommitBlockReturnsErrorIfBlockIsNotSequential(t *testing.T) {
	test.WithContext(func(ctx context.Context) {
		harness := newBlockStorageHarness().withSyncBroadcast(1).start(ctx)

		harness.commitBlock(ctx, builders.BlockPair().Build())

		_, err := harness.commitBlock(ctx, builders.BlockPair().WithHeight(1000).Build())
		require.EqualError(t, err, "block height is 1000, expected 2", "block height was mutate to be invalid, should return an error")
		require.EqualValues(t, 1, harness.numOfWrittenBlocks(), "only one block should have been written")
		harness.verifyMocks(t, 1)
	})
}

func TestCommitBlockWithSameTransactionTwicePanics(t *testing.T) {
	test.WithContext(func(ctx context.Context) {
		harness := newBlockStorageHarness().withSyncBroadcast(1).start(ctx)

		tx := builders.Transaction().Build()
		txReceipt := builders.TransactionReceipt().WithTransaction(tx.Transaction()).Build()

		block0 := builders.BlockPair().WithHeight(1).WithTimestampNow().WithTransaction(tx).WithReceipt(txReceipt).Build()
		block1 := builders.BlockPair().WithHeight(2).WithTimestampNow().WithTransaction(tx).WithReceipt(txReceipt).Build()

		txHash := digest.CalcTxHash(tx.Transaction())

		_, err := harness.commitBlock(ctx, block0)
		require.NoError(t, err)

		require.Panics(t, func() {
			harness.commitBlock(ctx, block1)
		})

		blockHeight := harness.storageAdapter.WaitForTransaction(ctx, txHash)
		require.EqualValues(t, 1, blockHeight)

		harness.verifyMocks(t, 1)
	})
}
