package coins

import (
	"github.com/phanhoc/wallet-bip44/coins/models"
	"math/big"
)

type SigningOptional struct {
	PreviousScript []byte
	RedeemScript   []byte
}

type Coiner interface {
	NewAccount(master string, accountIndex uint32) (*models.AccountInfo, error)
	GenerateNormalAddress(xPub string, index uint32, internal bool) (*models.AddressInfo, error)
	GenerateMultisigAddress(xPub []string, index, flagM, flagN uint32, internal bool) (*models.AddressInfo, error)
	// Name of the chain which supports this coin (eg, 'btc', 'eth')
	GetChain() string
	// Name of the coin family (eg. for tbtc, this would be btc)
	GetFamily() string
	// Human readable full name for the coin
	GetFullName() string
	BaseUnitsToBigUnits(baseUnits *big.Int) *big.Float
	BigUnitsToBaseUnits(bigUnits *big.Float) (*big.Int, error)
	// privKey format WIF for BTC, BCH
	SignTx(rawTx []byte, privKey string, optional SigningOptional) ([]byte, error)
}
