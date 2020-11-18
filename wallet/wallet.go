package wallet

import (
	"context"
	"sync"

	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/Secured-Finance/dione/sigs"
	_ "github.com/Secured-Finance/dione/sigs/ed25519" // enable ed25519 signatures
	"github.com/Secured-Finance/dione/types"
	"github.com/filecoin-project/go-address"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

const (
	KNamePrefix = "wallet-"
	KDefault    = "default"
)

type LocalWallet struct {
	keys     map[peer.ID]*Key
	keystore types.KeyStore

	lk sync.Mutex
}

type Default interface {
	GetDefault() (address.Address, error)
}

func NewWallet(keystore types.KeyStore) (*LocalWallet, error) {
	w := &LocalWallet{
		keys:     make(map[peer.ID]*Key),
		keystore: keystore,
	}

	return w, nil
}

func KeyWallet(keys ...*Key) *LocalWallet {
	m := make(map[peer.ID]*Key)
	for _, key := range keys {
		m[key.Address] = key
	}

	return &LocalWallet{
		keys: m,
	}
}

func (w *LocalWallet) WalletSign(ctx context.Context, addr peer.ID, msg []byte) (*types.Signature, error) {
	ki, err := w.findKey(addr)
	if err != nil {
		return nil, err
	}
	if ki == nil {
		return nil, xerrors.Errorf("failed to find private key of %w: %w", addr.String(), types.ErrKeyInfoNotFound)
	}

	return sigs.Sign(ActSigType(ki.Type), ki.PrivateKey, msg)
}

func (w *LocalWallet) findKey(addr peer.ID) (*Key, error) {
	w.lk.Lock()
	defer w.lk.Unlock()

	k, ok := w.keys[addr]
	if ok {
		return k, nil
	}
	if w.keystore == nil {
		logrus.Warn("findKey didn't find the key in in-memory wallet")
		return nil, nil
	}

	ki, err := w.tryFind(addr)
	if err != nil {
		if xerrors.Is(err, types.ErrKeyInfoNotFound) {
			return nil, nil
		}
		return nil, xerrors.Errorf("getting from keystore: %w", err)
	}
	k, err = NewKey(ki)
	if err != nil {
		return nil, xerrors.Errorf("decoding from keystore: %w", err)
	}
	w.keys[k.Address] = k
	return k, nil
}

func (w *LocalWallet) tryFind(addr peer.ID) (types.KeyInfo, error) {

	ki, err := w.keystore.Get(KNamePrefix + addr.String())
	if err == nil {
		return ki, err
	}

	if !xerrors.Is(err, types.ErrKeyInfoNotFound) {
		return types.KeyInfo{}, err
	}

	// We got an ErrKeyInfoNotFound error
	// Try again, this time with the testnet prefix

	tAddress, err := swapMainnetForTestnetPrefix(addr.String())
	if err != nil {
		return types.KeyInfo{}, err
	}
	logrus.Info("tAddress: ", tAddress)

	ki, err = w.keystore.Get(KNamePrefix + tAddress)
	if err != nil {
		return types.KeyInfo{}, err
	}
	logrus.Info("ki from tryFind: ", ki)

	// We found it with the testnet prefix
	// Add this KeyInfo with the mainnet prefix address string
	err = w.keystore.Put(KNamePrefix+addr.String(), ki)
	if err != nil {
		return types.KeyInfo{}, err
	}

	return ki, nil
}

func (w *LocalWallet) GetDefault() (peer.ID, error) {
	w.lk.Lock()
	defer w.lk.Unlock()

	ki, err := w.keystore.Get(KDefault)
	if err != nil {
		return "", xerrors.Errorf("failed to get default key: %w", err)
	}

	k, err := NewKey(ki)
	if err != nil {
		return "", xerrors.Errorf("failed to read default key from keystore: %w", err)
	}

	return k.Address, nil
}

func (w *LocalWallet) WalletNew(ctx context.Context, typ types.KeyType) (peer.ID, error) {
	w.lk.Lock()
	defer w.lk.Unlock()

	k, err := GenerateKey(typ)
	if err != nil {
		return "", err
	}

	if err := w.keystore.Put(KNamePrefix+k.Address.String(), k.KeyInfo); err != nil {
		return "", xerrors.Errorf("saving to keystore: %w", err)
	}
	w.keys[k.Address] = k

	_, err = w.keystore.Get(KDefault)
	if err != nil {
		if !xerrors.Is(err, types.ErrKeyInfoNotFound) {
			return "", err
		}

		if err := w.keystore.Put(KDefault, k.KeyInfo); err != nil {
			return "", xerrors.Errorf("failed to set new key as default: %w", err)
		}
	}

	return k.Address, nil
}

func (w *LocalWallet) WalletHas(ctx context.Context, addr peer.ID) (bool, error) {
	k, err := w.findKey(addr)
	if err != nil {
		return false, err
	}
	return k != nil, nil
}

func swapMainnetForTestnetPrefix(addr string) (string, error) {
	aChars := []rune(addr)
	prefixRunes := []rune(address.TestnetPrefix)
	if len(prefixRunes) != 1 {
		return "", xerrors.Errorf("unexpected prefix length: %d", len(prefixRunes))
	}

	aChars[0] = prefixRunes[0]
	return string(aChars), nil
}
