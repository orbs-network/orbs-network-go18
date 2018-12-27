package blockstorage

import (
	"context"
	"github.com/orbs-network/orbs-network-go/instrumentation/log"
	"github.com/orbs-network/orbs-network-go/instrumentation/trace"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/orbs-spec/types/go/protocol/gossipmessages"
	"github.com/orbs-network/orbs-spec/types/go/services"
	"github.com/orbs-network/orbs-spec/types/go/services/gossiptopics"
	"github.com/pkg/errors"
)

func (s *service) HandleBlockAvailabilityRequest(ctx context.Context, input *gossiptopics.BlockAvailabilityRequestInput) (*gossiptopics.EmptyOutput, error) {
	err := s.sourceHandleBlockAvailabilityRequest(ctx, input.Message)
	return nil, err
}

func (s *service) HandleBlockSyncRequest(ctx context.Context, input *gossiptopics.BlockSyncRequestInput) (*gossiptopics.EmptyOutput, error) {
	err := s.sourceHandleBlockSyncRequest(ctx, input.Message)
	return nil, err
}

func (s *service) sourceHandleBlockAvailabilityRequest(ctx context.Context, message *gossipmessages.BlockAvailabilityRequestMessage) error {
	logger := s.logger.WithTags(trace.LogFieldFrom(ctx))

	logger.Info("received block availability request",
		log.Stringable("petitioner", message.Sender.SenderNodeAddress()),
		log.Stringable("requested-first-block", message.SignedBatchRange.FirstBlockHeight()),
		log.Stringable("requested-last-block", message.SignedBatchRange.LastBlockHeight()),
		log.Stringable("requested-last-committed-block", message.SignedBatchRange.LastCommittedBlockHeight()))

	out, err := s.GetLastCommittedBlockHeight(ctx, &services.GetLastCommittedBlockHeightInput{})
	if err != nil {
		return err
	}
	lastCommittedBlockHeight := out.LastCommittedBlockHeight

	if lastCommittedBlockHeight <= message.SignedBatchRange.LastCommittedBlockHeight() {
		return nil
	}

	firstAvailableBlockHeight := primitives.BlockHeight(1)
	blockType := message.SignedBatchRange.BlockType()

	response := &gossiptopics.BlockAvailabilityResponseInput{
		RecipientNodeAddress: message.Sender.SenderNodeAddress(),
		Message: &gossipmessages.BlockAvailabilityResponseMessage{
			Sender: (&gossipmessages.SenderSignatureBuilder{
				SenderNodeAddress: s.config.NodeAddress(),
			}).Build(),
			SignedBatchRange: (&gossipmessages.BlockSyncRangeBuilder{
				BlockType:                blockType,
				LastBlockHeight:          lastCommittedBlockHeight,
				FirstBlockHeight:         firstAvailableBlockHeight,
				LastCommittedBlockHeight: lastCommittedBlockHeight,
			}).Build(),
		},
	}

	logger.Info("sending the response for availability request",
		log.Stringable("petitioner", response.RecipientNodeAddress),
		log.Stringable("first-available-block-height", response.Message.SignedBatchRange.FirstBlockHeight()),
		log.Stringable("last-available-block-height", response.Message.SignedBatchRange.LastBlockHeight()),
		log.Stringable("last-committed-available-block-height", response.Message.SignedBatchRange.LastCommittedBlockHeight()),
		log.Stringable("source", response.Message.Sender.SenderNodeAddress()),
	)

	_, err = s.gossip.SendBlockAvailabilityResponse(ctx, response)
	return err
}

func (s *service) sourceHandleBlockSyncRequest(ctx context.Context, message *gossipmessages.BlockSyncRequestMessage) error {
	logger := s.logger.WithTags(trace.LogFieldFrom(ctx))

	senderNodeAddress := message.Sender.SenderNodeAddress()
	blockType := message.SignedChunkRange.BlockType()
	firstRequestedBlockHeight := message.SignedChunkRange.FirstBlockHeight()
	lastRequestedBlockHeight := message.SignedChunkRange.LastBlockHeight()

	out, err := s.GetLastCommittedBlockHeight(ctx, &services.GetLastCommittedBlockHeightInput{})
	if err != nil {
		return err
	}
	lastCommittedBlockHeight := out.LastCommittedBlockHeight

	logger.Info("received block sync request",
		log.Stringable("petitioner", message.Sender.SenderNodeAddress()),
		log.Stringable("first-requested-block-height", firstRequestedBlockHeight),
		log.Stringable("last-requested-block-height", lastRequestedBlockHeight),
		log.Stringable("last-committed-block-height", lastCommittedBlockHeight))

	if lastCommittedBlockHeight <= firstRequestedBlockHeight {
		return errors.New("firstBlockHeight is greater or equal to lastCommittedBlockHeight")
	}

	if firstRequestedBlockHeight-lastCommittedBlockHeight > primitives.BlockHeight(s.config.BlockSyncBatchSize()-1) {
		lastRequestedBlockHeight = firstRequestedBlockHeight + primitives.BlockHeight(s.config.BlockSyncBatchSize()-1)
	}

	blocks, firstAvailableBlockHeight, lastAvailableBlockHeight, err := s.GetBlockSlice(firstRequestedBlockHeight, lastRequestedBlockHeight)
	if err != nil {
		return errors.Wrap(err, "block sync failed reading from block persistence")
	}

	logger.Info("sending blocks to another node via block sync",
		log.Stringable("petitioner", senderNodeAddress),
		log.Stringable("first-available-block-height", firstAvailableBlockHeight),
		log.Stringable("last-available-block-height", lastAvailableBlockHeight))

	response := &gossiptopics.BlockSyncResponseInput{
		RecipientNodeAddress: senderNodeAddress,
		Message: &gossipmessages.BlockSyncResponseMessage{
			Sender: (&gossipmessages.SenderSignatureBuilder{
				SenderNodeAddress: s.config.NodeAddress(),
			}).Build(),
			SignedChunkRange: (&gossipmessages.BlockSyncRangeBuilder{
				BlockType:                blockType,
				FirstBlockHeight:         firstAvailableBlockHeight,
				LastBlockHeight:          lastAvailableBlockHeight,
				LastCommittedBlockHeight: lastCommittedBlockHeight,
			}).Build(),
			BlockPairs: blocks,
		},
	}
	_, err = s.gossip.SendBlockSyncResponse(ctx, response)
	return err
}
