package sdk

import (
	"context"
	"math/big"
	"testing"
)

func TestTFC_BasicInfo(t *testing.T) {
	mockEth := NewMockEthereum()
	mockEth.Start()
	defer mockEth.Stop()

	amount := big.NewInt(1000)

	sdk := NewSDKWithBackend(mockEth.Backend)
	address, err := sdk.DeployTFCSync(context.Background(), []Address{PredefinedAccounts[0].Address()}, []*big.Int{amount}, PredefinedAccounts[0])
	if err != nil {
		t.Fatal(err)
	}
	tfc, err := sdk.TFC(address)
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
	address, err := sdk.DeployTFCSync(context.Background(), []Address{PredefinedAccounts[0].Address()}, []*big.Int{amount}, PredefinedAccounts[0])
	if err != nil {
		t.Fatal(err)
	}
	tfc, err := sdk.TFC(address)
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
