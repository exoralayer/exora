package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/exoralayer/exora/x/contracttoken/types"
)

func (m msgServer) CreateToken(ctx context.Context, msg *types.MsgCreateToken) (*types.MsgCreateTokenResponse, error) {
	contractAddr, err := m.addressCodec.StringToBytes(msg.ContractAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid contract address")
	}

	// Check if the contract address is already in use
	found := m.Keeper.HasToken(ctx, contractAddr)
	if found {
		return nil, errorsmod.Wrap(types.ErrTokenExists, "contract address already in use")
	}

	err = m.Keeper.SetToken(ctx, types.Token{
		ContractAddress:       msg.ContractAddress,
		BeforeSendHookEnabled: msg.BeforeSendHookEnabled,
	})
	if err != nil {
		return nil, err
	}

	// Consume gas for preventing spam
	params, err := m.Keeper.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if params.TokenCreationGas != 0 {
		sdkCtx.GasMeter().ConsumeGas(params.TokenCreationGas, "consume token creation gas")
	}

	// Denom
	denom, err := types.GetTokenDenom(msg.ContractAddress)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreateTokenResponse{
		Denom: denom,
	}, nil

}
