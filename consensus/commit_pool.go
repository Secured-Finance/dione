package consensus

import (
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/ethereum/go-ethereum/crypto"

	"go.dedis.ch/kyber/v3/sign/bls"

	"go.dedis.ch/kyber/v3/pairing"
	"go.dedis.ch/kyber/v3/pairing/bn256"

	"go.dedis.ch/kyber/v3"

	types2 "github.com/Secured-Finance/dione/consensus/types"
)

type CommitPool struct {
	mut            sync.RWMutex
	commitMsgs     map[string][]*types2.Message
	blsPubKeyCache map[string]map[peer.ID]string // consensus id -> {peer id -> public key}
	blsCurveSuite  pairing.Suite
	blsKeyGroup    kyber.Group
	blsPrivKey     kyber.Scalar
	blsPubKey      kyber.Point
}

func NewCommitPool(blsPrivKeyStr string) (*CommitPool, error) {
	if blsPrivKeyStr == "" {
		return nil, fmt.Errorf("bls private key is empty")
	}
	suite := bn256.NewSuite()
	keyGroup := suite.G2()
	blsPrivKey, err := stringToScalar(keyGroup, blsPrivKeyStr)
	blsPubKey := suite.G2().Point().Mul(blsPrivKey, nil)
	if err != nil {
		return nil, err
	}

	return &CommitPool{
		blsCurveSuite:  suite,
		blsKeyGroup:    keyGroup,
		blsPrivKey:     blsPrivKey,
		blsPubKey:      blsPubKey,
		commitMsgs:     map[string][]*types2.Message{},
		blsPubKeyCache: map[string]map[peer.ID]string{},
	}, nil
}

func stringToScalar(g kyber.Group, s string) (kyber.Scalar, error) {
	key, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	sc := g.Scalar()
	err = sc.UnmarshalBinary(key)
	if err != nil {
		return nil, err
	}
	return sc, nil
}

func stringToPoint(g kyber.Group, s string) (kyber.Point, error) {
	key, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	sc := g.Point()
	err = sc.UnmarshalBinary(key)
	if err != nil {
		return nil, err
	}
	return sc, nil
}

func pointToString(p kyber.Point) (string, error) {
	d, err := p.MarshalBinary()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(d), nil
}

func (cp *CommitPool) CreateCommit(prepareMsg *types2.Message) (*types2.Message, error) {
	var message types2.Message
	message.Type = types2.MessageTypeCommit
	newCMsg := prepareMsg.Payload

	msgHash := crypto.Keccak256(newCMsg.Task.Payload)
	blsSignature, err := cp.createBLSSignature(msgHash)
	if err != nil {
		return nil, err
	}
	newCMsg.Signature = blsSignature

	pkEncoded, err := pointToString(cp.blsPubKey)
	if err != nil {
		return nil, err
	}
	newCMsg.BLSPubKey = pkEncoded

	message.Payload = newCMsg
	return &message, nil
}

func (cp *CommitPool) createBLSSignature(msg []byte) ([]byte, error) {
	signature, err := bls.Sign(cp.blsCurveSuite, cp.blsPrivKey, msg)
	return signature, err
}

func (cp *CommitPool) IsExistingCommit(commitMsg *types2.Message) bool {
	cp.mut.RLock()
	defer cp.mut.RUnlock()

	consensusMessage := commitMsg.Payload
	var exists bool
	for _, v := range cp.commitMsgs[consensusMessage.ConsensusID] {
		if v.From == commitMsg.From {
			exists = true
		}
	}
	return exists
}

func (cp *CommitPool) IsValidCommit(commit *types2.Message) bool {
	cp.mut.Lock()
	defer cp.mut.Unlock()

	if len(commit.Payload.BLSPubKey) == 0 {
		return false
	}

	if _, ok := cp.blsPubKeyCache[commit.Payload.ConsensusID]; !ok {
		cp.blsPubKeyCache[commit.Payload.ConsensusID] = map[peer.ID]string{}
	}

	if _, ok := cp.blsPubKeyCache[commit.Payload.ConsensusID][commit.From]; ok {
		return false
	}

	blsPubKeyPoint, err := stringToPoint(cp.blsKeyGroup, commit.Payload.BLSPubKey)
	if err != nil {
		return false
	}

	msgHash := crypto.Keccak256(commit.Payload.Task.Payload)
	err = bls.Verify(cp.blsCurveSuite, blsPubKeyPoint, msgHash, commit.Payload.Signature)
	if err != nil {
		return false
	}

	cp.blsPubKeyCache[commit.Payload.ConsensusID][commit.From] = commit.Payload.BLSPubKey

	return true
}

func (cp *CommitPool) AddCommit(commit *types2.Message) {
	cp.mut.Lock()
	defer cp.mut.Unlock()

	consensusID := commit.Payload.ConsensusID
	if _, ok := cp.commitMsgs[consensusID]; !ok {
		cp.commitMsgs[consensusID] = []*types2.Message{}
	}

	cp.commitMsgs[consensusID] = append(cp.commitMsgs[consensusID], commit)
}

func (cp *CommitPool) CommitSize(consensusID string) int {
	cp.mut.RLock()
	defer cp.mut.RUnlock()

	if v, ok := cp.commitMsgs[consensusID]; ok {
		return len(v)
	}
	return 0
}

// @returns agg sig, pks of sigs, error
func (cp *CommitPool) AggregateSignatureForConsensus(consensusID string) ([]byte, [][]byte, error) {
	cp.mut.RLock()
	defer cp.mut.RUnlock()

	var publicKeys [][]byte
	var sigs [][]byte
	for _, v := range cp.commitMsgs[consensusID] {
		pkS, ok := cp.blsPubKeyCache[consensusID][v.From]
		if !ok {
			continue
		}

		pkRaw, err := hex.DecodeString(pkS)
		if err != nil {
			return nil, nil, err
		}

		publicKeys = append(publicKeys, pkRaw)
		sigs = append(sigs, v.Payload.Signature)
	}

	aggSig, err := bls.AggregateSignatures(cp.blsCurveSuite, sigs...)
	if err != nil {
		return nil, nil, err
	}

	return aggSig, publicKeys, nil
}
