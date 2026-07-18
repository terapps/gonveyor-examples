package stations

import "github.com/terapps/gonveyor"

type CreateContractInput struct {
	QuoteID     string
	ClientEmail string
	SignatureID string
	TxnID       string
	Amount      float64
}
type CreateContractOutput struct {
	ContractID  string
	ClientEmail string
}

var CreateContract = gonveyor.Define[CreateContractInput, CreateContractOutput]("create_contract")
