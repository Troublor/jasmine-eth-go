package sdk

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
)

type account struct {
	address    common.Address
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

type Hash string

type Address string

func (addr Address) Validate() bool {
	panic(UnimplementedError)
}

func (addr Address) Address() common.Address {
	return common.HexToAddress(string(addr))
}
