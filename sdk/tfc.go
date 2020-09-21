package sdk

import (
	"github.com/Troublor/jasmine-eth-go/token"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type TFC struct {
	sdk *SDK

	contract *token.TFCToken
}

/**
Create a new TFC instance by providing the sdk object and the address of TFC ERC20 contract
*/
func NewTFC(sdk *SDK, tfcAddress address) (tfc *TFC, err error) {
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
Returns the account balance with the provided address.
*/
func (tfc *TFC) BalanceOf(account address) (balance *big.Int, err error) {
	panic(UnimplementedError)
}

/**
Returns the amount which spender is still allowed to withdraw from owner.
*/
func (tfc *TFC) Allowance(owner address, spender address) (amount *big.Int, err error) {
	panic(UnimplementedError)
}

/* Send wrappers */

/**
Transfer the amount of balance from current account (specified in SDK) to the given "to" account.

This function requires privateKey has been set in SDK.
*/
func (tfc *TFC) Transfer(to address, amount *big.Int) (err error) {
	if tfc.sdk.account == nil {
		return NoPrivateKeyError
	}
	panic(UnimplementedError)
}

/**
Transfer the amount of balance from the given "from" account to the given "to" account.

This function requires privateKey has been set in SDK, which will be used to sign the ethereum transaction.
*/
func (tfc *TFC) TransferFrom(from address, to string, amount *big.Int) (err error) {
	if tfc.sdk.account == nil {
		return NoPrivateKeyError
	}
	panic(UnimplementedError)
}

/**
Allows spender to withdraw from the current account (specified in SDK) multiple times, up to the given amount.

This function requires privateKey has been set in SDK.
*/
func (tfc *TFC) Approve(spender address, amount *big.Int) (err error) {
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
func (tfc *TFC) Mint(to address, amount *big.Int) (err error) {
	if tfc.sdk.account == nil {
		return NoPrivateKeyError
	}
	panic(UnimplementedError)
}

/* Anonymous wrappers */
