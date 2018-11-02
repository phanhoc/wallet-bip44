package coins

import (
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/crypto"
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

func (e *Eth) getBaseFactor() *big.Float {
	return big.NewFloat(float64(10e18))
}
