package election

import (
	"banyan/identity"
)

type Static struct {
	master identity.NodeID
}

func NewStatic(master identity.NodeID) *Static {
	return &Static{
		master: master,
	}
}

func (st *Static) IsLeader(id identity.NodeID, height int, rank int) bool {
	return id == st.master
}

func (st *Static) FindLeaderFor(height int, rank int) identity.NodeID {
	return st.master
}
