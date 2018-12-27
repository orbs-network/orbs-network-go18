package inmemory

import (
	"context"
	"fmt"
	"github.com/orbs-network/orbs-network-go/bootstrap"
	"github.com/orbs-network/orbs-network-go/config"
	"github.com/orbs-network/orbs-network-go/crypto/digest"
	"github.com/orbs-network/orbs-network-go/crypto/keys"
	"github.com/orbs-network/orbs-network-go/instrumentation/log"
	"github.com/orbs-network/orbs-network-go/instrumentation/metric"
	ethereumAdapter "github.com/orbs-network/orbs-network-go/services/crosschainconnector/ethereum/adapter"
	"github.com/orbs-network/orbs-network-go/services/gossip/adapter"
	nativeProcessorAdapter "github.com/orbs-network/orbs-network-go/services/processor/native/adapter"
	"github.com/orbs-network/orbs-network-go/synchronization"
	"github.com/orbs-network/orbs-network-go/test"
	"github.com/orbs-network/orbs-network-go/test/harness/contracts"
	blockStorageAdapter "github.com/orbs-network/orbs-network-go/test/harness/services/blockstorage/adapter"
	harnessStateStorageAdapter "github.com/orbs-network/orbs-network-go/test/harness/services/statestorage/adapter"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/orbs-network/orbs-spec/types/go/protocol/client"
	"github.com/orbs-network/orbs-spec/types/go/services"
	"math"
)

type NetworkDriver interface {
	contracts.ContractAPI
	PublicApi(nodeIndex int) services.PublicApi
	Size() int
	WaitUntilReadyForTransactions(ctx context.Context)
}

type Network struct {
	Nodes     []*Node
	Logger    log.BasicLogger
	Transport adapter.Transport
}

type Node struct {
	index                             int
	name                              string
	config                            config.NodeConfig
	blockPersistence                  blockStorageAdapter.InMemoryBlockPersistence
	statePersistence                  harnessStateStorageAdapter.DumpingStatePersistence
	stateBlockHeightTracker           *synchronization.BlockTracker
	transactionPoolBlockHeightTracker *synchronization.BlockTracker
	nativeCompiler                    nativeProcessorAdapter.Compiler
	ethereumConnection                ethereumAdapter.EthereumConnection
	nodeLogic                         bootstrap.NodeLogic
	metricRegistry                    metric.Registry
}

func NewNetwork(logger log.BasicLogger, transport adapter.Transport) Network {
	return Network{Logger: logger, Transport: transport}
}

func (n *Network) AddNode(
	nodeKeyPair *keys.EcdsaSecp256K1KeyPair,
	cfg config.NodeConfig,
	blockPersistence blockStorageAdapter.InMemoryBlockPersistence,
	compiler nativeProcessorAdapter.Compiler,
	ethereumConnection ethereumAdapter.EthereumConnection,
	metricRegistry metric.Registry,
	logger log.BasicLogger) {

	node := &Node{}
	node.index = len(n.Nodes)
	node.name = fmt.Sprintf("%s", nodeKeyPair.PublicKey()[:3])
	node.config = cfg
	node.statePersistence = harnessStateStorageAdapter.NewDumpingStatePersistence(metricRegistry, logger)
	node.stateBlockHeightTracker = synchronization.NewBlockTracker(logger, 0, math.MaxUint16)
	node.transactionPoolBlockHeightTracker = synchronization.NewBlockTracker(logger, 0, math.MaxUint16)
	node.blockPersistence = blockPersistence
	node.nativeCompiler = compiler
	node.ethereumConnection = ethereumConnection
	node.metricRegistry = metricRegistry

	n.Nodes = append(n.Nodes, node)
}

func (n *Network) CreateAndStartNodes(ctx context.Context, numOfNodesToStart int) {
	for i, node := range n.Nodes {
		if i >= numOfNodesToStart {
			return
		}

		node.nodeLogic = bootstrap.NewNodeLogic(
			ctx,
			n.Transport,
			node.blockPersistence,
			node.statePersistence,
			node.stateBlockHeightTracker,
			node.transactionPoolBlockHeightTracker,
			node.nativeCompiler,
			n.Logger.WithTags(log.Node(node.name)),
			node.metricRegistry,
			node.config,
			node.ethereumConnection,
		)
	}
}

func (n *Node) GetPublicApi() services.PublicApi {
	return n.nodeLogic.PublicApi()
}

func (n *Node) GetCompiler() nativeProcessorAdapter.Compiler {
	return n.nativeCompiler
}

func (n *Node) WaitForTransactionInState(ctx context.Context, txHash primitives.Sha256) primitives.BlockHeight {
	blockHeight := n.blockPersistence.WaitForTransaction(ctx, txHash)
	err := n.stateBlockHeightTracker.WaitForBlock(ctx, blockHeight)
	if err != nil {
		test.DebugPrintGoroutineStacks() // since test timed out, help find deadlocked goroutines
		panic(fmt.Sprintf("statePersistence.WaitUntilCommittedBlockOfHeight failed: %s", err.Error()))
	}
	return blockHeight
}
func (n *Node) Started() bool {
	return n.nodeLogic != nil
}

func (n *Node) Destroy() {
	n.nodeLogic = nil
}

func (n *Network) PublicApi(nodeIndex int) services.PublicApi {
	return n.Nodes[nodeIndex].nodeLogic.PublicApi()
}

func (n *Network) GetBlockPersistence(nodeIndex int) blockStorageAdapter.InMemoryBlockPersistence {
	return n.Nodes[nodeIndex].blockPersistence
}

func (n *Network) GetStatePersistence(i int) harnessStateStorageAdapter.DumpingStatePersistence {
	return n.Nodes[i].statePersistence
}

func (n *Network) Size() int {
	return len(n.Nodes)
}

func (n *Network) SendTransaction(ctx context.Context, tx *protocol.SignedTransactionBuilder, nodeIndex int) (*client.SendTransactionResponse, primitives.Sha256) {
	n.assertStarted(nodeIndex)
	ch := make(chan *client.SendTransactionResponse)

	transactionRequestBuilder := &client.SendTransactionRequestBuilder{SignedTransaction: tx}
	txHash := digest.CalcTxHash(transactionRequestBuilder.SignedTransaction.Transaction.Build())

	go func() {
		defer close(ch)
		publicApi := n.Nodes[nodeIndex].GetPublicApi()
		output, err := publicApi.SendTransaction(ctx, &services.SendTransactionInput{
			ClientRequest: transactionRequestBuilder.Build(),
		})
		if output == nil {
			panic(fmt.Sprintf("error sending transaction: %v", err)) // TODO(https://github.com/orbs-network/orbs-network-go/issues/531): improve
		}

		select {
		case ch <- output.ClientResponse:
		case <-ctx.Done():
		}
	}()

	return <-ch, txHash
}

func (n *Network) SendTransactionInBackground(ctx context.Context, tx *protocol.SignedTransactionBuilder, nodeIndex int) {
	n.assertStarted(nodeIndex)

	go func() {
		publicApi := n.Nodes[nodeIndex].GetPublicApi()
		output, err := publicApi.SendTransaction(ctx, &services.SendTransactionInput{
			ClientRequest:     (&client.SendTransactionRequestBuilder{SignedTransaction: tx}).Build(),
			ReturnImmediately: 1,
		})
		if output == nil {
			panic(fmt.Sprintf("error sending transaction: %v", err)) // TODO(https://github.com/orbs-network/orbs-network-go/issues/531): improve
		}
	}()
}

func (n *Network) CallMethod(ctx context.Context, tx *protocol.TransactionBuilder, nodeIndex int) *client.CallMethodResponse {
	n.assertStarted(nodeIndex)

	ch := make(chan *client.CallMethodResponse)
	go func() {
		defer close(ch)
		publicApi := n.Nodes[nodeIndex].GetPublicApi()
		output, err := publicApi.CallMethod(ctx, &services.CallMethodInput{
			ClientRequest: (&client.CallMethodRequestBuilder{Transaction: tx}).Build(),
		})
		if output == nil {
			panic(fmt.Sprintf("error calling method: %v", err)) // TODO(https://github.com/orbs-network/orbs-network-go/issues/531): improve
		}
		select {
		case ch <- output.ClientResponse:
		case <-ctx.Done():
		}
	}()
	return <-ch
}

func (n *Network) assertStarted(nodeIndex int) {
	if !n.Nodes[nodeIndex].Started() {
		panic(fmt.Errorf("accessing a stopped node %d", nodeIndex))
	}
}

func (n *Network) WaitForTransactionInState(ctx context.Context, txHash primitives.Sha256) {
	for _, node := range n.Nodes {
		if node.Started() {
			h := node.WaitForTransactionInState(ctx, txHash)
			n.Logger.Info("WaitForTransactionInState found tx in state", log.BlockHeight(h), log.Node(node.name), log.Transaction(txHash))
		}
	}
}

func (n *Network) WaitUntilReadyForTransactions(ctx context.Context) {
	for _, node := range n.Nodes {
		if node.Started() {
			node.transactionPoolBlockHeightTracker.WaitForBlock(ctx, 1)
		}
	}
}
