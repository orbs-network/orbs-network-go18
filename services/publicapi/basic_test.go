package publicapi

import (
	"fmt"
	"github.com/orbs-network/orbs-network-go/config"
	"github.com/orbs-network/orbs-network-go/test/builders"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPublicApiSendTx_TranslateTransactionStatusToHttpCode(t *testing.T) {
	tests := []struct {
		name   string
		expect protocol.RequestStatus
		status protocol.TransactionStatus
	}{
		{"TRANSACTION_STATUS_RESERVED", protocol.REQUEST_STATUS_RESERVED, protocol.TRANSACTION_STATUS_RESERVED},
		{"TRANSACTION_STATUS_COMMITTED", protocol.REQUEST_STATUS_COMPLETED, protocol.TRANSACTION_STATUS_COMMITTED},
		{"TRANSACTION_STATUS_DUPLICATE_TRANSACTION_ALREADY_COMMITTED", protocol.REQUEST_STATUS_COMPLETED, protocol.TRANSACTION_STATUS_DUPLICATE_TRANSACTION_ALREADY_COMMITTED},
		{"TRANSACTION_STATUS_PENDING", protocol.REQUEST_STATUS_IN_PROCESS, protocol.TRANSACTION_STATUS_PENDING},
		{"TRANSACTION_STATUS_DUPLICATE_TRANSACTION_ALREADY_PENDING", protocol.REQUEST_STATUS_IN_PROCESS, protocol.TRANSACTION_STATUS_DUPLICATE_TRANSACTION_ALREADY_PENDING},
		{"TRANSACTION_STATUS_PRE_ORDER_VALID", protocol.REQUEST_STATUS_RESERVED, protocol.TRANSACTION_STATUS_PRE_ORDER_VALID},
		{"TRANSACTION_STATUS_NO_RECORD_FOUND", protocol.REQUEST_STATUS_NOT_FOUND, protocol.TRANSACTION_STATUS_NO_RECORD_FOUND},
		{"TRANSACTION_STATUS_REJECTED_UNSUPPORTED_VERSION", protocol.REQUEST_STATUS_REJECTED, protocol.TRANSACTION_STATUS_REJECTED_UNSUPPORTED_VERSION},
		{"TRANSACTION_STATUS_REJECTED_VIRTUAL_CHAIN_MISMATCH", protocol.REQUEST_STATUS_REJECTED, protocol.TRANSACTION_STATUS_REJECTED_VIRTUAL_CHAIN_MISMATCH},
		{"TRANSACTION_STATUS_REJECTED_TIMESTAMP_WINDOW_EXCEEDED", protocol.REQUEST_STATUS_REJECTED, protocol.TRANSACTION_STATUS_REJECTED_TIMESTAMP_WINDOW_EXCEEDED},
		{"TRANSACTION_STATUS_REJECTED_SIGNATURE_MISMATCH", protocol.REQUEST_STATUS_REJECTED, protocol.TRANSACTION_STATUS_REJECTED_SIGNATURE_MISMATCH},
		{"TRANSACTION_STATUS_REJECTED_UNKNOWN_SIGNER_SCHEME", protocol.REQUEST_STATUS_REJECTED, protocol.TRANSACTION_STATUS_REJECTED_UNKNOWN_SIGNER_SCHEME},
		{"TRANSACTION_STATUS_REJECTED_GLOBAL_PRE_ORDER", protocol.REQUEST_STATUS_REJECTED, protocol.TRANSACTION_STATUS_REJECTED_GLOBAL_PRE_ORDER},
		{"TRANSACTION_STATUS_REJECTED_VIRTUAL_CHAIN_PRE_ORDER", protocol.REQUEST_STATUS_REJECTED, protocol.TRANSACTION_STATUS_REJECTED_VIRTUAL_CHAIN_PRE_ORDER},
		{"TRANSACTION_STATUS_REJECTED_SMART_CONTRACT_PRE_ORDER", protocol.REQUEST_STATUS_REJECTED, protocol.TRANSACTION_STATUS_REJECTED_SMART_CONTRACT_PRE_ORDER},
		{"TRANSACTION_STATUS_REJECTED_TIMESTAMP_PRECEDES_NODE_TIME", protocol.REQUEST_STATUS_REJECTED, protocol.TRANSACTION_STATUS_REJECTED_TIMESTAMP_AHEAD_OF_NODE_TIME},
		{"TRANSACTION_STATUS_REJECTED_CONGESTION", protocol.REQUEST_STATUS_CONGESTION, protocol.TRANSACTION_STATUS_REJECTED_CONGESTION},
	}
	for i := range tests {
		currTest := tests[i] // this is so that we can run tests in parallel, see https://gist.github.com/posener/92a55c4cd441fc5e5e85f27bca008721
		t.Run(currTest.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, currTest.expect, translateTxStatusToResponseCode(currTest.status), fmt.Sprintf("%s was translated to %d", currTest.name, currTest.expect))
		})
	}
}

func TestPublicApi_TranslateExecutionStatusToHttpCode(t *testing.T) {
	tests := []struct {
		name   string
		expect protocol.RequestStatus
		status protocol.ExecutionResult
	}{
		{"EXECUTION_RESULT_RESERVED", protocol.REQUEST_STATUS_RESERVED, protocol.EXECUTION_RESULT_RESERVED},
		{"EXECUTION_RESULT_SUCCESS", protocol.REQUEST_STATUS_COMPLETED, protocol.EXECUTION_RESULT_SUCCESS},
		{"EXECUTION_RESULT_ERROR_SMART_CONTRACT", protocol.REQUEST_STATUS_COMPLETED, protocol.EXECUTION_RESULT_ERROR_SMART_CONTRACT},
		{"EXECUTION_RESULT_ERROR_INPUT", protocol.REQUEST_STATUS_REJECTED, protocol.EXECUTION_RESULT_ERROR_INPUT},
		{"EXECUTION_RESULT_ERROR_UNEXPECTED", protocol.REQUEST_STATUS_SYSTEM_ERROR, protocol.EXECUTION_RESULT_ERROR_UNEXPECTED},
	}
	for i := range tests {
		currTest := tests[i] // this is so that we can run tests in parallel, see https://gist.github.com/posener/92a55c4cd441fc5e5e85f27bca008721
		t.Run(currTest.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, currTest.expect, translateExecutionStatusToResponseCode(currTest.status), fmt.Sprintf("%s was translated to %d", currTest.name, currTest.expect))
		})
	}
}

func TestPublicApiBasic_IsRequestValidChain(t *testing.T) {
	cfg := config.ForPublicApiTests(6, 0)
	status := isTransactionRequestValid(cfg, builders.Transaction().WithVirtualChainId(6).Build().Transaction())

	require.EqualValues(t, protocol.TRANSACTION_STATUS_RESERVED, status, "virtual chain should be ok")
}

func TestPublicApiBasic_IsRequestValidChainNonValid(t *testing.T) {
	cfg := config.ForPublicApiTests(44, 0)
	status := isTransactionRequestValid(cfg, builders.Transaction().WithVirtualChainId(6).Build().Transaction())

	require.EqualValues(t, protocol.TRANSACTION_STATUS_REJECTED_VIRTUAL_CHAIN_MISMATCH, status, "virtual chain should be wrong")
}
