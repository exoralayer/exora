package custom

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
)

type CustomWasmModule struct {
	wasm.AppModuleBasic
	cdc codec.Codec
}

func (cm CustomWasmModule) DefaultGenesis() json.RawMessage {
	genesis := &wasmtypes.GenesisState{
		Params: wasmtypes.DefaultParams(),
	}

	genesis.Params.CodeUploadAccess = wasmtypes.AllowEverybody
	genesis.Params.InstantiateDefaultPermission = wasmtypes.AccessTypeEverybody

	return cm.cdc.MustMarshalJSON(genesis)
}
