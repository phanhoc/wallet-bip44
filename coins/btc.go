package coins

import (
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/phanhoc/wallet-bip44/coins/models"
)

type addressToKey struct {
	key        *btcec.PrivateKey
	compressed bool
}

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
		ecPubKey, branchNum, err := b.deriveEcPubKey(item, index, internal)
		if err != nil {
			return nil, err
		}
		path := fmt.Sprintf("%d/%d", branchNum, index)
		addressInfo, err := btcutil.NewAddressPubKey(ecPubKey.SerializeCompressed(), b.NetWork)
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
	return b.signTx(serializeTx, privKey, optional, b.NetWork)
	// var rawTx wire.MsgTx
	// if err := rawTx.Deserialize(bytes.NewBuffer(serializeTx)); err != nil {
	// 	return nil, fmt.Errorf("failed to deserialize transaction, err %v", err)
	// }
	// privWIF, err := btcutil.DecodeWIF(privKey)
	// if err != nil {
	// 	return nil, err
	// }
	// for i, tx := range rawTx.TxIn {
	// 	class, addresses, _, err := txscript.ExtractPkScriptAddrs(optional.PreviousScript, b.NetWork)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	getKey := mkGetKey(map[string]addressToKey{
	// 		addresses[0].String(): {key: privWIF.PrivKey, compressed: true},
	// 	})
	// 	getScript := mkGetScript(nil)
	// 	if class == txscript.ScriptHashTy {
	// 		getScript = mkGetScript(map[string][]byte{addresses[0].String(): optional.RedeemScript})
	// 	}
	// 	script, err := b.signAndCheck(b.NetWork, &rawTx, i, rawTx.TxOut[i].Value, optional.PreviousScript, getKey, getScript, tx.SignatureScript)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	rawTx.TxIn[i].SignatureScript = script
	// }
	//
	// buf := bytes.NewBuffer(make([]byte, 0, rawTx.SerializeSize()))
	// if err := rawTx.Serialize(buf); err != nil {
	// 	return nil, fmt.Errorf("failed to serialize transaction: %v", err)
	// }
	// return buf.Bytes(), nil
}

// func (b *Btc) signAndCheck(params *chaincfg.Params, tx *wire.MsgTx, index int, amount int64, pkScript []byte, keyDB txscript.KeyDB,
// 	scriptDB txscript.ScriptDB, prevScript []byte) ([]byte, error) {
// 	script, err := txscript.SignTxOutput(params, tx, index,
// 		pkScript, txscript.SigHashAll, keyDB, scriptDB, prevScript)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to signing transaction: %v", err)
// 	}
// 	// if len(b.PrevOutput[index].State.SignedBy) == int(b.PrevOutput[index].State.NRequire-1) {
// 	// 	if err := common.CheckScripts("check transaction script", tx, index, amount, script, pkScript); err != nil {
// 	// 		return nil, fmt.Errorf("failed to check signed transaction: %v", err)
// 	// 	}
// 	// }
//
// 	return script, nil
// }

// mkGetKey return Key for signing btc
func mkGetKey(keys map[string]addressToKey) txscript.KeyDB {
	return txscript.KeyClosure(func(addr btcutil.Address) (*btcec.PrivateKey, bool, error) {
		if keys != nil {
			if a2k, ok := keys[addr.EncodeAddress()]; ok {
				return a2k.key, a2k.compressed, nil
			}
		}

		return nil, false, errors.New("fail to get key")
	})
}

// mkGetScript return script for signing btc in case pay to script
func mkGetScript(scripts map[string][]byte) txscript.ScriptDB {
	return txscript.ScriptClosure(func(addr btcutil.Address) ([]byte, error) {
		if scripts != nil {
			if script, ok := scripts[addr.EncodeAddress()]; ok {
				return script, nil
			}
		}

		return nil, errors.New("fail to get script")
	})
}
