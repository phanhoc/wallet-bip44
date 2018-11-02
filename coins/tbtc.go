package coins

import "github.com/btcsuite/btcd/chaincfg"

type TBtc struct {
	Btc
}

func NewTBtc() *TBtc {
	return &TBtc{Btc{NetWork: &chaincfg.TestNet3Params}}
}