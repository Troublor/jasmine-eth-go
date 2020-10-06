package sdk

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"github.com/Troublor/jasmine-eth-go/token"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

type SDK struct {
	*provider

	// optional account info
	account *Account // default account
}

//NewSDK creates a new SDK instance with connection to Backend endpoint
func NewSDK(blockchainEndpoint string) (sdk *SDK, error error) {
	client, err := ethclient.Dial(blockchainEndpoint)
	if err != nil {
		return nil, err
	}
	return &SDK{
		provider: NewProvider(client),
	}, nil
}

func NewSDKWithBackend(backend Backend) (sdk *SDK) {
	return &SDK{
		provider: NewProvider(backend),
	}
}

//setDefaultAccount sets the default Account to sign ethereum transactions by providing its privateKey
func (sdk *SDK) SetDefaultAccount(privateKey string) (err error) {
	acc := &Account{}
	acc.privateKey, err = crypto.HexToECDSA(privateKey)
	if err != nil {
		return InvalidPrivateKeyError
	}
	var ok bool
	acc.publicKey, ok = acc.privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		return InvalidPrivateKeyError
	}
	acc.address = crypto.PubkeyToAddress(*acc.publicKey)
	sdk.account = acc
	return nil
}

func (sdk *SDK) RetrieveAccount(privateKey string) (account *Account, err error) {
	return retrieveAccount(privateKey)
}

func (sdk *SDK) CreateAccount() (account *Account) {
	return createAccount()
}

/**
DefaultAccount returns the current default Account in sdk (can be set via SetDefaultAccount())
*/
func (sdk *SDK) DefaultAccount() *Account {
	return sdk.account
}

func (sdk *SDK) DeployTFCSync(ctx context.Context, initialHolders []Address, initialSupplies []*big.Int, deployer *Account) (tfcAddress Address, err error) {
	tfcAddressCh, errCh := sdk.DeployTFC(ctx, initialHolders, initialSupplies, deployer)
	select {
	case addr := <-tfcAddressCh:
		return addr, nil
	case err := <-errCh:
		return "", err
	}
}

func (sdk *SDK) DeployTFC(ctx context.Context, initialHolders []Address, initialSupplies []*big.Int, deployer *Account) (tfcAddressCh chan Address, errCh chan error) {
	tfcAddressCh = make(chan Address, 1)
	errCh = make(chan error, 1)
	if len(initialHolders) != len(initialSupplies) {
		errCh <- errors.New("the number of initial holders is not equal to the number of supplies")
		return tfcAddressCh, errCh
	}
	auth := bind.NewKeyedTransactor(deployer.privateKey)
	nonce, err := sdk.backend.PendingNonceAt(context.Background(), deployer.address)
	if err != nil {
		errCh <- err
		return tfcAddressCh, errCh
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasPrice, err = sdk.backend.SuggestGasPrice(context.Background())
	if err != nil {
		errCh <- err
		return tfcAddressCh, errCh
	}
	holders := make([]common.Address, len(initialHolders))
	for i, holder := range initialHolders {
		if !holder.IsValid() {
			errCh <- InvalidAddressError
			return tfcAddressCh, errCh
		}
		holders[i] = common.HexToAddress(string(holder))
	}
	_, tx, _, err := token.DeployTFCToken(auth, sdk.backend, holders, initialSupplies)
	if err != nil {
		errCh <- err
		return tfcAddressCh, errCh
	}

	receiptCh, eCh := sdk.AsyncTransaction(ctx, tx.Hash(), ConfirmationRequirement)
	go func() {
		select {
		case receipt := <-receiptCh:
			tfcAddressCh <- Address(receipt.ContractAddress.Hex())
		case err := <-eCh:
			errCh <- err
		}
	}()
	return tfcAddressCh, errCh
}

func (sdk *SDK) DeployManagerSync(ctx context.Context, deployer *Account) (managerAddress Address, err error) {
	managerAddressCh, errCh := sdk.DeployManager(ctx, deployer)
	select {
	case addr := <-managerAddressCh:
		return addr, nil
	case err := <-errCh:
		return "", err
	}
}

func (sdk *SDK) DeployManager(ctx context.Context, deployer *Account) (managerAddressCh chan Address, errCh chan error) {
	managerAddressCh = make(chan Address, 1)
	errCh = make(chan error, 1)
	auth := bind.NewKeyedTransactor(deployer.privateKey)
	nonce, err := sdk.backend.PendingNonceAt(context.Background(), deployer.address)
	if err != nil {
		errCh <- err
		return managerAddressCh, errCh
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasPrice, err = sdk.backend.SuggestGasPrice(context.Background())
	if err != nil {
		errCh <- err
		return managerAddressCh, errCh
	}
	_, tx, _, err := token.DeployTFCManager(auth, sdk.backend)
	if err != nil {
		errCh <- err
		return managerAddressCh, errCh
	}

	receiptCh, eCh := sdk.AsyncTransaction(ctx, tx.Hash(), ConfirmationRequirement)
	go func() {
		select {
		case receipt := <-receiptCh:
			managerAddressCh <- Address(receipt.ContractAddress.Hex())
		case err := <-eCh:
			errCh <- err
		}
	}()
	return managerAddressCh, errCh
}

/**
Creates a new TFC instance based on current sdk.
This function is a wrapper of NewTFC()
*/
func (sdk *SDK) TFC(tfcAddress Address) (tfc *TFC, err error) {
	return NewTFC(sdk.backend, tfcAddress)
}

func (sdk *SDK) Manager(managerAddress Address) (manager *Manager, err error) {
	return NewManager(sdk.backend, managerAddress)
}
