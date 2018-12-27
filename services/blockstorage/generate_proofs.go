package blockstorage

import (
	"context"
	"github.com/orbs-network/orbs-network-go/crypto/digest"
	"github.com/orbs-network/orbs-network-go/crypto/merkle"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/orbs-network/orbs-spec/types/go/services"
	"github.com/pkg/errors"
)

func (s *service) GenerateReceiptProof(ctx context.Context, input *services.GenerateReceiptProofInput) (*services.GenerateReceiptProofOutput, error) {
	block, err := s.persistence.GetResultsBlock(input.BlockHeight)
	if err != nil {
		return nil, err
	}

	for i, txr := range block.TransactionReceipts {
		if txr.Txhash().Equal(input.Txhash) {

			proof, err := generateProof(block.TransactionReceipts, i)
			if err != nil {
				return nil, err
			}

			result := &services.GenerateReceiptProofOutput{
				Proof: (&protocol.ReceiptProofBuilder{
					Header:       protocol.ResultsBlockHeaderBuilderFromRaw(block.Header.Raw()),
					BlockProof:   protocol.ResultsBlockProofBuilderFromRaw(block.BlockProof.Raw()),
					ReceiptProof: proof,
				}).Build(),
			}
			return result, nil
		}
	}

	return nil, errors.Errorf("could not find transaction inside block %x", input.Txhash)

}

func generateProof(receipts []*protocol.TransactionReceipt, index int) (primitives.MerkleTreeProof, error) {
	hashes := digest.CalcReceiptHashes(receipts)
	proof, err := merkle.NewOrderedTree(hashes).GetProof(index)
	if err != nil {
		return nil, err
	}
	return merkle.FlattenOrderedTreeProof(proof), nil
}
