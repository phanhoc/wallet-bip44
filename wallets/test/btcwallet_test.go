package test

import (
	"encoding/hex"
	"github.com/phanhoc/wallet-bip44/wallets"
	"github.com/phanhoc/wallet-bip44/wallets/common"
	"github.com/tyler-smith/go-bip39"
	"testing"
)

func TestBtc_Mnemonic(t *testing.T) {
	btcWallet := wallets.NewBtcWallet(common.OneOfOne)
	mnem, err := btcWallet.NewMnemonic(common.JAPANESE)
	if err != nil {
		t.Fatalf("failed to generate mnemonic, err: %v", err)
	}
	t.Log(mnem)
}

func TestBtc_Seed(t *testing.T) {
	btcWallet := wallets.NewBtcWallet(common.OneOfOne)
	seed, err := btcWallet.NewSeed(common.JAPANESE, "つみき くふう てつや まとめ ひるま かんしゃ あてはまる ならび けってい むのう じてん かいよう とこや くちこみ たもつ たべる よごれる とらえる ぜんぶ ただしい こぐま ちへいせん おまいり あんまり")
	if err != nil {
		t.Fatalf("failed to generate mnemonic, err: %v", err)
	}
	t.Log(seed)
}

func TestBTC_NewMasterKey(t *testing.T) {
	seed := bip39.NewSeed("announce parent popular hybrid fine maid exile impulse unknown school castle wage hand impulse wing", "")
	t.Log(hex.EncodeToString(seed))
	btcWallet := wallets.NewBtcWallet(common.OneOfOne)
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

func TestBTC_NewAccount(t *testing.T) {
	seed := bip39.NewSeed("announce parent popular hybrid fine maid exile impulse unknown school castle wage hand impulse wing", "")
	t.Log(hex.EncodeToString(seed))
	btcWallet := wallets.NewBtcWallet(common.OneOfOne)
	masterKeyEx := "xprv9s21ZrQH143K26uUQxSgabsXyzC3VbzBsFs3dZuQLTbYgewoJbnvpUreykr8hThdNL238QdVTdhg7ns5jQW8veYDjuZRVqNppWEwc2wSuqB"
	accountKeyEx := "xprv9z9UCSxob6vZL7WkwjCPFJVj27R2gk11r4QD18wKH7LqrhiBXNg1ca3h7iyzULTCJa15JdyzPbArNCh9auTn8WDWEq2Fkf44sajGQQ2wYSs"
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

func TestBTC_NextAddressNormal(t *testing.T) {
	seed := bip39.NewSeed("announce parent popular hybrid fine maid exile impulse unknown school castle wage hand impulse wing", "")
	t.Log(hex.EncodeToString(seed))
	btcWallet := wallets.NewBtcWallet(common.OneOfOne)
	masterKeyEx := "xprv9s21ZrQH143K26uUQxSgabsXyzC3VbzBsFs3dZuQLTbYgewoJbnvpUreykr8hThdNL238QdVTdhg7ns5jQW8veYDjuZRVqNppWEwc2wSuqB"
	addressEx := "1Jgub4PAFFXabXemD1VZ143BSCBf1ax2S7"
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
