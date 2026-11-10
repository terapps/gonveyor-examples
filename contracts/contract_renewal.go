package contracts

// Workflow: Renouvellement de contrat
//
// Indépendant de quote_lifecycle (déclenché par contract_renewal_scan, voir
// contract_renewal_scan.go — un seul schedule récurrent, pas un par contrat), mais
// réutilise trois stations telles quelles — GenerateContractDoc, SendContractEmail,
// SyncCrmContract — câblées différemment ici via des chaînes Intake/After distinctes de celles de
// quote_lifecycle.go.
//
//   CheckContractRenewal ──> GenerateContractDoc ──┬──> SendContractEmail
//                                                   └──> SyncCrmContract

import (
	"github.com/terapps/gonveyor"
	st "github.com/terapps/gonveyor-examples/contracts/stations"
)

var ContractRenewal = gonveyor.New("contract_renewal",
	st.CheckContractRenewal, // root — dispatched via Seed at manifest time

	st.GenerateContractDoc.Intake(st.CheckContractRenewal, func(o st.CheckContractRenewalOutput, in *st.DocumentInput) {
		in.EntityID = o.ContractID
	}),

	st.SendContractEmail.
		Intake(st.CheckContractRenewal, func(o st.CheckContractRenewalOutput, in *st.SendEmailInput) {
			in.To = o.ClientEmail
			in.Vars = map[string]string{"renewal_url": o.RenewalURL}
		}).
		Intake(st.GenerateContractDoc, func(o st.DocumentOutput, in *st.SendEmailInput) {
			in.DocURLs = []string{o.DocURL}
		}),

	st.SyncCrmContract.
		Intake(st.CheckContractRenewal, func(o st.CheckContractRenewalOutput, in *st.SyncCrmInput) {
			in.EntityID = o.ContractID
		}).
		Intake(st.GenerateContractDoc, func(o st.DocumentOutput, in *st.SyncCrmInput) {
			in.DocURLs = []string{o.DocURL}
		}),
)

var RenewalTemplate = gonveyor.NewLaunchTemplate(ContractRenewal, func(p st.CheckContractRenewalInput) []gonveyor.ManifestOption {
	return []gonveyor.ManifestOption{
		st.CheckContractRenewal.Seed(p),
		st.GenerateContractDoc.Seed(st.DocumentInput{
			DocType: "renewal",
		}),
		st.SendContractEmail.Seed(st.SendEmailInput{
			Template: st.TemplateContractRenewal,
		}),
		st.SyncCrmContract.Seed(st.SyncCrmInput{
			EntityType: "contract",
		}),
	}
})
