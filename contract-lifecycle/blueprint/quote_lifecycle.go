package blueprint

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
	st "github.com/terapps/gonveyor-examples/contract-lifecycle/stations"
)

var QuoteLifecycle = gonveyor.New("quote_lifecycle",
	st.GenerateQuoteDoc, // root — dispatched via Seeds at manifest time
	gonveyor.Wire(st.AwaitSignature,
		gonveyor.After[struct{}](st.InitiateSignature),
	),

	gonveyor.Wire(st.AwaitPayment,
		gonveyor.After[struct{}](st.InitiatePayment),
	),

	gonveyor.Wire(st.InitiateSignature,
		gonveyor.Merge(st.GenerateQuoteDoc, func(outs []st.DocumentOutput, in *st.InitiateSignatureInput) {
			in.DocURLs = make([]string, len(outs))
			for i, o := range outs {
				in.DocURLs[i] = o.DocURL
			}
		}),
	),

	gonveyor.Wire(st.InitiatePayment,
		gonveyor.Merge(st.GenerateQuoteDoc, func(outs []st.DocumentOutput, in *st.InitiatePaymentInput) {
			in.DocURLs = make([]string, len(outs))
			for i, o := range outs {
				in.DocURLs[i] = o.DocURL
			}
		}),
	),

	gonveyor.Wire(st.SendQuoteEmail,
		gonveyor.Intake(st.InitiateSignature, func(o st.InitiateSignatureOutput, in *st.SendEmailInput) {
			in.Vars = map[string]string{"signature_url": o.SignatureURL}
		}),
	),

	gonveyor.Wire(st.SyncCrmQuote,
		gonveyor.Intake(st.InitiateSignature, func(o st.InitiateSignatureOutput, in *st.SyncCrmInput) {
			in.Metadata = map[string]string{"process_id": o.ProcessID}
		}),
	),

	// CreateContract waits for: email sent + crm synced + signature received + payment received
	gonveyor.Wire(st.CreateContract,
		gonveyor.After[st.CreateContractInput](st.SendQuoteEmail),
		gonveyor.After[st.CreateContractInput](st.SyncCrmQuote),
		gonveyor.Intake(st.AwaitSignature, func(p st.SignaturePayload, in *st.CreateContractInput) {
			in.SignatureID = p.SignatureID
		}),
		gonveyor.Intake(st.AwaitPayment, func(p st.PaymentPayload, in *st.CreateContractInput) {
			in.TxnID = p.TxnID
			in.Amount = p.Amount
		}),
	),

	gonveyor.Wire(st.GenerateContractDoc,
		gonveyor.Intake(st.CreateContract, func(o st.CreateContractOutput, in *st.DocumentInput) {
			in.EntityID = o.ContractID
		}),
	),

	gonveyor.Wire(st.BundleContractDocs,
		gonveyor.Merge(st.GenerateContractDoc, func(outs []st.DocumentOutput, in *st.BundleContractDocsInput) {
			in.DocURLs = make([]string, len(outs))
			for i, o := range outs {
				in.DocURLs[i] = o.DocURL
			}
			if len(outs) > 0 {
				in.ContractID = outs[0].EntityID
			}
		}),
	),

	gonveyor.Wire(st.SendContractEmail,
		gonveyor.Intake(st.BundleContractDocs, func(o st.BundleContractDocsOutput, in *st.SendEmailInput) {
			in.To = o.ClientEmail
			in.DocURLs = o.DocURLs
		}),
	),

	gonveyor.Wire(st.SyncCrmContract,
		gonveyor.Intake(st.BundleContractDocs, func(o st.BundleContractDocsOutput, in *st.SyncCrmInput) {
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
		gonveyor.Seeds(st.GenerateQuoteDoc, p.QuoteDocTypes, func(docType string, in *st.DocumentInput) {
			in.EntityID = p.QuoteID
			in.DocType = docType
		}),
		// Ambient context seeded into downstream nodes
		gonveyor.Seed(st.InitiateSignature, st.InitiateSignatureInput{
			QuoteID:     p.QuoteID,
			ClientEmail: p.ClientEmail,
		}),
		gonveyor.Seed(st.InitiatePayment, st.InitiatePaymentInput{
			QuoteID:     p.QuoteID,
			ClientEmail: p.ClientEmail,
			Amount:      p.Amount,
		}),
		gonveyor.Seed(st.SyncCrmQuote, st.SyncCrmInput{
			EntityType: "quote",
			EntityID:   p.QuoteID,
		}),
		gonveyor.Seed(st.SendQuoteEmail, st.SendEmailInput{
			To:       p.ClientEmail,
			Template: st.TemplateSignatureRequest,
		}),
		gonveyor.Seed(st.CreateContract, st.CreateContractInput{
			QuoteID:     p.QuoteID,
			ClientEmail: p.ClientEmail,
		}),
		gonveyor.Seed(st.BundleContractDocs, st.BundleContractDocsInput{
			ClientEmail: p.ClientEmail,
		}),
		gonveyor.Seed(st.SendContractEmail, st.SendEmailInput{
			Template: st.TemplateContractSigned,
		}),
		gonveyor.Seed(st.SyncCrmContract, st.SyncCrmInput{
			EntityType: "contract",
		}),
		// N contract documents dispatched in parallel after CreateContract
		gonveyor.Seeds(st.GenerateContractDoc, p.ContractDocTypes, func(docType string, in *st.DocumentInput) {
			in.DocType = docType
		}),
	}
})
