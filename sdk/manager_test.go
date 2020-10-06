package sdk

import (
	"context"
	"fmt"
	"math/big"
	"testing"
)

func TestManager_SignTFCClaim(t *testing.T) {
	mockEth := NewMockEthereum()
	mockEth.Start()
	defer mockEth.Stop()

	sdk := NewSDKWithBackend(mockEth.Backend)
	admin := PredefinedAccounts[0]
	user := PredefinedAccounts[2]

	//deploy manager
	address, err := sdk.DeployManagerSync(context.Background(), admin)
	if err != nil {
		t.Fatal(err)
	}

	manager, err := sdk.Manager(address)
	if err != nil {
		t.Fatal(err)
	}

	// sign claim
	sig, err := manager.SignTFCClaim(user.Address(), big.NewInt(1), big.NewInt(0), admin)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(admin.Address())
	fmt.Println(user.Address())
	fmt.Println(sig)

	sig = sig[:len(sig)-2] + "1C"

	// claim TFC using sig
	err = manager.ClaimTFCSync(context.Background(), big.NewInt(1), big.NewInt(0), sig, user)
	if err != nil {
		t.Fatal(err)
	}

	tfcAddr, err := manager.TFCAddress()
	if err != nil {
		t.Fatal(err)
	}

	tfc, err := sdk.TFC(tfcAddr)
	if err != nil {
		t.Fatal(err)
	}

	balance, err := tfc.BalanceOf(user.Address())
	if err != nil {
		t.Fatal(err)
	}

	if balance.Cmp(big.NewInt(1)) != 0 {
		t.Fatal("TFC claim failed")
	}
}
