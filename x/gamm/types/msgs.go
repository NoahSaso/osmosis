package types

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// constants
const (
	TypeMsgCreatePool              = "create_pool"
	TypeMsgSwapExactAmountIn       = "swap_exact_amount_in"
	TypeMsgSwapExactAmountOut      = "swap_exact_amount_out"
	TypeMsgJoinPool                = "join_pool"
	TypeMsgExitPool                = "exit_pool"
	TypeMsgJoinSwapExternAmountIn  = "join_swap_extern_amount_in"
	TypeMsgJoinSwapShareAmountOut  = "join_swap_share_amount_out"
	TypeMsgExitSwapExternAmountOut = "exit_swap_extern_amount_out"
	TypeMsgExitSwapShareAmountIn   = "exit_swap_share_amount_in"
	TypeMsgUpdateSwapfee           = "update_swapfee"
)

var _ sdk.Msg = &MsgCreatePool{}

func (msg MsgCreatePool) Route() string { return RouterKey }
func (msg MsgCreatePool) Type() string  { return TypeMsgCreatePool }
func (msg MsgCreatePool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	if len(msg.Records) == 0 {
		return ErrEmptyRecords
	}

	if len(msg.Records) == 1 {
		return ErrTooLittleRecords
	}

	// TODO: Add the limit of binding token to the pool params?
	if len(msg.Records) > 8 {
		return sdkerrors.Wrapf(ErrTooManyRecords, "%d", len(msg.Records))
	}

	for _, record := range msg.Records {
		if !record.Weight.IsPositive() {
			return sdkerrors.Wrap(ErrNotPositiveWeight, record.Weight.String())
		}

		if !record.Token.IsValid() || !record.Token.IsPositive() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, record.Token.String())
		}
	}

	if msg.PoolParams.Lock {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "can't create the locked pool")
	}

	err = msg.PoolParams.Validate()
	if err != nil {
		return err
	}

	return nil
}
func (msg MsgCreatePool) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}
func (msg MsgCreatePool) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

var _ sdk.Msg = &MsgSwapExactAmountIn{}

func (msg MsgSwapExactAmountIn) Route() string { return RouterKey }
func (msg MsgSwapExactAmountIn) Type() string  { return TypeMsgSwapExactAmountIn }
func (msg MsgSwapExactAmountIn) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	err = SwapAmountInRoutes(msg.Routes).Validate()
	if err != nil {
		return err
	}

	if !msg.TokenIn.IsValid() || !msg.TokenIn.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.TokenIn.String())
	}

	if !msg.TokenOutMinAmount.IsPositive() {
		return ErrNotPositiveCriteria
	}

	return nil
}
func (msg MsgSwapExactAmountIn) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}
func (msg MsgSwapExactAmountIn) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

var _ sdk.Msg = &MsgSwapExactAmountOut{}

func (msg MsgSwapExactAmountOut) Route() string { return RouterKey }
func (msg MsgSwapExactAmountOut) Type() string  { return TypeMsgSwapExactAmountOut }
func (msg MsgSwapExactAmountOut) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	err = SwapAmountOutRoutes(msg.Routes).Validate()
	if err != nil {
		return err
	}

	if !msg.TokenOut.IsValid() || !msg.TokenOut.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.TokenOut.String())
	}

	if !msg.TokenInMaxAmount.IsPositive() {
		return ErrNotPositiveCriteria
	}

	return nil
}
func (msg MsgSwapExactAmountOut) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}
func (msg MsgSwapExactAmountOut) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

var _ sdk.Msg = &MsgJoinPool{}

func (msg MsgJoinPool) Route() string { return RouterKey }
func (msg MsgJoinPool) Type() string  { return TypeMsgJoinPool }
func (msg MsgJoinPool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	if !msg.ShareOutAmount.IsPositive() {
		return sdkerrors.Wrap(ErrNotPositiveRequireAmount, msg.ShareOutAmount.String())
	}

	tokenInMaxs := sdk.Coins(msg.TokenInMaxs)
	if !tokenInMaxs.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, tokenInMaxs.String())
	}

	return nil
}
func (msg MsgJoinPool) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}
func (msg MsgJoinPool) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

var _ sdk.Msg = &MsgExitPool{}

func (msg MsgExitPool) Route() string { return RouterKey }
func (msg MsgExitPool) Type() string  { return TypeMsgExitPool }
func (msg MsgExitPool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	if !msg.ShareInAmount.IsPositive() {
		return sdkerrors.Wrap(ErrNotPositiveRequireAmount, msg.ShareInAmount.String())
	}

	tokenOutMins := sdk.Coins(msg.TokenOutMins)
	if !tokenOutMins.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, tokenOutMins.String())
	}

	return nil
}
func (msg MsgExitPool) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}
func (msg MsgExitPool) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

var _ sdk.Msg = &MsgJoinSwapExternAmountIn{}

func (msg MsgJoinSwapExternAmountIn) Route() string { return RouterKey }
func (msg MsgJoinSwapExternAmountIn) Type() string  { return TypeMsgJoinSwapExternAmountIn }
func (msg MsgJoinSwapExternAmountIn) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	if !msg.TokenIn.IsValid() || !msg.TokenIn.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.TokenIn.String())
	}

	if !msg.ShareOutMinAmount.IsPositive() {
		return sdkerrors.Wrap(ErrNotPositiveCriteria, msg.ShareOutMinAmount.String())
	}

	return nil
}
func (msg MsgJoinSwapExternAmountIn) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}
func (msg MsgJoinSwapExternAmountIn) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

var _ sdk.Msg = &MsgJoinSwapShareAmountOut{}

func (msg MsgJoinSwapShareAmountOut) Route() string { return RouterKey }
func (msg MsgJoinSwapShareAmountOut) Type() string  { return TypeMsgJoinSwapShareAmountOut }
func (msg MsgJoinSwapShareAmountOut) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	err = sdk.ValidateDenom(msg.TokenInDenom)
	if err != nil {
		return err
	}

	if !msg.ShareOutAmount.IsPositive() {
		return sdkerrors.Wrap(ErrNotPositiveRequireAmount, msg.ShareOutAmount.String())
	}

	if !msg.TokenInMaxAmount.IsPositive() {
		return sdkerrors.Wrap(ErrNotPositiveCriteria, msg.TokenInMaxAmount.String())
	}

	return nil
}
func (msg MsgJoinSwapShareAmountOut) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}
func (msg MsgJoinSwapShareAmountOut) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

var _ sdk.Msg = &MsgExitSwapExternAmountOut{}

func (msg MsgExitSwapExternAmountOut) Route() string { return RouterKey }
func (msg MsgExitSwapExternAmountOut) Type() string  { return TypeMsgExitSwapExternAmountOut }
func (msg MsgExitSwapExternAmountOut) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	if !msg.TokenOut.IsValid() || !msg.TokenOut.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.TokenOut.String())
	}

	if !msg.ShareInMaxAmount.IsPositive() {
		return sdkerrors.Wrap(ErrNotPositiveCriteria, msg.ShareInMaxAmount.String())
	}

	return nil
}
func (msg MsgExitSwapExternAmountOut) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}
func (msg MsgExitSwapExternAmountOut) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

var _ sdk.Msg = &MsgExitSwapShareAmountIn{}

func (msg MsgExitSwapShareAmountIn) Route() string { return RouterKey }
func (msg MsgExitSwapShareAmountIn) Type() string  { return TypeMsgExitSwapShareAmountIn }
func (msg MsgExitSwapShareAmountIn) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	err = sdk.ValidateDenom(msg.TokenOutDenom)
	if err != nil {
		return err
	}

	if !msg.ShareInAmount.IsPositive() {
		return sdkerrors.Wrap(ErrNotPositiveRequireAmount, msg.ShareInAmount.String())
	}

	if !msg.TokenOutMinAmount.IsPositive() {
		return sdkerrors.Wrap(ErrNotPositiveCriteria, msg.TokenOutMinAmount.String())
	}

	return nil
}
func (msg MsgExitSwapShareAmountIn) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}
func (msg MsgExitSwapShareAmountIn) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

var _ sdk.Msg = &MsgUpdateSwapFee{}

func (m MsgUpdateSwapFee) Route() string { return RouterKey }
func (m MsgUpdateSwapFee) Type() string  { return TypeMsgUpdateSwapfee }
func (m MsgUpdateSwapFee) ValidateBasic() error {
	if m.NewSwapFee.IsNegative() {
		return errors.New("swap fee could not be negative")
	}

	if m.NewSwapFee.GT(sdk.NewDec(1)) {
		return errors.New("swap fee could not exceed 100%")
	}

	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return err
	}
	if sender.Empty() {
		return errors.New("empty sender")
	}

	return nil
}
func (m MsgUpdateSwapFee) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgUpdateSwapFee) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}
