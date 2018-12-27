package test

import (
	"context"
	"github.com/orbs-network/go-mock"
	"github.com/orbs-network/orbs-network-go/config"
	"github.com/orbs-network/orbs-network-go/crypto/digest"
	"github.com/orbs-network/orbs-network-go/instrumentation/log"
	"github.com/orbs-network/orbs-network-go/instrumentation/metric"
	"github.com/orbs-network/orbs-network-go/services/transactionpool"
	testKeys "github.com/orbs-network/orbs-network-go/test/crypto/keys"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/orbs-network/orbs-spec/types/go/protocol/gossipmessages"
	"github.com/orbs-network/orbs-spec/types/go/services"
	"github.com/orbs-network/orbs-spec/types/go/services/gossiptopics"
	"github.com/orbs-network/orbs-spec/types/go/services/handlers"
	"time"
)

type harness struct {
	txpool             services.TransactionPool
	gossip             *gossiptopics.MockTransactionRelay
	vm                 *services.MockVirtualMachine
	trh                *handlers.MockTransactionResultsHandler
	lastBlockHeight    primitives.BlockHeight
	lastBlockTimestamp primitives.TimestampNano
	config             config.TransactionPoolConfig
}

var now = primitives.TimestampNano(time.Now().UnixNano())

var (
	thisNodeKeyPair  = testKeys.EcdsaSecp256K1KeyPairForTests(8)
	otherNodeKeyPair = testKeys.EcdsaSecp256K1KeyPairForTests(9)
)

func (h *harness) expectTransactionsToBeForwarded(sig primitives.EcdsaSecp256K1Sig, transactions ...*protocol.SignedTransaction) {

	h.gossip.When("BroadcastForwardedTransactions", mock.Any, &gossiptopics.ForwardedTransactionsInput{
		Message: &gossipmessages.ForwardedTransactionsMessage{
			Sender: (&gossipmessages.SenderSignatureBuilder{
				SenderNodeAddress: thisNodeKeyPair.NodeAddress(),
				Signature:         sig,
			}).Build(),
			SignedTransactions: transactions,
		},
	}).Return(&gossiptopics.EmptyOutput{}, nil).Times(1)
}

func (h *harness) expectNoTransactionsToBeForwarded() {
	h.gossip.Never("BroadcastForwardedTransactions", mock.Any, mock.Any)
}

func (h *harness) ignoringForwardMessages() {
	h.gossip.When("BroadcastForwardedTransactions", mock.Any, mock.Any).Return(&gossiptopics.EmptyOutput{}, nil).AtLeast(0)
}

func (h *harness) addNewTransaction(ctx context.Context, tx *protocol.SignedTransaction) (*services.AddNewTransactionOutput, error) {
	out, err := h.txpool.AddNewTransaction(ctx, &services.AddNewTransactionInput{
		SignedTransaction: tx,
	})

	return out, err
}

func (h *harness) addTransactions(ctx context.Context, txs ...*protocol.SignedTransaction) {
	for _, tx := range txs {
		h.addNewTransaction(ctx, tx)
	}
}

func (h *harness) reportTransactionsAsCommitted(ctx context.Context, transactions ...*protocol.SignedTransaction) (*services.CommitTransactionReceiptsOutput, error) {
	return h.txpool.CommitTransactionReceipts(ctx, &services.CommitTransactionReceiptsInput{
		LastCommittedBlockHeight: h.lastBlockHeight + 1,
		ResultsBlockHeader:       (&protocol.ResultsBlockHeaderBuilder{Timestamp: h.lastBlockTimestamp, BlockHeight: h.lastBlockHeight}).Build(), //TODO ResultsBlockHeader is too much info here, awaiting change in proto, see issue #121
		TransactionReceipts:      asReceipts(transactions),
	})

}

func (h *harness) verifyMocks() error {
	if _, err := h.gossip.Verify(); err != nil {
		return err
	}

	if _, err := h.trh.Verify(); err != nil {
		return err
	}

	if _, err := h.vm.Verify(); err != nil {
		return err
	}

	return nil
}

func (h *harness) handleForwardFrom(ctx context.Context, sender *testKeys.TestEcdsaSecp256K1KeyPair, transactions ...*protocol.SignedTransaction) {
	oneBigHash, _, _ := transactionpool.HashTransactions(transactions...)

	sig, err := digest.SignAsNode(sender.PrivateKey(), oneBigHash)
	if err != nil {
		panic(err)
	}

	h.txpool.HandleForwardedTransactions(ctx, &gossiptopics.ForwardedTransactionsInput{
		Message: &gossipmessages.ForwardedTransactionsMessage{
			Sender: (&gossipmessages.SenderSignatureBuilder{
				SenderNodeAddress: sender.NodeAddress(),
				Signature:         sig,
			}).Build(),
			SignedTransactions: transactions,
		},
	})
}
func (h *harness) expectTransactionResultsCallbackFor(transactions ...*protocol.SignedTransaction) {
	h.trh.When("HandleTransactionResults", mock.Any, &handlers.HandleTransactionResultsInput{
		BlockHeight:         h.lastBlockHeight,
		Timestamp:           h.lastBlockTimestamp,
		TransactionReceipts: asReceipts(transactions),
	}).Times(1).Return(&handlers.HandleTransactionResultsOutput{}, nil)
}

func (h *harness) expectTransactionErrorCallbackFor(tx *protocol.SignedTransaction, status protocol.TransactionStatus) {
	txHash := digest.CalcTxHash(tx.Transaction())
	h.trh.When("HandleTransactionError", mock.Any, mock.AnyIf("transaction error matching the given transaction", func(i interface{}) bool {
		tri := i.(*handlers.HandleTransactionErrorInput)
		return tri.Txhash.Equal(txHash) && tri.TransactionStatus == status
	})).Return(&handlers.HandleTransactionErrorOutput{}).Times(1)
}

func (h *harness) ignoringTransactionResults() {
	h.trh.When("HandleTransactionResults", mock.Any, mock.Any)
	h.trh.When("HandleTransactionError", mock.Any, mock.Any)
}

func (h *harness) assumeBlockStorageAtHeight(height primitives.BlockHeight) {
	h.lastBlockHeight = height
	h.lastBlockTimestamp = primitives.TimestampNano(time.Now().UnixNano())
}

func (h *harness) getTransactionsForOrdering(ctx context.Context, height primitives.BlockHeight, maxNumOfTransactions uint32) (*services.GetTransactionsForOrderingOutput, error) {
	return h.txpool.GetTransactionsForOrdering(ctx, &services.GetTransactionsForOrderingInput{
		BlockHeight:             height,
		MaxNumberOfTransactions: maxNumOfTransactions,
	})
}

func (h *harness) failPreOrderCheckFor(failOn func(tx *protocol.SignedTransaction) bool) {
	h.vm.Reset().When("TransactionSetPreOrder", mock.Any, mock.Any).Call(func(ctx context.Context, input *services.TransactionSetPreOrderInput) (*services.TransactionSetPreOrderOutput, error) {
		if input.BlockHeight != h.lastBlockHeight {
			log.GetLogger().Error("Invalid block height", log.Uint64("expected-block-height", h.lastBlockHeight.KeyForMap()), log.Uint64("actual-block-height", input.BlockHeight.KeyForMap()))
			panic("Invalid block height")
		}
		statuses := make([]protocol.TransactionStatus, len(input.SignedTransactions))
		for i, tx := range input.SignedTransactions {
			if failOn(tx) {
				statuses[i] = protocol.TRANSACTION_STATUS_REJECTED_SMART_CONTRACT_PRE_ORDER
			} else {
				statuses[i] = protocol.TRANSACTION_STATUS_PRE_ORDER_VALID
			}
		}
		return &services.TransactionSetPreOrderOutput{
			PreOrderResults: statuses,
		}, nil
	})
}

func (h *harness) passAllPreOrderChecks() {
	h.failPreOrderCheckFor(func(tx *protocol.SignedTransaction) bool {
		return false
	})
}
func (h *harness) goToBlock(ctx context.Context, height primitives.BlockHeight, timestamp primitives.TimestampNano) {
	h.ignoringTransactionResults()
	currentBlock := primitives.BlockHeight(0)
	for currentBlock <= height {
		out, _ := h.txpool.CommitTransactionReceipts(ctx, &services.CommitTransactionReceiptsInput{
			LastCommittedBlockHeight: currentBlock,
			ResultsBlockHeader:       (&protocol.ResultsBlockHeaderBuilder{BlockHeight: currentBlock, Timestamp: timestamp}).Build(),
		})
		currentBlock = out.NextDesiredBlockHeight
	}
	h.lastBlockHeight = height
}

func (h *harness) validateTransactionsForOrdering(ctx context.Context, blockHeight primitives.BlockHeight, txs ...*protocol.SignedTransaction) error {
	_, err := h.txpool.ValidateTransactionsForOrdering(ctx, &services.ValidateTransactionsForOrderingInput{
		BlockHeight:        blockHeight,
		SignedTransactions: txs,
	})

	return err
}

func newHarness(ctx context.Context) *harness {
	return newHarnessWithSizeLimit(ctx, 20*1024*1024)
}

func newHarnessWithSizeLimit(ctx context.Context, sizeLimit uint32) *harness {
	gossip := &gossiptopics.MockTransactionRelay{}
	gossip.When("RegisterTransactionRelayHandler", mock.Any).Return()

	virtualMachine := &services.MockVirtualMachine{}

	cfg := config.ForTransactionPoolTests(sizeLimit, thisNodeKeyPair)
	metricFactory := metric.NewRegistry()

	service := transactionpool.NewTransactionPool(ctx, gossip, virtualMachine, nil, cfg, log.GetLogger(), metricFactory)

	transactionResultHandler := &handlers.MockTransactionResultsHandler{}
	service.RegisterTransactionResultsHandler(transactionResultHandler)

	h := &harness{
		txpool:             service,
		gossip:             gossip,
		vm:                 virtualMachine,
		trh:                transactionResultHandler,
		lastBlockTimestamp: primitives.TimestampNano(time.Now().UnixNano()),
		config:             cfg,
	}

	h.passAllPreOrderChecks()

	h.goToBlock(ctx, 1, now)

	return h
}

func asReceipts(transactions transactionpool.Transactions) []*protocol.TransactionReceipt {
	var receipts []*protocol.TransactionReceipt
	for _, tx := range transactions {
		receipts = append(receipts, (&protocol.TransactionReceiptBuilder{
			Txhash: digest.CalcTxHash(tx.Transaction()),
		}).Build())
	}
	return receipts
}
