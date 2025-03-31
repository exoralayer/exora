package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewParams creates a new Params instance.
func NewParams() Params {
	return Params{}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams()
}

// Validate validates the set of params.
func (p Params) Validate() error {

	return nil
}

func validateDenomCreationFee(i interface{}) error {
	v, ok := i.(sdk.Coins)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.Validate() != nil {
		return fmt.Errorf("invalid denom creation fee: %+v", i)
	}

	return nil
}

func validateDenomCreationGasConsume(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateFeeCollectorAddress(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if len(v) == 0 {
		return nil
	}

	_, err := sdk.AccAddressFromBech32(v)
	if err != nil {
		return fmt.Errorf("invalid fee collector address: %w", err)
	}

	return nil
}

func validateWhitelistedHooks(i interface{}) error {
	hooks, ok := i.([]*WhitelistedHook)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	seenHooks := map[string]bool{}
	for _, hook := range hooks {
		hookStr := hook.String()
		if seenHooks[hookStr] {
			return fmt.Errorf("duplicate whitelisted hook: %s", hookStr)
		}
		seenHooks[hookStr] = true
		_, err := sdk.AccAddressFromBech32(hook.DenomCreator)
		if err != nil {
			return fmt.Errorf("invalid denom creator address: %w", err)
		}
	}
	return nil
}
