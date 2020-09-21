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

type address string

func (addr address) Validate() bool {
	panic(UnimplementedError)
}
