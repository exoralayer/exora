package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/exoralayer/exora/x/contracttoken/types"
)

func (k Keeper) HasToken(ctx context.Context, contract sdk.AccAddress) bool {
	_, err := k.GetToken(ctx, contract)
	return err == nil
}

func (k Keeper) GetToken(ctx context.Context, contract sdk.AccAddress) (types.Token, error) {
	token, err := k.Tokens.Get(ctx, contract)
	if err != nil {
		return types.Token{}, err
	}
	return token, nil
}

func (k Keeper) SetToken(ctx context.Context, token types.Token) error {
	contract, err := k.addressCodec.StringToBytes(token.ContractAddress)
	if err != nil {
		return err
	}
	return k.Tokens.Set(ctx, contract, token)
}

func (k Keeper) GetAllTokens(ctx context.Context) ([]types.Token, error) {
	tokens := []types.Token{}
	err := k.Tokens.Walk(ctx, nil, func(key sdk.AccAddress, value types.Token) (stop bool, err error) {
		tokens = append(tokens, value)
		return false, nil
	})
	return tokens, err
}
