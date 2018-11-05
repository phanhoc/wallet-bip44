package coins

import (
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/phanhoc/wallet-bip44/coins/models"
	"math/big"
)

type Eth struct {
	BaseCoin
	NetWork *chaincfg.Params
}

func NewEth() *Eth {
	ethChain := &chaincfg.MainNetParams
	ethChain.HDCoinType = 60
	return &Eth{NetWork: ethChain}
}

func (*Eth) GetChain() string {
	return "eth"
}

func (*Eth) GetFamily() string {
	return "eth"
}

func (*Eth) GetFullName() string {
	return "Ethereum"
}

func (e *Eth) GenerateNormalAddress(xPub string, index uint32, internal bool) (*models.AddressInfo, error) {
	ecPubKey, branchNum, err := e.deriveEcPubKey(xPub, index, internal)
	if err != nil {
		return nil, err
	}
	address := crypto.PubkeyToAddress(*ecPubKey.ToECDSA())
	return &models.AddressInfo{
		Address: address.Hex(),
		Type:    models.NORMAL_ADDRESS,
		// AddressRaw: address.(),
		Paths: []string{fmt.Sprintf("%d/%d", branchNum, index)},
	}, nil
}

func (e *Eth) GenerateMultisigAddress(xPub []string, index, flagM, flagN uint32, internal bool) (*models.AddressInfo, error) {
	return nil, errors.New("unsupported generate multisig address in coin")
}

func (e *Eth) NewAccount(master string, accountIndex uint32) (*models.AccountInfo, error) {
	return e.newAccount(master, e.NetWork.HDCoinType, accountIndex)
}

func (e *Eth) BaseUnitsToBigUnits(baseUnits *big.Int) *big.Float {
	return e.baseUnitsToBigUnits(baseUnits, e.getBaseFactor())
}

func (e *Eth) BigUnitsToBaseUnits(bigUnits *big.Float) (*big.Int, error) {
	return e.bigUnitsToBaseUnits(bigUnits, e.getBaseFactor())
}

func (e *Eth) SignTx(serializedTx []byte, privKey string, optional SigningOptional) ([]byte, error) {
	rawTx := new(types.Transaction)
	if err := rlp.DecodeBytes(serializedTx, rawTx); err != nil {
		return nil, fmt.Errorf("failed to decode raw transaction %v", err)
	}
	esPrvKey, err := crypto.ToECDSA([]byte(privKey))
	if err != nil {
		return nil, fmt.Errorf("failed to ecdsa key, err %v", err)
	}
	signTx, err := types.SignTx(rawTx, e.getSigner(), esPrvKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %v", err)
	}
	return rlp.EncodeToBytes(signTx)
}

func (e *Eth) getBaseFactor() *big.Float {
	return big.NewFloat(float64(10e18))
}

// GetSignerFromNetwork returns singer base on chain params
func (e *Eth) getSigner() types.Signer {
	return types.NewEIP155Signer(big.NewInt(1))
}
