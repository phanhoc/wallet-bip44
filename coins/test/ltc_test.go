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
	ltcChain := new(chaincfg.Params)
	ltcChain.HDCoinType = 2
	ltcChain.Bech32HRPSegwit = "ltc"
	ltcChain.PubKeyHashAddrID = 0x30        // starts with L
	ltcChain.ScriptHashAddrID = 0x32        // starts with M
	ltcChain.PrivateKeyID = 0xB0            // starts with 6 (uncompressed) or T (compressed)
	ltcChain.WitnessPubKeyHashAddrID = 0x06 // starts with p2
	ltcChain.WitnessScriptHashAddrID = 0x0A // starts with 7Xh

	// BIP32 hierarchical deterministic extended key magics
	ltcChain.HDPrivateKeyID = [4]byte{0x01, 0x9d, 0x9c, 0xfe} // starts with Lxprv
	ltcChain.HDPublicKeyID = [4]byte{0x01, 0x9d, 0xa4, 0x62}  // starts with Lxpub

	err := chaincfg.Register(ltcChain)

	if err != nil {
		t.Fatalf("failed to register, %v", err)
	}

	seed := bip39.NewSeed("announce parent popular hybrid fine maid exile impulse unknown school castle wage hand impulse wing", "")
	t.Log(hex.EncodeToString(seed))
	master, err := hdkeychain.NewMaster(seed, ltcChain)
	if err != nil {
		t.Fatalf("failed to new master, err: %v", err)
	}
	expectedAccount := "Ltpv774hUGKSWkkAbM1bMTHnUUDzRrxTCc7Y4iY5GoNosn37VxD9NSycro4bdHcCxWm7CfrLa1KqQ1rTEVUWwGbkbS4t5teUxpJ8RchiJ5r6GqC"
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
	ltcChain := new(chaincfg.Params)
	ltcChain.HDCoinType = 2
	ltcChain.Bech32HRPSegwit = "ltc"
	ltcChain.PubKeyHashAddrID = 0x30        // starts with L
	ltcChain.ScriptHashAddrID = 0x32        // starts with M
	ltcChain.PrivateKeyID = 0xB0            // starts with 6 (uncompressed) or T (compressed)
	ltcChain.WitnessPubKeyHashAddrID = 0x06 // starts with p2
	ltcChain.WitnessScriptHashAddrID = 0x0A // starts with 7Xh

	// BIP32 hierarchical deterministic extended key magics
	ltcChain.HDPrivateKeyID = [4]byte{0x01, 0x9d, 0x9c, 0xfe} // starts with Lxprv
	ltcChain.HDPublicKeyID = [4]byte{0x01, 0x9d, 0xa4, 0x62}  // starts with Lxpub
	chaincfg.Register(ltcChain)

	seed := bip39.NewSeed("announce parent popular hybrid fine maid exile impulse unknown school castle wage hand impulse wing", "")
	t.Log(hex.EncodeToString(seed))
	master, err := hdkeychain.NewMaster(seed, ltcChain)
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
