package acceptance

import (
	"context"
	"github.com/orbs-network/orbs-network-go/services/processor/native/repository/BenchmarkToken"
	"github.com/orbs-network/orbs-network-go/test/harness"
	"github.com/orbs-network/orbs-network-go/test/harness/services/gossip/adapter"
	"github.com/orbs-network/orbs-spec/types/go/protocol/consensus"
	"github.com/orbs-network/orbs-spec/types/go/protocol/gossipmessages"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestLeaderCommitsTransactionsAndSkipsInvalidOnesLeanHelix(t *testing.T) {
	harness.Network(t).
		WithNumNodes(4).
		WithConsensusAlgos(consensus.CONSENSUS_ALGO_TYPE_LEAN_HELIX).
		Start(func(ctx context.Context, network harness.TestNetworkDriver) {
			contract := network.GetBenchmarkTokenContract()
			t.Log("testing", network.Description()) // leader is nodeIndex 0, validator is nodeIndex 1
			_, txHash := contract.SendTransfer(ctx, 0, 17, 5, 6)

			network.WaitForTransactionInNodeState(ctx, txHash, 0)
			t.Log("finished waiting for tx")

			require.EqualValues(t, benchmarktoken.TOTAL_SUPPLY-17, contract.CallGetBalance(ctx, 0, 5), "getBalance result for the sender on gateway node")
			require.EqualValues(t, 17, contract.CallGetBalance(ctx, 0, 6), "getBalance result for the receiver on gateway node")
			t.Log("test done")
		})

}

func TestLeaderCommitsTransactionsAndSkipsInvalidOnes(t *testing.T) {
	harness.Network(t).Start(func(parent context.Context, network harness.TestNetworkDriver) {
		ctx, cancel := context.WithTimeout(parent, 1*time.Second)
		defer cancel()

		contract := network.GetBenchmarkTokenContract()
		contract.DeployBenchmarkToken(ctx, 5)

		t.Log("testing", network.Description()) // leader is nodeIndex 0, validator is nodeIndex 1

		_, txHash1 := contract.SendTransfer(ctx, 0, 17, 5, 6)
		contract.SendInvalidTransfer(ctx, 0, 5, 6)
		_, txHash2 := contract.SendTransfer(ctx, 0, 22, 5, 6)

		t.Log("waiting for leader blocks")

		network.WaitForTransactionInNodeState(ctx, txHash1, 0)
		network.WaitForTransactionInNodeState(ctx, txHash2, 0)
		require.EqualValues(t, benchmarktoken.TOTAL_SUPPLY-39, contract.CallGetBalance(ctx, 0, 5), "getBalance result on leader")
		require.EqualValues(t, 39, contract.CallGetBalance(ctx, 0, 6), "getBalance result on leader")

		t.Log("waiting for non leader blocks")

		network.WaitForTransactionInNodeState(ctx, txHash1, 1)
		network.WaitForTransactionInNodeState(ctx, txHash2, 1)
		require.EqualValues(t, benchmarktoken.TOTAL_SUPPLY-39, contract.CallGetBalance(ctx, 1, 5), "getBalance result on non leader")
		require.EqualValues(t, 39, contract.CallGetBalance(ctx, 1, 6), "getBalance result on non leader")

	})
}

func TestNonLeaderPropagatesTransactionsToLeader(t *testing.T) {
	harness.Network(t).Start(func(parent context.Context, network harness.TestNetworkDriver) {
		ctx, cancel := context.WithTimeout(parent, 1*time.Second)
		defer cancel()

		contract := network.GetBenchmarkTokenContract()
		contract.DeployBenchmarkToken(ctx, 5)

		t.Log("testing", network.Description()) // leader is nodeIndex 0, validator is nodeIndex 1

		pausedTxForwards := network.TransportTamperer().Pause(adapter.TransactionRelayMessage(gossipmessages.TRANSACTION_RELAY_FORWARDED_TRANSACTIONS))
		txHash := contract.SendTransferInBackground(ctx, 1, 17, 5, 6)

		if err := network.BlockPersistence(0).GetBlockTracker().WaitForBlock(ctx, 2); err != nil {
			t.Errorf("failed waiting for block on node 0: %s", err)
		}
		if err := network.BlockPersistence(1).GetBlockTracker().WaitForBlock(ctx, 2); err != nil {
			t.Errorf("failed waiting for block on node 1: %s", err)
		}

		require.EqualValues(t, 0, contract.CallGetBalance(ctx, 0, 6), "initial getBalance result on leader")
		require.EqualValues(t, 0, contract.CallGetBalance(ctx, 1, 6), "initial getBalance result on non leader")

		pausedTxForwards.Release(ctx)
		network.WaitForTransactionInNodeState(ctx, txHash, 0)
		require.EqualValues(t, 17, contract.CallGetBalance(ctx, 0, 6), "eventual getBalance result on leader")
		network.WaitForTransactionInNodeState(ctx, txHash, 1)
		require.EqualValues(t, 17, contract.CallGetBalance(ctx, 1, 6), "eventual getBalance result on non leader")

	})
}

func TestLeaderCommitsTwoTransactionsInOneBlock(t *testing.T) {
	harness.Network(t).Start(func(parent context.Context, network harness.TestNetworkDriver) {
		ctx, cancel := context.WithTimeout(parent, 1*time.Second)
		defer cancel()

		contract := network.GetBenchmarkTokenContract()
		contract.DeployBenchmarkToken(ctx, 5)

		t.Log("testing", network.Description()) // leader is nodeIndex 0, validator is nodeIndex 1

		txHash1 := contract.SendTransferInBackground(ctx, 0, 17, 5, 6)
		txHash2 := contract.SendTransferInBackground(ctx, 0, 22, 5, 6)

		t.Log("waiting for leader blocks")

		network.WaitForTransactionInNodeState(ctx, txHash1, 0)
		network.WaitForTransactionInNodeState(ctx, txHash2, 0)
		require.EqualValues(t, benchmarktoken.TOTAL_SUPPLY-39, contract.CallGetBalance(ctx, 0, 5), "getBalance result on leader")
		require.EqualValues(t, 39, contract.CallGetBalance(ctx, 0, 6), "getBalance result on leader")

		t.Log("waiting for non leader blocks")

		network.WaitForTransactionInNodeState(ctx, txHash1, 1)
		network.WaitForTransactionInNodeState(ctx, txHash2, 1)
		require.EqualValues(t, benchmarktoken.TOTAL_SUPPLY-39, contract.CallGetBalance(ctx, 1, 5), "getBalance result on non leader")
		require.EqualValues(t, 39, contract.CallGetBalance(ctx, 1, 6), "getBalance result on non leader")
	})
}
