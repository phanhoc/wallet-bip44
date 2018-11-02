package models

type AccountInfo struct {
	CoinType          uint32
	AccountIndex      uint32
	AccountPrivateKey string
	AccountPublicKey  string
	Path              string
}
