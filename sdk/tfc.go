package sdk

import (
	"github.com/Troublor/jasmine-go/token"
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
func (tfc *TFC) name() (name string) {
	panic(UnimplementedError)
}

/**
Returns the symbol of the token, i.e. TFC
*/
func (tfc *TFC) symbol() (symbol string) {
	panic(UnimplementedError)
}

/**
Returns the number of decimals the token uses - e.g. 8, means to divide the token amount by 100000000 to get its user representation.
*/
func (tfc *TFC) decimals() (decimals uint8) {
	panic(UnimplementedError)
}

/**
Returns the total token supply.
*/
func (tfc *TFC) totalSupply() (totalSupply *big.Int) {
	panic(UnimplementedError)
}

/**
Returns the account balance with the provided address.
*/
func (tfc *TFC) balanceOf(account address) (balance *big.Int, err error) {
	panic(UnimplementedError)
}

/**
Returns the amount which spender is still allowed to withdraw from owner.
*/
func (tfc *TFC) allowance(owner address, spender address) (amount *big.Int, err error) {
	panic(UnimplementedError)
}

/* Send wrappers */

/**
Transfer the amount of balance from current account (specified in SDK) to the given "to" account.

This function requires privateKey has been set in SDK.
*/
func (tfc *TFC) transfer(to address, amount *big.Int) (err error) {
	if tfc.sdk.account == nil {
		return NoPrivateKeyError
	}
	panic(UnimplementedError)
}

/**
Transfer the amount of balance from the given "from" account to the given "to" account.

This function requires privateKey has been set in SDK, which will be used to sign the ethereum transaction.
*/
func (tfc *TFC) transferFrom(from address, to string, amount *big.Int) (err error) {
	if tfc.sdk.account == nil {
		return NoPrivateKeyError
	}
	panic(UnimplementedError)
}

/**
Allows spender to withdraw from the current account (specified in SDK) multiple times, up to the given amount.

This function requires privateKey has been set in SDK.
*/
func (tfc *TFC) approve(spender address, amount *big.Int) (err error) {
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
func (tfc *TFC) mint(to address, amount *big.Int) (err error) {
	if tfc.sdk.account == nil {
		return NoPrivateKeyError
	}
	panic(UnimplementedError)
}

/* Anonymous wrappers */
