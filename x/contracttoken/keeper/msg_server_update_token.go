package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	"github.com/gluon-zone/gluon/x/contracttoken/types"
)

func (m msgServer) UpdateToken(ctx context.Context, msg *types.MsgUpdateToken) (*types.MsgUpdateTokenResponse, error) {
	contractAddr, err := m.addressCodec.StringToBytes(msg.ContractAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid contract address")
	}

	// Check if the contract address is already in use
	found, err := m.Keeper.HasToken(ctx, contractAddr)
	if err != nil {
		return nil, err
	}
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
