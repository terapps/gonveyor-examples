package blueprint

// Workflow: Renouvellement de contrat
//
// Indépendant de quote_lifecycle (déclenché séparément — voir pg.CreateSchedule pour un
// déclenchement récurrent natif, cmd/publisher schedule-contract-renewal), mais réutilise
// trois stations telles quelles — GenerateContractDoc, SendContractEmail, SyncCrmContract
// — câblées différemment ici via des Wire() distincts de ceux de quote_lifecycle.go.
//
//   CheckContractRenewal ──> GenerateContractDoc ──┬──> SendContractEmail
//                                                   └──> SyncCrmContract

import (
	"github.com/terapps/gonveyor"
	st "github.com/terapps/gonveyor-examples/contract-lifecycle/stations"
	"github.com/terapps/gonveyor/ledger"
)

var ContractRenewal = gonveyor.New("contract_renewal",
	st.CheckContractRenewal, // root — dispatched via Seed at manifest time

	gonveyor.Wire(st.GenerateContractDoc,
		gonveyor.Intake(st.CheckContractRenewal, func(o st.CheckContractRenewalOutput, in *st.DocumentInput) {
			in.EntityID = o.ContractID
		}),
	),

	gonveyor.Wire(st.SendContractEmail,
		gonveyor.Intake(st.CheckContractRenewal, func(o st.CheckContractRenewalOutput, in *st.SendEmailInput) {
			in.To = o.ClientEmail
			in.Vars = map[string]string{"renewal_url": o.RenewalURL}
		}),
		gonveyor.Intake(st.GenerateContractDoc, func(o st.DocumentOutput, in *st.SendEmailInput) {
			in.DocURLs = []string{o.DocURL}
		}),
	),

	gonveyor.Wire(st.SyncCrmContract,
		gonveyor.Intake(st.CheckContractRenewal, func(o st.CheckContractRenewalOutput, in *st.SyncCrmInput) {
			in.EntityID = o.ContractID
		}),
		gonveyor.Intake(st.GenerateContractDoc, func(o st.DocumentOutput, in *st.SyncCrmInput) {
			in.DocURLs = []string{o.DocURL}
		}),
	),
)

var RenewalLauncher = gonveyor.NewManifestBuilder(ContractRenewal, func(p st.CheckContractRenewalInput) []gonveyor.ManifestOption {
	return []gonveyor.ManifestOption{
		gonveyor.Seed(st.CheckContractRenewal, p),
		gonveyor.Seed(st.GenerateContractDoc, st.DocumentInput{
			DocType: "renewal",
		}),
		gonveyor.Seed(st.SendContractEmail, st.SendEmailInput{
			Template: st.TemplateContractRenewal,
		}),
		gonveyor.Seed(st.SyncCrmContract, st.SyncCrmInput{
			EntityType: "contract",
		}),
	}
})

func RenewalManifest(p st.CheckContractRenewalInput) (ledger.BlueprintManifest, error) {
	return RenewalLauncher.Manifest(p)
}
