package wallets

import (
	"fmt"
	"github.com/gcash/bchd/chaincfg"
	"github.com/gcash/bchutil/hdkeychain"
	"github.com/phanhoc/wallet-bip44/coins"
	"github.com/phanhoc/wallet-bip44/wallets/common"
)

type BchWallet struct {
	BaseWallet
	Network *chaincfg.Params
}

func NewBchWallet(walletType common.WalletType) Waller {
	return &BchWallet{
		BaseWallet{
			Coin:       coins.NewBch(),
			WalletType: walletType,
		},
		&chaincfg.MainNetParams,
	}
}

func (b *BchWallet) NewMaster(seed []byte) (string, error) {
	master, err := hdkeychain.NewMaster(seed, b.Network)
	if err != nil {
		return "", fmt.Errorf("failed to new master key, err %v", err)
	}
	return master.String(), nil
}
