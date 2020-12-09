package sdk

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"
)

func TestManager_TFCClaim(t *testing.T) {
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

	nonce, err := manager.GetUnusedNonce()
	if err != nil {
		t.Fatal(err)
	}

	// sign claim
	sig, err := manager.SignTFCClaim(user.Address(), big.NewInt(1), nonce, admin)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(admin.Address())
	fmt.Println(user.Address())
	fmt.Println(sig)

	// claim TFC using sig
	err = manager.ClaimTFCSync(context.Background(), big.NewInt(1), nonce, sig, user)
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

	// wait for confirmations
	// cancel ctx
	ctx, cancel := context.WithCancel(context.Background())
	doneCh, errCh := manager.UntilClaimTFCComplete(ctx, user.Address(), big.NewInt(1), nonce, sig, 1)
	time.Sleep(time.Millisecond * 100)
	cancel()
	select {
	case <-doneCh:
		t.Fatal()
	case err := <-errCh:
		if err != context.Canceled {
			t.Fatal()
		}
	}

	// zero confirmation requirement
	doneCh, errCh = manager.UntilClaimTFCComplete(context.Background(), user.Address(), big.NewInt(1), nonce, sig, 0)
	select {
	case <-doneCh:
	case err := <-errCh:
		t.Fatal(err)
	}

	// 6 confirmation requirement
	doneCh, errCh = manager.UntilClaimTFCComplete(context.Background(), user.Address(), big.NewInt(1), nonce, sig, 6)
	time.Sleep(time.Millisecond * 100)
	hasDone := func() (bool, error) {
		select {
		case <-doneCh:
			return true, nil
		case err := <-errCh:
			return true, err
		default:
			return false, nil
		}
	}
	for i := 0; i < 6; i++ {
		if done, err := hasDone(); done || err != nil {
			t.Fatal(err)
		}
		mockEth.Backend.Commit()
		time.Sleep(time.Millisecond * 100)
	}
	if done, err := hasDone(); !done || err != nil {
		t.Fatal(err)
	}
}
