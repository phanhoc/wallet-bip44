package test

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/phanhoc/wallet-bip44/coins"
	"github.com/tyler-smith/go-bip39"
	"testing"
)

func TestNewAccount(t *testing.T) {
	seed := bip39.NewSeed("announce parent popular hybrid fine maid exile impulse unknown school castle wage hand impulse wing", "")
	t.Log(hex.EncodeToString(seed))
	master, err := hdkeychain.NewMaster(seed, &chaincfg.TestNet3Params)
	if err != nil {
		t.Fatalf("failed to new master, err: %v", err)
	}
	expectedAccount := "tprv8fT91mKMRW1BVeCS9GtqqG6kN29CgZbaFFUxgNhheJ91iizKXysr7TRFeaxWAwrsLKNK7uQXnKtYWZEE2CWyebHPfFUChiXbXPnewdVx5So"
	tbtc := coins.NewTBtc()
	account, err := tbtc.NewAccount(master.String(), 0)
	if err != nil {
		t.Fatalf("failed to new master, err: %v", err)
	}
	t.Log(account.AccountPrivateKey)
	if account.AccountPrivateKey != expectedAccount {
		t.Fatalf("got: %s, expected: %s", account.AccountPrivateKey, expectedAccount)
	}
}

func TestBtc_GenerateNormalAddress(t *testing.T) {
	seed := bip39.NewSeed("announce parent popular hybrid fine maid exile impulse unknown school castle wage hand impulse wing", "")
	t.Log(hex.EncodeToString(seed))
	master, err := hdkeychain.NewMaster(seed, &chaincfg.TestNet3Params)
	if err != nil {
		t.Fatalf("failed to new master, err: %v", err)
	}
	expectedAddress := "mucfitzpxZMPYbw6X8j6R2S5YXBsiR1JLg"
	tbtc := coins.NewTBtc()
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

func TestTBtc_GenerateMultisigAddress(t *testing.T) {
	seed := bip39.NewSeed("announce parent popular hybrid fine maid exile impulse unknown school castle wage hand impulse wing", "")
	t.Log(hex.EncodeToString(seed))
	master, err := hdkeychain.NewMaster(seed, &chaincfg.TestNet3Params)
	if err != nil {
		t.Fatalf("failed to new master, err: %v", err)
	}
	expectedAddress := "2MxT3KJfogSPt7i93z61kRRrNP4GY78q4pv"
	tbtc := coins.NewTBtc()
	account_0, err := tbtc.NewAccount(master.String(), 0)
	if err != nil {
		t.Fatalf("failed to new master, err: %v", err)
	}
	account_1, err := tbtc.NewAccount(master.String(), 1)
	if err != nil {
		t.Fatalf("failed to new master, err: %v", err)
	}
	account_2, err := tbtc.NewAccount(master.String(), 2)
	if err != nil {
		t.Fatalf("failed to new master, err: %v", err)
	}
	xPub := make([]string, 0, 3)
	xPub = append(xPub, account_0.AccountPrivateKey)
	xPub = append(xPub, account_1.AccountPrivateKey)
	xPub = append(xPub, account_2.AccountPrivateKey)

	address, err := tbtc.GenerateMultisigAddress(xPub, 0, 2, 3, false)
	t.Log(address.Address)
	t.Log(hex.EncodeToString(address.RedeemScript))
	if expectedAddress != address.Address {
		t.Fatalf("failed to generate address, got %s - expected %s", address.Address, expectedAddress)
	}
}
