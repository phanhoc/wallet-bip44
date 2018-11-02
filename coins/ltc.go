package coins

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/phanhoc/wallet-bip44/coins/models"
	"github.com/pkg/errors"
)

type Ltc struct {
	Btc
}

func NewLtc() *Ltc {
	ltcChain := &chaincfg.MainNetParams
	ltcChain.HDCoinType = 2
	ltcChain.Bech32HRPSegwit = "ltc"
	ltcChain.PubKeyHashAddrID = 0x30        // starts with L
	ltcChain.ScriptHashAddrID = 0x32        // starts with M
	ltcChain.PrivateKeyID = 0xB0            // starts with 6 (uncompressed) or T (compressed)
	ltcChain.WitnessPubKeyHashAddrID = 0x06 // starts with p2
	ltcChain.WitnessScriptHashAddrID = 0x0A // starts with 7Xh

	// BIP32 hierarchical deterministic extended key magics
	ltcChain.HDPrivateKeyID = [4]byte{0x04, 0x88, 0xad, 0xe4} // starts with xprv
	ltcChain.HDPublicKeyID = [4]byte{0x04, 0x88, 0xb2, 0x1e}  // starts with xpub
	return &Ltc{
		Btc{NetWork: ltcChain},
	}
}

func (*Ltc) GetChain() string {
	return "ltc"
}

func (*Ltc) GetFamily() string {
	return "ltc"
}

func (*Ltc) GetFullName() string {
	return "Lite coin"
}

func (*Ltc) GenerateMultisigAddress(xPub []string, index, flagM, flagN uint32, internal bool) (*models.AddressInfo, error) {
	return nil, errors.New("currently not supported multisig for lite coin")
}
