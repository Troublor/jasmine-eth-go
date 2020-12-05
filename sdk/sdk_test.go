package sdk

import (
	"context"
	"math/big"
	"testing"
)

func TestSDK_DeployManager(t *testing.T) {
	mockEth := NewMockEthereum()
	mockEth.Start()
	defer mockEth.Stop()

	sdk := NewSDKWithBackend(mockEth.Backend)
	admin, err := sdk.RetrieveAccount("0x4f3edf983ac636a65a842ce7c78d9aa706d3b113bce9c46f30d7d21715b23b1d")
	if err != nil {
		t.Fatal(err)
	}
	addrCh, errCh := sdk.DeployManager(context.Background(), admin)
	select {
	case err := <-errCh:
		t.Fatal(err)
	case addr := <-addrCh:
		manager, err := sdk.Manager(addr)
		if err != nil {
			t.Fatal(err)
		}
		signer, err := manager.Signer()
		if err != nil {
			t.Fatal(err)
		}
		if signer != admin.Address() {
			t.Fatal("signer is not deployer")
		}
	}
}

func TestSDK_DeployTFC(t *testing.T) {
	mockEth := NewMockEthereum()
	mockEth.Start()
	defer mockEth.Stop()

	sdk := NewSDKWithBackend(mockEth.Backend)
	addressCh, errCh := sdk.DeployTFC(context.Background(), PredefinedAccounts[0])
	select {
	case err := <-errCh:
		t.Fatal(err)
	case address := <-addressCh:
		tfc, err := sdk.TFC(address)
		if err != nil {
			t.Fatal(err)
		}
		balance, err := tfc.BalanceOf(PredefinedAccounts[0].Address())
		if err != nil {
			t.Fatal(err)
		}
		if balance.Cmp(big.NewInt(0)) != 0 {
			t.Fatal("initial supply is not correct")
		}
	}
}
