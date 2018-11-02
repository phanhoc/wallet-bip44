package wallets

import (
	"github.com/phanhoc/wallet-bip44/coins"
	"github.com/phanhoc/wallet-bip44/coins/models"
	"github.com/phanhoc/wallet-bip44/wallets/common"
	"github.com/phanhoc/wallet-bip44/wallets/mnemonic"
)

type BaseWallet struct {
	Coin       coins.Coiner
	WalletType common.WalletType
}

func (*BaseWallet) NewMnemonic(language common.MnemonicLanguage) (string, error) {
	m := mnemonic.NewMnemonicWithLanguage(language)
	return m.GenerateMnemonic()
}

func (*BaseWallet) NewSeed(language common.MnemonicLanguage, mnem string) ([]byte, error) {
	m := mnemonic.NewMnemonicWithLanguage(language)
	return m.GenerateSeed(mnem)
}

func (*BaseWallet) Export(data []byte, filename, passphrase string) error {
	encrypt := common.NewEncrypter(filename, &passphrase)
	return encrypt.WriteEncryptedFile(data)
}

func (*BaseWallet) Encrypt(data []byte, passphrase string) ([]byte, error) {
	encrypt := common.NewEncrypterWriteBuffer(&passphrase)
	return encrypt.WriteToBuffer(data)
}

func (*BaseWallet) Decrypt(data []byte, passphrase string) ([]byte, error) {
	encrypt := common.NewEncrypterWriteBuffer(&passphrase)
	return encrypt.ReadEncrytedFromBuffer(data)
}

func (bw *BaseWallet) NewAccount(masterKey string, accountIndex uint32) (*models.AccountInfo, error) {
	return bw.Coin.NewAccount(masterKey, accountIndex)
}
func (bw *BaseWallet) NextAddress(xPub []string, optional AddressOptional) (*models.AddressInfo, error) {
	if bw.WalletType == common.OneOfOne {
		return bw.Coin.GenerateNormalAddress(xPub[0], optional.AddressIndex, optional.Internal)
	}
	return bw.Coin.GenerateMultisigAddress(xPub, optional.AddressIndex, optional.FlagM, optional.FlagN, optional.Internal)
}
