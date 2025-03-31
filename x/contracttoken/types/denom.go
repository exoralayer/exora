package types

import (
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleDenomPrefix = "contract"
)

// GetTokenDenom constructs a denom string for tokens created by contracttoken
// based on an input creator address and a subdenom
// The denom constructed is factory/{creator}/{subdenom}
func GetTokenDenom(contract string) (string, error) {
	denom := fmt.Sprintf("%s/%s", ModuleDenomPrefix, contract)
	return denom, sdk.ValidateDenom(denom)
}

func ParseTokenDenom(denom string) (sdk.AccAddress, error) {
	parts := strings.Split(denom, "/")
	if len(parts) != 2 || parts[0] != ModuleDenomPrefix {
		return nil, errorsmod.Wrapf(ErrInvalidDenom, "invalid denom format: %s", denom)
	}
	return sdk.AccAddressFromBech32(parts[1])
}
