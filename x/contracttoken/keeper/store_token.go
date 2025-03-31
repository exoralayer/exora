package keeper

import (
	"context"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/gluon-zone/gluon/x/contracttoken/types"
)

func (k Keeper) HasToken(ctx context.Context, contract sdk.AccAddress) (bool, error) {
	_, err, found := k.GetToken(ctx, contract)
	return found, err
}

func (k Keeper) GetToken(ctx context.Context, contract sdk.AccAddress) (types.Token, error, bool) {
	token, err := k.ContractTokens.Get(ctx, contract)
	if err == collections.ErrNotFound {
		return types.Token{}, nil, false
	}
	return token, err, true
}

func (k Keeper) SetToken(ctx context.Context, token types.Token) error {
	contract, err := k.addressCodec.StringToBytes(token.ContractAddress)
	if err != nil {
		return err
	}
	return k.ContractTokens.Set(ctx, contract, token)
}
