package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/gluon-zone/gluon/x/contracttoken/types"
)

func (m msgServer) CreateToken(ctx context.Context, msg *types.MsgCreateToken) (*types.MsgCreateTokenResponse, error) {
	params, err := m.Keeper.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// if TokenCreationGas is non-zero, consume the gas
	if params.TokenCreationGas != 0 {
		sdkCtx.GasMeter().ConsumeGas(params.TokenCreationGas, "consume token creation gas")
	}

	return &types.MsgCreateTokenResponse{}, nil

}
