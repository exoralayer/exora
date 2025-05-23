package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/exoralayer/exora/x/contracttoken/types"
)

func (m msgServer) UpdateToken(ctx context.Context, msg *types.MsgUpdateToken) (*types.MsgUpdateTokenResponse, error) {
	contractAddr, err := sdk.AccAddressFromBech32(msg.ContractAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid contract address")
	}

	// Check if the contract address is already in use
	found := m.Keeper.HasToken(ctx, contractAddr)
	if !found {
		return nil, errorsmod.Wrap(types.ErrTokenDoesNotExist, "token does not exist")
	}

	err = m.Keeper.SetToken(ctx, types.Token{
		ContractAddress:       msg.ContractAddress,
		BeforeSendHookEnabled: msg.BeforeSendHookEnabled,
	})
	if err != nil {
		return nil, err
	}

	return &types.MsgUpdateTokenResponse{}, nil
}
