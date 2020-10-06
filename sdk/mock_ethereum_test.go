package sdk

import (
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"testing"
	"time"
)

func TestMockBackend_SubscribeNewTransaction(t *testing.T) {
	backend := NewMockBackend()
	ch := make(chan *types.Transaction, 1)
	sub := backend.SubscribeNewTransaction(ch)
	defer sub.Unsubscribe()

	sdk := NewSDKWithBackend(backend)
	admin, _ := retrieveAccount(PredefinedPrivateKeys[0])
	_, _ = sdk.DeployManager(context.Background(), admin)
	select {
	case <-ch:
	case <-time.After(time.Second):
		t.Fatal("transaction is not fed")
	}
}

func TestMockEthereum_MineWhenTx(t *testing.T) {
	eth := NewMockEthereum()
	eth.Start()

	sdk := NewSDKWithBackend(eth.Backend)
	admin, _ := retrieveAccount(PredefinedPrivateKeys[0])
	_, err := sdk.DeployManagerSync(context.Background(), admin)
	if err != nil {
		t.Fatal(err)
	}
	_, err = sdk.DeployManagerSync(context.Background(), admin)
	if err != nil {
		t.Fatal(err)
	}
	addr, err := sdk.DeployManagerSync(context.Background(), admin)
	if err != nil {
		t.Fatal(err)
	}
	_, err = sdk.DeployManagerSync(context.Background(), admin)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(200 * time.Millisecond)
	eth.Stop()

	currentHead, err := eth.Backend.HeaderByNumber(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	eth.Backend.Commit()
	code, err := eth.Backend.CodeAt(context.Background(), addr.address(), currentHead.Number)
	if err != nil {
		t.Fatal(err)
	}
	if len(code) == 0 {
		t.Fatal("deploy reverted")
	}
}
