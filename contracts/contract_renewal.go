package contracts

// Workflow: Renouvellement de contrat
//
// Indépendant de quote_lifecycle (déclenché par contract_renewal_scan, voir
// contract_renewal_scan.go — un seul schedule récurrent, pas un par contrat), mais
// réutilise trois stations telles quelles — GenerateContractDoc, SendContractEmail,
// SyncCrmContract — câblées différemment ici via des chaînes Intake/After distinctes de celles de
// quote_lifecycle.go. La décision du client (signal) aiguille vers une branche When ;
// SyncCrmContract joint les deux (seule l'une a tourné, l'autre est skipped).
//
//   CheckContractRenewal ──> GenerateContractDoc ──> SendContractEmail ──> AwaitRenewalDecision
//     AwaitRenewalDecision ──When accepté──> RenewContract ──┬──> SyncCrmContract
//                          ──When refusé───> CloseContract ──┘

import (
	"fmt"

	"github.com/terapps/gonveyor"
	st "github.com/terapps/gonveyor-examples/contracts/stations"
)

func accepted(d st.RenewalDecision) bool { return d.Accepted }
func declined(d st.RenewalDecision) bool { return !d.Accepted }

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

	st.AwaitRenewalDecision.After(st.SendContractEmail),

	st.RenewContract.
		Intake(st.CheckContractRenewal, func(o st.CheckContractRenewalOutput, in *st.RenewContractInput) {
			in.ContractID = o.ContractID
		}).
		When(st.AwaitRenewalDecision, accepted),

	st.CloseContract.
		Intake(st.CheckContractRenewal, func(o st.CheckContractRenewalOutput, in *st.CloseContractInput) {
			in.ContractID = o.ContractID
		}).
		Intake(st.AwaitRenewalDecision, func(d st.RenewalDecision, in *st.CloseContractInput) {
			in.Reason = d.Reason
		}).
		When(st.AwaitRenewalDecision, declined),

	st.SyncCrmContract.
		Intake(st.CheckContractRenewal, func(o st.CheckContractRenewalOutput, in *st.SyncCrmInput) {
			in.EntityID = o.ContractID
		}).
		Intake(st.GenerateContractDoc, func(o st.DocumentOutput, in *st.SyncCrmInput) {
			in.DocURLs = []string{o.DocURL}
		}).
		Intake(st.AwaitRenewalDecision, func(d st.RenewalDecision, in *st.SyncCrmInput) {
			in.Metadata = map[string]string{"renewal_accepted": fmt.Sprint(d.Accepted)}
		}).
		After(st.RenewContract, st.CloseContract), // join: runs whichever branch ran
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
