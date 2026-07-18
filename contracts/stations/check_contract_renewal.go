package stations

import "github.com/terapps/gonveyor"

type CheckContractRenewalInput struct {
	ContractID  string `validate:"required"`
	ClientEmail string `validate:"required,email"`
}
type CheckContractRenewalOutput struct {
	ContractID  string
	ClientEmail string
	RenewalURL  string
}

// CheckContractRenewal is the root of the contract_renewal blueprint. SendContractEmail is
// reused as-is from quote_lifecycle — same station, different Wire().
var CheckContractRenewal = gonveyor.Define[CheckContractRenewalInput, CheckContractRenewalOutput]("check_contract_renewal")
