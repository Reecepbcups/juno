package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgRegisterRevenue{}
	_ sdk.Msg = &MsgCancelRevenue{}
	_ sdk.Msg = &MsgUpdateRevenue{}
)

const (
	TypeMsgRegisterRevenue = "register_revenue"
	TypeMsgCancelRevenue   = "cancel_revenue"
	TypeMsgUpdateRevenue   = "update_revenue"
)

// NewMsgRegisterRevenue creates new instance of MsgRegisterRevenue
func NewMsgRegisterRevenue(
	contract sdk.Address,
	deployer,
	withdrawer sdk.AccAddress,
) *MsgRegisterRevenue {
	withdrawerAddress := ""
	if withdrawer != nil {
		withdrawerAddress = withdrawer.String()
	}

	return &MsgRegisterRevenue{
		ContractAddress:   contract.String(),
		DeployerAddress:   deployer.String(),
		WithdrawerAddress: withdrawerAddress,
	}
}

// Route returns the name of the module
func (msg MsgRegisterRevenue) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgRegisterRevenue) Type() string { return TypeMsgRegisterRevenue }

// ValidateBasic runs stateless checks on the message
func (msg MsgRegisterRevenue) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.DeployerAddress); err != nil {
		return sdkerrors.Wrapf(err, "invalid deployer address %s", msg.DeployerAddress)
	}

	if _, err := sdk.AccAddressFromBech32(msg.ContractAddress); err != nil {
		return sdkerrors.Wrapf(err, "invalid contract address %s", msg.ContractAddress)
	}

	if msg.WithdrawerAddress != "" {
		if _, err := sdk.AccAddressFromBech32(msg.WithdrawerAddress); err != nil {
			return sdkerrors.Wrapf(err, "invalid withdraw address %s", msg.WithdrawerAddress)
		}
	}

	// if len(msg.Nonces) < 1 {
	// 	return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid nonces - empty array")
	// }

	// if len(msg.Nonces) > 20 {
	// 	return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid nonces - array length must be less than 20")
	// }

	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgRegisterRevenue) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgRegisterRevenue) GetSigners() []sdk.AccAddress {
	from := sdk.MustAccAddressFromBech32(msg.DeployerAddress)
	return []sdk.AccAddress{from}
}

// NewMsgCancelRevenue creates new instance of MsgCancelRevenue.
func NewMsgCancelRevenue(
	contract sdk.Address,
	deployer sdk.AccAddress,
) *MsgCancelRevenue {
	return &MsgCancelRevenue{
		ContractAddress: contract.String(),
		DeployerAddress: deployer.String(),
	}
}

// Route returns the message route for a MsgCancelRevenue.
func (msg MsgCancelRevenue) Route() string { return RouterKey }

// Type returns the message type for a MsgCancelRevenue.
func (msg MsgCancelRevenue) Type() string { return TypeMsgCancelRevenue }

// ValidateBasic runs stateless checks on the message
func (msg MsgCancelRevenue) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.DeployerAddress); err != nil {
		return sdkerrors.Wrapf(err, "invalid deployer address %s", msg.DeployerAddress)
	}

	if _, err := sdk.AccAddressFromBech32(msg.ContractAddress); err != nil {
		return sdkerrors.Wrapf(err, "invalid deployer address %s", msg.DeployerAddress)
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgCancelRevenue) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgCancelRevenue) GetSigners() []sdk.AccAddress {
	funder := sdk.MustAccAddressFromBech32(msg.DeployerAddress)
	return []sdk.AccAddress{funder}
}

// NewMsgUpdateRevenue creates new instance of MsgUpdateRevenue
func NewMsgUpdateRevenue(
	contract sdk.Address,
	deployer,
	withdraw sdk.AccAddress,
) *MsgUpdateRevenue {
	return &MsgUpdateRevenue{
		ContractAddress:   contract.String(),
		DeployerAddress:   deployer.String(),
		WithdrawerAddress: withdraw.String(),
	}
}

// Route returns the name of the module
func (msg MsgUpdateRevenue) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgUpdateRevenue) Type() string { return TypeMsgUpdateRevenue }

// ValidateBasic runs stateless checks on the message
func (msg MsgUpdateRevenue) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.DeployerAddress); err != nil {
		return sdkerrors.Wrapf(err, "invalid deployer address %s", msg.DeployerAddress)
	}

	if _, err := sdk.AccAddressFromBech32(msg.ContractAddress); err != nil {
		return sdkerrors.Wrapf(err, "invalid contract address %s", msg.ContractAddress)
	}

	if _, err := sdk.AccAddressFromBech32(msg.WithdrawerAddress); err != nil {
		return sdkerrors.Wrapf(err, "invalid withdraw address %s", msg.WithdrawerAddress)
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgUpdateRevenue) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgUpdateRevenue) GetSigners() []sdk.AccAddress {
	from := sdk.MustAccAddressFromBech32(msg.DeployerAddress)
	return []sdk.AccAddress{from}
}
