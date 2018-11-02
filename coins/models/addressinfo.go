package models

type AddressType string

const (
	NORMAL_ADDRESS   AddressType = "normal"
	MULTISIG_ADDRESS AddressType = "multisig"
)

type AddressVersion uint32

const (
	LEGACY_ADDRESS AddressVersion = 1
	CASH_ADDRESS   AddressVersion = 2
)

type AddressInfo struct {
	Version      string
	Prefix       string
	Address      string
	AddressRaw   []byte
	Paths        []string
	Type         AddressType
	RedeemScript []byte
}
