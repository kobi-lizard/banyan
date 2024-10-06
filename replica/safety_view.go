package replica

import (
	blockchain "banyan/blockchain_view"
	"banyan/pacemaker"
	"banyan/types"
)

type SafetyView interface {
	ProcessBlock(block *blockchain.Block) error
	ProcessVote(vote *blockchain.Vote)
	ProcessRemoteTmo(tmo *pacemaker.TMO)
	ProcessLocalTmo(view types.View)
	MakeProposal(view types.View, payloadSize int) *blockchain.Block
	GetChainStatus() string
}
