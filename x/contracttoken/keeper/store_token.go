package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/gluon-zone/gluon/x/contracttoken/types"
)

func (k Keeper) HasToken(ctx context.Context, contract sdk.AccAddress) bool {
	_, err := k.GetToken(ctx, contract)
	if err != nil {
		return false
	}
	return true
}

func (k Keeper) GetToken(ctx context.Context, contract sdk.AccAddress) (types.Token, error) {
	token, err := k.ContractTokens.Get(ctx, contract)
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
	return k.ContractTokens.Set(ctx, contract, token)
}
