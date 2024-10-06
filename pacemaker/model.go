package pacemaker

import (
	blockchain "banyan/blockchain_view"
	"banyan/crypto"
	"banyan/identity"
	"banyan/types"
)

type TMO struct {
	View   types.View
	NodeID identity.NodeID
	HighQC *blockchain.QC
}

type TC struct {
	types.View
	crypto.AggSig
	crypto.Signature
}

func NewTC(view types.View, requesters map[identity.NodeID]*TMO) *TC {
	// TODO: add crypto
	return &TC{View: view}
}
