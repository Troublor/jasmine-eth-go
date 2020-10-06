package sdk

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"math/big"
)

type MockEthereum struct {
	ctx     context.Context
	cancel  context.CancelFunc
	Backend *MockBackend
}

func NewMockEthereum() *MockEthereum {
	return &MockEthereum{Backend: NewMockBackend()}
}

func (eth *MockEthereum) Start() {
	eth.ctx, eth.cancel = context.WithCancel(context.Background())
	go func() {
		txCh := make(chan *types.Transaction, 1)
		sub := eth.Backend.SubscribeNewTransaction(txCh)
		defer sub.Unsubscribe()
		// mine block when there is transaction
		for {
			select {
			case <-eth.ctx.Done():
				return
			case <-txCh:
				eth.Backend.Commit()
			}
		}
	}()
}

func (eth *MockEthereum) Stop() {
	eth.cancel()
}

type MockBackend struct {
	*backends.SimulatedBackend

	newTxFeed event.Feed
}

func NewMockBackend() *MockBackend {
	balance := new(big.Int)
	balance.SetString("100000000000000000000", 10) // 100 eth in wei
	genesisAlloc := make(map[common.Address]core.GenesisAccount)
	for _, privateKey := range PredefinedPrivateKeys {
		account, _ := retrieveAccount(privateKey)
		genesisAlloc[account.address] = core.GenesisAccount{Balance: balance}
	}
	blockGasLimit := uint64(4712388)
	backend := backends.NewSimulatedBackend(genesisAlloc, blockGasLimit)
	mock := &MockBackend{SimulatedBackend: backend}
	return mock
}

func (b *MockBackend) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	// previous transactions must have been committed
	b.newTxFeed.Send(tx)
	err := b.SimulatedBackend.SendTransaction(ctx, tx)
	return err
}

func (b *MockBackend) SubscribeNewTransaction(ch chan *types.Transaction) event.Subscription {
	return b.newTxFeed.Subscribe(ch)
}
