package election

import (
	"banyan/identity"
	"banyan/types"
)

type Election interface {
	IsLeader(id identity.NodeID, height int, rank int) bool
	IsLeaderView(id identity.NodeID, view types.View) bool
	FindLeaderFor(height int, rank int) identity.NodeID
	FindLeaderForView(view types.View) identity.NodeID
}
