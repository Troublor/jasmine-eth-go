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

	minter, err := sdk.RetrieveAccount("0x96ca1b47bd2f7b6c1a3018e6038be291c9f5ff9556e5200f677c295693a31c60")
	if err != nil {
		t.Fatal(err)
	}

	amount := new(big.Int)
	amount.SetString("1000000000000000000", 10)

	fmt.Println("start")
	recipient, txHash, err := tfcContract.BridgeTFCExchangeAsync(
		context.Background(),
		"0xd551212792aa60482695c1e6eef52e8455a36e82d558583f77d8a572d3b67b77",
		amount,
		minter,
		6,
	)
	doneCh, errCh := tfcContract.UntilBridgeTFCExchangeComplete(context.Background(), txHash, 0)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Mint to " + recipient)
	select {
	case <-doneCh:
		fmt.Println("Mint done")
	case err = <-errCh:
		t.Fatal(err)
	}
}
