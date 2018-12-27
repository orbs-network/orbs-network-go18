package publicapi

import (
	"context"
	"github.com/orbs-network/orbs-network-go/instrumentation/log"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/orbs-network/orbs-spec/types/go/protocol/client"
	"github.com/orbs-network/orbs-spec/types/go/services"
	"github.com/pkg/errors"
	"time"
)

func (s *service) GetTransactionStatus(ctx context.Context, input *services.GetTransactionStatusInput) (*services.GetTransactionStatusOutput, error) {
	if input.ClientRequest == nil {
		err := errors.Errorf("error: missing input (client request is nil)")
		s.logger.Info("get transaction status received missing input", log.Error(err))
		return nil, err
	}

	s.logger.Info("get transaction status request received", log.String("flow", "checkpoint"), log.Transaction(input.ClientRequest.Txhash()))
	txHash := input.ClientRequest.Txhash()
	timestamp := input.ClientRequest.TransactionTimestamp()

	return s.getTransactionStatus(ctx, txHash, timestamp)
}

func (s *service) getTransactionStatus(ctx context.Context, txHash primitives.Sha256, timestamp primitives.TimestampNano) (*services.GetTransactionStatusOutput, error) {
	start := time.Now()
	defer s.metrics.getTransactionStatusTime.RecordSince(start)

	if txReceipt, err, ok := s.getFromTxPool(ctx, txHash, timestamp); ok {
		return toGetTxOutput(txReceipt), err
	}

	blockReceipt, err := s.getFromBlockStorage(ctx, txHash, timestamp)
	if err != nil {
		return nil, err
	}
	return toGetTxOutput(blockReceipt), err
}

func (s *service) getFromTxPool(ctx context.Context, txHash primitives.Sha256, timestamp primitives.TimestampNano) (*txResponse, error, bool) {
	txReceipt, err := s.transactionPool.GetCommittedTransactionReceipt(ctx, &services.GetCommittedTransactionReceiptInput{
		Txhash:               txHash,
		TransactionTimestamp: timestamp,
	})
	if err != nil {
		s.logger.Info("get transaction status failed in transactionPool", log.Error(err), log.String("flow", "checkpoint"), log.Transaction(txHash))
		return txStatusToTxResponse(txReceipt), err, true
	}
	if txReceipt.TransactionStatus != protocol.TRANSACTION_STATUS_NO_RECORD_FOUND {
		return txStatusToTxResponse(txReceipt), nil, true
	}
	return nil, nil, false
}

func txStatusToTxResponse(txStatus *services.GetCommittedTransactionReceiptOutput) *txResponse {
	return &txResponse{
		transactionStatus:  txStatus.TransactionStatus,
		transactionReceipt: txStatus.TransactionReceipt,
		blockHeight:        txStatus.BlockHeight,
		blockTimestamp:     txStatus.BlockTimestamp,
	}
}

func (s *service) getFromBlockStorage(ctx context.Context, txHash primitives.Sha256, timestamp primitives.TimestampNano) (*txResponse, error) {
	txReceipt, err := s.blockStorage.GetTransactionReceipt(ctx, &services.GetTransactionReceiptInput{
		Txhash:               txHash,
		TransactionTimestamp: timestamp,
	})
	if err != nil {
		s.logger.Info("get transaction status failed in blockStorage", log.Error(err), log.String("flow", "checkpoint"), log.Transaction(txHash))
		return nil, err
	}
	return blockToTxResponse(txReceipt), nil

}

func blockToTxResponse(bReceipt *services.GetTransactionReceiptOutput) *txResponse {
	status := protocol.TRANSACTION_STATUS_NO_RECORD_FOUND
	if bReceipt.TransactionReceipt != nil {
		status = protocol.TRANSACTION_STATUS_COMMITTED
	}
	return &txResponse{
		transactionStatus:  status,
		transactionReceipt: bReceipt.TransactionReceipt,
		blockHeight:        bReceipt.BlockHeight,
		blockTimestamp:     bReceipt.BlockTimestamp,
	}
}

func toGetTxOutput(transactionOutput *txResponse) *services.GetTransactionStatusOutput {
	var receiptForClient *protocol.TransactionReceiptBuilder = nil
	if transactionOutput.transactionReceipt != nil {
		receiptForClient = protocol.TransactionReceiptBuilderFromRaw(transactionOutput.transactionReceipt.Raw())
	}

	response := &client.GetTransactionStatusResponseBuilder{
		RequestStatus:      translateTxStatusToResponseCode(transactionOutput.transactionStatus),
		TransactionReceipt: receiptForClient,
		TransactionStatus:  transactionOutput.transactionStatus,
		BlockHeight:        transactionOutput.blockHeight,
		BlockTimestamp:     transactionOutput.blockTimestamp,
	}

	return &services.GetTransactionStatusOutput{ClientResponse: response.Build()}
}
