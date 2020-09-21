package sdk

import "errors"

var (
	UnimplementedError     = errors.New("unimplemented")
	NoPrivateKeyError      = errors.New("no private key is provided")
	InvalidPrivateKeyError = errors.New("invalid private key")
)
