package test

import (
	"context"
	"github.com/orbs-network/go-mock"
	"github.com/orbs-network/orbs-network-go/config"
	"github.com/orbs-network/orbs-network-go/instrumentation/log"
	"github.com/orbs-network/orbs-network-go/instrumentation/metric"
	"github.com/orbs-network/orbs-network-go/services/publicapi"
	"github.com/orbs-network/orbs-network-go/test/builders"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/orbs-network/orbs-spec/types/go/services"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

type harness struct {
	papi    services.PublicApi
	txpMock *services.MockTransactionPool
	bksMock *services.MockBlockStorage
	vmMock  *services.MockVirtualMachine
}

func newPublicApiHarness(ctx context.Context, txTimeout time.Duration) *harness {
	logger := log.GetLogger().WithOutput(log.NewFormattingOutput(os.Stdout, log.NewHumanReadableFormatter()))
	cfg := config.ForPublicApiTests(uint32(builders.DEFAULT_TEST_VIRTUAL_CHAIN_ID), txTimeout)
	txpMock := makeTxMock()
	vmMock := &services.MockVirtualMachine{}
	bksMock := &services.MockBlockStorage{}
	papi := publicapi.NewPublicApi(cfg, txpMock, vmMock, bksMock, logger, metric.NewRegistry())
	return &harness{
		papi:    papi,
		txpMock: txpMock,
		bksMock: bksMock,
		vmMock:  vmMock,
	}
}

func makeTxMock() *services.MockTransactionPool {
	txpMock := &services.MockTransactionPool{}
	txpMock.When("RegisterTransactionResultsHandler", mock.Any).Return(nil)
	return txpMock
}

func (h *harness) addTransactionReturnsAlreadyCommitted() {
	h.txpMock.When("AddNewTransaction", mock.Any, mock.Any).Return(&services.AddNewTransactionOutput{
		TransactionStatus:  protocol.TRANSACTION_STATUS_DUPLICATE_TRANSACTION_ALREADY_COMMITTED,
		TransactionReceipt: builders.TransactionReceipt().Build(),
	}).Times(1)
}

func (h *harness) onAddNewTransaction(f func()) {
	h.txpMock.When("AddNewTransaction", mock.Any, mock.Any).Times(1).
		Call(func(ctx context.Context, input *services.AddNewTransactionInput) (*services.AddNewTransactionOutput, error) {
			go func() {
				time.Sleep(1 * time.Millisecond)
				f()
			}()
			return &services.AddNewTransactionOutput{TransactionStatus: protocol.TRANSACTION_STATUS_PENDING}, nil
		})
}

func (h *harness) transactionIsPendingInPool() {
	h.txpMock.When("GetCommittedTransactionReceipt", mock.Any, mock.Any).Return(&services.GetCommittedTransactionReceiptOutput{
		TransactionStatus: protocol.TRANSACTION_STATUS_PENDING,
	}).Times(1)
	h.bksMock.Never("GetTransactionReceipt", mock.Any)
}

func (h *harness) transactionIsCommitedInPool() {
	h.txpMock.When("GetCommittedTransactionReceipt", mock.Any, mock.Any).Return(&services.GetCommittedTransactionReceiptOutput{
		TransactionStatus:  protocol.TRANSACTION_STATUS_COMMITTED,
		TransactionReceipt: builders.TransactionReceipt().Build(),
	}).Times(1)
	h.bksMock.Never("GetTransactionReceipt", mock.Any)
}

func (h *harness) transactionIsNotInPool() {
	h.txpMock.When("GetCommittedTransactionReceipt", mock.Any, mock.Any).Return(&services.GetCommittedTransactionReceiptOutput{
		TransactionStatus: protocol.TRANSACTION_STATUS_NO_RECORD_FOUND,
	}).Times(1)
}

func (h *harness) transactionIsNotInPoolIsInBlockStorage() {
	h.transactionIsNotInPool()
	h.bksMock.When("GetTransactionReceipt", mock.Any, mock.Any).Return(
		&services.GetTransactionReceiptOutput{
			TransactionReceipt: builders.TransactionReceipt().Build(),
		}).Times(1)
}

func (h *harness) runTransactionSuccess() {
	h.vmMock.When("RunLocalMethod", mock.Any, mock.Any).Times(1).
		Return(&services.RunLocalMethodOutput{
			CallResult:          protocol.EXECUTION_RESULT_SUCCESS,
			OutputArgumentArray: nil,
		})
}

func (h *harness) transactionHasProof() {
	h.transactionIsCommitedInPool()
	h.bksMock.When("GenerateReceiptProof", mock.Any, mock.Any).Return(
		&services.GenerateReceiptProofOutput{
			Proof: (&protocol.ReceiptProofBuilder{}).Build(),
		}).Times(1)
}

func (h *harness) transactionPendingNoProofCalled() {
	h.transactionIsPendingInPool()
	h.bksMock.Never("GenerateReceiptProof", mock.Any)
}

func (h *harness) getTransactionStatusFailed() {
	h.transactionIsNotInPool()
	h.bksMock.When("GetTransactionReceipt", mock.Any, mock.Any).Return(nil, errors.Errorf("stam")).Times(1)
	h.bksMock.Never("GenerateReceiptProof", mock.Any)
}

func (h *harness) verifyMocks(t *testing.T) {
	// contract test
	ok, errCalled := h.txpMock.Verify()
	require.True(t, ok, "txPool mock called incorrectly")
	require.NoError(t, errCalled, "error happened when it should not")
	ok, errCalled = h.bksMock.Verify()
	require.True(t, ok, "block storage mock called incorrectly")
	require.NoError(t, errCalled, "error happened when it should not")
	ok, errCalled = h.vmMock.Verify()
	require.True(t, ok, "virtual machine mock called incorrectly")
	require.NoError(t, errCalled, "error happened when it should not")
}
