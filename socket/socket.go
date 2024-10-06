package socket

import (
	"sync"
	"time"

	"banyan/identity"
	"banyan/log"
	"banyan/transport"
	"banyan/utils"
)

// Socket integrates all networking interface and fault injections
type Socket interface {

	// Send put message to outbound queue
	Send(to identity.NodeID, m interface{})

	// Broadcast send to all peers
	Broadcast(m interface{})

	// Recv receives a message
	Recv() interface{}

	Close()
}

type socket struct {
	silence   bool
	id        identity.NodeID
	addresses map[identity.NodeID]string
	nodes     map[identity.NodeID]transport.Transport

	lock sync.RWMutex // locking map nodes
}

// NewSocket return Socket interface instance given self NodeID, node list, transport and codec name
func NewSocket(id identity.NodeID, addrs map[identity.NodeID]string, silence bool) Socket {
	socket := &socket{
		silence:   silence,
		id:        id,
		addresses: addrs,
		nodes:     make(map[identity.NodeID]transport.Transport),
	}

	socket.nodes[id] = transport.NewTransport(addrs[id])
	socket.nodes[id].Listen()

	return socket
}

func (s *socket) Send(to identity.NodeID, m interface{}) {
	//log.Debugf("node %s send message %+v to %v", s.id, m, to)

	s.lock.RLock()
	t, exists := s.nodes[to]
	s.lock.RUnlock()
	if !exists {
		s.lock.RLock()
		address, ok := s.addresses[to]
		s.lock.RUnlock()
		if !ok {
			log.Errorf("socket does not have address of node %s", to)
			return
		}
		t = transport.NewTransport(address)
		err := utils.Retry(t.Dial, 100, time.Duration(50)*time.Millisecond)
		if err != nil {
			panic(err)
		}
		s.lock.Lock()
		s.nodes[to] = t
		s.lock.Unlock()
	}

	if !s.silence {
		t.Send(m)
	}
}

func (s *socket) Recv() interface{} {
	s.lock.RLock()
	t := s.nodes[s.id]
	s.lock.RUnlock()
	for {
		m := t.Recv()
		return m
	}
}

func (s *socket) Broadcast(m interface{}) {
	//log.Debugf("node %s broadcasting message %+v", s.id, m)
	for id := range s.addresses {
		if id == s.id {
			continue
		}
		s.Send(id, m)
	}
	//log.Debugf("node %s done  broadcasting message %+v", s.id, m)
}

func (s *socket) Close() {
	for _, t := range s.nodes {
		t.Close()
	}
}
