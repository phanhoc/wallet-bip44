package mnemonic

import (
	"github.com/phanhoc/wallet-bip44/wallets/common"
	"github.com/pkg/errors"
	"github.com/tyler-smith/go-bip39"
	"github.com/tyler-smith/go-bip39/wordlists"
)

// Please refer the link: https://iancoleman.io/bip39/ for purpose double check result

type Mnemonic struct {
	EntropySize int
	Password    string
}

func NewMnemonicWithDefaultOption() *Mnemonic {
	return &Mnemonic{EntropySize: common.DefaultEntropySize, Password: common.DefaultSeedPass}
}

func NewMnemonicWithLanguage(language common.MnemonicLanguage) *Mnemonic {
	bip39.SetWordList(loadWordList(language))
	return &Mnemonic{EntropySize: common.DefaultEntropySize, Password: common.DefaultSeedPass}
}

// New mnemonic follow the wordlists
func (m *Mnemonic) GenerateMnemonic() (string, error) {
	entropy, err := bip39.NewEntropy(m.EntropySize)
	if err != nil {
		return "", err
	}

	return bip39.NewMnemonic(entropy)
}

// Generate seed from mnemonic and pass( optional )
func (m *Mnemonic) GenerateSeed(mnemonic string) ([]byte, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, errors.New("invalidate mnemonic")
	}
	return bip39.NewSeed(mnemonic, m.Password), nil
}

// Get word list
func (m *Mnemonic) ListWord() []string {
	return bip39.GetWordList()
}

// loadWordList returns word lists base on language setting in the configuration
func loadWordList(language common.MnemonicLanguage) []string {
	switch language {
	case common.JAPANESE:
		return wordlists.Japanese
	case common.ITALIAN:
		return wordlists.Italian
	case common.KOREAN:
		return wordlists.Korean
	case common.SPANISH:
		return wordlists.Spanish
	default:
		return wordlists.English
	}
}
