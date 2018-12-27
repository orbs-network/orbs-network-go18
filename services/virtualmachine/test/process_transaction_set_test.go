package test

import (
	"context"
	"github.com/orbs-network/orbs-network-go/services/processor/native/repository/_Deployments"
	"github.com/orbs-network/orbs-network-go/test"
	"github.com/orbs-network/orbs-network-go/test/builders"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProcessTransactionSet_Success(t *testing.T) {
	test.WithContext(func(ctx context.Context) {
		h := newHarness()
		h.expectSystemContractCalled(deployments_systemcontract.CONTRACT_NAME, deployments_systemcontract.METHOD_GET_INFO, nil, uint32(protocol.PROCESSOR_TYPE_NATIVE)) // assume all contracts are deployed

		h.expectNativeContractMethodCalled("Contract1", "method1", func(executionContextId primitives.ExecutionContextId, inputArgs *protocol.MethodArgumentArray) (protocol.ExecutionResult, *protocol.MethodArgumentArray, error) {
			t.Log("Transaction 1: successful")
			return protocol.EXECUTION_RESULT_SUCCESS, builders.MethodArgumentsArray(), nil
		})
		h.expectNativeContractMethodCalled("Contract1", "method2", func(executionContextId primitives.ExecutionContextId, inputArgs *protocol.MethodArgumentArray) (protocol.ExecutionResult, *protocol.MethodArgumentArray, error) {
			t.Log("Transaction 2: successful")
			return protocol.EXECUTION_RESULT_SUCCESS, builders.MethodArgumentsArray(uint32(17), "hello", []byte{0x01, 0x02}), nil
		})

		results, outputArgs, _, _ := h.processTransactionSet(ctx, []*contractAndMethod{
			{"Contract1", "method1"},
			{"Contract1", "method2"},
		})
		require.Equal(t, results, []protocol.ExecutionResult{
			protocol.EXECUTION_RESULT_SUCCESS,
			protocol.EXECUTION_RESULT_SUCCESS,
		}, "processTransactionSet returned receipts should match")
		require.Equal(t, outputArgs, [][]byte{
			{},
			builders.PackedArgumentArrayEncode(uint32(17), "hello", []byte{0x01, 0x02}),
		}, "processTransactionSet returned output args should match")

		h.verifySystemContractCalled(t)
		h.verifyNativeContractMethodCalled(t)
	})
}

func TestProcessTransactionSet_WithErrors(t *testing.T) {
	test.WithContext(func(ctx context.Context) {
		h := newHarness()
		h.expectSystemContractCalled(deployments_systemcontract.CONTRACT_NAME, deployments_systemcontract.METHOD_GET_INFO, nil, uint32(protocol.PROCESSOR_TYPE_NATIVE)) // assume all contracts are deployed

		h.expectNativeContractMethodCalled("Contract1", "method1", func(executionContextId primitives.ExecutionContextId, inputArgs *protocol.MethodArgumentArray) (protocol.ExecutionResult, *protocol.MethodArgumentArray, error) {
			t.Log("Transaction 1: failed (contract error)")
			return protocol.EXECUTION_RESULT_ERROR_SMART_CONTRACT, builders.MethodArgumentsArray(), errors.New("contract error")
		})
		h.expectNativeContractMethodCalled("Contract1", "method2", func(executionContextId primitives.ExecutionContextId, inputArgs *protocol.MethodArgumentArray) (protocol.ExecutionResult, *protocol.MethodArgumentArray, error) {
			t.Log("Transaction 2: failed (unexpected error)")
			return protocol.EXECUTION_RESULT_ERROR_UNEXPECTED, builders.MethodArgumentsArray(), errors.New("unexpected error")
		})

		results, outputArgs, _, _ := h.processTransactionSet(ctx, []*contractAndMethod{
			{"Contract1", "method1"},
			{"Contract1", "method2"},
		})
		require.Equal(t, results, []protocol.ExecutionResult{
			protocol.EXECUTION_RESULT_ERROR_SMART_CONTRACT,
			protocol.EXECUTION_RESULT_ERROR_UNEXPECTED,
		}, "processTransactionSet returned receipts should match")
		require.Equal(t, outputArgs, [][]byte{
			{},
			{},
		}, "processTransactionSet returned output args should match")

		h.verifySystemContractCalled(t)
		h.verifyNativeContractMethodCalled(t)
	})
}
