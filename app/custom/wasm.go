package custom

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
)

type CustomWasmModule struct {
	wasm.AppModuleBasic
}

func (cm CustomWasmModule) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genesis := &wasmtypes.GenesisState{
		Params: wasmtypes.DefaultParams(),
	}

	genesis.Params.CodeUploadAccess = wasmtypes.AllowEverybody
	genesis.Params.InstantiateDefaultPermission = wasmtypes.AccessTypeEverybody

	return cdc.MustMarshalJSON(genesis)
}
