package replica

import (
	"encoding/gob"
	"fmt"
	"time"

	"go.uber.org/atomic"

	"banyan/blockchain"
	"banyan/config"
	"banyan/election"
	"banyan/identity"
	"banyan/local_timeout"
	"banyan/log"
	"banyan/message"
	"banyan/node"
	"banyan/protocol"
	"strconv"
)

type Replica struct {
	node.Node
	Safety
	election.Election
	lt              *local_timeout.LocalTimeout
	start           chan bool // signal to start the node
	isStarted       atomic.Bool
	isByz           bool
	strategy        string
	timer           *time.Timer // timeout for each rank
	committedBlocks chan *blockchain.Block
	forkedBlocks    chan *blockchain.Block
	eventChan       chan interface{}

	/* for monitoring node statistics */

	experimentStartTime time.Time
	experimentDuration  time.Duration

	allBlockLatency      []time.Duration
	myBlockLatency       []time.Duration
	allBlockTimes        []time.Duration
	lastBlockProposeTime time.Time
	oneBlockPayloadBytes int
	committedBlockNo     int
	lastHeightTime       time.Time
	experimentStarted    bool
}

// NewReplica creates a new replica instance
func NewReplica(id identity.NodeID, alg string, isByz bool) *Replica {
	r := new(Replica)
	r.Node = node.NewNode(id, isByz)
	if isByz {
		log.Infof("[%v] is Byzantine", r.ID())
	}
	r.Election = election.NewRotation(config.GetConfig().N)

	r.allBlockLatency = make([]time.Duration, 10000)
	r.myBlockLatency = make([]time.Duration, 10000)
	r.allBlockTimes = make([]time.Duration, 10000)
	r.experimentDuration = time.Second * time.Duration(config.GetConfig().ExperimentDuration)
	r.experimentStarted = false

	r.oneBlockPayloadBytes = config.GetConfig().PayloadSize
	r.isByz = isByz
	r.strategy = config.GetConfig().Strategy
	r.lt = local_timeout.NewLocalTimeout()
	r.start = make(chan bool)
	r.eventChan = make(chan interface{}, 100)
	r.committedBlocks = make(chan *blockchain.Block, 100)
	r.forkedBlocks = make(chan *blockchain.Block, 100)
	r.Register(blockchain.Block{}, r.HandleBlock)
	r.Register(blockchain.NotarizationShare{}, r.HandleNotarizationShare)
	r.Register(blockchain.FinalizationShare{}, r.HandleFinalizationShare)
	r.Register(message.Query{}, r.handleQuery)
	gob.Register(blockchain.Block{})
	gob.Register(blockchain.NotarizationShare{})
	gob.Register(blockchain.FinalizationShare{})

	switch alg {
	case "icc":
		r.Safety = protocol.NewIcc(r.Node, r.Election, r.lt, r.committedBlocks, r.forkedBlocks)
	case "banyan":
		r.Safety = protocol.NewBanyan(r.Node, r.Election, r.lt, r.committedBlocks, r.forkedBlocks, config.GetConfig().F, config.GetConfig().P)
	default:
		r.Safety = protocol.NewBanyan(r.Node, r.Election, r.lt, r.committedBlocks, r.forkedBlocks, config.GetConfig().F, config.GetConfig().P)
	}
	return r
}

/* Message Handlers */

func (r *Replica) HandleBlock(block blockchain.Block) {
	log.Debugf("[%v] received a block from %v, height is %v, id: %x, prevID: %x", r.ID(), block.Proposer, block.Height, block.ID, block.PrevID)
	r.eventChan <- block
}

func (r *Replica) HandleNotarizationShare(vote blockchain.NotarizationShare) {
	log.Debugf("[%v] received a N share frm %v, blockID is %x", r.ID(), vote.Voter, vote.BlockID)
	r.eventChan <- vote
}

func (r *Replica) HandleFinalizationShare(vote blockchain.FinalizationShare) {
	log.Debugf("[%v] received a F share frm %v, blockID is %x", r.ID(), vote.Voter, vote.BlockID)
	r.eventChan <- vote
}

// handleQuery replies a query with the statistics of the node
func (r *Replica) handleQuery(m message.Query) {
	r.startSignal()

	if !(r.experimentStarted && r.experimentStartTime.Add(r.experimentDuration).Before(time.Now())) {
		status := fmt.Sprintf("Committed blocks: %v.\n", r.committedBlockNo)
		m.Reply(message.QueryReply{Info: status})
		return
	}

	response := "blockPayloadSize\n"
	response += strconv.Itoa(r.oneBlockPayloadBytes) + "\n"

	response += "committedBlocks\n"
	response += strconv.Itoa(r.committedBlockNo) + "\n"

	response += "allBlockLatency\n"
	for i := 3; i <= r.committedBlockNo; i++ {
		response += strconv.Itoa(int(r.allBlockLatency[i].Milliseconds())) + ","
	}

	response += "\nproposerLatency\n"
	for i := 3; i <= r.committedBlockNo; i++ {
		response += strconv.Itoa(int(r.myBlockLatency[i].Milliseconds())) + ","
	}

	response += "\nblockTime\n"
	for i := 3; i <= r.committedBlockNo; i++ {
		response += strconv.Itoa(int(r.allBlockTimes[i].Milliseconds())) + ","
	}

	m.Reply(message.QueryReply{Info: response})
}

/* Processors */

func (r *Replica) processCommittedBlock(block *blockchain.Block) {
	if block.Height == 3 {
		r.experimentStartTime = time.Now()
		r.experimentStarted = true
	}
	if r.experimentStartTime.Add(r.experimentDuration).Before(time.Now()) && (block.Height > 3) {
		return
	}

	proposeTime := block.Timestamp
	if block.Height > 1 {
		r.allBlockTimes[block.Height] = proposeTime.Sub(r.lastBlockProposeTime)
	}
	now := time.Now()
	r.allBlockLatency[block.Height] = now.Sub(proposeTime)
	if block.Proposer == r.ID() {
		r.myBlockLatency[block.Height] = r.allBlockLatency[block.Height]
	}
	r.committedBlockNo++
	r.lastBlockProposeTime = proposeTime

	log.Infof("[%v] the block is committed, height: %v, id: %x", r.ID(), block.Height, block.ID)
}

func (r *Replica) processForkedBlock(block *blockchain.Block) {
	log.Infof("[%v] the block is forked, No. of transactions: %v, height: %v, id: %x", r.ID(), len(block.Payload), block.Height, block.ID)
}

func (r *Replica) proposeIfLeader(height int, rank int) {
	if !r.IsLeader(r.ID(), height, rank) {
		return
	}
	r.proposeBlock(height, rank)
}

func (r *Replica) proposeBlock(height int, rank int) {
	block := r.Safety.MakeProposal(height, rank, r.oneBlockPayloadBytes)
	block.Timestamp = time.Now()
	r.Broadcast(block)
	_ = r.Safety.ProcessBlock(block)
}

// ListenLocalEvent listens new height and timeout events
func (r *Replica) ListenLocalEvent() {
	block_production_height := 1
	block_production_rank := 0
	<-r.start
	r.proposeIfLeader(block_production_height, block_production_rank)
	r.lastHeightTime = time.Now()
	r.timer = time.NewTimer(r.lt.GetTimeoutDuration())
	for {
		r.timer.Reset(r.lt.GetTimeoutDuration())
	L:
		for {
			select {
			case new_height := <-r.lt.GetNewHeight():
				block_production_height = new_height
				block_production_rank = 0
				r.proposeIfLeader(block_production_height, block_production_rank)
				// measure round time
				now := time.Now()
				lasts := now.Sub(r.lastHeightTime)
				r.lastHeightTime = now
				log.Debugf("[%v] the last height lasted %v milliseconds", r.ID(), lasts.Milliseconds())
				break L
			case <-r.timer.C:
				block_production_rank += 1
				r.proposeIfLeader(block_production_height, block_production_rank)
				break L
			}
		}
	}
}

// ListenCommittedBlocks listens committed blocks and forked blocks from the protocols
func (r *Replica) ListenCommittedBlocks() {
	for {
		select {
		case committedBlock := <-r.committedBlocks:
			r.processCommittedBlock(committedBlock)
		case forkedBlock := <-r.forkedBlocks:
			r.processForkedBlock(forkedBlock)
		}
	}
}

func (r *Replica) startSignal() {
	if !r.isStarted.Load() {
		log.Debugf("[%v] is boosting", r.ID())
		r.isStarted.Store(true)
		r.start <- true
	}
}

// Starts event loop
func (r *Replica) Start() {
	go r.Run()

	go r.ListenLocalEvent()
	go r.ListenCommittedBlocks()
	for {
		event := <-r.eventChan
		r.startSignal()
		switch v := event.(type) {
		case blockchain.Block:
			r.Safety.ProcessBlock(&v)
		case blockchain.NotarizationShare:
			r.Safety.ProcessNotarizationShare(&v)
		case blockchain.FinalizationShare:
			r.Safety.ProcessFinalizationShare(&v)
		}
	}
}
