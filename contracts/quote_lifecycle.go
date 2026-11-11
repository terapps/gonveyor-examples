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
	st "github.com/terapps/gonveyor-examples/contracts/stations"
)

var QuoteLifecycle = gonveyor.New("quote_lifecycle",
	st.GenerateQuoteDoc, // root — dispatched via Seeds at manifest time
	st.AwaitSignature.After(st.InitiateSignature),

	st.AwaitPayment.After(st.InitiatePayment),

	st.InitiateSignature.Merge(st.GenerateQuoteDoc, func(outs []st.DocumentOutput, in *st.InitiateSignatureInput) {
		in.DocURLs = make([]string, len(outs))
		for i, o := range outs {
			in.DocURLs[i] = o.DocURL
		}
	}),

	st.InitiatePayment.Merge(st.GenerateQuoteDoc, func(outs []st.DocumentOutput, in *st.InitiatePaymentInput) {
		in.DocURLs = make([]string, len(outs))
		for i, o := range outs {
			in.DocURLs[i] = o.DocURL
		}
	}),

	st.SendQuoteEmail.Intake(st.InitiateSignature, func(o st.InitiateSignatureOutput, in *st.SendEmailInput) {
		in.Vars = map[string]string{"signature_url": o.SignatureURL}
	}),

	st.SyncCrmQuote.Intake(st.InitiateSignature, func(o st.InitiateSignatureOutput, in *st.SyncCrmInput) {
		in.Metadata = map[string]string{"process_id": o.ProcessID}
	}),

	// Refused signature or failed payment: their own branch runs, CreateContract and everything
	// after it are skipped, and the other wait is cancelled so the run ends.
	st.MarkQuoteDeclined.
		When(st.AwaitSignature, func(p st.SignaturePayload) bool { return p.Declined }).
		Intake(st.AwaitSignature, func(p st.SignaturePayload, in *st.QuoteOutcomeInput) { in.Reason = p.Reason }).
		Cancels(st.AwaitPayment),

	st.NotifyPaymentFailed.
		When(st.AwaitPayment, func(p st.PaymentPayload) bool { return p.Failed }).
		Intake(st.AwaitPayment, func(p st.PaymentPayload, in *st.QuoteOutcomeInput) { in.Reason = p.Reason }).
		Cancels(st.AwaitSignature),

	// CreateContract waits for: email sent + crm synced + a signed quote + a successful payment
	st.CreateContract.
		After(st.SendQuoteEmail, st.SyncCrmQuote).
		When(st.AwaitSignature, func(p st.SignaturePayload) bool { return !p.Declined }).
		When(st.AwaitPayment, func(p st.PaymentPayload) bool { return !p.Failed }).
		Intake(st.AwaitSignature, func(p st.SignaturePayload, in *st.CreateContractInput) {
			in.SignatureID = p.SignatureID
		}).
		Intake(st.AwaitPayment, func(p st.PaymentPayload, in *st.CreateContractInput) {
			in.TxnID = p.TxnID
			in.Amount = p.Amount
		}),

	st.GenerateContractDoc.Intake(st.CreateContract, func(o st.CreateContractOutput, in *st.DocumentInput) {
		in.EntityID = o.ContractID
	}),

	st.BundleContractDocs.Merge(st.GenerateContractDoc, func(outs []st.DocumentOutput, in *st.BundleContractDocsInput) {
		in.DocURLs = make([]string, len(outs))
		for i, o := range outs {
			in.DocURLs[i] = o.DocURL
		}
		if len(outs) > 0 {
			in.ContractID = outs[0].EntityID
		}
	}),

	st.SendContractEmail.Intake(st.BundleContractDocs, func(o st.BundleContractDocsOutput, in *st.SendEmailInput) {
		in.To = o.ClientEmail
		in.DocURLs = o.DocURLs
	}),

	st.SyncCrmContract.Intake(st.BundleContractDocs, func(o st.BundleContractDocsOutput, in *st.SyncCrmInput) {
		in.EntityID = o.ContractID
		in.DocURLs = o.DocURLs
	}),
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
		st.GenerateQuoteDoc.Seeds(p.QuoteDocTypes, func(docType string, in *st.DocumentInput) {
			in.EntityID = p.QuoteID
			in.DocType = docType
		}),
		// Ambient context seeded into downstream nodes
		st.InitiateSignature.Seed(st.InitiateSignatureInput{
			QuoteID:     p.QuoteID,
			ClientEmail: p.ClientEmail,
		}),
		st.MarkQuoteDeclined.Seed(st.QuoteOutcomeInput{QuoteID: p.QuoteID, ClientEmail: p.ClientEmail}),
		st.NotifyPaymentFailed.Seed(st.QuoteOutcomeInput{QuoteID: p.QuoteID, ClientEmail: p.ClientEmail}),
		st.InitiatePayment.Seed(st.InitiatePaymentInput{
			QuoteID:     p.QuoteID,
			ClientEmail: p.ClientEmail,
			Amount:      p.Amount,
		}),
		st.SyncCrmQuote.Seed(st.SyncCrmInput{
			EntityType: "quote",
			EntityID:   p.QuoteID,
		}),
		st.SendQuoteEmail.Seed(st.SendEmailInput{
			To:       p.ClientEmail,
			Template: st.TemplateSignatureRequest,
		}),
		st.CreateContract.Seed(st.CreateContractInput{
			QuoteID:     p.QuoteID,
			ClientEmail: p.ClientEmail,
		}),
		st.BundleContractDocs.Seed(st.BundleContractDocsInput{
			ClientEmail: p.ClientEmail,
		}),
		st.SendContractEmail.Seed(st.SendEmailInput{
			Template: st.TemplateContractSigned,
		}),
		st.SyncCrmContract.Seed(st.SyncCrmInput{
			EntityType: "contract",
		}),
		// N contract documents dispatched in parallel after CreateContract
		st.GenerateContractDoc.Seeds(p.ContractDocTypes, func(docType string, in *st.DocumentInput) {
			in.DocType = docType
		}),
	}
})
