package contracts

// Workflow: Renouvellement de contrat
//
// Indépendant de quote_lifecycle (déclenché par contract_renewal_scan, voir
// contract_renewal_scan.go — un seul schedule récurrent, pas un par contrat), mais
// réutilise trois stations telles quelles — GenerateContractDoc, SendContractEmail,
// SyncCrmContract — câblées différemment ici via des Wire() distincts de ceux de
// quote_lifecycle.go.
//
//   CheckContractRenewal ──> GenerateContractDoc ──┬──> SendContractEmail
//                                                   └──> SyncCrmContract

import (
	"github.com/terapps/gonveyor"
)

var ContractRenewal = gonveyor.New("contract_renewal",
	CheckContractRenewal, // root — dispatched via Seed at manifest time

	gonveyor.Wire(GenerateContractDoc,
		gonveyor.Intake(CheckContractRenewal, func(o CheckContractRenewalOutput, in *DocumentInput) {
			in.EntityID = o.ContractID
		}),
	),

	gonveyor.Wire(SendContractEmail,
		gonveyor.Intake(CheckContractRenewal, func(o CheckContractRenewalOutput, in *SendEmailInput) {
			in.To = o.ClientEmail
			in.Vars = map[string]string{"renewal_url": o.RenewalURL}
		}),
		gonveyor.Intake(GenerateContractDoc, func(o DocumentOutput, in *SendEmailInput) {
			in.DocURLs = []string{o.DocURL}
		}),
	),

	gonveyor.Wire(SyncCrmContract,
		gonveyor.Intake(CheckContractRenewal, func(o CheckContractRenewalOutput, in *SyncCrmInput) {
			in.EntityID = o.ContractID
		}),
		gonveyor.Intake(GenerateContractDoc, func(o DocumentOutput, in *SyncCrmInput) {
			in.DocURLs = []string{o.DocURL}
		}),
	),
)

var RenewalTemplate = gonveyor.NewLaunchTemplate(ContractRenewal, func(p CheckContractRenewalInput) []gonveyor.ManifestOption {
	return []gonveyor.ManifestOption{
		gonveyor.Seed(CheckContractRenewal, p),
		gonveyor.Seed(GenerateContractDoc, DocumentInput{
			DocType: "renewal",
		}),
		gonveyor.Seed(SendContractEmail, SendEmailInput{
			Template: TemplateContractRenewal,
		}),
		gonveyor.Seed(SyncCrmContract, SyncCrmInput{
			EntityType: "contract",
		}),
	}
})
