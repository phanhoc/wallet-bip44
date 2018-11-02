package common

type Coin string

const (
	BTC Coin = "btc"
	BCH Coin = "bch"
	ETH Coin = "eth"
)

const (
	MAIN_NET = "mainnet"
	TEST_NET = "testnet"
)

var True = Bool{true}
var False = Bool{false}

type Bool struct {
	value bool
}

func (b Bool) Pointer() *bool {
	return &b.value
}
