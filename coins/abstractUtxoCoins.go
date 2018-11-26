package coins

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
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

func (a *AbstractUtxoCoins) signTx(serializeTx []byte, privKey string, optional SigningOptional, net *chaincfg.Params) ([]byte, error) {
	var rawTx wire.MsgTx
	if err := rawTx.Deserialize(bytes.NewBuffer(serializeTx)); err != nil {
		return nil, fmt.Errorf("failed to deserialize transaction, err %v", err)
	}
	privWIF, err := btcutil.DecodeWIF(privKey)
	if err != nil {
		return nil, err
	}
	for i, tx := range rawTx.TxIn {
		class, addresses, _, err := txscript.ExtractPkScriptAddrs(optional.PreviousScript, net)
		if err != nil {
			return nil, err
		}
		getKey := mkGetKey(map[string]addressToKey{
			addresses[0].String(): {key: privWIF.PrivKey, compressed: true},
		})
		getScript := mkGetScript(nil)
		if class == txscript.ScriptHashTy {
			getScript = mkGetScript(map[string][]byte{addresses[0].String(): optional.RedeemScript})
		}
		script, err := a.signAndCheck(net, &rawTx, i, rawTx.TxOut[i].Value, optional.PreviousScript, getKey, getScript, tx.SignatureScript)
		if err != nil {
			return nil, err
		}
		rawTx.TxIn[i].SignatureScript = script
	}

	buf := bytes.NewBuffer(make([]byte, 0, rawTx.SerializeSize()))
	if err := rawTx.Serialize(buf); err != nil {
		return nil, fmt.Errorf("failed to serialize transaction: %v", err)
	}
	return buf.Bytes(), nil
}
func (*AbstractUtxoCoins) signAndCheck(params *chaincfg.Params, tx *wire.MsgTx, index int, amount int64, pkScript []byte, keyDB txscript.KeyDB,
	scriptDB txscript.ScriptDB, prevScript []byte) ([]byte, error) {
	script, err := txscript.SignTxOutput(params, tx, index,
		pkScript, txscript.SigHashAll, keyDB, scriptDB, prevScript)
	if err != nil {
		return nil, fmt.Errorf("failed to signing transaction: %v", err)
	}
	// if len(b.PrevOutput[index].State.SignedBy) == int(b.PrevOutput[index].State.NRequire-1) {
	// 	if err := common.CheckScripts("check transaction script", tx, index, amount, script, pkScript); err != nil {
	// 		return nil, fmt.Errorf("failed to check signed transaction: %v", err)
	// 	}
	// }

	return script, nil
}
