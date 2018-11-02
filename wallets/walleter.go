package wallets

import (
	"github.com/phanhoc/wallet-bip44/coins/models"
	"github.com/phanhoc/wallet-bip44/wallets/common"
)

type AddressOptional struct {
	AddressIndex uint32
	Internal     bool
	FlagM        uint32
	FlagN        uint32
}

type Waller interface {
	NewMnemonic(language common.MnemonicLanguage) (string, error)
	NewSeed(language common.MnemonicLanguage, mnemonic string) ([]byte, error)
	NewMaster(seed []byte) (string, error)
	NewAccount(masterKey string, accountIndex uint32) (*models.AccountInfo, error)
	NextAddress(xPub []string, optional AddressOptional) (*models.AddressInfo, error)
	Export(data []byte, filename, passphrase string) error
	Encrypt(data []byte, passphrase string) ([]byte, error)
	Decrypt(data []byte, passphrase string) ([]byte, error)
}
