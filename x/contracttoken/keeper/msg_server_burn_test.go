package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gluon-zone/gluon/x/contracttoken/keeper"
	"github.com/gluon-zone/gluon/x/contracttoken/types"
	"github.com/stretchr/testify/require"
)

func TestBurn(t *testing.T) {
	// Setup test keeper
	fixture := initFixture(t)
	msgServer := keeper.NewMsgServerImpl(fixture.keeper)

	// Generate test addresses
	contractAddr := sdk.AccAddress("test-contract-addr")

	// Register token
	token := types.Token{
		ContractAddress:       contractAddr.String(),
		BeforeSendHookEnabled: false,
	}
	err := fixture.keeper.Tokens.Set(fixture.ctx, contractAddr, token)
	require.NoError(t, err)

	// Set up mock expectations
	denom, err := types.GetTokenDenom(contractAddr.String())
	require.NoError(t, err)
	coin := sdk.NewCoin(denom, math.NewInt(100))
	coins := sdk.NewCoins(coin)

	fixture.mockBank.EXPECT().
		SendCoinsFromAccountToModule(fixture.ctx, contractAddr, types.ModuleName, coins).
		Return(nil)

	fixture.mockBank.EXPECT().
		BurnCoins(fixture.ctx, types.ModuleName, coins).
		Return(nil)

	tests := []struct {
		name          string
		msg           *types.MsgBurn
		expectError   bool
		errorContains string
	}{
		{
			name: "success: burn token",
			msg: &types.MsgBurn{
				ContractAddress: contractAddr.String(),
				Amount:          math.NewInt(100),
			},
			expectError: false,
		},
		{
			name: "error: invalid contract address",
			msg: &types.MsgBurn{
				ContractAddress: "invalid-address",
				Amount:          math.NewInt(100),
			},
			expectError:   true,
			errorContains: "invalid contract address",
		},
		{
			name: "error: token does not exist",
			msg: &types.MsgBurn{
				ContractAddress: sdk.AccAddress("non-existent-addr").String(),
				Amount:          math.NewInt(100),
			},
			expectError:   true,
			errorContains: "token does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := msgServer.Burn(fixture.ctx, tt.msg)
			if tt.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errorContains)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, resp)
		})
	}
}
