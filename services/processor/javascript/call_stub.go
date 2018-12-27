// +build !jsprocessor

package javascript

import (
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
)

func (s *service) processMethodCall(executionContextId primitives.ExecutionContextId, code string, methodName primitives.MethodName, args *protocol.MethodArgumentArray) (contractOutputArgs *protocol.MethodArgumentArray, contractOutputErr error, err error) {
	panic("Not implemented")
}
