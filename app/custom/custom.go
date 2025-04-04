package custom

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	"github.com/exoralayer/exora/app/consts"
)

func ReplaceCustomModules(
	manager module.BasicManager,
) {
	sdk.DefaultBondDenom = consts.Denom

	// bank
	oldBankModule, _ := manager[banktypes.ModuleName].(bank.AppModuleBasic)
	manager[banktypes.ModuleName] = CustomBankModule{
		AppModuleBasic: oldBankModule,
	}

	// wasm
	oldWasmModule, _ := manager[wasmtypes.ModuleName].(wasm.AppModuleBasic)
	manager[wasmtypes.ModuleName] = CustomWasmModule{
		AppModuleBasic: oldWasmModule,
	}
}
