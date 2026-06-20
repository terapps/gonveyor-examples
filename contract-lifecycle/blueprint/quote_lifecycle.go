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
	st.GenerateQuoteDoc, // root — dispatched via SplitWith at manifest time
	st.AwaitSignature,   // root signal — held until signature webhook
	st.AwaitPayment,     // root signal — held until payment webhook

	bp.Wire(st.InitiateSignature,
		gonveyor.Merge(st.GenerateQuoteDoc, func(outs []st.GenerateQuoteDocOutput, in *st.InitiateSignatureInput) {
			in.DocURLs = make([]string, len(outs))
			for i, o := range outs {
				in.DocURLs[i] = o.DocURL
			}
		}),
	),

	bp.Wire(st.SendQuoteEmail,
		gonveyor.Intake(st.InitiateSignature, func(o st.InitiateSignatureOutput, in *st.SendQuoteEmailInput) {
			in.SignatureURL = o.SignatureURL
		}),
	),

	bp.Wire(st.SyncCrmQuote,
		gonveyor.Intake(st.InitiateSignature, func(o st.InitiateSignatureOutput, in *st.SyncCrmQuoteInput) {
			in.ProcessID = o.ProcessID
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
		gonveyor.Intake(st.CreateContract, func(o st.CreateContractOutput, in *st.GenerateContractDocInput) {
			in.ContractID = o.ContractID
		}),
	),

	bp.Wire(st.SendContractEmail,
		gonveyor.Intake(st.CreateContract, func(o st.CreateContractOutput, in *st.SendContractEmailInput) {
			in.ClientEmail = o.ClientEmail
		}),
		gonveyor.Merge(st.GenerateContractDoc, func(outs []st.GenerateContractDocOutput, in *st.SendContractEmailInput) {
			in.DocURLs = make([]string, len(outs))
			for i, o := range outs {
				in.DocURLs[i] = o.DocURL
			}
		}),
	),

	bp.Wire(st.SyncCrmContract,
		gonveyor.Intake(st.CreateContract, func(o st.CreateContractOutput, in *st.SyncCrmContractInput) {
			in.ContractID = o.ContractID
		}),
		gonveyor.Merge(st.GenerateContractDoc, func(outs []st.GenerateContractDocOutput, in *st.SyncCrmContractInput) {
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
	ContractDocTypes []string // e.g. ["contract", "annex_a", "annex_b"]
}

func Manifest(p Params) (ledger.BlueprintManifest, error) {
	return QuoteLifecycle.Manifest(
		// N quote documents dispatched in parallel
		gonveyor.Seeds(st.GenerateQuoteDoc, p.QuoteDocTypes, func(docType string, in *st.GenerateQuoteDocInput) {
			in.QuoteID = p.QuoteID
			in.ClientEmail = p.ClientEmail
			in.DocType = docType
		}),
		// Ambient context threaded to nodes whose inputs aren't fully covered by Intake/Merge
		gonveyor.Seed(st.InitiateSignature, st.InitiateSignatureInput{
			QuoteID:     p.QuoteID,
			ClientEmail: p.ClientEmail,
		}),
		gonveyor.Seed(st.SyncCrmQuote, st.SyncCrmQuoteInput{QuoteID: p.QuoteID}),
		gonveyor.Seed(st.CreateContract, st.CreateContractInput{QuoteID: p.QuoteID}),
		gonveyor.Seed(st.SendContractEmail, st.SendContractEmailInput{ClientEmail: p.ClientEmail}),
		// N contract documents dispatched in parallel after CreateContract
		gonveyor.Seeds(st.GenerateContractDoc, p.ContractDocTypes, func(docType string, in *st.GenerateContractDocInput) {
			in.DocType = docType
		}),
	)
}
