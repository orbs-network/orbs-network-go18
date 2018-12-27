package gossip

import (
	"context"
	"github.com/orbs-network/orbs-network-go/instrumentation/log"
	"github.com/orbs-network/orbs-network-go/services/gossip/adapter"
	"github.com/orbs-network/orbs-network-go/services/gossip/codec"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/orbs-spec/types/go/protocol/gossipmessages"
	"github.com/orbs-network/orbs-spec/types/go/services/gossiptopics"
)

func (s *service) RegisterBlockSyncHandler(handler gossiptopics.BlockSyncHandler) {
	s.blockSyncHandlers = append(s.blockSyncHandlers, handler)
}

func (s *service) receivedBlockSyncMessage(ctx context.Context, header *gossipmessages.Header, payloads [][]byte) {
	switch header.BlockSync() {
	case gossipmessages.BLOCK_SYNC_AVAILABILITY_REQUEST:
		s.receivedBlockSyncAvailabilityRequest(ctx, header, payloads)
	case gossipmessages.BLOCK_SYNC_AVAILABILITY_RESPONSE:
		s.receivedBlockSyncAvailabilityResponse(ctx, header, payloads)
	case gossipmessages.BLOCK_SYNC_REQUEST:
		s.receivedBlockSyncRequest(ctx, header, payloads)
	case gossipmessages.BLOCK_SYNC_RESPONSE:
		s.receivedBlockSyncResponse(ctx, header, payloads)
	}
}

func (s *service) BroadcastBlockAvailabilityRequest(ctx context.Context, input *gossiptopics.BlockAvailabilityRequestInput) (*gossiptopics.EmptyOutput, error) {
	header := (&gossipmessages.HeaderBuilder{
		Topic:         gossipmessages.HEADER_TOPIC_BLOCK_SYNC,
		BlockSync:     gossipmessages.BLOCK_SYNC_AVAILABILITY_REQUEST,
		RecipientMode: gossipmessages.RECIPIENT_LIST_MODE_BROADCAST,
	}).Build()
	payloads, err := codec.EncodeBlockAvailabilityRequest(header, input.Message)
	if err != nil {
		return nil, err
	}
	return nil, s.transport.Send(ctx, &adapter.TransportData{
		SenderNodeAddress: s.config.NodeAddress(),
		RecipientMode:     gossipmessages.RECIPIENT_LIST_MODE_BROADCAST,
		Payloads:          payloads,
	})
}

func (s *service) receivedBlockSyncAvailabilityRequest(ctx context.Context, header *gossipmessages.Header, payloads [][]byte) {
	message, err := codec.DecodeBlockAvailabilityRequest(payloads)
	if err != nil {
		return
	}
	for _, l := range s.blockSyncHandlers {
		_, err := l.HandleBlockAvailabilityRequest(ctx, &gossiptopics.BlockAvailabilityRequestInput{Message: message})
		if err != nil {
			s.logger.Info("HandleBlockAvailabilityRequest failed", log.Error(err))
		}
	}
}

func (s *service) SendBlockAvailabilityResponse(ctx context.Context, input *gossiptopics.BlockAvailabilityResponseInput) (*gossiptopics.EmptyOutput, error) {
	header := (&gossipmessages.HeaderBuilder{
		Topic:                  gossipmessages.HEADER_TOPIC_BLOCK_SYNC,
		BlockSync:              gossipmessages.BLOCK_SYNC_AVAILABILITY_RESPONSE,
		RecipientMode:          gossipmessages.RECIPIENT_LIST_MODE_LIST,
		RecipientNodeAddresses: []primitives.NodeAddress{input.RecipientNodeAddress},
	}).Build()
	payloads, err := codec.EncodeBlockAvailabilityResponse(header, input.Message)
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

func (s *service) receivedBlockSyncAvailabilityResponse(ctx context.Context, header *gossipmessages.Header, payloads [][]byte) {
	message, err := codec.DecodeBlockAvailabilityResponse(payloads)
	if err != nil {
		return
	}
	for _, l := range s.blockSyncHandlers {
		_, err := l.HandleBlockAvailabilityResponse(ctx, &gossiptopics.BlockAvailabilityResponseInput{Message: message})
		if err != nil {
			s.logger.Info("HandleBlockAvailabilityResponse failed", log.Error(err))
		}
	}
}

func (s *service) SendBlockSyncRequest(ctx context.Context, input *gossiptopics.BlockSyncRequestInput) (*gossiptopics.EmptyOutput, error) {
	header := (&gossipmessages.HeaderBuilder{
		Topic:                  gossipmessages.HEADER_TOPIC_BLOCK_SYNC,
		BlockSync:              gossipmessages.BLOCK_SYNC_REQUEST,
		RecipientMode:          gossipmessages.RECIPIENT_LIST_MODE_LIST,
		RecipientNodeAddresses: []primitives.NodeAddress{input.RecipientNodeAddress},
	}).Build()
	payloads, err := codec.EncodeBlockSyncRequest(header, input.Message)
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

func (s *service) receivedBlockSyncRequest(ctx context.Context, header *gossipmessages.Header, payloads [][]byte) {
	message, err := codec.DecodeBlockSyncRequest(payloads)
	if err != nil {
		return
	}
	for _, l := range s.blockSyncHandlers {
		_, err := l.HandleBlockSyncRequest(ctx, &gossiptopics.BlockSyncRequestInput{Message: message})
		if err != nil {
			s.logger.Info("HandleBlockSyncRequest failed", log.Error(err))
		}
	}
}

func (s *service) SendBlockSyncResponse(ctx context.Context, input *gossiptopics.BlockSyncResponseInput) (*gossiptopics.EmptyOutput, error) {
	header := (&gossipmessages.HeaderBuilder{
		Topic:                  gossipmessages.HEADER_TOPIC_BLOCK_SYNC,
		BlockSync:              gossipmessages.BLOCK_SYNC_RESPONSE,
		RecipientMode:          gossipmessages.RECIPIENT_LIST_MODE_LIST,
		RecipientNodeAddresses: []primitives.NodeAddress{input.RecipientNodeAddress},
	}).Build()
	payloads, err := codec.EncodeBlockSyncResponse(header, input.Message)
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

func (s *service) receivedBlockSyncResponse(ctx context.Context, header *gossipmessages.Header, payloads [][]byte) {
	message, err := codec.DecodeBlockSyncResponse(payloads)
	if err != nil {
		return
	}
	for _, l := range s.blockSyncHandlers {
		_, err := l.HandleBlockSyncResponse(ctx, &gossiptopics.BlockSyncResponseInput{Message: message})
		if err != nil {
			s.logger.Info("HandleBlockSyncResponse failed", log.Error(err))
		}
	}
}
