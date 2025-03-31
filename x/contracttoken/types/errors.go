package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/contracttoken module sentinel errors
var (
	ErrInvalidSigner           = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrNotContract             = errors.Register(ModuleName, 1101, "not a contract address")
	ErrTokenExists             = errors.Register(ModuleName, 1102, "token already exists")
	ErrTokenDoesNotExist       = errors.Register(ModuleName, 1103, "token does not exist")
	ErrInvalidDenom            = errors.Register(ModuleName, 1104, "invalid denom")
	ErrTrackBeforeSendOutOfGas = errors.Register(ModuleName, 1105, "gas meter hit maximum limit")
)
