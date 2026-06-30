package domain

type FlowState string

const (
	StateMain           FlowState = "main"
	StateStarsFlow      FlowState = "stars_flow"
	StateStarsWho       FlowState = "stars_who"
	StateStarsCustom    FlowState = "stars_custom_amount"
	StateStarsPayment   FlowState = "stars_payment"
	StatePremiumFlow    FlowState = "premium_flow"
	StatePremiumWho     FlowState = "premium_who"
	StatePremiumPayment FlowState = "premium_payment"
	StateProfile        FlowState = "profile"
	StateReferral       FlowState = "referral_program"
	StateWithdrawFlow   FlowState = "withdraw_flow"
	StateWithdrawWho    FlowState = "withdraw_who"
	StateWithdrawAmount FlowState = "withdraw_amount"
)
