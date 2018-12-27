package adapter

import (
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
)

type ContractState map[string]*protocol.StateRecord
type ChainState map[primitives.ContractName]ContractState

type StatePersistence interface {
	Write(height primitives.BlockHeight, ts primitives.TimestampNano, root primitives.Sha256, diff ChainState) error
	Read(contract primitives.ContractName, key string) (*protocol.StateRecord, bool, error)
	ReadMetadata() (primitives.BlockHeight, primitives.TimestampNano, primitives.Sha256, error)
}
