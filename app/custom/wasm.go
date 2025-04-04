package custom

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
)

// CustomWasmModule は wasm.AppModuleBasic をラップし、
// DefaultGenesis をオーバーライドしてカスタムパラメータを設定します。
type CustomWasmModule struct {
	wasm.AppModuleBasic // wasm.AppModuleBasic を埋め込む
	// cdc codec.Codec // <- このフィールドは不要
}

// DefaultGenesis は module.AppModuleBasic インターフェースを実装します。
// ジェネシス状態を生成し、カスタムパラメータを設定します。
// cdc は引数として渡されます。
func (cm CustomWasmModule) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage { // ★ 引数で cdc codec.JSONCodec を受け取るように変更
	genesis := &wasmtypes.GenesisState{
		Params: wasmtypes.DefaultParams(),
	}

	genesis.Params.CodeUploadAccess = wasmtypes.AllowEverybody
	genesis.Params.InstantiateDefaultPermission = wasmtypes.AccessTypeEverybody

	// 引数で受け取った cdc を使用してマーシャリング
	return cdc.MustMarshalJSON(genesis) // ★ 引数の cdc を使用
}
