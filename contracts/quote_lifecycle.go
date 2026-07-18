package contracts

// Workflow: Création d'un Devis
//
// Phase 1 — Devis:
//   GenerateQuoteDoc * N ──> InitiateSignature ──> SendQuoteEmail ─┐
//                          │                   └──> SyncCrmQuote   │
//                          └──> InitiatePayment                    │
//   AwaitSignature (signal, after InitiateSignature) ──────────────┼──> CreateContract
//   AwaitPayment   (signal, after InitiatePayment)   ──────────────┘
//
// Phase 2 — Contrat:
//   CreateContract ──> GenerateContractDoc * N ──> SendContractEmail
//                                             └──> SyncCrmContract

import (
	"github.com/terapps/gonveyor"
)

var QuoteLifecycle = gonveyor.New("quote_lifecycle",
	GenerateQuoteDoc, // root — dispatched via Seeds at manifest time
	gonveyor.Wire(AwaitSignature,
		gonveyor.After[struct{}](InitiateSignature),
	),

	gonveyor.Wire(AwaitPayment,
		gonveyor.After[struct{}](InitiatePayment),
	),

	gonveyor.Wire(InitiateSignature,
		gonveyor.Merge(GenerateQuoteDoc, func(outs []DocumentOutput, in *InitiateSignatureInput) {
			in.DocURLs = make([]string, len(outs))
			for i, o := range outs {
				in.DocURLs[i] = o.DocURL
			}
		}),
	),

	gonveyor.Wire(InitiatePayment,
		gonveyor.Merge(GenerateQuoteDoc, func(outs []DocumentOutput, in *InitiatePaymentInput) {
			in.DocURLs = make([]string, len(outs))
			for i, o := range outs {
				in.DocURLs[i] = o.DocURL
			}
		}),
	),

	gonveyor.Wire(SendQuoteEmail,
		gonveyor.Intake(InitiateSignature, func(o InitiateSignatureOutput, in *SendEmailInput) {
			in.Vars = map[string]string{"signature_url": o.SignatureURL}
		}),
	),

	gonveyor.Wire(SyncCrmQuote,
		gonveyor.Intake(InitiateSignature, func(o InitiateSignatureOutput, in *SyncCrmInput) {
			in.Metadata = map[string]string{"process_id": o.ProcessID}
		}),
	),

	// CreateContract waits for: email sent + crm synced + signature received + payment received
	gonveyor.Wire(CreateContract,
		gonveyor.After[CreateContractInput](SendQuoteEmail),
		gonveyor.After[CreateContractInput](SyncCrmQuote),
		gonveyor.Intake(AwaitSignature, func(p SignaturePayload, in *CreateContractInput) {
			in.SignatureID = p.SignatureID
		}),
		gonveyor.Intake(AwaitPayment, func(p PaymentPayload, in *CreateContractInput) {
			in.TxnID = p.TxnID
			in.Amount = p.Amount
		}),
	),

	gonveyor.Wire(GenerateContractDoc,
		gonveyor.Intake(CreateContract, func(o CreateContractOutput, in *DocumentInput) {
			in.EntityID = o.ContractID
		}),
	),

	gonveyor.Wire(BundleContractDocs,
		gonveyor.Merge(GenerateContractDoc, func(outs []DocumentOutput, in *BundleContractDocsInput) {
			in.DocURLs = make([]string, len(outs))
			for i, o := range outs {
				in.DocURLs[i] = o.DocURL
			}
			if len(outs) > 0 {
				in.ContractID = outs[0].EntityID
			}
		}),
	),

	gonveyor.Wire(SendContractEmail,
		gonveyor.Intake(BundleContractDocs, func(o BundleContractDocsOutput, in *SendEmailInput) {
			in.To = o.ClientEmail
			in.DocURLs = o.DocURLs
		}),
	),

	gonveyor.Wire(SyncCrmContract,
		gonveyor.Intake(BundleContractDocs, func(o BundleContractDocsOutput, in *SyncCrmInput) {
			in.EntityID = o.ContractID
			in.DocURLs = o.DocURLs
		}),
	),
)

type Params struct {
	QuoteID          string   `validate:"required"`
	ClientEmail      string   `validate:"required,email"`
	Amount           float64  `validate:"gt=0"`
	QuoteDocTypes    []string `validate:"required,min=1"` // e.g. ["proposal", "pricing", "terms"]
	ContractDocTypes []string `validate:"required,min=1"` // e.g. ["contract", "annex_a"]
}

var QuoteLifecycleTemplate = gonveyor.NewLaunchTemplate(QuoteLifecycle, func(p Params) []gonveyor.ManifestOption {
	return []gonveyor.ManifestOption{
		// N quote documents dispatched in parallel
		gonveyor.Seeds(GenerateQuoteDoc, p.QuoteDocTypes, func(docType string, in *DocumentInput) {
			in.EntityID = p.QuoteID
			in.DocType = docType
		}),
		// Ambient context seeded into downstream nodes
		gonveyor.Seed(InitiateSignature, InitiateSignatureInput{
			QuoteID:     p.QuoteID,
			ClientEmail: p.ClientEmail,
		}),
		gonveyor.Seed(InitiatePayment, InitiatePaymentInput{
			QuoteID:     p.QuoteID,
			ClientEmail: p.ClientEmail,
			Amount:      p.Amount,
		}),
		gonveyor.Seed(SyncCrmQuote, SyncCrmInput{
			EntityType: "quote",
			EntityID:   p.QuoteID,
		}),
		gonveyor.Seed(SendQuoteEmail, SendEmailInput{
			To:       p.ClientEmail,
			Template: TemplateSignatureRequest,
		}),
		gonveyor.Seed(CreateContract, CreateContractInput{
			QuoteID:     p.QuoteID,
			ClientEmail: p.ClientEmail,
		}),
		gonveyor.Seed(BundleContractDocs, BundleContractDocsInput{
			ClientEmail: p.ClientEmail,
		}),
		gonveyor.Seed(SendContractEmail, SendEmailInput{
			Template: TemplateContractSigned,
		}),
		gonveyor.Seed(SyncCrmContract, SyncCrmInput{
			EntityType: "contract",
		}),
		// N contract documents dispatched in parallel after CreateContract
		gonveyor.Seeds(GenerateContractDoc, p.ContractDocTypes, func(docType string, in *DocumentInput) {
			in.DocType = docType
		}),
	}
})
