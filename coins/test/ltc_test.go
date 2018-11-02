package test

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/phanhoc/wallet-bip44/coins"
	"github.com/tyler-smith/go-bip39"
	"testing"
)

func TestLtc_NewAccount(t *testing.T) {
	ltcChain := &chaincfg.MainNetParams
	ltcChain.HDCoinType = 2
	ltcChain.Bech32HRPSegwit = "ltc"
	ltcChain.PubKeyHashAddrID = 0x30        // starts with L
	ltcChain.ScriptHashAddrID = 0x32        // starts with M
	ltcChain.PrivateKeyID = 0xB0            // starts with 6 (uncompressed) or T (compressed)
	ltcChain.WitnessPubKeyHashAddrID = 0x06 // starts with p2
	ltcChain.WitnessScriptHashAddrID = 0x0A // starts with 7Xh

	// BIP32 hierarchical deterministic extended key magics
	ltcChain.HDPrivateKeyID = [4]byte{0x04, 0x88, 0xad, 0xe4} // starts with xprv
	ltcChain.HDPublicKeyID = [4]byte{0x04, 0x88, 0xb2, 0x1e}  // starts with xpub
	seed := bip39.NewSeed("announce parent popular hybrid fine maid exile impulse unknown school castle wage hand impulse wing", "")
	t.Log(hex.EncodeToString(seed))
	master, err := hdkeychain.NewMaster(seed, ltcChain)
	if err != nil {
		t.Fatalf("failed to new master, err: %v", err)
	}
	expectedAccount := "xprv9xpaCu1436QpkNN1uZFAL41FEbVb7cSCejWXKJzNt8T2V9SGVpD4VxMwfX5C5o7vHZH4VVgKEw9oEzMSY9cmxKhY8KhJesF6dE5pEetoMA1"
	mString := master.String()
	ltc := coins.NewLtc()
	account, err := ltc.NewAccount(mString, 0)
	if err != nil {
		t.Fatalf("failed to new master, err: %v", err)
	}
	t.Log(account.AccountPrivateKey)
	if account.AccountPrivateKey != expectedAccount {
		t.Fatalf("got: %s, expected: %s", account.AccountPrivateKey, expectedAccount)
	}
}

func TestLtc_GenerateNormalAddress(t *testing.T) {
	seed := bip39.NewSeed("announce parent popular hybrid fine maid exile impulse unknown school castle wage hand impulse wing", "")
	t.Log(hex.EncodeToString(seed))
	master, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		t.Fatalf("failed to new master, err: %v", err)
	}
	expectedAddress := "LXBJmeqopk6bTYXcvUwH3TcTWM2gsZN2Ap"
	tbtc := coins.NewLtc()
	account, err := tbtc.NewAccount(master.String(), 0)
	if err != nil {
		t.Fatalf("failed to new master, err: %v", err)
	}
	t.Log(account.AccountPrivateKey)

	address, err := tbtc.GenerateNormalAddress(account.AccountPrivateKey, 0, false)
	t.Log(address.Address)
	if expectedAddress != address.Address {
		t.Fatalf("failed to generate address, got %s - expected %s", address.Address, expectedAddress)
	}
}
