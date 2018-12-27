package publicapi

import (
	"context"
	"github.com/orbs-network/orbs-network-go/crypto/digest"
	"github.com/orbs-network/orbs-network-go/instrumentation/log"
	"github.com/orbs-network/orbs-network-go/instrumentation/trace"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/orbs-network/orbs-spec/types/go/protocol/client"
	"github.com/orbs-network/orbs-spec/types/go/services"
	"github.com/pkg/errors"
	"time"
)

func (s *service) CallMethod(parentCtx context.Context, input *services.CallMethodInput) (*services.CallMethodOutput, error) {
	if input.ClientRequest == nil {
		err := errors.Errorf("error: missing input (client request is nil)")
		s.logger.Info("call method received via public api failed", log.Error(err))
		return nil, err
	}

	ctx := trace.NewContext(parentCtx, "PublicApi.CallMethod")
	tx := input.ClientRequest.Transaction()
	txHash := digest.CalcTxHash(tx)
	logger := s.logger.WithTags(trace.LogFieldFrom(ctx), log.Transaction(txHash), log.String("flow", "checkpoint"))

	if txStatus := isTransactionRequestValid(s.config, tx); txStatus != protocol.TRANSACTION_STATUS_RESERVED {
		err := errors.Errorf("error input %s", txStatus.String())
		logger.Info("call method received via public api", log.Error(err))
		return toCallMethodOutput(&services.RunLocalMethodOutput{CallResult: protocol.EXECUTION_RESULT_ERROR_INPUT}), err
	}
	logger.Info("call method request received via public api")

	start := time.Now()
	defer s.metrics.callMethodTime.RecordSince(start)

	result, err := s.virtualMachine.RunLocalMethod(ctx, &services.RunLocalMethodInput{
		Transaction: tx,
	})
	if err != nil {
		logger.Info("call method request failed", log.Error(err))
		return toCallMethodOutput(result), err
	}

	return toCallMethodOutput(result), nil
}

func toCallMethodOutput(output *services.RunLocalMethodOutput) *services.CallMethodOutput {
	response := &client.CallMethodResponseBuilder{
		RequestStatus:       translateExecutionStatusToResponseCode(output.CallResult),
		OutputArgumentArray: output.OutputArgumentArray,
		CallMethodResult:    output.CallResult,
		BlockHeight:         output.ReferenceBlockHeight,
		BlockTimestamp:      output.ReferenceBlockTimestamp,
	}

	return &services.CallMethodOutput{ClientResponse: response.Build()}
}
