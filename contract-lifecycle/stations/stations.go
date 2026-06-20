package stations

import "github.com/terapps/gonveyor/blueprint"

// --- Shared types ---

type DocumentInput struct {
	EntityID string // QuoteID or ContractID
	DocType  string // e.g. "proposal", "pricing", "contract", "annex_a"
}
type DocumentOutput struct {
	DocURL  string
	DocType string
}

type EmailTemplate string

const (
	TemplateSignatureRequest EmailTemplate = "signature_request"
	TemplateContractSigned   EmailTemplate = "contract_signed"
)

type SendEmailInput struct {
	To       string
	Template EmailTemplate
	Vars     map[string]string // template variables
	DocURLs  []string          // attachments / links
}
type SendEmailOutput struct{}

type SyncCrmInput struct {
	EntityType string // "quote" or "contract"
	EntityID   string
	Metadata   map[string]string
	DocURLs    []string
}
type SyncCrmOutput struct{}

// --- Phase 1: Devis ---

type InitiateSignatureInput struct {
	QuoteID     string
	ClientEmail string
	DocURLs     []string
}
type InitiateSignatureOutput struct {
	ProcessID    string
	SignatureURL string
}

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
	ClientEmail string
	SignatureID string
	TxnID       string
	Amount      float64
}
type CreateContractOutput struct {
	ContractID  string
	ClientEmail string
}

// --- Station definitions ---

// GenerateQuoteDoc and GenerateContractDoc share DocumentInput/Output and the same handler.
// They are distinct stations because a node key must be unique within a blueprint.
var GenerateQuoteDoc    = blueprint.Define[DocumentInput, DocumentOutput]("generate_quote_doc")
var InitiateSignature   = blueprint.Define[InitiateSignatureInput, InitiateSignatureOutput]("initiate_signature")
var SendQuoteEmail      = blueprint.Define[SendEmailInput, SendEmailOutput]("send_quote_email")
var SyncCrmQuote        = blueprint.Define[SyncCrmInput, SyncCrmOutput]("sync_crm_quote")
var AwaitSignature      = blueprint.Signal[SignaturePayload]("await_signature")
var AwaitPayment        = blueprint.Signal[PaymentPayload]("await_payment")
var CreateContract      = blueprint.Define[CreateContractInput, CreateContractOutput]("create_contract")
var GenerateContractDoc = blueprint.Define[DocumentInput, DocumentOutput]("generate_contract_doc")
var SendContractEmail   = blueprint.Define[SendEmailInput, SendEmailOutput]("send_contract_email")
var SyncCrmContract     = blueprint.Define[SyncCrmInput, SyncCrmOutput]("sync_crm_contract")
