package wasmbinding

import (
	contractmanagerkeeper "github.com/neutron-org/neutron/v5/x/contractmanager/keeper"
	contractmanagertypes "github.com/neutron-org/neutron/v5/x/contractmanager/types"
	feerefunderkeeper "github.com/neutron-org/neutron/v5/x/feerefunder/keeper"
	icacontrollerkeeper "github.com/neutron-org/neutron/v5/x/interchaintxs/keeper"
	tokenfactorykeeper "github.com/neutron-org/neutron/v5/x/tokenfactory/keeper"
)

type QueryPlugin struct {
	contractmanagerQueryServer contractmanagertypes.QueryServer
	feeRefunderKeeper          *feerefunderkeeper.Keeper
	icaControllerKeeper        *icacontrollerkeeper.Keeper
	tokenFactoryKeeper         *tokenfactorykeeper.Keeper
}

// NewQueryPlugin returns a reference to a new QueryPlugin.
func NewQueryPlugin(
	contractmanagerKeeper *contractmanagerkeeper.Keeper,
	icaControllerKeeper *icacontrollerkeeper.Keeper,
	feeRefunderKeeper *feerefunderkeeper.Keeper,
	tfk *tokenfactorykeeper.Keeper,
) *QueryPlugin {
	return &QueryPlugin{
		contractmanagerQueryServer: contractmanagerkeeper.NewQueryServerImpl(*contractmanagerKeeper),
		feeRefunderKeeper:          feeRefunderKeeper,
		icaControllerKeeper:        icaControllerKeeper,
		tokenFactoryKeeper:         tfk,
	}
}
