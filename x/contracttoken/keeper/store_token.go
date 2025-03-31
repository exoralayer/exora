package keeper

import (
	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/gluon-zone/gluon/x/contracttoken/types"
)

func (k Keeper) GetToken(ctx sdk.Context, contract sdk.AccAddress) (types.Token, error, bool) {
	token, err := k.ContractTokens.Get(ctx, contract)
	if err == collections.ErrNotFound {
		return types.Token{}, nil, false
	}
	return token, err, true
}

func (k Keeper) SetToken(ctx sdk.Context, token types.Token) error {
	contract, err := k.addressCodec.StringToBytes(token.ContractAddress)
	if err != nil {
		return err
	}
	return k.ContractTokens.Set(ctx, contract, token)
}
