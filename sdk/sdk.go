package sdk

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type SDK struct {
	client *ethclient.Client // connection client to ethereum

	// optional account info
	account *account // default account

	// components
	tfc *TFC // TFC ERC20 token
}

//NewSDK creates a new SDK instance with connection to blockchain endpoint
func NewSDK(blockchainEndpoint string) (sdk *SDK, error error) {
	client, err := ethclient.Dial(blockchainEndpoint)
	if err != nil {
		return nil, err
	}
	return &SDK{
		client: client,
	}, nil
}

//setDefaultAccount sets the default account to sign ethereum transactions by providing its privateKey
func (sdk *SDK) SetDefaultAccount(privateKey string) (err error) {
	acc := &account{}
	acc.privateKey, err = crypto.HexToECDSA(privateKey)
	if err != nil {
		return InvalidPrivateKeyError
	}
	var ok bool
	acc.publicKey, ok = acc.privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		return InvalidPrivateKeyError
	}
	acc.address = crypto.PubkeyToAddress(*acc.publicKey)
	sdk.account = acc
	return nil
}

/**
DefaultAccount returns the current default account in sdk (can be set via SetDefaultAccount())
*/
func (sdk *SDK) DefaultAccount() *account {
	return sdk.account
}
