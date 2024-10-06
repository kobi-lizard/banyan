package replica

import (
	"banyan/blockchain"
)

type Safety interface {
	ProcessBlock(block *blockchain.Block) error
	ProcessNotarizationShare(vote *blockchain.NotarizationShare)
	ProcessFinalizationShare(vote *blockchain.FinalizationShare)
	MakeProposal(height int, rank int, payloadSize int) *blockchain.Block
}
