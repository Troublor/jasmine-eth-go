# TFC-ERC20 Exchange Specification

1. Construct sdk object
```go
endpoint := "wss://rinkeby.infura.io/ws/v3/e8e5b9ad18ad4daeb0e01a522a989d66"
sdk, err := NewSDK(endpoint)
if err != nil {
    panic(err)
}
```
2. Construct tfcERC20Contract object
```go
erc20ContractAddress := Address("0x401Ef2b876Db2608e4A353800BBaD1E3e3Ea8B46")
erc20Contract, err := sdk.TFC(erc20ContractAddress)
if err != nil {
    panic(err)
}
```
3. Retrieve bridge account
```go
bridgeAccountPrivateKey := "0x96ca1b47bd2f7b6c1a3018e6038be291c9f5ff9556e5200f677c295693a31c60"
bridgeAccount, err := sdk.RetrieveAccount(bridgeAccountPrivateKey)
if err != nil {
    panic(err)
}
```

## Estimate Transaction Fee

Ethereum transaction fee is calculated by `gas * gasPrice`, where `gas` is the measurement of how much computation the transaction needs. `gasPrice` is the price of each `gas` in Ethereum native currency `wei`. 

The estimatation will estimate the `gas` needed by the transaction, and get the latest `gasPrice` on the network. 

### Inputs
1. `recipient` address
2. exchange tfc `amount`
3. `bridgeAccount`
4. `minGas`: if `estimatedGas < mimGas`, then `estimatedGas = minGas`.
5. `transactionFeeRate`: the rate of interest we take from each exchange transaction. `estimatedGas = estimatedGas * (1 + transactionFeeRate)`.

### Outputs
1. `requiredTransferAmount`: the amount of `wei` need to be transferred as transaction fee. `requiredTransferAmount = estimatedGas * gasPrice * (1 + transactionFeeRate)`. 
2. `estimatedGas`
3. `gasPrice`
4. `err`

### Error Handling
Error should not happen in usual cases. 
If there is error, it usually the serious issues on the server side.

### Usage
```go
recipient := Address("0x90F8bf6A479f320ead074411a4B0e7944Ea8c9C1")
amount, _ := new(big.Int).SetString("1000000000000000000", 10) // 1 TFC
minGas := 60000 // this amount of gas will cover all cases of transaction to exchange TFC ERC20
transactionFeeRate := 0.1
requiredTransferAmount, estimatedGas, gasPrice, err := tfcContract.EstimateTFCExchangeFee(context.Background(), recipient, amount, bridgeAccount, minGas, transactionFeeRate)
if err != nil {
    panic(err)
}
```

## Check Fee Deposit Transaction Hash

### Inputs
1. `depositTransactionHash`: transaction hash provided by users which is the transaction to transfer fee to bridge account.
2. `bridgeAccountAddress`
3. `transactionConfirmationRequirement`: the number of confirmations required for `depositTransaction` to be confirmed. 

### Outputs
1. `recipient`: the sender of `depositTransaction`, the address which will received TFC ERC20.
2. `depositAmount`: the amount of `wei` deposit in the `depositTransaction`.
3. `err`

### Error Handling
1. `UnknowTransactionHashErr`: if the `depositTransaction` cannot be found by the hash.
2. `UnconfirmedTransactionErr`: if the `depositTransaction` are not confirmed.
3. `InvalidDepositErr`: if the `depositTransaction` is not sending fee to bridge account.
4. other unusual errors

### Usage
```go
depositTransactionHash := "0x0e87e93aa08fd149f4f66e6939543b220b2ac77697f786c0ca5e4e88022c564d"
transactionConfirmationRequirement := 6
recipient, depositAmount, err := tfcContract.CheckTransactionFeeDeposit(context.Background(), depositTransactionHash, bridgeAccount.Address(), transactionConfirmationRequirement)
if err != nil {
    panic(err)
}
```

## Send ERC20 Mint Transaction

### Inputs
1. `recipient`: the address which will received TFC ERC20.
2. exchange tfc `amount`
3. `bridgeAccount`
4. `depositAmount`: the amount of `wei` deposit in the `depositTransaction`.
5. `estimatedGas`
6. `gasPrice`
7. `transactionFeeRate`: the rate of interest we take from each exchange transaction.

### Outputs
1. `txHash`: hash of the mint transaction
2. `err`

### Error Handling
1. `InsufficientBalanceErr`: if the balance of `bridgeAccount` is less than the given `depositAmount`. 
   This is usually should not happen, if this happens, there might be a bug in the program. 
2. `InsufficientGasErr`: if the provided `estimatedGas` is not enough for the transaction.
   This is usually due to the changes in `recipient` account which causes the gas requirement of transaction changes, i.e., fault of users.
3. `InsufficientTransactionFeeErr`: if the `depositAmount` is not enough to pay for the transaction fee (`estimatedGas * gasPrice`). 
4. other unusual errors

### Usage
```go
txHash, err := tfcContract.SendMintTransaction(
     context.Background(),
     recipient,
     amount,
     bridgeAccount,
     depositAmount,
     estimatedGas,
     gasPrice,
     transactionFeeRate,
 )
 if err != nil {
     panic(err)
 }
```

## Wait until TFC ERC20 Mint Transaction is Confirmed

### Inputs
1. `transactionHash`: transaction hash of ERC20 mint transaction.
2. `transactionConfirmationRequirement`: the number of confirmations required for ERC20 mint transaction to be confirmed.

### Outputs
1. `doneCh`: will be closed when the transaction is confirmed.
2. `errCh`: will be fed error when there is error.

### Error Handling
1. unusual errors

### Usage
```go
transactionConfirmationRequirement := 6
doneCh, errCh := tfcContract.UntilBridgeTFCExchangeComplete(context.Background(), txHash, transactionConfirmationRequirement)
 select {
 case <-doneCh:
     fmt.Println("Mint done")
 case err = <-errCh:
     t.Fatal(err)
 }
```