package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gluon-zone/gluon/x/contracttoken/keeper"
	"github.com/gluon-zone/gluon/x/contracttoken/types"
	"github.com/stretchr/testify/require"
)

func TestCreateToken(t *testing.T) {
	// Generate test addresses
	contractAddr := sdk.AccAddress("test-contract-addr")

	tests := []struct {
		name          string
		msg           *types.MsgCreateToken
		setupState    func(*fixture)
		expectError   bool
		errorContains string
	}{
		{
			name: "success: create token",
			msg: &types.MsgCreateToken{
				ContractAddress:       contractAddr.String(),
				BeforeSendHookEnabled: false,
			},
			setupState:  func(f *fixture) {},
			expectError: false,
		},
		{
			name: "error: invalid contract address",
			msg: &types.MsgCreateToken{
				ContractAddress:       "invalid-address",
				BeforeSendHookEnabled: false,
			},
			setupState:    func(f *fixture) {},
			expectError:   true,
			errorContains: "invalid contract address",
		},
		{
			name: "error: token already exists",
			msg: &types.MsgCreateToken{
				ContractAddress:       contractAddr.String(),
				BeforeSendHookEnabled: false,
			},
			setupState: func(f *fixture) {
				// Create a token first
				token := types.Token{
					ContractAddress:       contractAddr.String(),
					BeforeSendHookEnabled: false,
				}
				err := f.keeper.ContractTokens.Set(f.ctx, contractAddr, token)
				require.NoError(t, err)
			},
			expectError:   true,
			errorContains: "contract address already in use",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create new fixture for this test case
			fixture := initFixture(t)
			msgServer := keeper.NewMsgServerImpl(fixture.keeper)

			// Setup mock and state for this test case
			tt.setupState(fixture)

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
