package stations

import "github.com/terapps/gonveyor/blueprint"

// --- Phase 1: Devis ---

type GenerateQuoteDocInput struct {
	QuoteID     string
	ClientEmail string
	DocType     string
}
type GenerateQuoteDocOutput struct {
	DocURL  string
	DocType string
}

type InitiateSignatureInput struct {
	QuoteID     string
	ClientEmail string
	DocURLs     []string
}
type InitiateSignatureOutput struct {
	ProcessID   string
	SignatureURL string
}

type SendQuoteEmailInput struct {
	ClientEmail string
	SignatureURL string
}
type SendQuoteEmailOutput struct{}

type SyncCrmQuoteInput struct {
	QuoteID   string
	ProcessID string
}
type SyncCrmQuoteOutput struct{}

// --- Signals ---

type SignaturePayload struct {
	SignatureID string
	SignedAt    string
}

type PaymentPayload struct {
	TxnID  string
	Amount float64
	PaidAt string
}

// --- Phase 2: Contrat ---

type CreateContractInput struct {
	QuoteID     string
	SignatureID string
	TxnID       string
	Amount      float64
}
type CreateContractOutput struct {
	ContractID  string
	ClientEmail string
}

type GenerateContractDocInput struct {
	ContractID string
	DocType    string
}
type GenerateContractDocOutput struct {
	DocURL  string
	DocType string
}

type SendContractEmailInput struct {
	ClientEmail string
	DocURLs     []string
}
type SendContractEmailOutput struct{}

type SyncCrmContractInput struct {
	ContractID string
	DocURLs    []string
}
type SyncCrmContractOutput struct{}

// --- Station definitions ---

var GenerateQuoteDoc    = blueprint.Define[GenerateQuoteDocInput, GenerateQuoteDocOutput]("generate_quote_doc")
var InitiateSignature   = blueprint.Define[InitiateSignatureInput, InitiateSignatureOutput]("initiate_signature")
var SendQuoteEmail      = blueprint.Define[SendQuoteEmailInput, SendQuoteEmailOutput]("send_quote_email")
var SyncCrmQuote        = blueprint.Define[SyncCrmQuoteInput, SyncCrmQuoteOutput]("sync_crm_quote")
var AwaitSignature      = blueprint.Signal[SignaturePayload]("await_signature")
var AwaitPayment        = blueprint.Signal[PaymentPayload]("await_payment")
var CreateContract      = blueprint.Define[CreateContractInput, CreateContractOutput]("create_contract")
var GenerateContractDoc = blueprint.Define[GenerateContractDocInput, GenerateContractDocOutput]("generate_contract_doc")
var SendContractEmail   = blueprint.Define[SendContractEmailInput, SendContractEmailOutput]("send_contract_email")
var SyncCrmContract     = blueprint.Define[SyncCrmContractInput, SyncCrmContractOutput]("sync_crm_contract")
