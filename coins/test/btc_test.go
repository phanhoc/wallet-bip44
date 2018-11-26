package test

import (
	"bytes"
	"encoding/hex"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
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
	if err != nil {
		t.Fatalf("failed to generate multisig, err: %v", err)
	}
	t.Log(address.Address)
	t.Log(hex.EncodeToString(address.RedeemScript))
	if expectedAddress != address.Address {
		t.Fatalf("failed to generate address, got %s - expected %s", address.Address, expectedAddress)
	}
}

func TestExtract(t *testing.T) {
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

	txHex := "0100000000010178299dc3409b7d9e7a11d648ed1bac525ef0724065580d2730c5c9a2d0e425e60200000000ffffffff02e981cd2f02000000160014d34468064a13c08c37582123223876d74d796495da5f3e240000000017a9149f5d9e0a0fc7eff817667a0d95773eb45ddf0710870247304402204b4dec3fae70eaa9968ff304c50c4941648c81152a7063bd62520d5e5e84cb89022014e373fc1d1e8bd7e6f947bb0073995f148ba0438b96f1cde68e19ba53b43d1c012103177fe9afb07a931846485623a5994baaa8078d075cbb7245ac716be6245f24d200000000"
	txByte, err := hex.DecodeString(txHex)
	if err != nil {
		t.Fatalf("failed to decode hex transaction, err %v", err)
	}
	var rawTx wire.MsgTx
	if err := rawTx.Deserialize(bytes.NewBuffer(txByte)); err != nil {
		t.Fatalf("failed to decode raw transaction, err %v", err)
	}
	class, addresses, n, err := txscript.ExtractPkScriptAddrs(rawTx.TxOut[0].PkScript, ltcChain)
	if err != nil {
		t.Fatalf("failed to decode hex transaction, err %v", err)
	}
	t.Log(n)
	t.Log(class)
	t.Log(addresses)
}
