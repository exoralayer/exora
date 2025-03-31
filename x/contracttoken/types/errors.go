package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/contracttoken module sentinel errors
var (
	ErrInvalidSigner           = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrTokenExists             = errors.Register(ModuleName, 1101, "token already exists")
	ErrTokenDoesNotExist       = errors.Register(ModuleName, 1102, "token does not exist")
	ErrInvalidDenom            = errors.Register(ModuleName, 1103, "invalid denom")
	ErrTrackBeforeSendOutOfGas = errors.Register(ModuleName, 1104, "gas meter hit maximum limit")
)
