package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	"github.com/gluon-zone/gluon/x/contracttoken/types"
)

func (m msgServer) SetMetadata(ctx context.Context, msg *types.MsgSetMetadata) (*types.MsgSetMetadataResponse, error) {
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
	if denom != msg.Metadata.Base {
		return nil, errorsmod.Wrap(types.ErrInvalidDenom, "base denom in the metadata does not match the denom")
	}

	m.Keeper.bankKeeper.SetDenomMetaData(ctx, msg.Metadata)

	return &types.MsgSetMetadataResponse{}, nil
}
