package blueprint

import (
	"github.com/terapps/gonveyor"
	"github.com/terapps/gonveyor/blueprint"
	"github.com/terapps/gonveyor/ledger"
)

// --- types ---

type PrepareContractInput struct {
	ClientID    string
	ContractRef string
}
type PrepareContractOutput struct {
	ContractID string
	Amount     float64
}

type PaymentPayload struct {
	TxnID  string
	Amount float64
}

type FinalizeInput struct {
	ContractID string
	TxnID      string
	Amount     float64
}
type FinalizeOutput struct {
	SignedURL string
}

// --- stations ---

var PrepareContract = blueprint.Define[PrepareContractInput, PrepareContractOutput]("prepare_contract")
var AwaitPayment    = blueprint.Signal[PaymentPayload]("await_payment")
var FinalizeContract = blueprint.Define[FinalizeInput, FinalizeOutput]("finalize_contract")

// ContractFlow: prepare → await_payment (signal) → finalize
//
//	prepare_contract ──> [await_payment] ──> finalize_contract
//	                          ↑
//	                   payment webhook calls SendSignal
var ContractFlow = blueprint.New("contract_flow",
	PrepareContract,
	AwaitPayment,
	blueprint.Wire(FinalizeContract,
		gonveyor.After[FinalizeInput](PrepareContract),
		gonveyor.Intake(AwaitPayment, func(p PaymentPayload, in *FinalizeInput) {
			in.TxnID = p.TxnID
			in.Amount = p.Amount
		}),
		gonveyor.Intake(PrepareContract, func(o PrepareContractOutput, in *FinalizeInput) {
			in.ContractID = o.ContractID
		}),
	),
)

func ContractManifest(clientID, contractRef string) (ledger.BlueprintManifest, error) {
	return ContractFlow.Manifest(
		gonveyor.Seed(PrepareContract, PrepareContractInput{
			ClientID:    clientID,
			ContractRef: contractRef,
		}),
	)
}
