package test

import (
	"encoding/hex"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gcash/bchd/chaincfg"
	"github.com/gcash/bchd/txscript"
	"github.com/gcash/bchutil/hdkeychain"
	"github.com/phanhoc/wallet-bip44/coins"
	"github.com/tyler-smith/go-bip39"
	"testing"
)

func TestNewAccountBCH(t *testing.T) {
	seed := bip39.NewSeed("announce parent popular hybrid fine maid exile impulse unknown school castle wage hand impulse wing", "")
	t.Log(hex.EncodeToString(seed))
	master, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		t.Fatalf("failed to new master, err: %v", err)
	}
	expectedAccount := "xprv9zBgSSmb5TdBEsxUQ43Dc9MQM9ve1xo1Qp4ZMdNSVV78Wukbp9nA2PHfitgQKkHUuuK8J4GcfQxFyeMcnMaod6boWdtxhYmN1Rshi8oVdFA"
	tbtc := coins.NewBch()
	account, err := tbtc.NewAccount(master.String(), 0)
	if err != nil {
		t.Fatalf("failed to new master, err: %v", err)
	}
	t.Log(spew.Sdump(account))
	if account.AccountPrivateKey != expectedAccount {
		t.Fatalf("got: %s, expected: %s", account.AccountPrivateKey, expectedAccount)
	}
}

func TestBtc_GenerateNormalAddressBCH(t *testing.T) {
	seed := bip39.NewSeed("announce parent popular hybrid fine maid exile impulse unknown school castle wage hand impulse wing", "")
	t.Log(hex.EncodeToString(seed))
	master, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		t.Fatalf("failed to new master, err: %v", err)
	}
	expectedAddress := "qza4l4rgr3emd4rjy9zxdhenwvjhkphprvcasc58ct"
	tbtc := coins.NewBch()
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

func TestBch_GenerateMultisigAddress(t *testing.T) {
	seed := bip39.NewSeed("announce parent popular hybrid fine maid exile impulse unknown school castle wage hand impulse wing", "")
	t.Log(hex.EncodeToString(seed))
	master, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		t.Fatalf("failed to new master, err: %v", err)
	}
	expectedAddress := "pr46qq772t3j7xnkua4k0qk5z3vwrzlt5ysq2395z2"
	tbtc := coins.NewBch()
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
	if err != nil {
		t.Fatalf("failed to generate multisig, err: %v", err)
	}
	t.Log(address.Address)
	t.Log(hex.EncodeToString(address.RedeemScript))
	_, addresses, _, _ := txscript.ExtractPkScriptAddrs(address.RedeemScript, &chaincfg.MainNetParams)
	fmt.Println(spew.Sdump(addresses))
	fmt.Println(addresses[0].EncodeAddress())

	if expectedAddress != address.Address {
		t.Fatalf("failed to generate address, got %s - expected %s", address.Address, expectedAddress)
	}

}
