package keeper

import (
	"context"
	"encoding/json"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/gluon-zone/gluon/x/contracttoken/types"
)

func (k Keeper) BeforeSendHook(ctx context.Context, from sdk.AccAddress, to sdk.AccAddress, amount sdk.Coins) error {
	for _, coin := range amount {
		contract, err := types.ParseTokenDenom(coin.Denom)
		if err != nil {
			continue
		}

		token, err, found := k.GetToken(ctx, contract)
		if !found {
			continue
		}
		if err != nil {
			return err
		}

		if token.BeforeSendHookEnabled {
			err = k.ExecuteBeforeSend(ctx, contract, from, to, coin.Amount)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (k Keeper) ExecuteBeforeSend(ctx context.Context, contract sdk.AccAddress, from sdk.AccAddress, to sdk.AccAddress, amount math.Int) error {
	msg := BeforeSendMsg{
		From:   from.String(),
		To:     to.String(),
		Amount: amount,
	}
	msgJson, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	_, err = k.wasmKeeper.Execute(sdkCtx, contract, from, msgJson, sdk.NewCoins())
	if err != nil {
		return err
	}

	return nil
}

type BeforeSendMsg struct {
	From   string   `json:"from"`
	To     string   `json:"to"`
	Amount math.Int `json:"amount"`
}
