package sdk

import (
	"fmt"
	"testing"
)

func TestVersion(t *testing.T) {
	fmt.Println(Version())
	newSDK, _ := NewSDK("http://localhost:8545")
	fmt.Println(newSDK.Version())
}
