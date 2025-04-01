package custom

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

type CustomBankModule struct {
	bank.AppModule
	cdc codec.Codec
}

func (cm CustomBankModule) DefaultGenesis() json.RawMessage {
	genesis := banktypes.DefaultGenesisState()

	metadata := banktypes.Metadata{
		Description: "The native token of the Gluon network.",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    "uglu",
				Exponent: 0,
				Aliases:  []string{"microglu"},
			},
			{
				Denom:    "glu",
				Exponent: 6,
			},
		},
		Base:    "uglu",
		Display: "glu",
		Name:    "Gluon GLU",
		Symbol:  "GLU",
	}

	genesis.DenomMetadata = append(genesis.DenomMetadata, metadata)

	return cm.cdc.MustMarshalJSON(genesis)
}
