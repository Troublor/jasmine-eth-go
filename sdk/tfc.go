package sdk

import (
	"context"
	"github.com/Troublor/jasmine-eth-go/token"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type TFC struct {
	sdk *SDK

	contract *token.TFCToken
}

func DeployTFC(sdk *SDK) (tfcAddress Address, txHash Hash, err error) {
	auth := bind.NewKeyedTransactor(sdk.account.privateKey)
	nonce, err := sdk.client.PendingNonceAt(context.Background(), sdk.account.address)
	if err != nil {
		return "", "", err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice, err = sdk.client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", "", err
	}
	address, tx, _, err := token.DeployTFCToken(auth, sdk.client, big.NewInt(1000))
	if err != nil {
		return "", "", err
	}
	return Address(address.Hex()), Hash(tx.Hash().Hex()), nil
}

/**
Create a new TFC instance by providing the sdk object and the Address of TFC ERC20 contract
*/
func NewTFC(sdk *SDK, tfcAddress Address) (tfc *TFC, err error) {
	tfc = &TFC{
		sdk: sdk,
	}
	tfc.contract, err = token.NewTFCToken(common.HexToAddress(string(tfcAddress)), sdk.client)
	if err != nil {
		return nil, err
	}
	return tfc, nil
}

/* Call wrappers */

/**
Returns the name of the token, i.e. TFCToken
*/
func (tfc *TFC) Name() (name string) {
	panic(UnimplementedError)
}

/**
Returns the symbol of the token, i.e. TFC
*/
func (tfc *TFC) Symbol() (symbol string) {
	panic(UnimplementedError)
}

/**
Returns the number of decimals the token uses - e.g. 8, means to divide the token amount by 100000000 to get its user representation.
*/
func (tfc *TFC) Decimals() (decimals uint8) {
	panic(UnimplementedError)
}

/**
Returns the total token supply.
*/
func (tfc *TFC) TotalSupply() (totalSupply *big.Int) {
	panic(UnimplementedError)
}

/**
Returns the account balance with the provided Address.
*/
func (tfc *TFC) BalanceOf(account Address) (balance *big.Int, err error) {
	panic(UnimplementedError)
}

/**
Returns the amount which spender is still allowed to withdraw from owner.
*/
func (tfc *TFC) Allowance(owner Address, spender Address) (amount *big.Int, err error) {
	panic(UnimplementedError)
}

/* Send wrappers */

/**
Transfer the amount of balance from current account (specified in SDK) to the given "to" account.

This function requires privateKey has been set in SDK.
*/
func (tfc *TFC) Transfer(to Address, amount *big.Int) (err error) {
	if tfc.sdk.account == nil {
		return NoPrivateKeyError
	}
	panic(UnimplementedError)
}

/**
Transfer the amount of balance from the given "from" account to the given "to" account.

This function requires privateKey has been set in SDK, which will be used to sign the ethereum transaction.
*/
func (tfc *TFC) TransferFrom(from Address, to string, amount *big.Int) (err error) {
	if tfc.sdk.account == nil {
		return NoPrivateKeyError
	}
	panic(UnimplementedError)
}

/**
Allows spender to withdraw from the current account (specified in SDK) multiple times, up to the given amount.

This function requires privateKey has been set in SDK.
*/
func (tfc *TFC) Approve(spender Address, amount *big.Int) (err error) {
	if tfc.sdk.account == nil {
		return NoPrivateKeyError
	}
	panic(UnimplementedError)
}

/**
Generate the amount of tokens and put them in the balance of the given "to" account.
This function can only be called by account (specified in SDK) which has MINTER_ROLE of smart contract.

This function requires privateKey has been set in SDK.
*/
func (tfc *TFC) Mint(to Address, amount *big.Int) (err error) {
	if tfc.sdk.account == nil {
		return NoPrivateKeyError
	}
	panic(UnimplementedError)
}

/* Anonymous wrappers */
