package test

import (
	"encoding/hex"
	"github.com/phanhoc/wallet-bip44/wallets"
	"github.com/phanhoc/wallet-bip44/wallets/common"
	"github.com/tyler-smith/go-bip39"
	"testing"
)

func TestBCH_NewMasterKey(t *testing.T) {
	seed := bip39.NewSeed("announce parent popular hybrid fine maid exile impulse unknown school castle wage hand impulse wing", "")
	t.Log(hex.EncodeToString(seed))
	btcWallet := wallets.NewBchWallet(common.OneOfOne)
	masterKeyEx := "xprv9s21ZrQH143K26uUQxSgabsXyzC3VbzBsFs3dZuQLTbYgewoJbnvpUreykr8hThdNL238QdVTdhg7ns5jQW8veYDjuZRVqNppWEwc2wSuqB"
	masterKey, err := btcWallet.NewMaster(seed)
	if err != nil {
		t.Fatalf("failed to new master key, err %v", err)
	}
	t.Log(masterKey)
	if masterKeyEx != masterKey {
		t.Fatalf("got %s - ex %s", masterKey, masterKeyEx)
	}
}

func TestBCH_NewAccount(t *testing.T) {
	seed := bip39.NewSeed("announce parent popular hybrid fine maid exile impulse unknown school castle wage hand impulse wing", "")
	t.Log(hex.EncodeToString(seed))
	btcWallet := wallets.NewBchWallet(common.OneOfOne)
	masterKeyEx := "xprv9s21ZrQH143K26uUQxSgabsXyzC3VbzBsFs3dZuQLTbYgewoJbnvpUreykr8hThdNL238QdVTdhg7ns5jQW8veYDjuZRVqNppWEwc2wSuqB"
	accountKeyEx := "xprv9zBgSSmb5TdBEsxUQ43Dc9MQM9ve1xo1Qp4ZMdNSVV78Wukbp9nA2PHfitgQKkHUuuK8J4GcfQxFyeMcnMaod6boWdtxhYmN1Rshi8oVdFA"
	masterKey, err := btcWallet.NewMaster(seed)
	if err != nil {
		t.Fatalf("failed to new master key, err %v", err)
	}
	t.Log(masterKey)
	if masterKeyEx != masterKey {
		t.Fatalf("got %s - ex %s", masterKey, masterKeyEx)
	}
	account0, err := btcWallet.NewAccount(masterKey, 0)
	if err != nil {
		t.Fatalf("failed to new account, err %v", err)
	}
	if accountKeyEx != account0.AccountPrivateKey {
		t.Fatalf("failed to create new account, got %s - expected %s", account0.AccountPrivateKey, accountKeyEx)
	}
}

func TestBCH_NextAddressNormal(t *testing.T) {
	seed := bip39.NewSeed("announce parent popular hybrid fine maid exile impulse unknown school castle wage hand impulse wing", "")
	t.Log(hex.EncodeToString(seed))
	btcWallet := wallets.NewBchWallet(common.OneOfOne)
	masterKeyEx := "xprv9s21ZrQH143K26uUQxSgabsXyzC3VbzBsFs3dZuQLTbYgewoJbnvpUreykr8hThdNL238QdVTdhg7ns5jQW8veYDjuZRVqNppWEwc2wSuqB"
	addressEx := "qza4l4rgr3emd4rjy9zxdhenwvjhkphprvcasc58ct"
	masterKey, err := btcWallet.NewMaster(seed)
	if err != nil {
		t.Fatalf("failed to new master key, err %v", err)
	}
	t.Log(masterKey)
	if masterKeyEx != masterKey {
		t.Fatalf("got %s - ex %s", masterKey, masterKeyEx)
	}
	account0, err := btcWallet.NewAccount(masterKey, 0)
	if err != nil {
		t.Fatalf("failed to new account, err %v", err)
	}
	addressInfo, err := btcWallet.NextAddress([]string{account0.AccountPublicKey}, wallets.AddressOptional{AddressIndex: 0, Internal: false})
	if err != nil {
		t.Fatalf("failed to next address, err %v", err)
	}
	t.Log(addressInfo.Address)
	if addressEx != addressInfo.Address {
		t.Fatalf("failed to generate address, got %s - expected %s", addressInfo.Address, addressEx)
	}
}
