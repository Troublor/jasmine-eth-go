package sdk

import (
	"context"
	"fmt"
	"github.com/Troublor/jasmine-eth-go/token"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type TFC struct {
	*provider

	contract *token.TFCToken
}

/**
Create a new TFC instance by providing the sdk object and the Address of TFC ERC20 contract
*/
func NewTFC(backend Backend, tfcAddress Address) (tfc *TFC, err error) {
	tfc = &TFC{
		provider: NewProvider(backend),
	}
	tfc.contract, err = token.NewTFCToken(common.HexToAddress(string(tfcAddress)), backend)
	if err != nil {
		return nil, err
	}
	return tfc, nil
}

/* Call wrappers */

/**
Returns the name of the token, i.e. TFCToken
*/
func (tfc *TFC) Name() (name string, err error) {
	return tfc.contract.Name(nil)
}

/**
Returns the symbol of the token, i.e. TFC
*/
func (tfc *TFC) Symbol() (symbol string, err error) {
	return tfc.contract.Symbol(nil)
}

/**
Returns the number of decimals the token uses - e.g. 8, means to divide the token amount by 100000000 to get its user representation.
*/
func (tfc *TFC) Decimals() (decimals uint8, err error) {
	return tfc.contract.Decimals(nil)
}

/**
Returns the total token supply.
*/
func (tfc *TFC) TotalSupply() (totalSupply *big.Int, err error) {
	return tfc.contract.TotalSupply(nil)
}

/**
Returns the Account balance with the provided Address.
*/
func (tfc *TFC) BalanceOf(address Address) (balance *big.Int, err error) {
	if !address.IsValid() {
		return nil, InvalidAddressError
	}
	return tfc.contract.BalanceOf(nil, common.HexToAddress(string(address)))
}

/**
Returns the amount which spender is still allowed to withdraw from owner.
*/
func (tfc *TFC) Allowance(owner Address, spender Address) (amount *big.Int, err error) {
	if !owner.IsValid() {
		return nil, InvalidAddressError
	}
	if !spender.IsValid() {
		return nil, InvalidAddressError
	}
	return tfc.contract.Allowance(nil, common.HexToAddress(string(owner)), common.HexToAddress(string(spender)))
}

/* Send wrappers */

/**
Transfer the amount of balance from current Account (specified in SDK) to the given "to" Account.

This function requires privateKey has been set in SDK.
*/
func (tfc *TFC) Transfer(ctx context.Context, to Address, amount *big.Int, sender *Account) (doneCh chan interface{}, errCh chan error) {
	doneCh = make(chan interface{}, 0)
	errCh = make(chan error, 0)
	auth := bind.NewKeyedTransactor(sender.privateKey)
	tx, err := tfc.contract.Transfer(auth, to.address(), amount)
	if err != nil {
		errCh <- err
		return doneCh, errCh
	}
	receiptCh, eCh := tfc.AsyncTransaction(ctx, tx.Hash(), ConfirmationRequirement)
	go func() {
		select {
		case <-receiptCh:
			close(doneCh)
		case err := <-eCh:
			errCh <- err
		}
	}()
	return doneCh, errCh
}

func (tfc *TFC) TransferSync(ctx context.Context, to Address, amount *big.Int, sender *Account) (err error) {
	doneCh, errCh := tfc.Transfer(ctx, to, amount, sender)
	select {
	case <-doneCh:
		return nil
	case err := <-errCh:
		return err
	}
}

/**
Transfer the amount of balance from the given "from" Account to the given "to" Account.

This function requires privateKey has been set in SDK, which will be used to sign the ethereum transaction.
*/
func (tfc *TFC) TransferFrom(ctx context.Context, from Address, to string, amount *big.Int, sender *Account) (doneCh chan interface{}, errCh chan error) {
	panic(UnimplementedError)
}

/**
Allows spender to withdraw from the current Account (specified in SDK) multiple times, up to the given amount.

This function requires privateKey has been set in SDK.
*/
func (tfc *TFC) Approve(ctx context.Context, spender Address, amount *big.Int, sender *Account) (doneCh chan interface{}, errCh chan error) {
	panic(UnimplementedError)
}

/**
Generate the amount of tokens and put them in the balance of the given "to" Account.
This function can only be called by Account (specified in SDK) which has MINTER_ROLE of smart contract.

This function requires privateKey has been set in SDK.
*/
func (tfc *TFC) Mint(ctx context.Context, to Address, amount *big.Int, sender *Account) (doneCh chan interface{}, errCh chan error) {
	doneCh = make(chan interface{}, 0)
	errCh = make(chan error, 0)
	auth := bind.NewKeyedTransactor(sender.privateKey)
	tx, err := tfc.contract.Mint(auth, to.address(), amount)
	if err != nil {
		errCh <- err
		return doneCh, errCh
	}
	receiptCh, eCh := tfc.AsyncTransaction(ctx, tx.Hash(), ConfirmationRequirement)
	go func() {
		select {
		case <-receiptCh:
			close(doneCh)
		case err := <-eCh:
			errCh <- err
		}
	}()
	return doneCh, errCh
}

func (tfc *TFC) MintSync(ctx context.Context, to Address, amount *big.Int, sender *Account) (err error) {
	doneCh, errCh := tfc.Mint(ctx, to, amount, sender)
	select {
	case <-doneCh:
		return nil
	case err := <-errCh:
		return err
	}
}

/* Anonymous wrappers */

func (tfc *TFC) BridgeTFCExchange(ctx context.Context, depositTransactionHash string, amount *big.Int, minter *Account, depositTransactionConfirmationRequirement int) (recipient Address, transactionHashErr error, doneCh chan interface{}, errCh chan error) {
	chainID, err := tfc.backend.NetworkID(context.Background())
	if err != nil {
		return "", err, nil, nil
	}

	tx, pending, err := tfc.backend.TransactionByHash(ctx, common.HexToHash(depositTransactionHash))
	if err == ethereum.NotFound {
		return "", UnknownTransactionHashErr, nil, nil
	} else if err != nil {
		return "", err, nil, nil
	}
	if pending {
		return "", UnconfirmedTransactionErr, nil, nil
	}

	msg, err := tx.AsMessage(types.NewEIP155Signer(chainID))
	if err != nil {
		return "", err, nil, nil
	}

	receipt, err := tfc.backend.TransactionReceipt(ctx, common.HexToHash(depositTransactionHash))
	if err == ethereum.NotFound {
		return "", UnknownTransactionHashErr, nil, nil
	} else if err != nil {
		return "", err, nil, nil
	}
	// check if receipt is on canonical chain
	blockHash := receipt.BlockHash
	canonicalBlock, err := tfc.backend.BlockByNumber(ctx, receipt.BlockNumber)
	if err != nil {
		return "", err, nil, nil
	}
	if blockHash != canonicalBlock.Hash() {
		return "", UnconfirmedTransactionErr, nil, nil
	}
	currentBlock, err := tfc.backend.BlockByNumber(ctx, nil)
	if err != nil {
		return "", err, nil, nil
	}
	if currentBlock.Number().Sub(currentBlock.Number(), receipt.BlockNumber).Cmp(big.NewInt(int64(depositTransactionConfirmationRequirement))) < 0 {
		return "", UnconfirmedTransactionErr, nil, nil
	}
	recipient = Address(msg.From().Hex())
	// transaction confirmed
	fmt.Println("mint")
	doneCh, errCh = tfc.Mint(ctx, recipient, amount, minter)
	return recipient, nil, doneCh, errCh
}
