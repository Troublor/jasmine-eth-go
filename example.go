package jasmine_eth_go

import (
	"context"
	"errors"
	"fmt"
	"github.com/Troublor/jasmine-eth-go/sdk"
	"math/big"
)

func checkErr(err error) {
	if err != nil {
		checkErr(err)
	}
}

func depositTxHashUsed(txHash string) bool {
	return true
}
func main() {
	const transactionFeeRate = 0.1

	erc20ContractAddress := sdk.Address("0x401Ef2b876Db2608e4A353800BBaD1E3e3Ea8B46")
	sdkObject, err := sdk.NewSDK("wss://rinkeby.infura.io/ws/v3/e8e5b9ad18ad4daeb0e01a522a989d66")
	if err != nil {
		checkErr(err)
	}

	tfcContract, err := sdkObject.TFC(erc20ContractAddress)
	if err != nil {
		checkErr(err)
	}

	bridgeAccount, err := sdkObject.RetrieveAccount("0x96ca1b47bd2f7b6c1a3018e6038be291c9f5ff9556e5200f677c295693a31c60")
	if err != nil {
		checkErr(err)
	}

	recipient := sdk.Address("0x90F8bf6A479f320ead074411a4B0e7944Ea8c9C1")
	amount := new(big.Int)
	amount.SetString("1000000000000000000", 10)

	requiredTransferAmount, estimatedGas, gasPrice, err := tfcContract.EstimateTFCExchangeFee(context.Background(), recipient, amount, bridgeAccount, 0, transactionFeeRate)
	if err != nil {
		checkErr(err)
	}
	fmt.Println("required transfer amount", requiredTransferAmount.Uint64())
	fmt.Println("estimated gas", estimatedGas)
	fmt.Println("gas price", gasPrice.Uint64())

	depositTransactionHash := "0x0e87e93aa08fd149f4f66e6939543b220b2ac77697f786c0ca5e4e88022c564d"
	if depositTxHashUsed(depositTransactionHash) {
		panic(errors.New("deposit tx used"))
	}
	recipient, depositAmount, err := tfcContract.CheckTransactionFeeDeposit(context.Background(), depositTransactionHash, bridgeAccount.Address(), 6)
	if err != nil {
		checkErr(err)
	}

	if depositAmount.Cmp(requiredTransferAmount) < 0 {
		panic(errors.New("deposit amount too small"))
	}

	fmt.Println("recipient", recipient)
	fmt.Println("deposit amount", depositAmount.Uint64())
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
		checkErr(err)
	}
	fmt.Println("txHash", txHash)

	doneCh, errCh := tfcContract.UntilBridgeTFCExchangeComplete(context.Background(), txHash, 2)
	fmt.Println("Mint to " + recipient)

	select {
	case <-doneCh:
		fmt.Println("Mint done")
	case err = <-errCh:
		checkErr(err)
	}
}
