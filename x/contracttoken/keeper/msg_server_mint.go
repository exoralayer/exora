package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gluon-zone/gluon/x/contracttoken/types"
)

func (m msgServer) Mint(ctx context.Context, msg *types.MsgMint) (*types.MsgMintResponse, error) {
	contractAddr, err := m.addressCodec.StringToBytes(msg.ContractAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid contract address")
	}

	// Check if the contract address is already in use
	found := m.Keeper.HasToken(ctx, contractAddr)
	if !found {
		return nil, errorsmod.Wrap(types.ErrTokenDoesNotExist, "token does not exist")
	}

	var recipient sdk.AccAddress
	if msg.Recipient != "" {
		recipient, err = m.addressCodec.StringToBytes(msg.Recipient)
		if err != nil {
			return nil, errorsmod.Wrap(err, "invalid recipient address")
		}
	} else {
		recipient = contractAddr
	}

	// Denom
	denom, err := types.GetTokenDenom(msg.ContractAddress)
	if err != nil {
		return nil, err
	}
	coin := sdk.NewCoin(denom, msg.Amount)
	coins := sdk.NewCoins(coin)

	// Mint
	err = m.Keeper.bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	if err != nil {
		return nil, err
	}

	// Send to the recipient
	err = m.Keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, coins)
	if err != nil {
		return nil, err
	}

	return &types.MsgMintResponse{}, nil
}
