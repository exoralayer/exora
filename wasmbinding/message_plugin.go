package wasmbinding

import (
	"encoding/json"
	"fmt"

	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	channeltypesv8 "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"

	//nolint:staticcheck

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/CosmWasm/wasmvm/v2/types"

	"gluon/wasmbinding/bindings"

	contractmanagerkeeper "github.com/neutron-org/neutron/v5/x/contractmanager/keeper"
	contractmanagertypes "github.com/neutron-org/neutron/v5/x/contractmanager/types"
	ictxkeeper "github.com/neutron-org/neutron/v5/x/interchaintxs/keeper"
	ictxtypes "github.com/neutron-org/neutron/v5/x/interchaintxs/types"
	tokenfactorykeeper "github.com/neutron-org/neutron/v5/x/tokenfactory/keeper"
	tokenfactorytypes "github.com/neutron-org/neutron/v5/x/tokenfactory/types"
	transferwrapperkeeper "github.com/neutron-org/neutron/v5/x/transfer/keeper"
	transferwrappertypes "github.com/neutron-org/neutron/v5/x/transfer/types"
)

type CustomMessenger struct {
	Keeper                     ictxkeeper.Keeper
	Wrapped                    wasmkeeper.Messenger
	ContractmanagerMsgServer   contractmanagertypes.MsgServer
	ContractmanagerQueryServer contractmanagertypes.QueryServer
	Ictxmsgserver              ictxtypes.MsgServer
	TokenFactory               *tokenfactorykeeper.Keeper
	transferKeeper             *transferwrapperkeeper.KeeperTransferWrapper
}

var _ wasmkeeper.Messenger = (*CustomMessenger)(nil)

func CustomMessageDecorator(
	contractmanagerKeeper *contractmanagerkeeper.Keeper,
	ictx *ictxkeeper.Keeper,
	tokenFactoryKeeper *tokenfactorykeeper.Keeper,
	transferKeeper *transferwrapperkeeper.KeeperTransferWrapper,
) func(messenger wasmkeeper.Messenger) wasmkeeper.Messenger {
	return func(old wasmkeeper.Messenger) wasmkeeper.Messenger {
		return &CustomMessenger{
			Keeper:                     *ictx,
			Wrapped:                    old,
			ContractmanagerMsgServer:   contractmanagerkeeper.NewMsgServerImpl(*contractmanagerKeeper),
			ContractmanagerQueryServer: contractmanagerkeeper.NewQueryServerImpl(*contractmanagerKeeper),
			Ictxmsgserver:              ictxkeeper.NewMsgServerImpl(*ictx),
			TokenFactory:               tokenFactoryKeeper,
			transferKeeper:             transferKeeper,
		}
	}
}

func (m *CustomMessenger) DispatchMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg) ([]sdk.Event, [][]byte, [][]*types.Any, error) {
	// Return early if msg.Custom is nil
	if msg.Custom == nil {
		return m.Wrapped.DispatchMsg(ctx, contractAddr, contractIBCPortID, msg)
	}

	var contractMsg bindings.NeutronMsg
	if err := json.Unmarshal(msg.Custom, &contractMsg); err != nil {
		ctx.Logger().Debug("json.Unmarshal: failed to decode incoming custom cosmos message",
			"from_address", contractAddr.String(),
			"message", string(msg.Custom),
			"error", err,
		)
		return nil, nil, nil, errors.Wrap(err, "failed to decode incoming custom cosmos message")
	}

	// Dispatch the message based on its type by checking each possible field
	if contractMsg.SubmitTx != nil {
		return m.submitTx(ctx, contractAddr, contractMsg.SubmitTx)
	}
	if contractMsg.RegisterInterchainAccount != nil {
		return m.registerInterchainAccount(ctx, contractAddr, contractMsg.RegisterInterchainAccount)
	}
	if contractMsg.IBCTransfer != nil {
		return m.ibcTransfer(ctx, contractAddr, *contractMsg.IBCTransfer)
	}

	if contractMsg.CreateDenom != nil {
		return m.createDenom(ctx, contractAddr, contractMsg.CreateDenom)
	}
	if contractMsg.MintTokens != nil {
		return m.mintTokens(ctx, contractAddr, contractMsg.MintTokens)
	}
	if contractMsg.SetBeforeSendHook != nil {
		return m.setBeforeSendHook(ctx, contractAddr, contractMsg.SetBeforeSendHook)
	}
	if contractMsg.ChangeAdmin != nil {
		return m.changeAdmin(ctx, contractAddr, contractMsg.ChangeAdmin)
	}
	if contractMsg.BurnTokens != nil {
		return m.burnTokens(ctx, contractAddr, contractMsg.BurnTokens)
	}
	if contractMsg.ForceTransfer != nil {
		return m.forceTransfer(ctx, contractAddr, contractMsg.ForceTransfer)
	}
	if contractMsg.SetDenomMetadata != nil {
		return m.setDenomMetadata(ctx, contractAddr, contractMsg.SetDenomMetadata)
	}

	// If none of the conditions are met, forward the message to the wrapped handler
	return m.Wrapped.DispatchMsg(ctx, contractAddr, contractIBCPortID, msg)
}

func (m *CustomMessenger) ibcTransfer(ctx sdk.Context, contractAddr sdk.AccAddress, ibcTransferMsg transferwrappertypes.MsgTransfer) ([]sdk.Event, [][]byte, [][]*types.Any, error) {
	ibcTransferMsg.Sender = contractAddr.String()

	response, err := m.transferKeeper.Transfer(ctx, &ibcTransferMsg)
	if err != nil {
		ctx.Logger().Debug("transferServer.Transfer: failed to transfer",
			"from_address", contractAddr.String(),
			"msg", ibcTransferMsg,
			"error", err,
		)
		return nil, nil, nil, errors.Wrap(err, "failed to execute IBCTransfer")
	}

	data, err := json.Marshal(response)
	if err != nil {
		ctx.Logger().Error("json.Marshal: failed to marshal MsgTransferResponse response to JSON",
			"from_address", contractAddr.String(),
			"msg", response,
			"error", err,
		)
		return nil, nil, nil, errors.Wrap(err, "marshal json failed")
	}

	ctx.Logger().Debug("ibcTransferMsg completed",
		"from_address", contractAddr.String(),
		"msg", ibcTransferMsg,
	)

	anyResp, err := types.NewAnyWithValue(response)
	if err != nil {
		return nil, nil, nil, errors.Wrapf(err, "failed to convert {%T} to Any", response)
	}
	msgResponses := [][]*types.Any{{anyResp}}
	return nil, [][]byte{data}, msgResponses, nil
}

func (m *CustomMessenger) submitTx(ctx sdk.Context, contractAddr sdk.AccAddress, submitTx *bindings.SubmitTx) ([]sdk.Event, [][]byte, [][]*types.Any, error) {
	response, err := m.performSubmitTx(ctx, contractAddr, submitTx)
	if err != nil {
		ctx.Logger().Debug("performSubmitTx: failed to submit interchain transaction",
			"from_address", contractAddr.String(),
			"connection_id", submitTx.ConnectionId,
			"interchain_account_id", submitTx.InterchainAccountId,
			"error", err,
		)
		return nil, nil, nil, errors.Wrap(err, "failed to submit interchain transaction")
	}

	data, err := json.Marshal(response)
	if err != nil {
		ctx.Logger().Error("json.Marshal: failed to marshal submitTx response to JSON",
			"from_address", contractAddr.String(),
			"connection_id", submitTx.ConnectionId,
			"interchain_account_id", submitTx.InterchainAccountId,
			"error", err,
		)
		return nil, nil, nil, errors.Wrap(err, "marshal json failed")
	}

	ctx.Logger().Debug("interchain transaction submitted",
		"from_address", contractAddr.String(),
		"connection_id", submitTx.ConnectionId,
		"interchain_account_id", submitTx.InterchainAccountId,
	)

	anyResp, err := types.NewAnyWithValue(response)
	if err != nil {
		return nil, nil, nil, errors.Wrapf(err, "failed to convert {%T} to Any", response)
	}
	msgResponses := [][]*types.Any{{anyResp}}
	return nil, [][]byte{data}, msgResponses, nil
}

// createDenom creates a new token denom
func (m *CustomMessenger) createDenom(ctx sdk.Context, contractAddr sdk.AccAddress, createDenom *bindings.CreateDenom) ([]sdk.Event, [][]byte, [][]*types.Any, error) {
	err := PerformCreateDenom(m.TokenFactory, ctx, contractAddr, createDenom)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "perform create denom")
	}
	return nil, nil, nil, nil
}

// PerformCreateDenom is used with createDenom to create a token denom; validates the msgCreateDenom.
func PerformCreateDenom(f *tokenfactorykeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, createDenom *bindings.CreateDenom) error {
	msgServer := tokenfactorykeeper.NewMsgServerImpl(*f)

	msgCreateDenom := tokenfactorytypes.NewMsgCreateDenom(contractAddr.String(), createDenom.Subdenom)

	// Create denom
	_, err := msgServer.CreateDenom(
		ctx,
		msgCreateDenom,
	)
	if err != nil {
		return errors.Wrap(err, "creating denom")
	}
	return nil
}

// createDenom forces a transfer of a tokenFactory token
func (m *CustomMessenger) forceTransfer(ctx sdk.Context, contractAddr sdk.AccAddress, forceTransfer *bindings.ForceTransfer) ([]sdk.Event, [][]byte, [][]*types.Any, error) {
	err := PerformForceTransfer(m.TokenFactory, ctx, contractAddr, forceTransfer)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "perform force transfer")
	}
	return nil, nil, nil, nil
}

// PerformForceTransfer is used with forceTransfer to force a tokenfactory token transfer; validates the msgForceTransfer.
func PerformForceTransfer(f *tokenfactorykeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, forceTransfer *bindings.ForceTransfer) error {
	msgServer := tokenfactorykeeper.NewMsgServerImpl(*f)

	msgForceTransfer := tokenfactorytypes.NewMsgForceTransfer(contractAddr.String(), sdk.NewInt64Coin(forceTransfer.Denom, forceTransfer.Amount.Int64()), forceTransfer.TransferFromAddress, forceTransfer.TransferToAddress)

	// Force Transfer
	_, err := msgServer.ForceTransfer(
		ctx,
		msgForceTransfer,
	)
	if err != nil {
		return errors.Wrap(err, "forcing transfer")
	}
	return nil
}

// setDenomMetadata sets a metadata for a tokenfactory denom
func (m *CustomMessenger) setDenomMetadata(ctx sdk.Context, contractAddr sdk.AccAddress, setDenomMetadata *bindings.SetDenomMetadata) ([]sdk.Event, [][]byte, [][]*types.Any, error) {
	err := PerformSetDenomMetadata(m.TokenFactory, ctx, contractAddr, setDenomMetadata)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "perform set denom metadata")
	}
	return nil, nil, nil, nil
}

// PerformSetDenomMetadata is used with setDenomMetadata to set a metadata for a tokenfactory denom; validates the msgSetDenomMetadata.
func PerformSetDenomMetadata(f *tokenfactorykeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, setDenomMetadata *bindings.SetDenomMetadata) error {
	msgServer := tokenfactorykeeper.NewMsgServerImpl(*f)

	msgSetDenomMetadata := tokenfactorytypes.NewMsgSetDenomMetadata(contractAddr.String(), setDenomMetadata.Metadata)

	// Set denom metadata
	_, err := msgServer.SetDenomMetadata(
		ctx,
		msgSetDenomMetadata,
	)
	if err != nil {
		return errors.Wrap(err, "setting denom metadata")
	}
	return nil
}

// mintTokens mints tokens of a specified denom to an address.
func (m *CustomMessenger) mintTokens(ctx sdk.Context, contractAddr sdk.AccAddress, mint *bindings.MintTokens) ([]sdk.Event, [][]byte, [][]*types.Any, error) {
	err := PerformMint(m.TokenFactory, ctx, contractAddr, mint)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "perform mint")
	}
	return nil, nil, nil, nil
}

// setBeforeSendHook sets before send hook for a specified denom.
func (m *CustomMessenger) setBeforeSendHook(ctx sdk.Context, contractAddr sdk.AccAddress, set *bindings.SetBeforeSendHook) ([]sdk.Event, [][]byte, [][]*types.Any, error) {
	err := PerformSetBeforeSendHook(m.TokenFactory, ctx, contractAddr, set)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to perform set before send hook")
	}
	return nil, nil, nil, nil
}

// PerformMint used with mintTokens to validate the mint message and mint through token factory.
func PerformMint(f *tokenfactorykeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, mint *bindings.MintTokens) error {
	rcpt, err := parseAddress(mint.MintToAddress)
	if err != nil {
		return err
	}

	coin := sdk.Coin{Denom: mint.Denom, Amount: mint.Amount}
	sdkMsg := tokenfactorytypes.NewMsgMintTo(contractAddr.String(), coin, rcpt.String())

	// Mint through token factory / message server
	msgServer := tokenfactorykeeper.NewMsgServerImpl(*f)
	_, err = msgServer.Mint(ctx, sdkMsg)
	if err != nil {
		return errors.Wrap(err, "minting coins from message")
	}

	return nil
}

func PerformSetBeforeSendHook(f *tokenfactorykeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, set *bindings.SetBeforeSendHook) error {
	sdkMsg := tokenfactorytypes.NewMsgSetBeforeSendHook(contractAddr.String(), set.Denom, set.ContractAddr)

	// SetBeforeSendHook through token factory / message server
	msgServer := tokenfactorykeeper.NewMsgServerImpl(*f)
	_, err := msgServer.SetBeforeSendHook(ctx, sdkMsg)
	if err != nil {
		return errors.Wrap(err, "set before send from message")
	}

	return nil
}

// changeAdmin changes the admin.
func (m *CustomMessenger) changeAdmin(ctx sdk.Context, contractAddr sdk.AccAddress, changeAdmin *bindings.ChangeAdmin) ([]sdk.Event, [][]byte, [][]*types.Any, error) {
	err := ChangeAdmin(m.TokenFactory, ctx, contractAddr, changeAdmin)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to change admin")
	}

	return nil, nil, nil, nil
}

// ChangeAdmin is used with changeAdmin to validate changeAdmin messages and to dispatch.
func ChangeAdmin(f *tokenfactorykeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, changeAdmin *bindings.ChangeAdmin) error {
	newAdminAddr, err := parseAddress(changeAdmin.NewAdminAddress)
	if err != nil {
		return err
	}

	changeAdminMsg := tokenfactorytypes.NewMsgChangeAdmin(contractAddr.String(), changeAdmin.Denom, newAdminAddr.String())

	msgServer := tokenfactorykeeper.NewMsgServerImpl(*f)
	_, err = msgServer.ChangeAdmin(ctx, changeAdminMsg)
	if err != nil {
		return errors.Wrap(err, "failed changing admin from message")
	}
	return nil
}

// burnTokens burns tokens.
func (m *CustomMessenger) burnTokens(ctx sdk.Context, contractAddr sdk.AccAddress, burn *bindings.BurnTokens) ([]sdk.Event, [][]byte, [][]*types.Any, error) {
	err := PerformBurn(m.TokenFactory, ctx, contractAddr, burn)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "perform burn")
	}

	return nil, nil, nil, nil
}

// PerformBurn performs token burning after validating tokenBurn message.
func PerformBurn(f *tokenfactorykeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, burn *bindings.BurnTokens) error {
	coin := sdk.Coin{Denom: burn.Denom, Amount: burn.Amount}
	sdkMsg := tokenfactorytypes.NewMsgBurnFrom(contractAddr.String(), coin, burn.BurnFromAddress)

	// Burn through token factory / message server
	msgServer := tokenfactorykeeper.NewMsgServerImpl(*f)
	_, err := msgServer.Burn(ctx, sdkMsg)
	if err != nil {
		return errors.Wrap(err, "burning coins from message")
	}

	return nil
}

// GetFullDenom is a function, not method, so the message_plugin can use it
func GetFullDenom(contract, subDenom string) (string, error) {
	// Address validation
	if _, err := parseAddress(contract); err != nil {
		return "", err
	}

	fullDenom, err := tokenfactorytypes.GetTokenDenom(contract, subDenom)
	if err != nil {
		return "", errors.Wrap(err, "validate sub-denom")
	}

	return fullDenom, nil
}

// parseAddress parses address from bech32 string and verifies its format.
func parseAddress(addr string) (sdk.AccAddress, error) {
	parsed, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return nil, errors.Wrap(err, "address from bech32")
	}

	err = sdk.VerifyAddressFormat(parsed)
	if err != nil {
		return nil, errors.Wrap(err, "verify address format")
	}

	return parsed, nil
}

func (m *CustomMessenger) performSubmitTx(ctx sdk.Context, contractAddr sdk.AccAddress, submitTx *bindings.SubmitTx) (*ictxtypes.MsgSubmitTxResponse, error) {
	tx := ictxtypes.MsgSubmitTx{
		FromAddress:         contractAddr.String(),
		ConnectionId:        submitTx.ConnectionId,
		Memo:                submitTx.Memo,
		InterchainAccountId: submitTx.InterchainAccountId,
		Timeout:             submitTx.Timeout,
		Fee:                 submitTx.Fee,
	}
	for _, msg := range submitTx.Msgs {
		tx.Msgs = append(tx.Msgs, &types.Any{
			TypeUrl: msg.TypeURL,
			Value:   msg.Value,
		})
	}

	response, err := m.Ictxmsgserver.SubmitTx(ctx, &tx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to submit interchain transaction")
	}

	return response, nil
}

func (m *CustomMessenger) registerInterchainAccount(ctx sdk.Context, contractAddr sdk.AccAddress, reg *bindings.RegisterInterchainAccount) ([]sdk.Event, [][]byte, [][]*types.Any, error) {
	response, err := m.performRegisterInterchainAccount(ctx, contractAddr, reg)
	if err != nil {
		ctx.Logger().Debug("performRegisterInterchainAccount: failed to register interchain account",
			"from_address", contractAddr.String(),
			"connection_id", reg.ConnectionId,
			"interchain_account_id", reg.InterchainAccountId,
			"error", err,
		)
		return nil, nil, nil, errors.Wrap(err, "failed to register interchain account")
	}

	data, err := json.Marshal(response)
	if err != nil {
		ctx.Logger().Error("json.Marshal: failed to marshal register interchain account response to JSON",
			"from_address", contractAddr.String(),
			"connection_id", reg.ConnectionId,
			"interchain_account_id", reg.InterchainAccountId,
			"error", err,
		)
		return nil, nil, nil, errors.Wrap(err, "marshal json failed")
	}

	ctx.Logger().Debug("registered interchain account",
		"from_address", contractAddr.String(),
		"connection_id", reg.ConnectionId,
		"interchain_account_id", reg.InterchainAccountId,
	)

	anyResp, err := types.NewAnyWithValue(response)
	if err != nil {
		return nil, nil, nil, errors.Wrapf(err, "failed to convert {%T} to Any", response)
	}
	msgResponses := [][]*types.Any{{anyResp}}
	return nil, [][]byte{data}, msgResponses, nil
}

func (m *CustomMessenger) performRegisterInterchainAccount(ctx sdk.Context, contractAddr sdk.AccAddress, reg *bindings.RegisterInterchainAccount) (*ictxtypes.MsgRegisterInterchainAccountResponse, error) {
	// parse incoming ordering. If nothing passed, use ORDERED by default
	var orderValue channeltypes.Order
	if reg.Ordering == "" {
		orderValue = channeltypes.ORDERED
	} else {
		orderValueInt, ok := channeltypes.Order_value[reg.Ordering]

		if !ok {
			return nil, fmt.Errorf("failed to register interchain account: incorrect order value passed: %s", reg.Ordering)
		}
		orderValue = channeltypes.Order(orderValueInt)
	}

	msg := ictxtypes.MsgRegisterInterchainAccount{
		FromAddress:         contractAddr.String(),
		ConnectionId:        reg.ConnectionId,
		InterchainAccountId: reg.InterchainAccountId,
		RegisterFee:         getRegisterFee(reg.RegisterFee),
		Ordering:            channeltypesv8.Order(orderValue),
	}

	response, err := m.Ictxmsgserver.RegisterInterchainAccount(ctx, &msg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to register interchain account")
	}

	return response, nil
}

func getRegisterFee(fee sdk.Coins) sdk.Coins {
	if fee == nil {
		return make(sdk.Coins, 0)
	}
	return fee
}
