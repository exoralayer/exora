//nolint:revive,stylecheck  // if we change the names of var-naming things here, we harm some kind of mapping.
package bindings

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	feetypes "github.com/neutron-org/neutron/v5/x/feerefunder/types"
	transferwrappertypes "github.com/neutron-org/neutron/v5/x/transfer/types"
)

// ProtobufAny is a hack-struct to serialize protobuf Any message into JSON object
type ProtobufAny struct {
	TypeURL string `json:"type_url"`
	Value   []byte `json:"value"`
}

// NeutronMsg is used like a sum type to hold one of custom Neutron messages.
// Follow https://github.com/neutron-org/neutron-sdk/blob/main/packages/neutron-sdk/src/bindings/msg.rs
// for more information.
type NeutronMsg struct {
	// Interchain txs
	SubmitTx                  *SubmitTx                  `json:"submit_tx,omitempty"`
	RegisterInterchainAccount *RegisterInterchainAccount `json:"register_interchain_account,omitempty"`

	// Token factory types
	/// Contracts can create denoms, namespaced under the contract's address.
	/// A contract may create any number of independent sub-denoms.
	CreateDenom *CreateDenom `json:"create_denom,omitempty"`
	/// Contracts can change the admin of a denom that they are the admin of.
	ChangeAdmin *ChangeAdmin `json:"change_admin,omitempty"`
	/// Contracts can mint native tokens for an existing factory denom
	/// that they are the admin of.
	MintTokens *MintTokens `json:"mint_tokens,omitempty"`
	/// Contracts can burn native tokens for an existing factory denom
	/// that they are the admin of.
	/// Currently, the burn from address must be the admin contract.
	BurnTokens *BurnTokens `json:"burn_tokens,omitempty"`
	/// Contracts can set before send hook for an existing factory denom
	///	that they are the admin of.
	///	Currently, the set before hook call should be performed from address that must be the admin contract.
	SetBeforeSendHook *SetBeforeSendHook `json:"set_before_send_hook,omitempty"`
	/// Force transferring of a specific denom is only allowed for the creator of the denom registered during CreateDenom.
	ForceTransfer *ForceTransfer `json:"force_transfer,omitempty"`
	/// Setting of metadata for a specific denom is only allowed for the admin of the denom.
	/// It allows the overwriting of the denom metadata in the bank module.
	SetDenomMetadata *SetDenomMetadata `json:"set_denom_metadata,omitempty"`

	// Transfer
	IBCTransfer *transferwrappertypes.MsgTransfer `json:"ibc_transfer,omitempty"`
}

// SubmitTx submits interchain transaction on a remote chain.
type SubmitTx struct {
	ConnectionId        string        `json:"connection_id"`
	InterchainAccountId string        `json:"interchain_account_id"`
	Msgs                []ProtobufAny `json:"msgs"`
	Memo                string        `json:"memo"`
	Timeout             uint64        `json:"timeout"`
	Fee                 feetypes.Fee  `json:"fee"`
}

// RegisterInterchainAccount creates account on remote chain.
type RegisterInterchainAccount struct {
	ConnectionId        string    `json:"connection_id"`
	InterchainAccountId string    `json:"interchain_account_id"`
	RegisterFee         sdk.Coins `json:"register_fee,omitempty"`
	Ordering            string    `json:"ordering,omitempty"`
}

// RegisterInterchainAccountResponse holds response for RegisterInterchainAccount.
type RegisterInterchainAccountResponse struct {
	ChannelId string `json:"channel_id"`
	PortId    string `json:"port_id"`
}

// CreateDenom creates a new factory denom, of denomination:
// factory/{creating contract address}/{Subdenom}
// Subdenom can be of length at most 44 characters, in [0-9a-zA-Z./]
// The (creating contract address, subdenom) pair must be unique.
// The created denom's admin is the creating contract address,
// but this admin can be changed using the ChangeAdmin binding.
type CreateDenom struct {
	Subdenom string `json:"subdenom"`
}

// ChangeAdmin changes the admin for a factory denom.
// If the NewAdminAddress is empty, the denom has no admin.
type ChangeAdmin struct {
	Denom           string `json:"denom"`
	NewAdminAddress string `json:"new_admin_address"`
}

type MintTokens struct {
	Denom         string   `json:"denom"`
	Amount        math.Int `json:"amount"`
	MintToAddress string   `json:"mint_to_address"`
}

type BurnTokens struct {
	Denom  string   `json:"denom"`
	Amount math.Int `json:"amount"`
	// BurnFromAddress must be set to "" for now.
	BurnFromAddress string `json:"burn_from_address"`
}

// SetBeforeSendHook Allowing to assign a CosmWasm contract to call with a BeforeSend hook for a specific denom is only
// allowed for the creator of the denom registered during CreateDenom.
type SetBeforeSendHook struct {
	Denom        string `json:"denom"`
	ContractAddr string `json:"contract_addr"`
}

// SetDenomMetadata is sets the denom's bank metadata
type SetDenomMetadata struct {
	banktypes.Metadata
}

// ForceTransfer forces transferring of a specific denom is only allowed for the creator of the denom registered during CreateDenom.
type ForceTransfer struct {
	Denom               string   `json:"denom"`
	Amount              math.Int `json:"amount"`
	TransferFromAddress string   `json:"transfer_from_address"`
	TransferToAddress   string   `json:"transfer_to_address"`
}
