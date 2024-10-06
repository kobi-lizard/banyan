package main

import (
	"banyan"
	"flag"
	"strconv"
	"sync"

	"banyan/config"
	"banyan/crypto"
	"banyan/identity"
	"banyan/log"
	"banyan/replica"
)

var algorithm = flag.String("algorithm", "hotstuff", "BFT consensus algorithm")
var id = flag.String("id", "", "NodeID of the node")
var simulation = flag.Bool("sim", false, "simulation mode")

func initReplica(id identity.NodeID, isByz bool) {
	log.Infof("node %v starting...", id)
	if isByz {
		log.Infof("node %v is Byzantine", id)
	}

	switch *algorithm {
	case "icc":
		r := replica.NewReplica(id, *algorithm, isByz)
		r.Start()
	case "banyan":
		r := replica.NewReplica(id, *algorithm, isByz)
		r.Start()
	default:
		r := replica.NewReplicaView(id, *algorithm, isByz)
		r.Start()
	}
}

func main() {
	banyan.Init()
	// the private and public keys are generated here
	errCrypto := crypto.SetKeys()
	if errCrypto != nil {
		log.Fatal("Could not generate keys:", errCrypto)
	}
	if *simulation {
		var wg sync.WaitGroup
		wg.Add(1)
		config.Simulation()
		for id := range config.GetConfig().Addrs {
			isByz := false
			if id.Node() > config.GetConfig().N-config.GetConfig().ByzNo {
				isByz = true
			}
			go initReplica(id, isByz)
		}
		wg.Wait()
	} else {
		isByz := false
		i, _ := strconv.Atoi(*id)
		if i > config.GetConfig().N-config.GetConfig().ByzNo {
			isByz = true
		}
		initReplica(identity.NodeID(*id), isByz)
	}
}
