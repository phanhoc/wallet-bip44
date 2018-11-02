package wallets

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/phanhoc/wallet-bip44/coins"
	common2 "github.com/phanhoc/wallet-bip44/wallets/common"
)

type EthWallet struct {
	BaseWallet
	Network *chaincfg.Params
}

func NewEthWallet(walletType common2.WalletType) Waller {
	chainParam := &chaincfg.MainNetParams
	chainParam.HDCoinType = 60
	return &EthWallet{
		BaseWallet{
			Coin:       coins.NewEth(),
			WalletType: walletType,
		},
		chainParam,
	}
}

func (e *EthWallet) NewMaster(seed []byte) (string, error) {
	master, err := hdkeychain.NewMaster(seed, e.Network)
	if err != nil {
		return "", fmt.Errorf("failed to new master key, err %v", err)
	}
	return master.String(), nil
}
