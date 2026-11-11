package stations

import "github.com/terapps/gonveyor"

type RenewalDecision struct {
	Accepted bool
	Reason   string
}

// AwaitRenewalDecision is a signal — the client's answer, sent via Gonductor.Send.
var AwaitRenewalDecision = gonveyor.Signal[RenewalDecision]("await_renewal_decision")

type RenewContractInput struct {
	ContractID string `validate:"required"`
}
type RenewContractOutput struct {
	ContractID string
}

// RenewContract runs only if the decision is accepted (When).
var RenewContract = gonveyor.Define[RenewContractInput, RenewContractOutput]("renew_contract")

type CloseContractInput struct {
	ContractID string `validate:"required"`
	Reason     string
}
type CloseContractOutput struct{}

// CloseContract runs only if the decision is declined (When).
var CloseContract = gonveyor.Define[CloseContractInput, CloseContractOutput]("close_contract")
