package adapter

import (
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/orbs-network/orbs-network-go/instrumentation/log"
	"math/big"
	"sync"
)

type EthereumRpcConnection struct {
	connectorCommon

	config ethereumAdapterConfig

	mu struct {
		sync.Mutex
		client *ethclient.Client
	}
}

func NewEthereumRpcConnection(config ethereumAdapterConfig, logger log.BasicLogger) *EthereumRpcConnection {
	rpc := &EthereumRpcConnection{
		config: config,
	}
	rpc.logger = logger.WithTags(log.String("adapter", "ethereum"))
	rpc.getContractCaller = func() (caller EthereumCaller, e error) {
		return rpc.dialIfNeededAndReturnClient()
	}
	return rpc
}

func (rpc *EthereumRpcConnection) dial() error {
	rpc.mu.Lock()
	defer rpc.mu.Unlock()
	if rpc.mu.client != nil {
		return nil
	}
	if client, err := ethclient.Dial(rpc.config.EthereumEndpoint()); err != nil {
		return err
	} else {
		rpc.mu.client = client
	}
	return nil
}

func (rpc *EthereumRpcConnection) dialIfNeededAndReturnClient() (*ethclient.Client, error) {
	if rpc.mu.client == nil {
		if err := rpc.dial(); err != nil {
			return nil, err
		}
	}
	return rpc.mu.client, nil
}

func (rpc *EthereumRpcConnection) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	client, err := rpc.dialIfNeededAndReturnClient()
	if err != nil {
		return nil, err
	}

	return client.HeaderByNumber(ctx, number)
}
