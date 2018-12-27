package gossip

import (
	"context"
	"github.com/orbs-network/orbs-network-go/instrumentation/log"
	"github.com/orbs-network/orbs-network-go/instrumentation/trace"
	"github.com/orbs-network/orbs-network-go/services/gossip/adapter"
	"github.com/orbs-network/orbs-network-go/services/gossip/codec"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/orbs-spec/types/go/protocol/consensus"
	"github.com/orbs-network/orbs-spec/types/go/protocol/gossipmessages"
	"github.com/orbs-network/orbs-spec/types/go/services/gossiptopics"
)

func (s *service) RegisterBenchmarkConsensusHandler(handler gossiptopics.BenchmarkConsensusHandler) {
	s.benchmarkConsensusHandlers = append(s.benchmarkConsensusHandlers, handler)
}

func (s *service) receivedBenchmarkConsensusMessage(ctx context.Context, header *gossipmessages.Header, payloads [][]byte) {
	switch header.BenchmarkConsensus() {
	case consensus.BENCHMARK_CONSENSUS_COMMIT:
		s.receivedBenchmarkConsensusCommit(ctx, header, payloads)
	case consensus.BENCHMARK_CONSENSUS_COMMITTED:
		s.receivedBenchmarkConsensusCommitted(ctx, header, payloads)
	}
}

func (s *service) BroadcastBenchmarkConsensusCommit(ctx context.Context, input *gossiptopics.BenchmarkConsensusCommitInput) (*gossiptopics.EmptyOutput, error) {
	header := (&gossipmessages.HeaderBuilder{
		Topic:              gossipmessages.HEADER_TOPIC_BENCHMARK_CONSENSUS,
		BenchmarkConsensus: consensus.BENCHMARK_CONSENSUS_COMMIT,
		RecipientMode:      gossipmessages.RECIPIENT_LIST_MODE_BROADCAST,
	}).Build()

	payloads, err := codec.EncodeBenchmarkConsensusCommitMessage(header, input.Message)
	if err != nil {
		return nil, err
	}

	return nil, s.transport.Send(ctx, &adapter.TransportData{
		SenderNodeAddress: s.config.NodeAddress(),
		RecipientMode:     gossipmessages.RECIPIENT_LIST_MODE_BROADCAST,
		Payloads:          payloads,
	})
}

func (s *service) receivedBenchmarkConsensusCommit(ctx context.Context, header *gossipmessages.Header, payloads [][]byte) {
	logger := s.logger.WithTags(trace.LogFieldFrom(ctx))
	message, err := codec.DecodeBenchmarkConsensusCommitMessage(payloads)
	if err != nil {
		logger.Info("HandleBenchmarkConsensusCommit failed to decode block pair", log.Error(err))
		return
	}

	for _, l := range s.benchmarkConsensusHandlers {
		_, err := l.HandleBenchmarkConsensusCommit(ctx, &gossiptopics.BenchmarkConsensusCommitInput{Message: message})
		if err != nil {
			logger.Info("HandleBenchmarkConsensusCommit failed", log.Error(err))
		}
	}
}

func (s *service) SendBenchmarkConsensusCommitted(ctx context.Context, input *gossiptopics.BenchmarkConsensusCommittedInput) (*gossiptopics.EmptyOutput, error) {
	header := (&gossipmessages.HeaderBuilder{
		Topic:                  gossipmessages.HEADER_TOPIC_BENCHMARK_CONSENSUS,
		BenchmarkConsensus:     consensus.BENCHMARK_CONSENSUS_COMMITTED,
		RecipientMode:          gossipmessages.RECIPIENT_LIST_MODE_LIST,
		RecipientNodeAddresses: []primitives.NodeAddress{input.RecipientNodeAddress},
	}).Build()
	payloads, err := codec.EncodeBenchmarkConsensusCommittedMessage(header, input.Message)
	if err != nil {
		return nil, err
	}

	return nil, s.transport.Send(ctx, &adapter.TransportData{
		SenderNodeAddress:      s.config.NodeAddress(),
		RecipientMode:          gossipmessages.RECIPIENT_LIST_MODE_LIST,
		RecipientNodeAddresses: []primitives.NodeAddress{input.RecipientNodeAddress},
		Payloads:               payloads,
	})
}

func (s *service) receivedBenchmarkConsensusCommitted(ctx context.Context, header *gossipmessages.Header, payloads [][]byte) {
	message, err := codec.DecodeBenchmarkConsensusCommittedMessage(payloads)
	if err != nil {
		return
	}

	for _, l := range s.benchmarkConsensusHandlers {
		_, err := l.HandleBenchmarkConsensusCommitted(ctx, &gossiptopics.BenchmarkConsensusCommittedInput{Message: message})
		if err != nil {
			s.logger.Info("HandleBenchmarkConsensusCommitted failed", log.Error(err))
		}
	}
}
