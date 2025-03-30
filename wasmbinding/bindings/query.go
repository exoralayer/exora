package bindings

import (
	"github.com/cosmos/cosmos-sdk/types/query"

	contractmanagertypes "github.com/neutron-org/neutron/v5/x/contractmanager/types"
	feerefundertypes "github.com/neutron-org/neutron/v5/x/feerefunder/types"
)

// NeutronQuery contains neutron custom queries.
type NeutronQuery struct {
	// Contractmanager queries
	// Query all failures for address
	Failures *Failures `json:"failures,omitempty"`
	// MinIbcFee
	MinIbcFee *QueryMinIbcFeeRequest `json:"min_ibc_fee,omitempty"`
	// Interchain account address for specified ConnectionID and OwnerAddress
	InterchainAccountAddress *QueryInterchainAccountAddressRequest `json:"interchain_account_address,omitempty"`
	// Token Factory queries
	// Given a subdenom minted by a contract via `NeutronMsg::MintTokens`,
	// returns the full denom as used by `BankMsg::Send`.
	FullDenom *FullDenom `json:"full_denom,omitempty"`
	// Returns the admin of a denom, if the denom is a Token Factory denom.
	DenomAdmin *DenomAdmin `json:"denom_admin,omitempty"`
	// Returns the before send hook if it was set before
	BeforeSendHook *BeforeSendHook `json:"before_send_hook,omitempty"`
}

type Failures struct {
	Address    string             `json:"address"`
	Pagination *query.PageRequest `json:"pagination,omitempty"`
}

type FailuresResponse struct {
	Failures []contractmanagertypes.Failure `json:"failures"`
}

type QueryMinIbcFeeRequest struct{}

type QueryMinIbcFeeResponse struct {
	MinFee feerefundertypes.Fee `json:"min_fee"`
}

type QueryInterchainAccountAddressRequest struct {
	// owner_address is the owner of the interchain account on the controller chain
	OwnerAddress string `json:"owner_address,omitempty"`
	// interchain_account_id is an identifier of your interchain account from which you want to execute msgs
	InterchainAccountID string `json:"interchain_account_id,omitempty"`
	// connection_id is an IBC connection identifier between Neutron and remote chain
	ConnectionID string `json:"connection_id,omitempty"`
}

// Query response for an interchain account address
type QueryInterchainAccountAddressResponse struct {
	// The corresponding interchain account address on the host chain
	InterchainAccountAddress string `json:"interchain_account_address,omitempty"`
}

type FullDenom struct {
	CreatorAddr string `json:"creator_addr"`
	Subdenom    string `json:"subdenom"`
}

type FullDenomResponse struct {
	Denom string `json:"denom"`
}

type DenomAdmin struct {
	Subdenom string `json:"subdenom"`
}

type DenomAdminResponse struct {
	Admin string `json:"admin"`
}

type BeforeSendHook struct {
	Denom string `json:"denom"`
}

type BeforeSendHookResponse struct {
	ContractAddr string `json:"contract_addr"`
}
