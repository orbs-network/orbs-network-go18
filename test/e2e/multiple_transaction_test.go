package e2e

import (
	"github.com/orbs-network/orbs-client-sdk-go/codec"
	"github.com/orbs-network/orbs-network-go/crypto/keys"
	"github.com/orbs-network/orbs-network-go/test"
	"github.com/orbs-network/orbs-network-go/test/builders"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNetworkCommitsMultipleTransactions(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E tests in short mode")
	}

	runMultipleTimes(t, func(t *testing.T) {

		h := newHarness()
		lt := time.Now()
		printTestTime(t, "started", &lt)

		h.waitUntilTransactionPoolIsReady(t)
		printTestTime(t, "first block committed", &lt)

		transferTo, _ := keys.GenerateEd25519Key()
		targetAddress := builders.AddressFor(transferTo)

		// send 3 transactions with total of 70
		amounts := []uint64{15, 22, 33}
		txIds := []string{}
		for _, amount := range amounts {
			printTestTime(t, "send transaction - start", &lt)
			response, txId, err := h.sendTransaction(OwnerOfAllSupply, "BenchmarkToken", "transfer", uint64(amount), []byte(targetAddress))
			printTestTime(t, "send transaction - end", &lt)

			txIds = append(txIds, txId)
			require.NoError(t, err, "transaction for amount %d should not return error", amount)
			require.Equal(t, codec.TRANSACTION_STATUS_COMMITTED, response.TransactionStatus)
			require.Equal(t, codec.EXECUTION_RESULT_SUCCESS, response.ExecutionResult)
		}

		// get statuses and receipt proofs
		for _, txId := range txIds {
			printTestTime(t, "get status - start", &lt)
			response, err := h.getTransactionStatus(txId)
			printTestTime(t, "get status - end", &lt)

			require.NoError(t, err, "get status for txid %s should not return error", txId)
			require.Equal(t, codec.TRANSACTION_STATUS_COMMITTED, response.TransactionStatus)
			require.Equal(t, codec.EXECUTION_RESULT_SUCCESS, response.ExecutionResult)

			printTestTime(t, "get receipt proof - start", &lt)
			proofResponse, err := h.getTransactionReceiptProof(txId)
			printTestTime(t, "get receipt proof - end", &lt)

			require.NoError(t, err, "get receipt proof for txid %s should not return error", txId)
			require.Equal(t, codec.TRANSACTION_STATUS_COMMITTED, proofResponse.TransactionStatus)
			require.True(t, len(proofResponse.PackedProof) > 20, "packed receipt proof for txid %s should return at least 20 bytes", txId)
			require.True(t, len(proofResponse.PackedReceipt) > 10, "packed receipt for txid %s should return at least 10 bytes", txId)
		}

		// check balance
		ok := test.Eventually(test.EVENTUALLY_DOCKER_E2E_TIMEOUT, func() bool {
			printTestTime(t, "call method - start", &lt)
			response, err := h.callMethod(transferTo, "BenchmarkToken", "getBalance", []byte(targetAddress))
			printTestTime(t, "call method - end", &lt)

			if err == nil && response.ExecutionResult == codec.EXECUTION_RESULT_SUCCESS {
				return response.OutputArguments[0] == uint64(70)
			}
			return false
		})

		require.True(t, ok, "getBalance should return total amount")
		printTestTime(t, "done", &lt)

	})
}
