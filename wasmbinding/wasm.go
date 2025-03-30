package wasmbinding

import (
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"

	contractmanagerkeeper "github.com/neutron-org/neutron/v5/x/contractmanager/keeper"
	feerefunderkeeper "github.com/neutron-org/neutron/v5/x/feerefunder/keeper"
	interchaintxsmodulekeeper "github.com/neutron-org/neutron/v5/x/interchaintxs/keeper"
	tokenfactorykeeper "github.com/neutron-org/neutron/v5/x/tokenfactory/keeper"
	transferkeeper "github.com/neutron-org/neutron/v5/x/transfer/keeper"
)

// RegisterCustomPlugins returns wasmkeeper.Option that we can use to connect handlers for implemented custom queries and messages to the App
func RegisterCustomPlugins(
	contractmanagerKeeper *contractmanagerkeeper.Keeper,
	feeRefunderKeeper *feerefunderkeeper.Keeper,
	ictxKeeper *interchaintxsmodulekeeper.Keeper,
	tfk *tokenfactorykeeper.Keeper,
	transfer *transferkeeper.KeeperTransferWrapper,
) []wasmkeeper.Option {
	wasmQueryPlugin := NewQueryPlugin(contractmanagerKeeper, ictxKeeper, feeRefunderKeeper, tfk)

	queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom: CustomQuerier(wasmQueryPlugin),
	})
	messagePluginOpt := wasmkeeper.WithMessageHandlerDecorator(
		CustomMessageDecorator(contractmanagerKeeper, ictxKeeper, tfk, transfer),
	)

	return []wasmkeeper.Option{
		queryPluginOpt,
		messagePluginOpt,
	}
}
