package consensus

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/ipfs/go-log"
	consensus "github.com/libp2p/go-libp2p-consensus"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"

	raft "github.com/hashicorp/raft"
	libp2praft "github.com/libp2p/go-libp2p-raft"
)

var raftTmpFolder = "raft-consensus"
var raftQuiet = true

type RaftConsensus struct {
	Logger    *log.ZapEventLogger
	Raft      *raft.Raft
	Consensus *libp2praft.Consensus
	State     consensus.State
	Transport *raft.NetworkTransport
}

type RaftState struct {
	Msg string
}

func NewRaftConsensus() *RaftConsensus {
	raftConsensus := &RaftConsensus{
		Logger: log.Logger("rendezvous"),
	}

	return raftConsensus

}

func (raftConsensus *RaftConsensus) NewState(value string) {
	raftConsensus.State = &RaftState{value}
}

func (raftConsensus *RaftConsensus) NewConsensus(op consensus.Op) {
	if op != nil {
		raftConsensus.Consensus = libp2praft.NewOpLog(&RaftState{}, op)
	} else {
		raftConsensus.Consensus = libp2praft.NewConsensus(&RaftState{"i am not consensuated"})
	}
}

func (raftConsensus *RaftConsensus) GetConsensusState(consensus *libp2praft.Consensus) (consensus.State, error) {
	var err error
	raftConsensus.State, err = raftConsensus.Consensus.GetCurrentState()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return raftConsensus.State, nil
}

func (raftConsensus *RaftConsensus) UpdateState(value string) error {
	raftConsensus.NewState(value)

	// CommitState() blocks until the state has been
	// agreed upon by everyone
	agreedState, err := raftConsensus.Consensus.CommitState(raftConsensus.State)
	if err != nil {
		raftConsensus.Logger.Warn("Failed to commit new state", err)
		return err
	}
	if agreedState == nil {
		fmt.Println("agreedState is nil: commited on a non-leader?")
		return err
	}
	return nil
}

func (raftConsensus *RaftConsensus) Shutdown() {
	err := raftConsensus.Raft.Shutdown().Error()
	if err != nil {
		raftConsensus.Logger.Fatal(err)
	}
}

func (raftConsensus *RaftConsensus) WaitForLeader(r *raft.Raft) {
	obsCh := make(chan raft.Observation, 1)
	observer := raft.NewObserver(obsCh, false, nil)
	r.RegisterObserver(observer)
	defer r.DeregisterObserver(observer)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	ticker := time.NewTicker(time.Second / 2)
	defer ticker.Stop()
	for {
		select {
		case obs := <-obsCh:
			switch obs.Data.(type) {
			case raft.RaftState:
				if r.Leader() != "" {
					return
				}
			}
		case <-ticker.C:
			if r.Leader() != "" {
				return
			}
		case <-ctx.Done():
			raftConsensus.Logger.Fatal("timed out waiting for Leader")
		}
	}
}

func (raftConsensus *RaftConsensus) MakeConsensus(h host.Host, pids []peer.ID, op consensus.Op) {
	raftConsensus.NewConsensus(op)

	// -- Create Raft servers configuration
	servers := make([]raft.Server, len(pids))
	for i, pid := range pids {
		servers[i] = raft.Server{
			Suffrage: raft.Voter,
			ID:       raft.ServerID(pid.Pretty()),
			Address:  raft.ServerAddress(pid.Pretty()),
		}
	}
	serverConfig := raft.Configuration{
		Servers: servers,
	}

	// -- Create LibP2P transports Raft
	transport, err := libp2praft.NewLibp2pTransport(h, 2*time.Second)
	if err != nil {
		raftConsensus.Logger.Fatal(err)
	}
	raftConsensus.Transport = transport

	// -- Configuration
	config := raft.DefaultConfig()
	if raftQuiet {
		config.LogOutput = ioutil.Discard
		config.Logger = nil
	}
	config.LocalID = raft.ServerID(h.ID().Pretty())

	// -- SnapshotStore
	snapshots, err := raft.NewFileSnapshotStore(raftTmpFolder, 3, nil)
	if err != nil {
		raftConsensus.Logger.Fatal(err)
	}

	// -- Log store and stable store: we use inmem.
	logStore := raft.NewInmemStore()

	// -- Boostrap everything if necessary
	bootstrapped, err := raft.HasExistingState(logStore, logStore, snapshots)
	if err != nil {
		raftConsensus.Logger.Fatal(err)
	}

	if !bootstrapped {
		raft.BootstrapCluster(config, logStore, logStore, snapshots, transport, serverConfig)
	} else {
		raftConsensus.Logger.Info("Already initialized!!")
	}

	raft, err := raft.NewRaft(config, raftConsensus.Consensus.FSM(), logStore, logStore, snapshots, transport)
	if err != nil {
		raftConsensus.Logger.Fatal(err)
	}
	raftConsensus.Raft = raft
}

func (raftConsensus *RaftConsensus) StartConsensus(host host.Host, peers []peer.ID) {
	raftConsensus.MakeConsensus(host, peers, nil)

	// Create the actors using the Raft nodes
	actor := libp2praft.NewActor(raftConsensus.Raft)

	// Set the actors so that we can CommitState() and GetCurrentState()
	raftConsensus.Consensus.SetActor(actor)

	// Provide some time for leader election
	time.Sleep(5 * time.Second)

	// Run the 1000 updates on the leader
	// Barrier() will wait until updates have been applied
	if actor.IsLeader() {
		value := "new value"
		if err := raftConsensus.UpdateState(value); err != nil {
			raftConsensus.Logger.Fatal("Failed to update state", err)
		}
	} else {
		raftConsensus.WaitForLeader(raftConsensus.Raft)
	}
}
