package coins

import (
	"errors"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/golangcrypto/ripemd160"
	"math/big"
)

const (
	PurposePath = "m/44'/"

	ExternalBranch uint32 = 0

	// internalBranch is the child number to use when performing BIP0044
	// style hierarchical deterministic key derivation for the internal
	// branch.
	InternalBranch uint32 = 1
)

type AbstractUtxoCoins struct {
	BaseCoin
}

// Convert a currency amount represented in base units (satoshi, wei, atoms, drops, stroops)
// to big units (btc, eth, rmg, xrp, xlm)
func (a *AbstractUtxoCoins) BaseUnitsToBigUnits(baseUnits *big.Int) *big.Float {
	return a.baseUnitsToBigUnits(baseUnits, a.getBaseFactor())
}

// Convert a currency amount represented in big units (btc, eth, rmg, xrp, xlm)
// to base units (satoshi, wei, atoms, drops, stroops)
func (a *AbstractUtxoCoins) BigUnitsToBaseUnits(bigUnits *big.Float) (*big.Int, error) {
	return a.bigUnitsToBaseUnits(bigUnits, a.getBaseFactor())
}

func (a *AbstractUtxoCoins) getBaseFactor() *big.Float {
	return big.NewFloat(float64(10e8))
}

func (a *AbstractUtxoCoins) utxoDeriveEcPubKey(acctKey string, index uint32, internal bool) ([]byte, uint32, error) {
	ecpubKey, branch, err := a.deriveEcPubKey(acctKey, index, internal)
	if err != nil {
		return nil, 0, err
	}
	pubKeyHash := btcutil.Hash160(ecpubKey.SerializeCompressed())
	if len(pubKeyHash) != ripemd160.Size {
		return nil, 0, errors.New("pubKeyHash must be 20 bytes")
	}
	return pubKeyHash, branch, err
}
