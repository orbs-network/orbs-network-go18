package test

import (
	"context"
	"github.com/orbs-network/orbs-network-go/services/processor/native"
	"github.com/orbs-network/orbs-network-go/services/processor/native/repository/_Deployments"
	"github.com/orbs-network/orbs-network-go/test"
	"github.com/orbs-network/orbs-network-go/test/builders"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSdkService_CallMethodFailingCall(t *testing.T) {
	test.WithContext(func(ctx context.Context) {
		h := newHarness()
		h.expectSystemContractCalled(deployments_systemcontract.CONTRACT_NAME, deployments_systemcontract.METHOD_GET_INFO, nil, uint32(protocol.PROCESSOR_TYPE_NATIVE)) // assume all contracts are deployed

		h.expectNativeContractMethodCalled("Contract1", "method1", func(executionContextId primitives.ExecutionContextId, inputArgs *protocol.MethodArgumentArray) (protocol.ExecutionResult, *protocol.MethodArgumentArray, error) {
			t.Log("CallMethod on failing contract")
			_, err := h.handleSdkCall(ctx, executionContextId, native.SDK_OPERATION_NAME_SERVICE, "callMethod", "FailingContract", "method1", builders.MethodArgumentsArray().Raw())
			require.Error(t, err, "handleSdkCall should fail")
			return protocol.EXECUTION_RESULT_SUCCESS, builders.MethodArgumentsArray(), nil
		})
		h.expectNativeContractMethodCalled("FailingContract", "method1", func(executionContextId primitives.ExecutionContextId, inputArgs *protocol.MethodArgumentArray) (protocol.ExecutionResult, *protocol.MethodArgumentArray, error) {
			return protocol.EXECUTION_RESULT_ERROR_UNEXPECTED, builders.MethodArgumentsArray(), errors.New("call error")
		})

		h.processTransactionSet(ctx, []*contractAndMethod{
			{"Contract1", "method1"},
		})

		h.verifySystemContractCalled(t)
		h.verifyNativeContractMethodCalled(t)
	})
}

func TestSdkService_CallMethodMaintainsAddressSpaceUnderSameContract(t *testing.T) {
	test.WithContext(func(ctx context.Context) {
		h := newHarness()
		h.expectSystemContractCalled(deployments_systemcontract.CONTRACT_NAME, deployments_systemcontract.METHOD_GET_INFO, nil, uint32(protocol.PROCESSOR_TYPE_NATIVE)) // assume all contracts are deployed

		h.expectNativeContractMethodCalled("Contract1", "method1", func(executionContextId primitives.ExecutionContextId, inputArgs *protocol.MethodArgumentArray) (protocol.ExecutionResult, *protocol.MethodArgumentArray, error) {
			t.Log("Write to key in first contract")
			_, err := h.handleSdkCall(ctx, executionContextId, native.SDK_OPERATION_NAME_STATE, "write", []byte{0x01}, []byte{0x02, 0x03})
			require.NoError(t, err, "handleSdkCall should succeed")

			t.Log("CallMethod on a the same contract")
			_, err = h.handleSdkCall(ctx, executionContextId, native.SDK_OPERATION_NAME_SERVICE, "callMethod", "Contract1", "method2", builders.MethodArgumentsArray().Raw())
			require.NoError(t, err, "handleSdkCall should succeed")

			return protocol.EXECUTION_RESULT_SUCCESS, builders.MethodArgumentsArray(), nil
		})
		h.expectNativeContractMethodCalled("Contract1", "method2", func(executionContextId primitives.ExecutionContextId, inputArgs *protocol.MethodArgumentArray) (protocol.ExecutionResult, *protocol.MethodArgumentArray, error) {
			t.Log("Read the same key in the first contract")
			res, err := h.handleSdkCall(ctx, executionContextId, native.SDK_OPERATION_NAME_STATE, "read", []byte{0x01})
			require.NoError(t, err, "handleSdkCall should not fail")
			require.Equal(t, []byte{0x02, 0x03}, res[0].BytesValue(), "handleSdkCall result should be equal")
			return protocol.EXECUTION_RESULT_SUCCESS, builders.MethodArgumentsArray(), nil
		})
		h.expectStateStorageNotRead()

		h.processTransactionSet(ctx, []*contractAndMethod{
			{"Contract1", "method1"},
		})

		h.verifySystemContractCalled(t)
		h.verifyNativeContractMethodCalled(t)
		h.verifyStateStorageRead(t)
	})
}

func TestSdkService_CallMethodChangesAddressSpaceBetweenContracts(t *testing.T) {
	test.WithContext(func(ctx context.Context) {
		h := newHarness()
		h.expectSystemContractCalled(deployments_systemcontract.CONTRACT_NAME, deployments_systemcontract.METHOD_GET_INFO, nil, uint32(protocol.PROCESSOR_TYPE_NATIVE)) // assume all contracts are deployed

		h.expectNativeContractMethodCalled("Contract1", "method1", func(executionContextId primitives.ExecutionContextId, inputArgs *protocol.MethodArgumentArray) (protocol.ExecutionResult, *protocol.MethodArgumentArray, error) {
			t.Log("Write to key in first contract")
			_, err := h.handleSdkCall(ctx, executionContextId, native.SDK_OPERATION_NAME_STATE, "write", []byte{0x01}, []byte{0x02, 0x03})
			require.NoError(t, err, "handleSdkCall should succeed")

			t.Log("CallMethod on a different contract")
			_, err = h.handleSdkCall(ctx, executionContextId, native.SDK_OPERATION_NAME_SERVICE, "callMethod", "Contract2", "method1", builders.MethodArgumentsArray().Raw())
			require.NoError(t, err, "handleSdkCall should succeed")

			t.Log("Read the same key in the first contract after the call")
			res, err := h.handleSdkCall(ctx, executionContextId, native.SDK_OPERATION_NAME_STATE, "read", []byte{0x01})
			require.NoError(t, err, "handleSdkCall should not fail")
			require.Equal(t, []byte{0x02, 0x03}, res[0].BytesValue(), "handleSdkCall result should be equal")

			return protocol.EXECUTION_RESULT_SUCCESS, builders.MethodArgumentsArray(), nil
		})
		h.expectNativeContractMethodCalled("Contract2", "method1", func(executionContextId primitives.ExecutionContextId, inputArgs *protocol.MethodArgumentArray) (protocol.ExecutionResult, *protocol.MethodArgumentArray, error) {
			t.Log("Read the same key in the second contract")
			res, err := h.handleSdkCall(ctx, executionContextId, native.SDK_OPERATION_NAME_STATE, "read", []byte{0x01})
			require.NoError(t, err, "handleSdkCall should not fail")
			require.Equal(t, []byte{0x04, 0x05, 0x06}, res[0].BytesValue(), "handleSdkCall result should be equal")
			return protocol.EXECUTION_RESULT_SUCCESS, builders.MethodArgumentsArray(), nil
		})
		h.expectStateStorageRead(11, "Contract2", []byte{0x01}, []byte{0x04, 0x05, 0x06})

		h.processTransactionSet(ctx, []*contractAndMethod{
			{"Contract1", "method1"},
		})

		h.verifySystemContractCalled(t)
		h.verifyNativeContractMethodCalled(t)
		h.verifyStateStorageRead(t)
	})
}

func TestSdkService_CallMethodWithSystemPermissions(t *testing.T) {
	test.WithContext(func(ctx context.Context) {
		h := newHarness()
		h.expectSystemContractCalled(deployments_systemcontract.CONTRACT_NAME, deployments_systemcontract.METHOD_GET_INFO, nil, uint32(protocol.PROCESSOR_TYPE_NATIVE)) // assume all contracts are deployed

		h.expectNativeContractMethodCalled("Contract1", "method1", func(executionContextId primitives.ExecutionContextId, inputArgs *protocol.MethodArgumentArray) (protocol.ExecutionResult, *protocol.MethodArgumentArray, error) {
			t.Log("CallMethod on a different contract with system permissions")
			_, err := h.handleSdkCallWithSystemPermissions(ctx, executionContextId, native.SDK_OPERATION_NAME_SERVICE, "callMethod", "Contract2", "method1", builders.MethodArgumentsArray().Raw())
			require.NoError(t, err, "handleSdkCall should succeed")

			return protocol.EXECUTION_RESULT_SUCCESS, builders.MethodArgumentsArray(), nil
		})
		h.expectNativeContractMethodCalledWithSystemPermissions("Contract2", "method1", func(executionContextId primitives.ExecutionContextId) (protocol.ExecutionResult, *protocol.MethodArgumentArray, error) {
			return protocol.EXECUTION_RESULT_SUCCESS, builders.MethodArgumentsArray(), nil
		})

		h.processTransactionSet(ctx, []*contractAndMethod{
			{"Contract1", "method1"},
		})

		h.verifySystemContractCalled(t)
		h.verifyNativeContractMethodCalled(t)
	})
}

func TestSdkService_CallMethodWithMultipleArguments(t *testing.T) {
	test.WithContext(func(ctx context.Context) {
		h := newHarness()
		h.expectSystemContractCalled(deployments_systemcontract.CONTRACT_NAME, deployments_systemcontract.METHOD_GET_INFO, nil, uint32(protocol.PROCESSOR_TYPE_NATIVE)) // assume all contracts are deployed

		h.expectNativeContractMethodCalled("Contract1", "method1", func(executionContextId primitives.ExecutionContextId, inputArgs *protocol.MethodArgumentArray) (protocol.ExecutionResult, *protocol.MethodArgumentArray, error) {
			t.Log("CallMethod with multiple arguments")
			sdkCallInputArgs := builders.MethodArgumentsArray(uint64(17), "hello").Raw()
			res, err := h.handleSdkCall(ctx, executionContextId, native.SDK_OPERATION_NAME_SERVICE, "callMethod", "Contract2", "method1", sdkCallInputArgs)
			require.NoError(t, err, "handleSdkCall should not fail")
			require.Equal(t, builders.MethodArgumentsArray(uint64(18), "hello2").Raw(), res[0].BytesValue(), "handleSdkCall result should be equal")
			return protocol.EXECUTION_RESULT_SUCCESS, builders.MethodArgumentsArray(), nil
		})
		h.expectNativeContractMethodCalled("Contract2", "method1", func(executionContextId primitives.ExecutionContextId, inputArgs *protocol.MethodArgumentArray) (protocol.ExecutionResult, *protocol.MethodArgumentArray, error) {
			inputArgsIterator := inputArgs.ArgumentsIterator()
			arg0 := inputArgsIterator.NextArguments().Uint64Value() + 1
			arg1 := inputArgsIterator.NextArguments().StringValue() + "2"
			return protocol.EXECUTION_RESULT_SUCCESS, builders.MethodArgumentsArray(arg0, arg1), nil
		})

		h.processTransactionSet(ctx, []*contractAndMethod{
			{"Contract1", "method1"},
		})

		h.verifySystemContractCalled(t)
		h.verifyNativeContractMethodCalled(t)
	})
}
