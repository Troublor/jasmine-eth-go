package sdk

import (
	"crypto/ecdsa"
	"errors"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"testing"
)

/**
Generate a ethereum Account.
If any error happens inside, the third return value, err, will be non-nil.
Otherwise err will be nil, with the first and second return value being Address and privateKey of the generated Account.
*/
func GenerateEthAccount() (address string, privateKey string, err error) {
	privateKeyECDSA, err := crypto.GenerateKey()
	if err != nil {
		return "", "", err
	}
	privateKeyBytes := crypto.FromECDSA(privateKeyECDSA)
	privateKey = hexutil.Encode(privateKeyBytes)
	publicKeyECDSA, ok := privateKeyECDSA.Public().(*ecdsa.PublicKey)
	if !ok {
		return "", "", errors.New("generate Account failed")
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	publicKey := hexutil.Encode(publicKeyBytes)
	address = crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	return publicKey, privateKey, nil
}

func checkError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
