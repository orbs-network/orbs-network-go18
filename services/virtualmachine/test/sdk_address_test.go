package test

import (
	"context"
	"fmt"
	"github.com/orbs-network/orbs-network-go/crypto/hash"
	"github.com/orbs-network/orbs-network-go/services/processor/native"
	"github.com/orbs-network/orbs-network-go/services/processor/native/repository/_Deployments"
	"github.com/orbs-network/orbs-network-go/test"
	"github.com/orbs-network/orbs-network-go/test/builders"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSdkAddress_GetSignerAddressWithoutContext(t *testing.T) {
	test.WithContext(func(ctx context.Context) {
		h := newHarness()

		_, err := h.handleSdkCall(ctx, 999, native.SDK_OPERATION_NAME_ADDRESS, "getSignerAddress")
		require.Error(t, err, "handleSdkCall should fail")
	})
}

func TestSdkAddress_GetSignerAddressWithoutSignerFails(t *testing.T) {
	test.WithContext(func(ctx context.Context) {
		h := newHarness()
		h.expectSystemContractCalled(deployments_systemcontract.CONTRACT_NAME, deployments_systemcontract.METHOD_GET_INFO, nil, uint32(protocol.PROCESSOR_TYPE_NATIVE)) // assume all contracts are deployed

		h.expectStateStorageBlockHeightRequested(12)
		h.expectNativeContractMethodCalled("Contract1", "method1", func(executionContextId primitives.ExecutionContextId, inputArgs *protocol.MethodArgumentArray) (protocol.ExecutionResult, *protocol.MethodArgumentArray, error) {
			_, err := h.handleSdkCall(ctx, executionContextId, native.SDK_OPERATION_NAME_ADDRESS, "getSignerAddress")
			fmt.Println(err)
			require.Error(t, err, "handleSdkCall should fail since not signed")
			return protocol.EXECUTION_RESULT_SUCCESS, builders.MethodArgumentsArray(), nil
		})

		// runLocalMethod in harness uses nil as the Signer
		h.runLocalMethod(ctx, "Contract1", "method1")

		h.verifySystemContractCalled(t)
		h.verifyStateStorageBlockHeightRequested(t)
		h.verifyNativeContractMethodCalled(t)
	})
}

func TestSdkAddress_GetSignerAddressDoesNotChangeWithContractCalls(t *testing.T) {
	test.WithContext(func(ctx context.Context) {
		h := newHarness()
		h.expectSystemContractCalled(deployments_systemcontract.CONTRACT_NAME, deployments_systemcontract.METHOD_GET_INFO, nil, uint32(protocol.PROCESSOR_TYPE_NATIVE)) // assume all contracts are deployed

		var signerAddressRes []byte

		h.expectNativeContractMethodCalled("Contract1", "method1", func(executionContextId primitives.ExecutionContextId, inputArgs *protocol.MethodArgumentArray) (protocol.ExecutionResult, *protocol.MethodArgumentArray, error) {
			t.Log("GetSignerAddress in the first contract")
			res, err := h.handleSdkCall(ctx, executionContextId, native.SDK_OPERATION_NAME_ADDRESS, "getSignerAddress")
			require.NoError(t, err, "handleSdkCall should succeed")
			require.Equal(t, hash.RIPEMD160_HASH_SIZE_BYTES, len(res[0].BytesValue()), "signer address should be a valid address")
			signerAddressRes = res[0].BytesValue()

			t.Log("CallMethod on a different contract")
			_, err = h.handleSdkCall(ctx, executionContextId, native.SDK_OPERATION_NAME_SERVICE, "callMethod", "Contract2", "method1", builders.MethodArgumentsArray().Raw())
			require.NoError(t, err, "handleSdkCall should succeed")

			return protocol.EXECUTION_RESULT_SUCCESS, builders.MethodArgumentsArray(), nil
		})
		h.expectNativeContractMethodCalled("Contract2", "method1", func(executionContextId primitives.ExecutionContextId, inputArgs *protocol.MethodArgumentArray) (protocol.ExecutionResult, *protocol.MethodArgumentArray, error) {
			t.Log("GetSignerAddress in the second contract")
			res, err := h.handleSdkCall(ctx, executionContextId, native.SDK_OPERATION_NAME_ADDRESS, "getSignerAddress")
			require.NoError(t, err, "handleSdkCall should succeed")
			require.Equal(t, hash.RIPEMD160_HASH_SIZE_BYTES, len(res[0].BytesValue()), "signer address should be a valid address")
			require.Equal(t, signerAddressRes, res[0].BytesValue(), "signer address should be equal to the first call")
			return protocol.EXECUTION_RESULT_SUCCESS, builders.MethodArgumentsArray(), nil
		})

		h.processTransactionSet(ctx, []*contractAndMethod{
			{"Contract1", "method1"},
		})

		h.verifySystemContractCalled(t)
		h.verifyNativeContractMethodCalled(t)
		h.verifyStateStorageRead(t)
	})
}

func TestSdkAddress_GetCallerAddressWithoutSignerFails(t *testing.T) {
	test.WithContext(func(ctx context.Context) {
		h := newHarness()
		h.expectSystemContractCalled(deployments_systemcontract.CONTRACT_NAME, deployments_systemcontract.METHOD_GET_INFO, nil, uint32(protocol.PROCESSOR_TYPE_NATIVE)) // assume all contracts are deployed

		h.expectStateStorageBlockHeightRequested(12)
		h.expectNativeContractMethodCalled("Contract1", "method1", func(executionContextId primitives.ExecutionContextId, inputArgs *protocol.MethodArgumentArray) (protocol.ExecutionResult, *protocol.MethodArgumentArray, error) {
			_, err := h.handleSdkCall(ctx, executionContextId, native.SDK_OPERATION_NAME_ADDRESS, "getCallerAddress")
			fmt.Println(err)
			require.Error(t, err, "handleSdkCall should fail since not signed")
			return protocol.EXECUTION_RESULT_SUCCESS, builders.MethodArgumentsArray(), nil
		})

		// runLocalMethod in harness uses nil as the Signer
		h.runLocalMethod(ctx, "Contract1", "method1")

		h.verifySystemContractCalled(t)
		h.verifyStateStorageBlockHeightRequested(t)
		h.verifyNativeContractMethodCalled(t)
	})
}

func TestSdkAddress_GetCallerAddressChangesWithContractCalls(t *testing.T) {
	test.WithContext(func(ctx context.Context) {
		h := newHarness()
		h.expectSystemContractCalled(deployments_systemcontract.CONTRACT_NAME, deployments_systemcontract.METHOD_GET_INFO, nil, uint32(protocol.PROCESSOR_TYPE_NATIVE)) // assume all contracts are deployed

		var initialCallerAddress []byte
		var firstCallerAddress []byte

		h.expectNativeContractMethodCalled("Contract1", "method1", func(executionContextId primitives.ExecutionContextId, inputArgs *protocol.MethodArgumentArray) (protocol.ExecutionResult, *protocol.MethodArgumentArray, error) {
			t.Log("GetCallerAddress in the first contract (1)")
			res, err := h.handleSdkCall(ctx, executionContextId, native.SDK_OPERATION_NAME_ADDRESS, "getCallerAddress")
			require.NoError(t, err, "handleSdkCall should succeed")
			require.Equal(t, hash.RIPEMD160_HASH_SIZE_BYTES, len(res[0].BytesValue()), "caller address should be a valid address")
			initialCallerAddress = res[0].BytesValue()

			t.Log("CallMethod on a different contract (1->1.2)")
			_, err = h.handleSdkCall(ctx, executionContextId, native.SDK_OPERATION_NAME_SERVICE, "callMethod", "Contract2", "method1", builders.MethodArgumentsArray().Raw())
			require.NoError(t, err, "handleSdkCall should succeed")

			t.Log("CallMethod on a different contract (1->1.4)")
			_, err = h.handleSdkCall(ctx, executionContextId, native.SDK_OPERATION_NAME_SERVICE, "callMethod", "Contract4", "method1", builders.MethodArgumentsArray().Raw())
			require.NoError(t, err, "handleSdkCall should succeed")

			return protocol.EXECUTION_RESULT_SUCCESS, builders.MethodArgumentsArray(), nil
		})
		h.expectNativeContractMethodCalled("Contract2", "method1", func(executionContextId primitives.ExecutionContextId, inputArgs *protocol.MethodArgumentArray) (protocol.ExecutionResult, *protocol.MethodArgumentArray, error) {
			t.Log("GetCallerAddress in the second contract (1.2)")
			res, err := h.handleSdkCall(ctx, executionContextId, native.SDK_OPERATION_NAME_ADDRESS, "getCallerAddress")
			require.NoError(t, err, "handleSdkCall should succeed")
			require.Equal(t, hash.RIPEMD160_HASH_SIZE_BYTES, len(res[0].BytesValue()), "caller address should be a valid address")
			require.NotEqual(t, initialCallerAddress, res[0].BytesValue(), "called address should be different from the initial call")
			firstCallerAddress = res[0].BytesValue()

			t.Log("CallMethod on a different contract (1.2->1.2.3)")
			_, err = h.handleSdkCall(ctx, executionContextId, native.SDK_OPERATION_NAME_SERVICE, "callMethod", "Contract3", "method1", builders.MethodArgumentsArray().Raw())
			require.NoError(t, err, "handleSdkCall should succeed")

			return protocol.EXECUTION_RESULT_SUCCESS, builders.MethodArgumentsArray(), nil
		})
		h.expectNativeContractMethodCalled("Contract3", "method1", func(executionContextId primitives.ExecutionContextId, inputArgs *protocol.MethodArgumentArray) (protocol.ExecutionResult, *protocol.MethodArgumentArray, error) {
			t.Log("GetCallerAddress in the third contract (1.2.3)")
			res, err := h.handleSdkCall(ctx, executionContextId, native.SDK_OPERATION_NAME_ADDRESS, "getCallerAddress")
			require.NoError(t, err, "handleSdkCall should succeed")
			require.Equal(t, hash.RIPEMD160_HASH_SIZE_BYTES, len(res[0].BytesValue()), "caller address should be a valid address")
			require.NotEqual(t, initialCallerAddress, res[0].BytesValue(), "called address should be different from the initial call")
			require.NotEqual(t, firstCallerAddress, res[0].BytesValue(), "called address should be different from the first call")
			return protocol.EXECUTION_RESULT_SUCCESS, builders.MethodArgumentsArray(), nil
		})
		h.expectNativeContractMethodCalled("Contract4", "method1", func(executionContextId primitives.ExecutionContextId, inputArgs *protocol.MethodArgumentArray) (protocol.ExecutionResult, *protocol.MethodArgumentArray, error) {
			t.Log("GetCallerAddress in the fourth contract (1.4)")
			res, err := h.handleSdkCall(ctx, executionContextId, native.SDK_OPERATION_NAME_ADDRESS, "getCallerAddress")
			require.NoError(t, err, "handleSdkCall should succeed")
			require.Equal(t, hash.RIPEMD160_HASH_SIZE_BYTES, len(res[0].BytesValue()), "caller address should be a valid address")
			require.NotEqual(t, initialCallerAddress, res[0].BytesValue(), "called address should be different from the initial call")
			require.Equal(t, firstCallerAddress, res[0].BytesValue(), "called address should be equal to the first call")
			return protocol.EXECUTION_RESULT_SUCCESS, builders.MethodArgumentsArray(), nil
		})

		h.processTransactionSet(ctx, []*contractAndMethod{
			{"Contract1", "method1"},
		})

		h.verifySystemContractCalled(t)
		h.verifyNativeContractMethodCalled(t)
		h.verifyStateStorageRead(t)
	})
}
