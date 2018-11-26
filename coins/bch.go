package coins

import (
	"fmt"
	"github.com/gcash/bchd/chaincfg"
	"github.com/gcash/bchd/txscript"
	"github.com/gcash/bchutil"
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
	pubKeyHash, branchNum, err := b.utxoDeriveEcPubKey(xPub, index, internal)
	if err != nil {
		return nil, fmt.Errorf("failed to drive ecrypt public key, err %v", err)
	}

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
		ecPubKey, branchNum, err := b.deriveEcPubKey(item, index, internal)
		if err != nil {
			return nil, err
		}
		path := fmt.Sprintf("%d/%d", branchNum, index)
		addressInfo, err := bchutil.NewAddressPubKey(ecPubKey.SerializeCompressed(), b.NetWork)
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
}

func (b *Bch) SignTx(serializeTx []byte, privKey string, optional SigningOptional) ([]byte, error) {
	// return b.signTx(serializeTx, privKey, optional, b.NetWork)
	return nil, nil
}
