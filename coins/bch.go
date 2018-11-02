package coins

import (
	"fmt"
	"github.com/gcash/bchd/chaincfg"
	"github.com/gcash/bchd/txscript"
	"github.com/gcash/bchutil"
	"github.com/gcash/bchutil/hdkeychain"
	"github.com/phanhoc/wallet-bip44/coins/models"
)

type Bch struct {
	AbstractUtxoCoins
	NetWork *chaincfg.Params
}

func NewBch() *Bch {
	return &Bch{NetWork: &chaincfg.MainNetParams}
}

func (*Bch) GetChain() string {
	return "bch"
}

func (*Bch) GetFamily() string {
	return "bch"
}

func (*Bch) GetFullName() string {
	return "Bitcoin Cash"
}

func (b *Bch) GenerateNormalAddress(xPub string, index uint32, internal bool) (*models.AddressInfo, error) {
	// hdxPubKey, err := hdkeychain.NewKeyFromString(xPub)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to new public key from string, err %v", err)
	// }
	pubKeyHash, branchNum, err := b.utxoDeriveEcPubKey(xPub, index, internal)
	if err != nil {
		return nil, fmt.Errorf("failed to drive ecrypt public key, err %v", err)
	}
	// pubKeyHash := bchutil.Hash160(ecPubKey.SerializeCompressed())
	// if len(pubKeyHash) != ripemd160.Size {
	// 	return nil, errors.New("pubKeyHash must be 20 bytes")
	// }

	address, err := bchutil.NewAddressPubKeyHash(pubKeyHash, b.NetWork)
	if err != nil {
		return nil, err
	}
	fmt.Println("-----------", address.EncodeAddress())

	return &models.AddressInfo{
		Address:    address.String(),
		Prefix:     "bitcoincash",
		Type:       models.NORMAL_ADDRESS,
		AddressRaw: address.ScriptAddress(),
		Paths:      []string{fmt.Sprintf("%d/%d", branchNum, index)},
	}, nil
}

func (b *Bch) GenerateMultisigAddress(xPub []string, index, flagM, flagN uint32, internal bool) (*models.AddressInfo, error) {
	paths := make([]string, 0, len(xPub))
	addresses := make([]*bchutil.AddressPubKey, 0, len(xPub))
	for _, item := range xPub {
		// hdxPubKey, err := hdkeychain.NewKeyFromString(item)
		// if err != nil {
		// 	return nil, fmt.Errorf("failed to new public key from string, err %v", err)
		// }
		ecPubKey, branchNum, err := b.utxoDeriveEcPubKey(item, index, internal)
		if err != nil {
			return nil, err
		}
		path := fmt.Sprintf("%d/%d", branchNum, index)
		addressInfo, err := bchutil.NewAddressPubKey(ecPubKey, b.NetWork)
		if err != nil {
			return nil, err
		}
		paths = append(paths, path)
		addresses = append(addresses, addressInfo)
	}
	redeemScript, err := txscript.MultiSigScript(addresses, int(flagM))
	if err != nil {
		return nil, err
	}

	multisigAddr, err := bchutil.NewAddressScriptHash(redeemScript, b.NetWork)
	if err != nil {
		return nil, err
	}

	return &models.AddressInfo{
		Address:      multisigAddr.String(),
		Type:         models.MULTISIG_ADDRESS,
		Paths:        paths,
		RedeemScript: redeemScript,
	}, nil
}

func (b *Bch) NewAccount(master string, accountIndex uint32) (*models.AccountInfo, error) {
	return b.newAccount(master, b.NetWork.HDCoinType, accountIndex)
	// masterKey, err := hdkeychain.NewKeyFromString(master)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to new master key from string, err %v", err)
	// }
	// purpose, err := masterKey.Child(44 + hdkeychain.HardenedKeyStart)
	// if err != nil {
	// 	return nil, err
	// }
	// // Derive the coin type key as a child of the purpose key.
	// coinTypeKey, err := purpose.Child(b.NetWork.HDCoinType + hdkeychain.HardenedKeyStart)
	// if err != nil {
	// 	return nil, err
	// }
	// defer coinTypeKey.Zero()
	//
	// acctKeyPriv, err := coinTypeKey.Child(accountIndex + hdkeychain.HardenedKeyStart)
	// if err != nil {
	// 	return nil, err
	// }
	// // Ensure the branch keys can be derived for the provided seed according
	// // to BIP0044.
	// if err := bchCheckBranchKeys(acctKeyPriv); err != nil {
	// 	// The seed is unusable if the any of the children in the
	// 	// required hierarchy can't be derived due to invalid child.
	// 	if err == hdkeychain.ErrInvalidChild {
	// 		return nil, errors.New("the provided seed is unusable")
	// 	}
	//
	// 	return nil, err
	// }
	// accountPath := fmt.Sprintf("%s%d'/%d'", PurposePath, b.NetWork.HDCoinType, accountIndex)
	// accountInfo := &models.AccountInfo{
	// 	CoinType:          b.NetWork.HDCoinType,
	// 	AccountIndex:      accountIndex,
	// 	AccountPrivateKey: acctKeyPriv.String(),
	// 	Path:              accountPath,
	// }
	// return accountInfo, nil
}

func bchCheckBranchKeys(acctKey *hdkeychain.ExtendedKey) error {
	// Derive the external branch as the first child of the account key.
	if _, err := acctKey.Child(ExternalBranch); err != nil {
		return err
	}

	// Derive the external branch as the second child of the account key.
	_, err := acctKey.Child(InternalBranch)
	return err
}

// func (b *Bch) deriveEcPubKey(acctKey *hdkeychain.ExtendedKey, index uint32, internal bool) (*bchec.PublicKey, uint32, error) {
// 	branch := ExternalBranch
// 	if internal {
// 		branch = InternalBranch
// 	}
// 	acctPubKey := acctKey
// 	var err error
// 	if acctKey.IsPrivate() {
// 		acctPubKey, err = acctKey.Neuter()
// 		if err != nil {
// 			return nil, 0, err
// 		}
// 	}
//
// 	// Derive the appropriate branch key and ensure it is zeroed when done.
// 	branchKey, err := acctPubKey.Child(branch)
// 	if err != nil {
// 		return nil, 0, fmt.Errorf("failed to derive extended key branch %d", branch)
// 	}
// 	defer branchKey.Zero() // Ensure branch key is zeroed when done.
//
// 	privKey, err := branchKey.Child(index)
// 	if err != nil {
// 		return nil, 0, fmt.Errorf("failed to generate child %d", index)
// 	}
// 	defer privKey.Zero()
// 	pubKey, err := privKey.Neuter()
// 	if err != nil {
// 		return nil, 0, fmt.Errorf("failed to derive pubKey")
// 	}
// 	pubKey.SetNet(b.NetWork)
// 	ecpubKey, err := pubKey.ECPubKey()
//
// 	return ecpubKey, branch, err
// }
