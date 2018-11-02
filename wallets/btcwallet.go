package wallets

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/phanhoc/wallet-bip44/coins"
	"github.com/phanhoc/wallet-bip44/wallets/common"
)

type BtcWallet struct {
	BaseWallet
	Network *chaincfg.Params
}

func NewBtcWallet(walletType common.WalletType) Waller {
	return &BtcWallet{
		BaseWallet{
			Coin:       coins.NewBtc(),
			WalletType: walletType,
		},
		&chaincfg.MainNetParams,
	}
}

func (b *BtcWallet) NewMaster(seed []byte) (string, error) {
	master, err := hdkeychain.NewMaster(seed, b.Network)
	if err != nil {
		return "", fmt.Errorf("failed to new master key, err %v", err)
	}
	return master.String(), nil
}

// func (b *BtcWallet) NewAccount(masterKey string, accountIndex uint32) (*models.AccountInfo, error) {
// 	return b.Coin.NewAccount(masterKey, accountIndex)
// }
// func (b *BtcWallet) NextAddress(xPub []string, optional AddressOptional) (*models.AddressInfo, error) {
// 	if b.WalletType == common.OneOfOne {
// 		return b.Coin.GenerateNormalAddress(xPub[0], optional.AddressIndex, optional.Internal)
// 	}
// 	return b.Coin.GenerateMultisigAddress(xPub, optional.AddressIndex, optional.FlagM, optional.FlagN, optional.Internal)
// }
