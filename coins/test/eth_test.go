package test

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/davecgh/go-spew/spew"
	"github.com/phanhoc/wallet-bip44/coins"
	"github.com/tyler-smith/go-bip39"
	"testing"
)

func TestNewAccountETH(t *testing.T) {
	seed := bip39.NewSeed("announce parent popular hybrid fine maid exile impulse unknown school castle wage hand impulse wing", "")
	t.Log(hex.EncodeToString(seed))
	master, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		t.Fatalf("failed to new master, err: %v", err)
	}
	expectedAccount := "xprv9zS1ivPzUYQ8rrZXRzdkzgurFjqP2HZ5dmt4vGjoFLzXKRJLY4o192mf36DhaCEJFN9HnUNUZMSCdQan19FqDPdGioRbf2nCNpKmJ2SJt86"
	eth := coins.NewEth()
	account, err := eth.NewAccount(master.String(), 0)
	if err != nil {
		t.Fatalf("failed to new master, err: %v", err)
	}
	t.Log(spew.Sdump(account))
	if account.AccountPrivateKey != expectedAccount {
		t.Fatalf("got: %s, expected: %s", account.AccountPrivateKey, expectedAccount)
	}
}

func TestBtc_GenerateNormalAddressEth(t *testing.T) {
	seed := bip39.NewSeed("announce parent popular hybrid fine maid exile impulse unknown school castle wage hand impulse wing", "")
	t.Log(hex.EncodeToString(seed))
	master, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		t.Fatalf("failed to new master, err: %v", err)
	}
	expectedAddress := "0x6d3D71E7A30dc0f7BE1DF407fA361240f63038Ec"
	tbtc := coins.NewEth()
	account, err := tbtc.NewAccount(master.String(), 0)
	if err != nil {
		t.Fatalf("failed to new master, err: %v", err)
	}
	t.Log(account.AccountPrivateKey)

	address, err := tbtc.GenerateNormalAddress(account.AccountPrivateKey, 0, false)
	t.Log(spew.Sdump(address))
	if expectedAddress != address.Address {
		t.Fatalf("failed to generate address, got %s - expected %s", address.Address, expectedAddress)
	}
}