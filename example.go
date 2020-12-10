package jasmine_eth_go

import (
	"context"
	"fmt"
	"github.com/Troublor/jasmine-eth-go/sdk"
	"math/big"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
func main() {
	sdkObject, err := sdk.NewSDK("ws://localhost:8545")
	checkErr(err)

	// get admin account using private key
	adminAccount, err := sdkObject.RetrieveAccount("0x11cb04ef3d5b276da031e0410d9425726187739cbe54cdedd5401911e7428df3")
	checkErr(err)

	manager, err := sdkObject.Manager("0xb402822CC243E8f86E28c2F79c67DAcD14A9cc01")
	checkErr(err)

	// data from message of client
	var recipient sdk.Address
	var amount *big.Int

	// generate an used nonce
	var nonce *big.Int

	// sign message
	signature, err := manager.SignTFCClaim(recipient, amount, nonce, adminAccount)
	checkErr(err)
	fmt.Println(signature)

	// wait for 6 block confirmations
	confirmationRequirement := 6
	doneCh, errCh := manager.UntilClaimTFCComplete(context.Background(), recipient, amount, nonce, confirmationRequirement)
	select {
	case <-doneCh:
		fmt.Println("TFC Claim confirmed")
	case err = <-errCh:
		panic(err)
	}
}
