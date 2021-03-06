package sdk

import (
	"context"
	"github.com/Troublor/jasmine-eth-go/token"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	solsha3 "github.com/offchainlabs/go-solidity-sha3"
	"github.com/status-im/keycard-go/hexutils"
	"math/big"
	"strings"
)

type Manager struct {
	backend  Backend // connection backend to ethereum
	provider *provider

	address  common.Address
	contract *token.TFCManager
}

/**
Create a new TFC instance by providing the sdk object and the Address of TFC ERC20 contract
*/
func NewManager(backend Backend, managerAddress Address) (manager *Manager, err error) {
	manager = &Manager{
		address:  managerAddress.address(),
		backend:  backend,
		provider: NewProvider(backend),
	}
	manager.contract, err = token.NewTFCManager(common.HexToAddress(string(managerAddress)), backend)
	if err != nil {
		return nil, err
	}
	return manager, nil
}

/* Call wrappers */

/**
Returns the name of the token, i.e. TFCToken
*/
func (manager *Manager) TFCAddress() (tfcAddress Address, err error) {
	addr, err := manager.contract.TfcToken(nil)
	if err != nil {
		return "", nil
	}
	return Address(addr.Hex()), err
}

func (manager *Manager) IsNonceUsed(nonce *big.Int) (used bool, err error) {
	return manager.contract.UsedNonces(nil, nonce)
}

func (manager *Manager) Signer() (signerAddress Address, err error) {
	addr, err := manager.contract.Signer(nil)
	if err != nil {
		return "", nil
	}
	return Address(addr.Hex()), err
}

func (manager *Manager) GetUnusedNonce() (nonce *big.Int, err error) {
	for nonce = big.NewInt(0); ; nonce.Add(nonce, big.NewInt(1)) {
		used, err := manager.IsNonceUsed(nonce)
		if err != nil {
			return nil, err
		}
		if !used {
			return nonce, nil
		}
	}
}

func (manager *Manager) SignTFCClaim(recipient Address, amount *big.Int, nonce *big.Int, signer *Account) (signature string, err error) {
	hash := solsha3.SoliditySHA3(
		[]string{"address", "uint256", "uint256", "address"},
		[]interface{}{
			recipient.address().Hex(),
			amount.String(),
			nonce.String(),
			manager.address.Hex(),
		},
	)

	hash = solsha3.SoliditySHA3(
		[]string{"string", "bytes32"},
		[]interface{}{
			"\x19Ethereum Signed Message:\n32",
			hash,
		},
	)

	sig, err := crypto.Sign(hash[:], signer.privateKey)
	if err != nil {
		return "", err
	}
	// weird Ethereum quirk
	sig[64] += 27

	signature = "0x" + hexutils.BytesToHex(sig)
	return signature, nil
}

func (manager *Manager) ClaimTFC(ctx context.Context, amount *big.Int, nonce *big.Int, signature string, claimer *Account) (doneCh chan interface{}, errCh chan error) {
	doneCh = make(chan interface{}, 1)
	errCh = make(chan error, 1)
	auth := bind.NewKeyedTransactor(claimer.privateKey)
	if strings.HasPrefix(signature, "0x") {
		signature = signature[2:]
	}
	tx, err := manager.contract.ClaimTFC(auth, amount, nonce, hexutils.HexToBytes(signature))
	if err != nil {
		errCh <- err
		return nil, errCh
	}
	receiptCh, eCh := manager.provider.AsyncTransaction(ctx, tx.Hash(), ConfirmationRequirement)
	go func() {
		select {
		case <-receiptCh:
			close(doneCh)
		case err := <-eCh:
			errCh <- err
		}
	}()
	return doneCh, errCh
}

func (manager *Manager) ClaimTFCSync(ctx context.Context, amount *big.Int, nonce *big.Int, signature string, claimer *Account) (err error) {
	doneCh, errCh := manager.ClaimTFC(ctx, amount, nonce, signature, claimer)
	select {
	case <-doneCh:
		return nil
	case err := <-errCh:
		return err
	}
}

func (manager *Manager) UntilClaimTFCComplete(ctx context.Context, recipient Address, amount *big.Int, nonce *big.Int, confirmationRequirement int) (doneCh chan interface{}, errCh chan error) {
	doneCh = make(chan interface{}, 1)
	errCh = make(chan error, 1)
	go func() {
		query := ethereum.FilterQuery{
			Addresses: []common.Address{manager.address},
		}
		logsCh := make(chan types.Log)
		sub, err := manager.backend.SubscribeFilterLogs(ctx, query, logsCh)
		if err != nil {
			errCh <- err
			return
		}
		defer sub.Unsubscribe()
		contractAbi, err := abi.JSON(strings.NewReader(token.TFCManagerABI))
		if err != nil {
			errCh <- err
			return
		}

		isTargetEvent := func(log types.Log) bool {
			event := struct {
				Recipient common.Address
				Amount    *big.Int
				Nonce     *big.Int
				Sig       []byte
			}{}
			err := contractAbi.Unpack(&event, "ClaimTFC", log.Data)
			if err != nil {
				return false
			}
			event.Recipient = common.BytesToAddress(log.Topics[1].Bytes())
			event.Amount = new(big.Int).SetBytes(log.Topics[2].Bytes())
			event.Nonce = new(big.Int).SetBytes(log.Topics[3].Bytes())
			return event.Recipient == recipient.address() &&
				event.Amount.Cmp(amount) == 0 &&
				event.Nonce.Cmp(nonce) == 0
		}

		var eventWaitingCtx context.Context
		var eventWaitingCancel context.CancelFunc = func() {}
		waitTxConfirm := func(txHash common.Hash) {
			receiptCh, eCh := manager.provider.AsyncTransaction(eventWaitingCtx, txHash, confirmationRequirement)
			select {
			case <-ctx.Done():
				if ctx.Err() != context.Canceled {
					errCh <- ctx.Err()
				}
			case err := <-eCh:
				errCh <- err
			case <-receiptCh:
				close(doneCh)
			}
		}

		// check past events
		pastLogs, err := manager.backend.FilterLogs(ctx, query)
		for _, log := range pastLogs {
			if !log.Removed && isTargetEvent(log) {
				eventWaitingCtx, eventWaitingCancel = context.WithCancel(ctx)
				go waitTxConfirm(log.TxHash)
			}
		}

		for {
			select {
			case <-doneCh:
				eventWaitingCancel()
				return
			case <-ctx.Done():
				errCh <- ctx.Err()
				eventWaitingCancel()
				return
			case err := <-sub.Err():
				errCh <- err
				eventWaitingCancel()
				return
			case log := <-logsCh:
				if isTargetEvent(log) {
					if log.Removed {
						eventWaitingCancel()
						continue
					}
					eventWaitingCtx, eventWaitingCancel = context.WithCancel(ctx)
					go waitTxConfirm(log.TxHash)
				}
			}
		}

	}()
	return doneCh, errCh
}
