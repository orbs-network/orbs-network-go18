package adapter

import (
	"github.com/orbs-network/orbs-network-go/instrumentation/log"
	"github.com/orbs-network/orbs-network-go/instrumentation/metric"
	"github.com/orbs-network/orbs-network-go/services/statestorage/adapter"
)

type DumpingStatePersistence interface {
	adapter.StatePersistence
	Dump() string
}

type TestStatePersistence struct {
	*adapter.InMemoryStatePersistence
}

func NewDumpingStatePersistence(metric metric.Registry, log log.BasicLogger) *TestStatePersistence {
	result := &TestStatePersistence{
		InMemoryStatePersistence: adapter.NewInMemoryStatePersistence(metric),
	}
	return result
}
