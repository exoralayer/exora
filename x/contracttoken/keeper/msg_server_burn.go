package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/gluon-zone/gluon/x/contracttoken/types"
)

func (m msgServer) Burn(ctx context.Context, msg *types.MsgBurn) (*types.MsgBurnResponse, error) {
	contractAddr, err := m.addressCodec.StringToBytes(msg.ContractAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid contract address")
	}

	// Check if the contract address is already in use
	found := m.Keeper.HasToken(ctx, contractAddr)
	if !found {
		return nil, errorsmod.Wrap(types.ErrTokenDoesNotExist, "token does not exist")
	}

	// Denom
	denom, err := types.GetTokenDenom(msg.ContractAddress)
	if err != nil {
		return nil, err
	}
	coin := sdk.NewCoin(denom, msg.Amount)
	coins := sdk.NewCoins(coin)

	// Send to the module address
	err = m.Keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, contractAddr, types.ModuleName, coins)
	if err != nil {
		return nil, err
	}

	// Burn
	err = m.Keeper.bankKeeper.BurnCoins(ctx, types.ModuleName, coins)
	if err != nil {
		return nil, err
	}

	return &types.MsgBurnResponse{}, nil
}
