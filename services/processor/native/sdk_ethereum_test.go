package native

import (
	"context"
	sdkContext "github.com/orbs-network/orbs-contract-sdk/go/context"
	"github.com/orbs-network/orbs-network-go/test/builders"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/orbs-network/orbs-spec/types/go/services/handlers"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"testing"
)

// this example represents respones from "say" function with data "etherworld"
var examplePackedOutput = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 104, 101, 108, 108, 111, 32, 101, 116, 104, 101, 114, 119, 111, 114, 108, 100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

func TestSdkEthereum_CallMethod(t *testing.T) {
	s := createEthereumSdk()

	var out string
	sampleABI := `[{"inputs":[],"name":"say","outputs":[{"name":"","type":"string"}],"type":"function"}]`
	s.SdkEthereumCallMethod(EXAMPLE_CONTEXT, sdkContext.PERMISSION_SCOPE_SYSTEM, "ExampleAddress", sampleABI, "say", &out)

	require.Equal(t, "hello etherworld", out, "did not get the expected return value from ethereum call")
}

func TestSdkEthereum_GetTransactionLog(t *testing.T) {
	s := createEthereumSdk()

	var out string
	sampleABI := `[{"inputs":[{"name":"sentence","type":"string"}],"name":"said","type":"event"}]`
	s.SdkEthereumGetTransactionLog(EXAMPLE_CONTEXT, sdkContext.PERMISSION_SCOPE_SYSTEM, "ExampleAddress", sampleABI, "0x88df016429689c079f3b2f6ad39fa052532c56795b733da78a91ebe6a713944b", "said", &out)

	require.Equal(t, "hello etherworld", out, "did not get the expected return value from transaction log")
}

func createEthereumSdk() *service {
	return &service{sdkHandler: &contractSdkEthereumCallHandlerStub{}}
}

type contractSdkEthereumCallHandlerStub struct{}

func (c *contractSdkEthereumCallHandlerStub) HandleSdkCall(ctx context.Context, input *handlers.HandleSdkCallInput) (*handlers.HandleSdkCallOutput, error) {
	if input.PermissionScope != protocol.PERMISSION_SCOPE_SYSTEM {
		panic("permissions passed to SDK are incorrect")
	}
	switch input.MethodName {
	case "callMethod":
		return &handlers.HandleSdkCallOutput{
			OutputArguments: builders.MethodArguments(examplePackedOutput),
		}, nil
	case "getTransactionLog":
		return &handlers.HandleSdkCallOutput{
			OutputArguments: builders.MethodArguments(examplePackedOutput),
		}, nil
	default:
		return nil, errors.New("unknown method")
	}
}
