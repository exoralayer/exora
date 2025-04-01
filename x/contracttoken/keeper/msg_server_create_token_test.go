package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gluon-zone/gluon/x/contracttoken/keeper"
	"github.com/gluon-zone/gluon/x/contracttoken/types"
	"github.com/stretchr/testify/require"
)

func TestCreateToken(t *testing.T) {
	// Setup test keeper
	fixture := initFixture(t)
	msgServer := keeper.NewMsgServerImpl(fixture.keeper)

	// Generate test addresses
	contractAddr := sdk.AccAddress("test-contract-addr")

	// Set up mock expectations
	fixture.mockWasm.EXPECT().
		HasContractInfo(fixture.ctx, contractAddr).
		Return(true)

	tests := []struct {
		name          string
		msg           *types.MsgCreateToken
		expectError   bool
		errorContains string
	}{
		{
			name: "success: create token",
			msg: &types.MsgCreateToken{
				ContractAddress:       contractAddr.String(),
				BeforeSendHookEnabled: false,
			},
			expectError: false,
		},
		{
			name: "error: invalid contract address",
			msg: &types.MsgCreateToken{
				ContractAddress:       "invalid-address",
				BeforeSendHookEnabled: false,
			},
			expectError:   true,
			errorContains: "invalid contract address",
		},
		{
			name: "error: not a contract",
			msg: &types.MsgCreateToken{
				ContractAddress:       sdk.AccAddress("not-contract").String(),
				BeforeSendHookEnabled: false,
			},
			expectError:   true,
			errorContains: "contract address is not a contract",
		},
		{
			name: "error: token already exists",
			msg: &types.MsgCreateToken{
				ContractAddress:       contractAddr.String(),
				BeforeSendHookEnabled: false,
			},
			expectError:   true,
			errorContains: "contract address already in use",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := msgServer.CreateToken(fixture.ctx, tt.msg)
			if tt.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errorContains)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, resp)

			// Verify token was created
			token, err := fixture.keeper.ContractTokens.Get(fixture.ctx, contractAddr)
			require.NoError(t, err)
			require.Equal(t, tt.msg.ContractAddress, token.ContractAddress)
			require.Equal(t, tt.msg.BeforeSendHookEnabled, token.BeforeSendHookEnabled)
		})
	}
}
