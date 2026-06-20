package blueprint

// Workflow: Création d'un Devis
//
// Phase 1 — Devis:
//   GenerateQuoteDoc * N ──> InitiateSignature ──> SendQuoteEmail ──┐
//                                              └──> SyncCrmQuote   ──┤
//   AwaitSignature (signal) ────────────────────────────────────────┼──> CreateContract
//   AwaitPayment   (signal) ────────────────────────────────────────┘
//
// Phase 2 — Contrat:
//   CreateContract ──> GenerateContractDoc * N ──> SendContractEmail
//                                             └──> SyncCrmContract

import (
	"github.com/terapps/gonveyor"
	bp "github.com/terapps/gonveyor/blueprint"
	"github.com/terapps/gonveyor/ledger"
	st "github.com/terapps/gonveyor-examples/contract-lifecycle/stations"
)

var QuoteLifecycle = bp.New("quote_lifecycle",
	st.GenerateQuoteDoc, // root — dispatched via Seeds at manifest time
	st.AwaitSignature,   // root signal — held until signature webhook
	st.AwaitPayment,     // root signal — held until payment webhook

	bp.Wire(st.InitiateSignature,
		gonveyor.Merge(st.GenerateQuoteDoc, func(outs []st.DocumentOutput, in *st.InitiateSignatureInput) {
			in.DocURLs = make([]string, len(outs))
			for i, o := range outs {
				in.DocURLs[i] = o.DocURL
			}
		}),
	),

	bp.Wire(st.SendQuoteEmail,
		gonveyor.Intake(st.InitiateSignature, func(o st.InitiateSignatureOutput, in *st.SendEmailInput) {
			in.Vars = map[string]string{"signature_url": o.SignatureURL}
		}),
	),

	bp.Wire(st.SyncCrmQuote,
		gonveyor.Intake(st.InitiateSignature, func(o st.InitiateSignatureOutput, in *st.SyncCrmInput) {
			in.Metadata = map[string]string{"process_id": o.ProcessID}
		}),
	),

	// CreateContract waits for: email sent + crm synced + signature received + payment received
	bp.Wire(st.CreateContract,
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

	bp.Wire(st.GenerateContractDoc,
		gonveyor.Intake(st.CreateContract, func(o st.CreateContractOutput, in *st.DocumentInput) {
			in.EntityID = o.ContractID
		}),
	),

	bp.Wire(st.SendContractEmail,
		gonveyor.Intake(st.CreateContract, func(o st.CreateContractOutput, in *st.SendEmailInput) {
			in.To = o.ClientEmail
		}),
		gonveyor.Merge(st.GenerateContractDoc, func(outs []st.DocumentOutput, in *st.SendEmailInput) {
			in.DocURLs = make([]string, len(outs))
			for i, o := range outs {
				in.DocURLs[i] = o.DocURL
			}
		}),
	),

	bp.Wire(st.SyncCrmContract,
		gonveyor.Intake(st.CreateContract, func(o st.CreateContractOutput, in *st.SyncCrmInput) {
			in.EntityID = o.ContractID
		}),
		gonveyor.Merge(st.GenerateContractDoc, func(outs []st.DocumentOutput, in *st.SyncCrmInput) {
			in.DocURLs = make([]string, len(outs))
			for i, o := range outs {
				in.DocURLs[i] = o.DocURL
			}
		}),
	),
)

type Params struct {
	QuoteID          string
	ClientEmail      string
	QuoteDocTypes    []string // e.g. ["proposal", "pricing", "terms"]
	ContractDocTypes []string // e.g. ["contract", "annex_a"]
}

func Manifest(p Params) (ledger.BlueprintManifest, error) {
	return QuoteLifecycle.Manifest(
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
	)
}
