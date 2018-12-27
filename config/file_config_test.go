package config

import (
	"github.com/orbs-network/orbs-network-go/test/crypto/keys"
	"github.com/orbs-network/orbs-spec/types/go/protocol/consensus"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestFileConfigConstructor(t *testing.T) {
	cfg, err := newEmptyFileConfig(`{}`)

	require.NotNil(t, cfg)
	require.NoError(t, err)
}

func TestFileConfigSetUint32(t *testing.T) {
	cfg, err := newEmptyFileConfig(`{"block-sync-batch-size": 999}`)

	require.NotNil(t, cfg)
	require.NoError(t, err)
	require.EqualValues(t, 999, cfg.BlockSyncBatchSize())
}

func TestFileConfigSetDuration(t *testing.T) {
	cfg, err := newEmptyFileConfig(`{"block-sync-collect-response-timeout": "10m"}`)

	require.NotNil(t, cfg)
	require.NoError(t, err)
	require.EqualValues(t, 10*time.Minute, cfg.BlockSyncCollectResponseTimeout())
}

func TestSetNodeAddress(t *testing.T) {
	cfg, err := newEmptyFileConfig(`{"node-address": "a328846cd5b4979d68a8c58a9bdfeee657b34de7"}`)

	keyPair := keys.EcdsaSecp256K1KeyPairForTests(0)

	require.NotNil(t, cfg)
	require.NoError(t, err)
	require.EqualValues(t, keyPair.NodeAddress(), cfg.NodeAddress())
}

func TestSetNodePrivateKey(t *testing.T) {
	cfg, err := newEmptyFileConfig(`{"node-private-key": "901a1a0bfbe217593062a054e561e708707cb814a123474c25fd567a0fe088f8"}`)

	keyPair := keys.EcdsaSecp256K1KeyPairForTests(0)

	require.NotNil(t, cfg)
	require.NoError(t, err)
	require.EqualValues(t, keyPair.PrivateKey(), cfg.NodePrivateKey())
}

func TestSetConstantConsensusLeader(t *testing.T) {
	cfg, err := newEmptyFileConfig(`{"constant-consensus-leader": "d27e2e7398e2582f63d0800330010b3e58952ff6"}`)

	keyPair := keys.EcdsaSecp256K1KeyPairForTests(1)

	require.NotNil(t, cfg)
	require.NoError(t, err)
	require.EqualValues(t, keyPair.NodeAddress(), cfg.ConstantConsensusLeader())
}

func TestSetActiveConsensusAlgo(t *testing.T) {
	cfg, err := newEmptyFileConfig(`{"active-consensus-algo": 999}`)

	require.NotNil(t, cfg)
	require.NoError(t, err)
	require.EqualValues(t, 999, cfg.ActiveConsensusAlgo())
}

func TestSetFederationNodes(t *testing.T) {
	cfg, err := newEmptyFileConfig(`{
	"federation-nodes": [
    {"address":"a328846cd5b4979d68a8c58a9bdfeee657b34de7","ip":"192.168.199.2","port":4400},
    {"address":"d27e2e7398e2582f63d0800330010b3e58952ff6","ip":"192.168.199.3","port":4400},
    {"address":"6e2cb55e4cbe97bf5b1e731d51cc2c285d83cbf9","ip":"192.168.199.4","port":4400}
	]
}`)

	require.NotNil(t, cfg)
	require.NoError(t, err)
	require.EqualValues(t, 3, len(cfg.FederationNodes(0)))

	keyPair := keys.EcdsaSecp256K1KeyPairForTests(0)

	node1 := &hardCodedFederationNode{
		nodeAddress: keyPair.NodeAddress(),
	}

	require.EqualValues(t, node1, cfg.FederationNodes(0)[keyPair.NodeAddress().KeyForMap()])
}

func TestSetGossipPeers(t *testing.T) {
	cfg, err := newEmptyFileConfig(`{
	"federation-nodes": [
    {"address":"a328846cd5b4979d68a8c58a9bdfeee657b34de7","ip":"192.168.199.2","port":4400},
    {"address":"d27e2e7398e2582f63d0800330010b3e58952ff6","ip":"192.168.199.3","port":4400},
    {"address":"6e2cb55e4cbe97bf5b1e731d51cc2c285d83cbf9","ip":"192.168.199.4","port":4400}
	]
}`)

	require.NotNil(t, cfg)
	require.NoError(t, err)
	require.EqualValues(t, 3, len(cfg.GossipPeers(0)))

	keyPair := keys.EcdsaSecp256K1KeyPairForTests(0)

	node1 := &hardCodedGossipPeer{
		gossipEndpoint: "192.168.199.2",
		gossipPort:     4400,
	}

	require.EqualValues(t, node1, cfg.GossipPeers(0)[keyPair.NodeAddress().KeyForMap()])
}

func TestSetGossipPort(t *testing.T) {
	cfg, err := newEmptyFileConfig(`{"gossip-port": 4500}`)

	require.NotNil(t, cfg)
	require.NoError(t, err)
	require.EqualValues(t, 4500, cfg.GossipListenPort())
}

func TestMergeWithFileConfig(t *testing.T) {
	nodes := make(map[string]FederationNode)
	keyPair := keys.EcdsaSecp256K1KeyPairForTests(2)

	cfg := ForAcceptanceTestNetwork(nodes,
		keyPair.NodeAddress(),
		consensus.CONSENSUS_ALGO_TYPE_BENCHMARK_CONSENSUS, 30, 100)

	require.EqualValues(t, 0, len(cfg.FederationNodes(0)))

	cfg.MergeWithFileConfig(`
{
	"block-sync-batch-size": 999,
	"block-sync-collect-response-timeout": "10m",
	"node-address": "a328846cd5b4979d68a8c58a9bdfeee657b34de7",
	"node-private-key": "901a1a0bfbe217593062a054e561e708707cb814a123474c25fd567a0fe088f8",
	"constant-consensus-leader": "a328846cd5b4979d68a8c58a9bdfeee657b34de7",
	"active-consensus-algo": 999,
	"gossip-port": 4500,
	"federation-nodes": [
    {"address":"a328846cd5b4979d68a8c58a9bdfeee657b34de7","ip":"192.168.199.2","port":4400},
    {"address":"d27e2e7398e2582f63d0800330010b3e58952ff6","ip":"192.168.199.3","port":4400},
    {"address":"6e2cb55e4cbe97bf5b1e731d51cc2c285d83cbf9","ip":"192.168.199.4","port":4400}
	]
}
`)

	newKeyPair := keys.EcdsaSecp256K1KeyPairForTests(0)

	require.EqualValues(t, 3, len(cfg.FederationNodes(0)))
	require.EqualValues(t, newKeyPair.NodeAddress(), cfg.NodeAddress())
}

func TestConfig_EthereumEndpoint(t *testing.T) {
	cfg, err := newEmptyFileConfig(`{"ethereum-endpoint":"http://172.31.1.100:8545"}`)
	require.NoError(t, err)

	require.EqualValues(t, "http://172.31.1.100:8545", cfg.EthereumEndpoint())
}
