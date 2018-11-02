package coins

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/phanhoc/wallet-bip44/coins/models"
)

type Btc struct {
	AbstractUtxoCoins
	NetWork *chaincfg.Params
}

func NewBtc() *Btc {
	return &Btc{NetWork: &chaincfg.MainNetParams}
}

func (*Btc) GetChain() string {
	return "btc"
}

func (*Btc) GetFamily() string {
	return "btc"
}

func (*Btc) GetFullName() string {
	return "Bitcoin"
}

func (b *Btc) GenerateNormalAddress(xPub string, index uint32, internal bool) (*models.AddressInfo, error) {
	pubKeyHash, branchNum, err := b.utxoDeriveEcPubKey(xPub, index, internal)
	if err != nil {
		return nil, fmt.Errorf("failed to drive ecrypt public key, err %v", err)
	}

	address, err := btcutil.NewAddressPubKeyHash(pubKeyHash, b.NetWork)
	if err != nil {
		return nil, err
	}
	fmt.Println("-----------", address.EncodeAddress())

	return &models.AddressInfo{
		Address:    address.String(),
		Type:       models.NORMAL_ADDRESS,
		AddressRaw: address.ScriptAddress(),
		Paths:      []string{fmt.Sprintf("%d/%d", branchNum, index)},
	}, nil
}

func (b *Btc) GenerateMultisigAddress(xPub []string, index, flagM, flagN uint32, internal bool) (*models.AddressInfo, error) {
	paths := make([]string, 0, len(xPub))
	addresses := make([]*btcutil.AddressPubKey, 0, len(xPub))
	for _, item := range xPub {
		ecPubKey, branchNum, err := b.utxoDeriveEcPubKey(item, index, internal)
		if err != nil {
			return nil, err
		}
		path := fmt.Sprintf("%d/%d", branchNum, index)
		addressInfo, err := btcutil.NewAddressPubKey(ecPubKey, b.NetWork)
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

	multisigAddr, err := btcutil.NewAddressScriptHash(redeemScript, b.NetWork)
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

func (b *Btc) NewAccount(master string, accountIndex uint32) (*models.AccountInfo, error) {
	return b.newAccount(master, b.NetWork.HDCoinType, accountIndex)
}

func (b *Btc) SignTx(serializeTx []byte, privKey string, optional SigningOptional) ([]byte, error) {
	// var rawTx wire.MsgTx
	// if err := rawTx.Deserialize(bytes.NewBuffer(serializeTx)); err != nil {
	// 	return nil, fmt.Errorf("failed to deserialize transaction, err %v", err)
	// }
	// getKey := mkGetKey(map[string]addressToKey{
	// 	address: {key: ecPrvKey, compressed: true},
	// })
	//
	// getScript := mkGetScript(nil)
	// if class == txscript.ScriptHashTy {
	// 	redeemScript, err := hex.DecodeString(prevInfo.RedeemScript)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to decode redeem script %v", err)
	// 	}
	// 	getScript = mkGetScript(map[string][]byte{addresses[0].String(): redeemScript})
	// }
	// script, err := b.signAndCheck(chainParams, &rawTx, i, prevInfo.Amount, prevInfo.ScriptPubKey, getKey, getScript, txIn.SignatureScript)
	// if err != nil {
	// 	return nil, err
	// }
	return nil, nil
}

//
// func (b *Btc) deriveEcPubKey(acctKey *hdkeychain.ExtendedKey, index uint32, internal bool) (*btcec.PublicKey, uint32, error) {
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
