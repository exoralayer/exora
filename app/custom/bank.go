package custom

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

type CustomBankModule struct {
	bank.AppModuleBasic
}

func (cm CustomBankModule) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genesis := banktypes.DefaultGenesisState()

	metadata := banktypes.Metadata{
		Description: "The native token of the Exora network.",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    "uxora",
				Exponent: 0,
				Aliases:  []string{"microxora"},
			},
			{
				Denom:    "xora",
				Exponent: 6,
			},
		},
		Base:    "uxora",
		Display: "xora",
		Name:    "Exora XORA",
		Symbol:  "XORA",
	}

	genesis.DenomMetadata = append(genesis.DenomMetadata, metadata)

	return cdc.MustMarshalJSON(genesis)
}
