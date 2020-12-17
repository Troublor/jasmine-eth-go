package sdk

import (
	"context"
	"fmt"
	"github.com/Troublor/jasmine-eth-go/token"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"strings"
)

type TFC struct {
	*provider

	address  Address
	contract *token.TFCToken
}

/**
Create a new TFC instance by providing the sdk object and the Address of TFC ERC20 contract
*/
func NewTFC(backend Backend, tfcAddress Address) (tfc *TFC, err error) {
	tfc = &TFC{
		provider: NewProvider(backend),
		address:  tfcAddress,
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
	doneCh, errCh = tfc.Mint(ctx, recipient, amount, minter)
	return recipient, nil, doneCh, errCh
}

func (tfc *TFC) EstimateTFCExchangeFee(ctx context.Context, recipient Address, amount *big.Int, bridgeAccount *Account, minGas uint64, transactionFeeRate float64) (requiredTransferAmount *big.Int, estimatedGas uint64, gasPrice *big.Int, err error) {
	gasPrice, err = tfc.backend.SuggestGasPrice(ctx)
	if err != nil {
		return nil, 0, nil, err
	}
	parsedABI, err := abi.JSON(strings.NewReader(token.TFCTokenABI))
	if err != nil {
		return nil, 0, nil, err
	}
	// encode input
	input, err := parsedABI.Pack("mint", recipient.address(), amount)
	if err != nil {
		return nil, 0, nil, err
	}
	// estimate gas
	// Gas estimation cannot succeed without code for method invocations
	if code, err := tfc.backend.PendingCodeAt(ctx, tfc.address.address()); err != nil {
		return nil, 0, nil, err
	} else if len(code) == 0 {
		return nil, 0, nil, bind.ErrNoCode
	}
	// If the contract surely has code (or code is not needed), estimate the transaction
	tfcAddress := tfc.address.address()
	msg := ethereum.CallMsg{From: bridgeAccount.address, To: &tfcAddress, GasPrice: gasPrice, Value: big.NewInt(0), Data: input}
	estimatedGas, err = tfc.backend.EstimateGas(ctx, msg)
	if err != nil && strings.Contains(err.Error(), "insufficient funds") {
		estimatedGas = 60000 // if estimate gas fails due to bridge account has no balance, assign a default safe gasLimit for ERC20 mint
	} else if err != nil {
		return nil, 0, nil, fmt.Errorf("failed to estimate gas needed: %v", err)
	}
	if estimatedGas < minGas {
		estimatedGas = minGas
	}
	// calculate fee with rate
	fee := new(big.Int).Mul(gasPrice, big.NewInt(int64(estimatedGas)))
	interestedFee := new(big.Float).Mul(new(big.Float).SetInt(fee), big.NewFloat(1+transactionFeeRate))
	requiredTransferAmount, accuracy := interestedFee.Int(nil)
	if accuracy == big.Below &&
		interestedFee.Cmp(new(big.Float).SetInt(requiredTransferAmount)) > 0 {
		// add one gas
		requiredTransferAmount = requiredTransferAmount.Add(requiredTransferAmount, gasPrice)
	}
	return requiredTransferAmount, estimatedGas, gasPrice, nil
}

func (tfc *TFC) CheckTransactionFeeDeposit(ctx context.Context, depositTransactionHash string, bridgeAccountAddress Address, depositTransactionConfirmationRequirement int) (recipient Address, depositAmount *big.Int, err error) {
	chainID, err := tfc.backend.NetworkID(context.Background())
	if err != nil {
		return "", nil, err
	}

	tx, pending, err := tfc.backend.TransactionByHash(ctx, common.HexToHash(depositTransactionHash))
	if err == ethereum.NotFound {
		return "", nil, UnknownTransactionHashErr
	} else if err != nil {
		return "", nil, err
	}
	if pending {
		return "", nil, UnconfirmedTransactionErr
	}

	msg, err := tx.AsMessage(types.NewEIP155Signer(chainID))
	if err != nil {
		return "", nil, err
	}

	if *msg.To() != bridgeAccountAddress.address() {
		return "", nil, InvalidDepositErr
	}

	receipt, err := tfc.backend.TransactionReceipt(ctx, common.HexToHash(depositTransactionHash))
	if err == ethereum.NotFound {
		return "", nil, UnknownTransactionHashErr
	} else if err != nil {
		return "", nil, err
	}
	// check if receipt is on canonical chain
	blockHash := receipt.BlockHash
	canonicalBlock, err := tfc.backend.BlockByNumber(ctx, receipt.BlockNumber)
	if err != nil {
		return "", nil, err
	}
	if blockHash != canonicalBlock.Hash() {
		return "", nil, UnconfirmedTransactionErr
	}
	currentBlock, err := tfc.backend.BlockByNumber(ctx, nil)
	if err != nil {
		return "", nil, err
	}
	if currentBlock.Number().Sub(currentBlock.Number(), receipt.BlockNumber).Cmp(big.NewInt(int64(depositTransactionConfirmationRequirement))) < 0 {
		return "", nil, UnconfirmedTransactionErr
	}
	recipient = Address(msg.From().Hex())
	return recipient, tx.Value(), nil
}

func (tfc *TFC) SendMintTransaction(ctx context.Context, recipient Address, amount *big.Int, minter *Account, depositAmount *big.Int, estimatedGas uint64, gasPrice *big.Int, transactionFeeRate float64) (mintTransactionHash string, err error) {
	// get the fee received from user
	receivedFee := depositAmount
	// make sure the minter account has at least receivedFee amount of ETH
	balance, err := tfc.backend.BalanceAt(ctx, minter.address, nil)
	if err != nil {
		return "", err
	}
	if balance.Cmp(receivedFee) < 0 {
		return "", InsufficientBalanceErr
	}
	// deduct transaction fee rate
	txFee := new(big.Float).Quo(new(big.Float).SetInt(receivedFee), big.NewFloat(1+transactionFeeRate))

	// encode input
	parsedABI, err := abi.JSON(strings.NewReader(token.TFCTokenABI))
	if err != nil {
		return "", err
	}
	// pack mint transaction input
	input, err := parsedABI.Pack("mint", recipient.address(), amount)
	if err != nil {
		return "", err
	}

	// estimate gas if estimatedGas == 0
	// Gas estimation cannot succeed without code for method invocations
	if code, err := tfc.backend.PendingCodeAt(ctx, tfc.address.address()); err != nil {
		return "", err
	} else if len(code) == 0 {
		return "", bind.ErrNoCode
	}
	// If the contract surely has code (or code is not needed), estimate the transaction
	tfcAddress := tfc.address.address()
	msg := ethereum.CallMsg{From: minter.address, To: &tfcAddress, GasPrice: gasPrice, Value: big.NewInt(0), Data: input}
	gas, err := tfc.backend.EstimateGas(ctx, msg)
	if err != nil {
		return "", fmt.Errorf("failed to estimate gas needed: %v", err)
	}
	if estimatedGas == 0 {
		estimatedGas = gas
	} else if estimatedGas > 0 && estimatedGas < gas {
		return "", InsufficientGasErr
	}

	// if gasPrice is zero or nil, calculate gasPrice using estimatedGas and txFee
	if gasPrice == nil || gasPrice.Cmp(big.NewInt(0)) == 0 {
		// calculate gasPrice
		calculatedGasPrice := new(big.Float).Quo(txFee, big.NewFloat(float64(estimatedGas)))
		gasPrice, _ = calculatedGasPrice.Int(nil)
	}

	if txFee.Cmp(new(big.Float).Mul(big.NewFloat(float64(estimatedGas)), new(big.Float).SetInt(gasPrice))) < 0 {
		return "", InsufficientTransactionFeeErr
	}

	// send mint transaction
	auth := bind.NewKeyedTransactor(minter.privateKey)
	auth.GasLimit = estimatedGas
	auth.GasPrice = gasPrice
	tx, err := tfc.contract.Mint(auth, recipient.address(), amount)
	if err != nil {
		return "", err
	}
	return tx.Hash().Hex(), nil
}

func (tfc *TFC) BridgeTFCExchangeAsync(ctx context.Context, depositTransactionHash string, amount *big.Int, minter *Account, depositTransactionConfirmationRequirement int) (recipient Address, mintTransactionHash string, err error) {
	chainID, err := tfc.backend.NetworkID(context.Background())
	if err != nil {
		return "", "", err
	}

	tx, pending, err := tfc.backend.TransactionByHash(ctx, common.HexToHash(depositTransactionHash))
	if err == ethereum.NotFound {
		return "", "", UnknownTransactionHashErr
	} else if err != nil {
		return "", "", err
	}
	if pending {
		return "", "", UnconfirmedTransactionErr
	}

	msg, err := tx.AsMessage(types.NewEIP155Signer(chainID))
	if err != nil {
		return "", "", err
	}

	receipt, err := tfc.backend.TransactionReceipt(ctx, common.HexToHash(depositTransactionHash))
	if err == ethereum.NotFound {
		return "", "", UnknownTransactionHashErr
	} else if err != nil {
		return "", "", err
	}
	// check if receipt is on canonical chain
	blockHash := receipt.BlockHash
	canonicalBlock, err := tfc.backend.BlockByNumber(ctx, receipt.BlockNumber)
	if err != nil {
		return "", "", err
	}
	if blockHash != canonicalBlock.Hash() {
		return "", "", UnconfirmedTransactionErr
	}
	currentBlock, err := tfc.backend.BlockByNumber(ctx, nil)
	if err != nil {
		return "", "", err
	}
	if currentBlock.Number().Sub(currentBlock.Number(), receipt.BlockNumber).Cmp(big.NewInt(int64(depositTransactionConfirmationRequirement))) < 0 {
		return "", "", UnconfirmedTransactionErr
	}
	recipient = Address(msg.From().Hex())

	// send mint transaction
	auth := bind.NewKeyedTransactor(minter.privateKey)
	tx, err = tfc.contract.Mint(auth, recipient.address(), amount)
	if err != nil {
		return "", "", err
	}
	return recipient, tx.Hash().Hex(), nil
}

func (tfc *TFC) UntilBridgeTFCExchangeComplete(ctx context.Context, mintTransactionHash string, confirmationRequirement int) (doneCh chan interface{}, errCh chan error) {
	doneCh = make(chan interface{}, 0)
	errCh = make(chan error, 0)
	receiptCh, eCh := tfc.AsyncTransaction(ctx, common.HexToHash(mintTransactionHash), confirmationRequirement)
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
