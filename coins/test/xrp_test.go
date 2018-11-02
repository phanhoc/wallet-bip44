package test

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/phanhoc/wallet-bip44/coins"
	"github.com/tyler-smith/go-bip39"
	"testing"
)

func TestXrp_GenerateMultisigAddress(t *testing.T) {
	seed := bip39.NewSeed("announce parent popular hybrid fine maid exile impulse unknown school castle wage hand impulse wing", "")
	t.Log(hex.EncodeToString(seed))
	master, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		t.Fatalf("failed to new master, err: %v", err)
	}
	expectedAddress := "r38r1gszFSsYKmmgp7b86jKEFU9ywg1WfK"
	tbtc := coins.NewXrp()
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
