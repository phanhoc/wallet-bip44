package common

type WalletStatus int64

const (
	// Default entropy size for mnemonic
	DefaultEntropySize = 256
	// Default seed pass. it used to generate seed from mnemonic( BIP39 ). Don't change if determined
	DefaultSeedPass = ""
)

type MnemonicLanguage string

// List Mnemonic language support
const (
	ENGLISH  MnemonicLanguage = "EN"
	JAPANESE                  = "JP"
	FRENCH                    = "FR"
	ITALIAN                   = "IT"
	KOREAN                    = "KR"
	SPANISH                   = "ES"
)

const (
	WalletStatusNew WalletStatus = iota + 1
	WalletStatusCreating
	WalletStatusError
	WalletStatusReady
	WalletStatusDisable
)

// Support wallet types as 1 of 1, 2 of 3 and 3 of 5
type WalletType uint32

const (
	OneOfOne WalletType = iota + 1
	TwoOfThree
	ThreeOfFive
)
