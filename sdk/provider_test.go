package sdk

import (
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"strings"
	"testing"
	"time"
)

func prepareEthTransferTransaction(backend *MockBackend, from *Account, to *Account, amount *big.Int) (signedTx *types.Transaction) {
	nonce, _ := backend.PendingNonceAt(context.Background(), from.address)
	gasLimit := uint64(21000) // in units
	gasPrice, _ := backend.SuggestGasPrice(context.Background())
	tx := types.NewTransaction(nonce, to.address, amount, gasLimit, gasPrice, nil)
	chainId := backend.Blockchain().Config().ChainID
	signedTx, _ = types.SignTx(tx, types.NewEIP155Signer(chainId), from.privateKey)
	return signedTx
}

func TestProvider_AsyncTransaction_nonblock(t *testing.T) {
	backend := NewMockBackend()
	provider := NewProvider(backend)
	signedTx := prepareEthTransferTransaction(backend, PredefinedAccounts[0], PredefinedAccounts[1], big.NewInt(1000000000000000000))
	err := backend.SendTransaction(context.Background(), signedTx)
	if err != nil {
		t.Fatal(err)
	}
	provider.AsyncTransaction(context.Background(), signedTx.Hash(), 1)
}

func TestProvider_AsyncTransaction_context_cancel(t *testing.T) {
	backend := NewMockBackend()
	provider := NewProvider(backend)
	signedTx := prepareEthTransferTransaction(backend, PredefinedAccounts[0], PredefinedAccounts[1], big.NewInt(1000000000000000000))
	err := backend.SendTransaction(context.Background(), signedTx)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	receiptCh, errCh := provider.AsyncTransaction(ctx, signedTx.Hash(), 1)
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()
	select {
	case <-receiptCh:
		t.Fatal("tx should not be mined")
	case err = <-errCh:
		if err != context.Canceled {
			t.Fatal(err)
		}
	}
}

func TestProvider_AsyncTransaction_receiptCh(t *testing.T) {
	backend := NewMockBackend()
	provider := NewProvider(backend)
	signedTx := prepareEthTransferTransaction(backend, PredefinedAccounts[0], PredefinedAccounts[1], big.NewInt(1000000000000000000))
	err := backend.SendTransaction(context.Background(), signedTx)
	if err != nil {
		t.Fatal(err)
	}
	receiptCh, errCh := provider.AsyncTransaction(context.Background(), signedTx.Hash(), 0)
	go func() {
		time.Sleep(100 * time.Millisecond)
		// mine block
		backend.Commit()
	}()
	select {
	case receipt := <-receiptCh:
		if receipt.TxHash != signedTx.Hash() {
			t.Fatal("tx hash not match")
		}
	case err = <-errCh:
		t.Fatal(err)
	}
}

func TestProvider_AsyncTransaction_already_confirmed(t *testing.T) {
	backend := NewMockBackend()
	provider := NewProvider(backend)
	signedTx := prepareEthTransferTransaction(backend, PredefinedAccounts[0], PredefinedAccounts[1], big.NewInt(1000000000000000000))
	err := backend.SendTransaction(context.Background(), signedTx)
	if err != nil {
		t.Fatal(err)
	}
	receiptCh, errCh := provider.AsyncTransaction(context.Background(), signedTx.Hash(), 0)
	// mine block
	backend.Commit()
	select {
	case receipt := <-receiptCh:
		if receipt.TxHash != signedTx.Hash() {
			t.Fatal("tx hash not match")
		}
	case err = <-errCh:
		t.Fatal(err)
	}
}

func TestProvider_AsyncTransaction_negative_confirmation_number(t *testing.T) {
	backend := NewMockBackend()
	provider := NewProvider(backend)
	signedTx := prepareEthTransferTransaction(backend, PredefinedAccounts[0], PredefinedAccounts[1], big.NewInt(1000000000000000000))
	err := backend.SendTransaction(context.Background(), signedTx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		e := recover()
		if !strings.Contains(e.(error).Error(), "confirmation number must be non-negative") {
			t.Fatal(e)
		}
	}()
	_, _ = provider.AsyncTransaction(context.Background(), signedTx.Hash(), -1)
}

func TestProvider_AsyncTransaction_confirmation(t *testing.T) {
	backend := NewMockBackend()
	provider := NewProvider(backend)
	signedTx := prepareEthTransferTransaction(backend, PredefinedAccounts[0], PredefinedAccounts[1], big.NewInt(1000000000000000000))
	err := backend.SendTransaction(context.Background(), signedTx)
	if err != nil {
		t.Fatal(err)
	}
	receiptCh, errCh := provider.AsyncTransaction(context.Background(), signedTx.Hash(), 1)
	select {
	case <-receiptCh:
		t.Fatal("should not get receipt when tx is not executed")
	case err = <-errCh:
		t.Fatal(err)
	default:
	}
	// mine one block
	backend.Commit()
	time.Sleep(100 * time.Millisecond)

	select {
	case <-receiptCh:
		t.Fatal("should not get receipt when confirmation number is 0")
	case err = <-errCh:
		t.Fatal(err)
	default:
	}
	// mine one block
	backend.Commit()
	time.Sleep(100 * time.Millisecond)

	select {
	case <-receiptCh:
	case err = <-errCh:
		t.Fatal(err)
	default:
		t.Fatal("should get receipt when confirmation requirement is met")
	}
}
