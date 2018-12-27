package virtualmachine

import (
	"context"
	"github.com/orbs-network/orbs-network-go/instrumentation/log"
	"github.com/orbs-network/orbs-network-go/instrumentation/trace"
	"github.com/orbs-network/orbs-network-go/services/processor/native"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/orbs-network/orbs-spec/types/go/services"
	"github.com/orbs-network/orbs-spec/types/go/services/handlers"
	"github.com/pkg/errors"
)

var LogTag = log.Service("virtual-machine")

type service struct {
	stateStorage         services.StateStorage
	processors           map[protocol.ProcessorType]services.Processor
	crosschainConnectors map[protocol.CrosschainConnectorType]services.CrosschainConnector
	logger               log.BasicLogger

	contexts *executionContextProvider
}

func NewVirtualMachine(
	stateStorage services.StateStorage,
	processors map[protocol.ProcessorType]services.Processor,
	crosschainConnectors map[protocol.CrosschainConnectorType]services.CrosschainConnector,
	logger log.BasicLogger,
) services.VirtualMachine {

	s := &service{
		processors:           processors,
		crosschainConnectors: crosschainConnectors,
		stateStorage:         stateStorage,
		logger:               logger.WithTags(LogTag),

		contexts: newExecutionContextProvider(),
	}

	for _, processor := range processors {
		processor.RegisterContractSdkCallHandler(s)
	}

	return s
}

func (s *service) RunLocalMethod(ctx context.Context, input *services.RunLocalMethodInput) (*services.RunLocalMethodOutput, error) {
	logger := s.logger.WithTags(trace.LogFieldFrom(ctx))

	blockHeight, blockTimestamp, err := s.getRecentBlockHeight(ctx)
	if err != nil {
		return &services.RunLocalMethodOutput{
			CallResult:              protocol.EXECUTION_RESULT_ERROR_UNEXPECTED,
			OutputArgumentArray:     []byte{},
			ReferenceBlockHeight:    blockHeight,
			ReferenceBlockTimestamp: blockTimestamp,
		}, err
	}

	logger.Info("running local method", log.Stringable("contract", input.Transaction.ContractName()), log.Stringable("method", input.Transaction.MethodName()), log.BlockHeight(blockHeight))
	callResult, outputArgs, outputEvents, err := s.runMethod(ctx, blockHeight, blockTimestamp, input.Transaction, protocol.ACCESS_SCOPE_READ_ONLY, nil)
	if outputArgs == nil {
		outputArgs = (&protocol.MethodArgumentArrayBuilder{}).Build()
	}
	if outputEvents == nil {
		outputEvents = (&protocol.EventsArrayBuilder{}).Build()
	}

	return &services.RunLocalMethodOutput{
		CallResult:              callResult,
		OutputEventsArray:       outputEvents.RawEventsArray(),
		OutputArgumentArray:     outputArgs.RawArgumentsArray(),
		ReferenceBlockHeight:    blockHeight,
		ReferenceBlockTimestamp: blockTimestamp,
	}, err
}

func (s *service) ProcessTransactionSet(ctx context.Context, input *services.ProcessTransactionSetInput) (*services.ProcessTransactionSetOutput, error) {
	logger := s.logger.WithTags(trace.LogFieldFrom(ctx))
	previousBlockHeight := input.BlockHeight - 1      // our contracts rely on this block's state for execution (TODO(https://github.com/orbs-network/orbs-spec/issues/123): maybe move the minus one to processTransactionSet so it would match the timestamp)
	blockTimestamp := primitives.TimestampNano(0x777) // TODO(v1): replace placeholder timestamp 0x777 with actual value (probably needs to be added to input in protos

	logger.Info("processing transaction set", log.Int("num-transactions", len(input.SignedTransactions)))
	receipts, stateDiffs := s.processTransactionSet(ctx, previousBlockHeight, blockTimestamp, input.SignedTransactions)

	return &services.ProcessTransactionSetOutput{
		TransactionReceipts: receipts,
		ContractStateDiffs:  stateDiffs,
	}, nil
}

func (s *service) TransactionSetPreOrder(ctx context.Context, input *services.TransactionSetPreOrderInput) (*services.TransactionSetPreOrderOutput, error) {
	logger := s.logger.WithTags(trace.LogFieldFrom(ctx))

	statuses := make([]protocol.TransactionStatus, len(input.SignedTransactions))
	// TODO(v1) sometimes we get value of ffffffffffffffff
	previousBlockHeight := input.BlockHeight - 1      // our contracts rely on this block's state for execution (TODO(https://github.com/orbs-network/orbs-spec/issues/123): maybe move the minus one to callGlobalPreOrderSystemContract so it would match the timestamp)
	blockTimestamp := primitives.TimestampNano(0x777) // TODO(v1): replace placeholder timestamp 0x777 with actual value (probably needs to be added to input in protos)

	// check subscription
	err := s.callGlobalPreOrderSystemContract(ctx, previousBlockHeight, blockTimestamp)
	if err != nil {
		for i := 0; i < len(statuses); i++ {
			statuses[i] = protocol.TRANSACTION_STATUS_REJECTED_GLOBAL_PRE_ORDER
		}
	} else {
		// check signatures
		err = s.verifyTransactionSignatures(input.SignedTransactions, statuses)
	}

	if err != nil {
		logger.Info("performed pre order checks", log.Error(err), log.BlockHeight(previousBlockHeight), log.Int("num-statuses", len(statuses)))
	} else {
		logger.Info("performed pre order checks", log.BlockHeight(previousBlockHeight), log.Int("num-statuses", len(statuses)))
	}

	return &services.TransactionSetPreOrderOutput{
		PreOrderResults: statuses,
	}, err
}

func (s *service) HandleSdkCall(ctx context.Context, input *handlers.HandleSdkCallInput) (*handlers.HandleSdkCallOutput, error) {
	var output []*protocol.MethodArgument
	var err error

	executionContext := s.contexts.loadExecutionContext(input.ContextId)
	if executionContext == nil {
		return nil, errors.Errorf("invalid execution context %s", input.ContextId)
	}

	switch input.OperationName {
	case native.SDK_OPERATION_NAME_STATE:
		output, err = s.handleSdkStateCall(ctx, executionContext, input.MethodName, input.InputArguments, input.PermissionScope)
	case native.SDK_OPERATION_NAME_SERVICE:
		output, err = s.handleSdkServiceCall(ctx, executionContext, input.MethodName, input.InputArguments, input.PermissionScope)
	case native.SDK_OPERATION_NAME_EVENTS:
		output, err = s.handleSdkEventsCall(ctx, executionContext, input.MethodName, input.InputArguments, input.PermissionScope)
	case native.SDK_OPERATION_NAME_ETHEREUM:
		output, err = s.handleSdkEthereumCall(ctx, executionContext, input.MethodName, input.InputArguments, input.PermissionScope)
	case native.SDK_OPERATION_NAME_ADDRESS:
		output, err = s.handleSdkAddressCall(ctx, executionContext, input.MethodName, input.InputArguments, input.PermissionScope)
	default:
		return nil, errors.Errorf("unknown SDK call operation: %s", input.OperationName)
	}

	if err != nil {
		return nil, err
	}

	return &handlers.HandleSdkCallOutput{
		OutputArguments: output,
	}, nil
}
