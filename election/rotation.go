package election

import (
	"banyan/identity"
	"banyan/types"
)

type Rotation struct {
	peerNo int
}

func NewRotation(peerNo int) *Rotation {
	return &Rotation{
		peerNo: peerNo,
	}
}

func (r *Rotation) IsLeader(id identity.NodeID, height int, rank int) bool {
	return (uint64(height-1)+uint64(rank))%uint64(r.peerNo) == uint64(id.Node()-1)
}

func (r *Rotation) IsLeaderView(id identity.NodeID, view types.View) bool {
	return uint64(view-1)%uint64(r.peerNo) == uint64(id.Node()-1)
}

func (r *Rotation) FindLeaderFor(height int, rank int) identity.NodeID {
	return identity.NewNodeID((height+rank-1)%r.peerNo + 1)
}

func (r *Rotation) FindLeaderForView(view types.View) identity.NodeID {
	return identity.NewNodeID(int(view-1)%r.peerNo + 1)
}
