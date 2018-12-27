package test

import (
	"context"
	"github.com/orbs-network/orbs-network-go/test"
	"github.com/orbs-network/orbs-network-go/test/builders"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/orbs-network/orbs-spec/types/go/services"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProcessCall_Errors(t *testing.T) {
	tests := []struct {
		name           string
		input          *services.ProcessCallInput
		expectedError  bool
		expectedResult protocol.ExecutionResult
		expectedOutput *protocol.MethodArgumentArray
	}{
		{
			name:           "ThatThrowsError",
			input:          processCallInput().WithMethod("BenchmarkContract", "throw").Build(),
			expectedError:  true,
			expectedResult: protocol.EXECUTION_RESULT_ERROR_SMART_CONTRACT,
			expectedOutput: builders.MethodArgumentsArray("example error returned by contract"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.WithContext(func(ctx context.Context) {
				h := newHarness()

				output, err := h.service.ProcessCall(ctx, tt.input)
				if tt.expectedError {
					require.Error(t, err, "call should fail")
					require.Equal(t, tt.expectedOutput, output.OutputArgumentArray, "call return args should be equal")
				} else {
					require.NoError(t, err, "call should succeed")
					require.Equal(t, tt.expectedOutput, output.OutputArgumentArray, "call return args should be equal")
				}
				require.Equal(t, tt.expectedResult, output.CallResult, "call result should be equal")
			})
		})
	}
}
