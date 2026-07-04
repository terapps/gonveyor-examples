package stations

import "github.com/terapps/gonveyor"

// --- Shared types ---

type DocumentInput struct {
	EntityID string // QuoteID or ContractID
	DocType  string // e.g. "proposal", "pricing", "contract", "annex_a"
}
type DocumentOutput struct {
	DocURL   string
	DocType  string
	EntityID string
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

type InitiatePaymentInput struct {
	QuoteID     string
	ClientEmail string
	Amount      float64
	DocURLs     []string
}
type InitiatePaymentOutput struct {
	ProcessID  string
	PaymentURL string
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

type BundleContractDocsInput struct {
	ContractID  string
	ClientEmail string
	DocURLs     []string
}
type BundleContractDocsOutput struct {
	ContractID  string
	ClientEmail string
	DocURLs     []string
}

// --- Station definitions ---

// GenerateQuoteDoc and GenerateContractDoc share DocumentInput/Output and the same handler.
// They are distinct stations because a node key must be unique within a blueprint.
var GenerateQuoteDoc    = gonveyor.Define[DocumentInput, DocumentOutput]("generate_quote_doc", gonveyor.WithRoute("tasks.document"))
var InitiateSignature   = gonveyor.Define[InitiateSignatureInput, InitiateSignatureOutput]("initiate_signature")
var InitiatePayment     = gonveyor.Define[InitiatePaymentInput, InitiatePaymentOutput]("initiate_payment")
var SendQuoteEmail      = gonveyor.Define[SendEmailInput, SendEmailOutput]("send_quote_email")
var SyncCrmQuote        = gonveyor.Define[SyncCrmInput, SyncCrmOutput]("sync_crm_quote")
var AwaitSignature      = gonveyor.Signal[SignaturePayload]("await_signature")
var AwaitPayment        = gonveyor.Signal[PaymentPayload]("await_payment")
var CreateContract      = gonveyor.Define[CreateContractInput, CreateContractOutput]("create_contract")
var GenerateContractDoc  = gonveyor.Define[DocumentInput, DocumentOutput]("generate_contract_doc", gonveyor.WithRoute("tasks.document"))
var BundleContractDocs  = gonveyor.Define[BundleContractDocsInput, BundleContractDocsOutput]("bundle_contract_docs")
var SendContractEmail   = gonveyor.Define[SendEmailInput, SendEmailOutput]("send_contract_email")
var SyncCrmContract     = gonveyor.Define[SyncCrmInput, SyncCrmOutput]("sync_crm_contract")
