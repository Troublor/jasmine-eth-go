package sdk

import (
	"context"
	"errors"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type provider struct {
	backend Backend
}

func NewProvider(backend Backend) *provider {
	return &provider{backend: backend}
}

func (p *provider) getConfirmationCount(ctx context.Context, receipt *types.Receipt) (count int, err error) {
	receiptBlockNumber := receipt.BlockNumber
	receiptBlock, err := p.backend.BlockByNumber(ctx, receiptBlockNumber)
	if err != nil {
		return math.MinInt32, err
	}
	if receiptBlock.Hash() != receipt.BlockHash {
		// the receipt has been removed from history (due to reorg)
		return -1, nil
	}
	// the receipt is in canonical chain
	currentHeader, err := p.backend.HeaderByNumber(ctx, nil)
	if err != nil {
		return math.MinInt32, err
	}
	currentNumber := currentHeader.Number
	c := big.NewInt(0).Sub(currentNumber, receiptBlockNumber)
	if !c.IsInt64() || c.Int64() > math.MaxInt32 {
		// if confirmation count is too large
		return math.MaxInt32, nil
	} else {
		return int(c.Int64()), nil
	}
}

func (p *provider) AsyncTransaction(ctx context.Context, txHash common.Hash, confirmationNumber int) (receiptCh chan *types.Receipt, errCh chan error) {
	if confirmationNumber < 0 {
		panic(errors.New("confirmation number must be non-negative"))
	}
	receiptCh = make(chan *types.Receipt, 1)
	errCh = make(chan error, 1)
	go func() {
		// listen to new headers
		headerCh := make(chan *types.Header, confirmationNumber)
		headerSub, err := p.backend.SubscribeNewHead(ctx, headerCh)
		if err != nil {
			// there is some error when try to subscribe new head
			errCh <- err
			return
		}
		defer headerSub.Unsubscribe()

		var receipt *types.Receipt

		checkTransaction := func() (waitForMoreBlocks bool) {
			// try to fetch receipt
			receipt, err = p.backend.TransactionReceipt(ctx, txHash)
			if err == ethereum.NotFound || receipt == nil {
				// transaction is still pending
				return true
			}
			if err != nil {
				errCh <- err
				return false
			}
			// the transaction has already been mined
			// check confirmation count
			confirmationCount, err := p.getConfirmationCount(ctx, receipt)
			if err != nil {
				// there is some error when try to get confirmation count
				errCh <- err
				return false
			}
			if confirmationCount >= confirmationNumber {
				// confirmation requirement achieved
				receiptCh <- receipt
				return false
			}
			return true
		}

		// check if transaction is already confirmed
		if waitForMoreBlocks := checkTransaction(); !waitForMoreBlocks {
			return
		}
		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case <-headerCh:
				waitForMoreBlocks := checkTransaction()
				if !waitForMoreBlocks {
					return
				}
			}
		}
	}()
	return receiptCh, errCh
}
