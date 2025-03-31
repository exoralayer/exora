package types

// DONTCOVER

import (
	"fmt"

	"cosmossdk.io/errors"
)

// x/tokenfactory module sentinel errors
var (
	ErrInvalidSigner                = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrDenomExists                  = errors.Register(ModuleName, 2, "attempting to create a denom that already exists (has bank metadata)")
	ErrUnauthorized                 = errors.Register(ModuleName, 3, "unauthorized account")
	ErrInvalidDenom                 = errors.Register(ModuleName, 4, "invalid denom")
	ErrInvalidCreator               = errors.Register(ModuleName, 5, "invalid creator")
	ErrInvalidAuthorityMetadata     = errors.Register(ModuleName, 6, "invalid authority metadata")
	ErrInvalidGenesis               = errors.Register(ModuleName, 7, "invalid genesis")
	ErrSubdenomTooLong              = errors.Register(ModuleName, 8, fmt.Sprintf("subdenom too long, max length is %d bytes", MaxSubdenomLength))
	ErrCreatorTooLong               = errors.Register(ModuleName, 9, fmt.Sprintf("creator too long, max length is %d bytes", MaxCreatorLength))
	ErrDenomDoesNotExist            = errors.Register(ModuleName, 10, "denom does not exist")
	ErrBurnFromModuleAccount        = errors.Register(ModuleName, 11, "burning from Module Account is not allowed")
	ErrTrackBeforeSendOutOfGas      = errors.Register(ModuleName, 12, "gas meter hit maximum limit")
	ErrInvalidHookContractAddress   = errors.Register(ModuleName, 13, "invalid hook contract address")
	ErrBeforeSendHookNotWhitelisted = errors.Register(ModuleName, 14, "beforeSendHook is not whitelisted")
)
