package blockchain

import (
	"io"
	"math/rand"
	"time"

	"banyan/crypto"
	"banyan/identity"
	"banyan/types"
)

type Block struct {
	types.View
	QC        *QC
	Proposer  identity.NodeID
	Timestamp time.Time
	Payload   []byte
	PrevID    crypto.Identifier
	Sig       crypto.Signature
	ID        crypto.Identifier
	Ts        time.Duration
}

type rawBlock struct {
	types.View
	QC          *QC
	Proposer    identity.NodeID
	PayloadHash crypto.Identifier
	PrevID      crypto.Identifier
	Sig         crypto.Signature
	ID          crypto.Identifier
}

// MakeBlock creates an unsigned block
func MakeBlock(view types.View, qc *QC, prevID crypto.Identifier, proposer identity.NodeID, blockByteSize int, r *rand.Rand) *Block {
	b := new(Block)
	b.View = view
	b.Proposer = proposer
	b.QC = qc
	b.Payload = generateRandomPayload(blockByteSize, r)
	b.PrevID = prevID
	b.makeID(proposer)
	return b
}

func (b *Block) makeID(nodeID identity.NodeID) {
	raw := &rawBlock{
		View:     b.View,
		QC:       b.QC,
		Proposer: b.Proposer,
		PrevID:   b.PrevID,
	}
	raw.PayloadHash = crypto.MakeID(b.Payload)
	b.ID = crypto.MakeID(raw)
	// TODO: uncomment the following
	b.Sig, _ = crypto.PrivSign(crypto.IDToByte(b.ID), nodeID, nil)
}

func generateRandomPayload(size int, r *rand.Rand) []byte {
	payload := make([]byte, size)
	if _, err := io.ReadFull(r, payload); err != nil {
		panic(err)
	}
	return payload
}
