package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) BeforeSend(ctx context.Context, contract sdk.AccAddress, from sdk.AccAddress, to sdk.AccAddress, amount sdk.Coins) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	_, err := k.wasmKeeper.Execute(sdkCtx, contract, from, []byte{}, sdk.NewCoins())
	if err != nil {
		return err
	}

	return nil
}
