package coins

import (
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/phanhoc/wallet-bip44/coins/models"
	"github.com/rubblelabs/ripple/crypto"
	"golang.org/x/crypto/ripemd160"
	"math/big"
)

type Xrp struct {
	BaseCoin
	NetWork *chaincfg.Params
}

func NewXrp() *Xrp {
	xrpChain := &chaincfg.MainNetParams
	xrpChain.HDCoinType = 144
	return &Xrp{NetWork: xrpChain}
}

func (*Xrp) GetChain() string {
	return "xrp"
}

func (*Xrp) GetFamily() string {
	return "xrp"
}

func (*Xrp) GetFullName() string {
	return "Ripple"
}

func (x *Xrp) GenerateNormalAddress(xPub string, index uint32, internal bool) (*models.AddressInfo, error) {
	ecPubKey, branchNum, err := x.deriveEcPubKey(xPub, index, internal)
	if err != nil {
		return nil, fmt.Errorf("failed to drive ecrypt public key, err %v", err)
	}
	pubKeyHash := crypto.Sha256RipeMD160(ecPubKey.SerializeCompressed())
	if len(pubKeyHash) != ripemd160.Size {
		return nil, errors.New("pubKeyHash must be 20 bytes")
	}
	address, err := crypto.NewAccountId(pubKeyHash)
	if err != nil {
		return nil, fmt.Errorf("failed to generate address, err %v", err)
	}

	return &models.AddressInfo{
		Address:    address.String(),
		Type:       models.NORMAL_ADDRESS,
		AddressRaw: address.Payload(),
		Paths:      []string{fmt.Sprintf("%d/%d", branchNum, index)},
	}, nil
}

func (x *Xrp) GenerateMultisigAddress(xPub []string, index, flagM, flagN uint32, internal bool) (*models.AddressInfo, error) {
	return nil, nil
}

func (x *Xrp) NewAccount(master string, accountIndex uint32) (*models.AccountInfo, error) {
	return x.newAccount(master, x.NetWork.HDCoinType, accountIndex)
}

func (x *Xrp) BaseUnitsToBigUnits(baseUnits *big.Int) *big.Float {
	return x.baseUnitsToBigUnits(baseUnits, x.getBaseFactor())
}

func (x *Xrp) BigUnitsToBaseUnits(bigUnits *big.Float) (*big.Int, error) {
	return x.bigUnitsToBaseUnits(bigUnits, x.getBaseFactor())
}

func (x *Xrp) getBaseFactor() *big.Float {
	return big.NewFloat(float64(1))
}