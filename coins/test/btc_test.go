package test

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/phanhoc/wallet-bip44/coins"
	"github.com/tyler-smith/go-bip39"
	"testing"
)

func TestNewAccountBTC(t *testing.T) {
	seed := bip39.NewSeed("announce parent popular hybrid fine maid exile impulse unknown school castle wage hand impulse wing", "")
	t.Log(hex.EncodeToString(seed))
	master, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		t.Fatalf("failed to new master, err: %v", err)
	}
	expectedAccount := "xprv9z9UCSxob6vZL7WkwjCPFJVj27R2gk11r4QD18wKH7LqrhiBXNg1ca3h7iyzULTCJa15JdyzPbArNCh9auTn8WDWEq2Fkf44sajGQQ2wYSs"
	tbtc := coins.NewBtc()
	account, err := tbtc.NewAccount(master.String(), 0)
	if err != nil {
		t.Fatalf("failed to new master, err: %v", err)
	}
	t.Log(account.AccountPrivateKey)
	if account.AccountPrivateKey != expectedAccount {
		t.Fatalf("got: %s, expected: %s", account.AccountPrivateKey, expectedAccount)
	}
}

func TestBtc_GenerateNormalAddressBTC(t *testing.T) {
	seed := bip39.NewSeed("announce parent popular hybrid fine maid exile impulse unknown school castle wage hand impulse wing", "")
	t.Log(hex.EncodeToString(seed))
	master, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		t.Fatalf("failed to new master, err: %v", err)
	}
	expectedAddress := "1Jgub4PAFFXabXemD1VZ143BSCBf1ax2S7"
	tbtc := coins.NewBtc()
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

func TestBtc_GenerateMultisigAddress(t *testing.T) {
	seed := bip39.NewSeed("announce parent popular hybrid fine maid exile impulse unknown school castle wage hand impulse wing", "")
	t.Log(hex.EncodeToString(seed))
	master, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		t.Fatalf("failed to new master, err: %v", err)
	}
	expectedAddress := "3NB8ZTiBpBnza7ueKwXvtmkQVS3miJZTiW"
	tbtc := coins.NewBtc()
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
