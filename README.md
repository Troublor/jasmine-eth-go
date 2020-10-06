# Ethereum Client of Jasmine Project (Go Implementation)

The api documentation can be found at [pkg.go.dev](https://pkg.go.dev/github.com/Troublor/jasmine-eth-go/sdk). 

## Usage

To sign a TFC claim message: 

First instantiate a new SDK object. 
```go
sdk := NewSDK("ws://3.125.17.119:8546") // connect to the dev blockchain running on server 9523
```

Then instantiate an TFC Manager object using the Manager contract address
```go
manager := sdk.Manager(managerAddress)
```

Retrieve admin account using private key:
```go
admin := sdk.RetrieveAccount(privateKey)
```

Sign a TFC claim message which will mint a certain amount of TFC token for a recipient address;
```go
nonce, err := manager.GetUnusedNonce()
signature := manager.SignTFCClaim(context.Background(), recipientAddress, amount, nonce, admin)
```

The signature string can be given to user to claim TFC tokens by themselves. 