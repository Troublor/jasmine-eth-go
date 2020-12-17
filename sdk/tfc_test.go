package sdk

import (
	"context"
	"fmt"
	"math/big"
	"testing"
)

func TestTFC_BasicInfo(t *testing.T) {
	mockEth := NewMockEthereum()
	mockEth.Start()
	defer mockEth.Stop()

	amount := big.NewInt(1000)

	sdk := NewSDKWithBackend(mockEth.Backend)
	address, err := sdk.DeployTFCSync(context.Background(), PredefinedAccounts[0])
	if err != nil {
		t.Fatal(err)
	}
	tfc, err := sdk.TFC(address)
	if err != nil {
		t.Fatal(err)
	}

	err = tfc.MintSync(context.Background(), PredefinedAccounts[0].Address(), amount, PredefinedAccounts[0])
	if err != nil {
		t.Fatal(err)
	}

	balance, err := tfc.BalanceOf(PredefinedAccounts[0].Address())
	if err != nil {
		t.Fatal(err)
	}
	if balance.Cmp(amount) != 0 {
		t.Fatal("initial supply is incorrect")
	}

	totalSupply, err := tfc.TotalSupply()
	checkError(t, err)
	if totalSupply.Cmp(amount) != 0 {
		t.Fatal("total supply is incorrect")
	}
}

func TestTFC_Transfer(t *testing.T) {
	mockEth := NewMockEthereum()
	mockEth.Start()
	defer mockEth.Stop()

	amount := big.NewInt(1000)

	sdk := NewSDKWithBackend(mockEth.Backend)
	address, err := sdk.DeployTFCSync(context.Background(), PredefinedAccounts[0])
	if err != nil {
		t.Fatal(err)
	}
	tfc, err := sdk.TFC(address)
	if err != nil {
		t.Fatal(err)
	}

	err = tfc.MintSync(context.Background(), PredefinedAccounts[0].Address(), amount, PredefinedAccounts[0])
	if err != nil {
		t.Fatal(err)
	}

	err = tfc.TransferSync(context.Background(), PredefinedAccounts[1].Address(), amount, PredefinedAccounts[0])
	if err != nil {
		t.Fatal(err)
	}

	balance0, err := tfc.BalanceOf(PredefinedAccounts[0].Address())
	if err != nil {
		t.Fatal(err)
	}
	balance1, err := tfc.BalanceOf(PredefinedAccounts[1].Address())
	if err != nil {
		t.Fatal(err)
	}
	if balance0.Cmp(big.NewInt(0)) != 0 || balance1.Cmp(amount) != 0 {
		t.Fatal("transfer does not work")
	}
}

func TestTFC_BridgeTFCExchange(t *testing.T) {
	erc20ContractAddress := Address("0x401Ef2b876Db2608e4A353800BBaD1E3e3Ea8B46")
	sdk, err := NewSDK("wss://rinkeby.infura.io/ws/v3/e8e5b9ad18ad4daeb0e01a522a989d66")
	if err != nil {
		t.Fatal(err)
	}

	tfcContract, err := sdk.TFC(erc20ContractAddress)
	if err != nil {
		t.Fatal(err)
	}

	bridgeAccount, err := sdk.RetrieveAccount("0x96ca1b47bd2f7b6c1a3018e6038be291c9f5ff9556e5200f677c295693a31c60")
	if err != nil {
		t.Fatal(err)
	}

	recipient := Address("0x90F8bf6A479f320ead074411a4B0e7944Ea8c9C1")
	amount := new(big.Int)
	amount.SetString("1000000000000000000", 10)

	requiredTransferAmount, estimatedGas, gasPrice, err := tfcContract.EstimateTFCExchangeFee(context.Background(), recipient, amount, bridgeAccount, 0, 0.1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("required transfer amount", requiredTransferAmount.Uint64())
	fmt.Println("estimated gas", estimatedGas)
	fmt.Println("gas price", gasPrice.Uint64())

	depositTransactionHash := "0x0e87e93aa08fd149f4f66e6939543b220b2ac77697f786c0ca5e4e88022c564d"
	recipient, depositAmount, err := tfcContract.CheckTransactionFeeDeposit(context.Background(), depositTransactionHash, bridgeAccount.Address(), 6)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("recipient", recipient)
	fmt.Println("deposit amount", depositAmount.Uint64())
	txHash, err := tfcContract.SendMintTransaction(
		context.Background(),
		recipient,
		amount,
		bridgeAccount,
		depositAmount,
		0,
		big.NewInt(157000000000),
		0.2,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("txHash", txHash)

	doneCh, errCh := tfcContract.UntilBridgeTFCExchangeComplete(context.Background(), txHash, 0)
	fmt.Println("Mint to " + recipient)

	select {
	case <-doneCh:
		fmt.Println("Mint done")
	case err = <-errCh:
		t.Fatal(err)
	}
}
