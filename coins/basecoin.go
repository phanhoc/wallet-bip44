package coins

import (
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/phanhoc/wallet-bip44/coins/models"
	"math/big"
)

type BaseCoin struct {
}

func (*BaseCoin) SignTx(rawTx []byte, privKey string, optional SigningOptional) ([]byte, error) {
	return nil, nil
}

// Convert a currency amount represented in base units (satoshi, wei, atoms, drops, stroops)
// to big units (btc, eth, rmg, xrp, xlm)
func (*BaseCoin) baseUnitsToBigUnits(baseUnits *big.Int, baseFactor *big.Float) *big.Float {
	return new(big.Float).Quo(big.NewFloat(float64(baseUnits.Int64())), baseFactor)
}

// Convert a currency amount represented in big units (btc, eth, rmg, xrp, xlm)
// to base units (satoshi, wei, atoms, drops, stroops)
func (*BaseCoin) bigUnitsToBaseUnits(bigUnits *big.Float, baseFactor *big.Float) (*big.Int, error) {
	res := new(big.Float).Mul(bigUnits, baseFactor)
	resInt64, accuracy := res.Int64()
	if accuracy.String() != "0" {
		return nil, fmt.Errorf("non-integer output resulted from multiplying %s by %s", bigUnits.String(), baseFactor.String())
	}
	return new(big.Int).SetInt64(resInt64), nil
}

func (*BaseCoin) newAccount(master string, coinType, accountIndex uint32) (*models.AccountInfo, error) {
	masterKey, err := hdkeychain.NewKeyFromString(master)
	if err != nil {
		return nil, fmt.Errorf("failed to new master key from string, err %v", err)
	}
	purpose, err := masterKey.Child(44 + hdkeychain.HardenedKeyStart)
	if err != nil {
		return nil, err
	}
	// Derive the coin type key as a child of the purpose key.
	coinTypeKey, err := purpose.Child(coinType + hdkeychain.HardenedKeyStart)
	if err != nil {
		return nil, err
	}
	defer coinTypeKey.Zero()

	acctKeyPriv, err := coinTypeKey.Child(accountIndex + hdkeychain.HardenedKeyStart)
	if err != nil {
		return nil, err
	}
	// Ensure the branch keys can be derived for the provided seed according
	// to BIP0044.
	if err := checkBranchKeys(acctKeyPriv); err != nil {
		// The seed is unusable if the any of the children in the
		// required hierarchy can't be derived due to invalid child.
		if err == hdkeychain.ErrInvalidChild {
			return nil, errors.New("the provided seed is unusable")
		}

		return nil, err
	}
	acctPubKey, err := acctKeyPriv.Neuter()
	if err != nil {
		return nil, err
	}

	accountPath := fmt.Sprintf("%s%d'/%d'", PurposePath, coinType, accountIndex)
	accountInfo := &models.AccountInfo{
		CoinType:          coinType,
		AccountIndex:      accountIndex,
		AccountPrivateKey: acctKeyPriv.String(),
		AccountPublicKey:  acctPubKey.String(),
		Path:              accountPath,
	}
	return accountInfo, nil
}

func checkBranchKeys(acctKey *hdkeychain.ExtendedKey) error {
	// Derive the external branch as the first child of the account key.
	if _, err := acctKey.Child(ExternalBranch); err != nil {
		return err
	}

	// Derive the external branch as the second child of the account key.
	_, err := acctKey.Child(InternalBranch)
	return err
}

func (*BaseCoin) deriveEcPubKey(acctKey string, index uint32, internal bool) (*btcec.PublicKey, uint32, error) {
	accountKey, err := hdkeychain.NewKeyFromString(acctKey)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to new account key from string, err %v", err)
	}
	branch := ExternalBranch
	if internal {
		branch = InternalBranch
	}
	acctPubKey := accountKey
	if accountKey.IsPrivate() {
		acctPubKey, err = accountKey.Neuter()
		if err != nil {
			return nil, 0, err
		}
	}

	// Derive the appropriate branch key and ensure it is zeroed when done.
	branchKey, err := acctPubKey.Child(branch)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to derive extended key branch %d", branch)
	}
	defer branchKey.Zero() // Ensure branch key is zeroed when done.

	privKey, err := branchKey.Child(index)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to generate child %d", index)
	}
	defer privKey.Zero()
	pubKey, err := privKey.Neuter()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to derive pubKey")
	}
	// pubKey.SetNet(net)
	ecpubKey, err := pubKey.ECPubKey()

	return ecpubKey, branch, err
}
