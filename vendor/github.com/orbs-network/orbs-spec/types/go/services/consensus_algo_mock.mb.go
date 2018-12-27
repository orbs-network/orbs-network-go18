// AUTO GENERATED FILE (by membufc proto compiler v0.0.21)
package services

import (
	"github.com/orbs-network/go-mock"
	"github.com/orbs-network/orbs-spec/types/go/services/gossiptopics"
	"github.com/orbs-network/orbs-spec/types/go/services/handlers"
)

/////////////////////////////////////////////////////////////////////////////
// service ConsensusAlgo

type MockConsensusAlgo struct {
	mock.Mock
	handlers.MockConsensusBlocksHandler
}

/////////////////////////////////////////////////////////////////////////////
// service ConsensusAlgoLeanHelix

type MockConsensusAlgoLeanHelix struct {
	mock.Mock
	MockConsensusAlgo
	gossiptopics.MockLeanHelixHandler
}

/////////////////////////////////////////////////////////////////////////////
// service ConsensusAlgoBenchmark

type MockConsensusAlgoBenchmark struct {
	mock.Mock
	MockConsensusAlgo
	gossiptopics.MockBenchmarkConsensusHandler
}
