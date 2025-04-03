package custom

import (
	"github.com/cosmos/cosmos-sdk/codec"
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
	cdc codec.Codec,
) {
	sdk.DefaultBondDenom = consts.Denom

	// bank
	oldBankModule, _ := manager[banktypes.ModuleName].(bank.AppModule)
	manager[banktypes.ModuleName] = CustomBankModule{
		AppModule: oldBankModule,
		cdc:       cdc,
	}

	// wasm
	oldWasmModule, _ := manager[wasmtypes.ModuleName].(wasm.AppModule)
	manager[wasmtypes.ModuleName] = CustomWasmModule{
		AppModule: oldWasmModule,
		cdc:       cdc,
	}
}
