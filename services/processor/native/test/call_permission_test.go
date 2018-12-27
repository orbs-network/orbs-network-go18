package test

import (
	"context"
	"github.com/orbs-network/orbs-network-go/test"
	"github.com/orbs-network/orbs-spec/types/go/services"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProcessCall_Permissions(t *testing.T) {
	tests := []struct {
		name             string
		input            *services.ProcessCallInput
		expectedError    bool
		expectedSdkWrite bool
	}{
		{
			name:          "UnknownMethodFails",
			input:         processCallInput().WithUnknownMethod().Build(),
			expectedError: true,
		},
		{
			name:          "PublicMethodSucceeds",
			input:         processCallInput().WithPublicMethod().Build(),
			expectedError: false,
		},
		{
			name:          "InternalMethodFails",
			input:         processCallInput().WithInternalMethod().Build(),
			expectedError: true,
		},
		{
			name:          "SystemMethodFails",
			input:         processCallInput().WithSystemMethod().Build(),
			expectedError: true,
		},
		{
			name:             "SystemMethodUnderSystemPermissionsSucceeds",
			input:            processCallInput().WithSystemMethod().WithSystemPermissions().Build(),
			expectedError:    false,
			expectedSdkWrite: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.WithContext(func(ctx context.Context) {
				h := newHarness()
				if tt.expectedSdkWrite {
					h.expectSdkCallMadeWithStateWrite(nil, nil)
				}

				_, err := h.service.ProcessCall(ctx, tt.input)
				if tt.expectedError {
					require.Error(t, err, "call should fail")
				} else {
					require.NoError(t, err, "call should succeed")
				}
			})
		})
	}
}
