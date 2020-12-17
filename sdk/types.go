package sdk

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"strings"
)

type Backend interface {
	bind.ContractBackend
	ethereum.ChainReader
	ethereum.ChainStateReader
	ethereum.TransactionReader
	NetworkID(ctx context.Context) (*big.Int, error)
}

type Account struct {
	address    common.Address
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func retrieveAccount(privateKey string) (account *Account, err error) {
	if strings.HasPrefix(privateKey, "0x") {
		privateKey = privateKey[2:]
	}
	acc := &Account{}
	acc.privateKey, err = crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, InvalidPrivateKeyError
	}
	var ok bool
	acc.publicKey, ok = acc.privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		return nil, InvalidPrivateKeyError
	}
	acc.address = crypto.PubkeyToAddress(*acc.publicKey)
	return acc, nil
}

func createAccount() (account *Account) {
	privateKey, _ := crypto.GenerateKey()
	account, _ = retrieveAccount(hexutil.Encode(crypto.FromECDSA(privateKey)))
	return account
}

func (acc *Account) Address() Address {
	return Address(acc.address.Hex())
}

func (acc *Account) PrivateKey() string {
	return string(crypto.FromECDSA(acc.privateKey))
}

type Hash string

type Address string

func (addr Address) IsValid() bool {
	return common.IsHexAddress(string(addr))
}

func (addr Address) address() common.Address {
	return common.HexToAddress(string(addr))
}
