package types

import

// DefaultGenesis returns the default genesis state
host "github.com/cosmos/ibc-go/v10/modules/core/24-host"

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(), PortId: PortID,
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if err := host.PortIdentifierValidator(gs.PortId); err != nil {
		return err
	}

	return gs.Params.Validate()
}
