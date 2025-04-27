package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/exoralayer/exora/x/contracttoken/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) Token(ctx context.Context, req *types.QueryTokenRequest) (*types.QueryTokenResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	contractAddr, err := sdk.AccAddressFromBech32(req.ContractAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid contract address")
	}

	token, err := q.k.GetToken(ctx, contractAddr)
	if err != nil {
		return nil, err
	}
	return &types.QueryTokenResponse{Token: token}, nil
}

func (q queryServer) Tokens(ctx context.Context, req *types.QueryTokensRequest) (*types.QueryTokensResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	tokens, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.Tokens,
		req.Pagination,
		func(key sdk.AccAddress, value types.Token) (types.Token, error) {
			return value, nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryTokensResponse{Tokens: tokens, Pagination: pageRes}, nil
}
