package handler

import (
	"context"
	"log/slog"

	"github.com/terapps/gonveyor"
	clbp "github.com/terapps/gonveyor-examples/contract-lifecycle/blueprint"
	st "github.com/terapps/gonveyor-examples/contract-lifecycle/stations"
)

// dueContract stubs what a real CRM/contracts query would return: contracts whose end
// date falls within the renewal reminder window. A real implementation would query that
// database directly — scanning live data on every run means there's nothing to keep in
// sync in gonveyor when a renewal date changes.
type dueContract struct {
	ContractID  string
	ClientEmail string
}

func findDueContracts() []dueContract {
	return []dueContract{
		{ContractID: "contract-42", ClientEmail: "client42@example.com"},
	}
}

// NewScanContractRenewals returns the ScanContractRenewals handler, closing over gc to
// launch one contract_renewal sub-blueprint per contract found due — a direct Launch,
// not the launch_requests mailbox, since the scan already has the ManifestBuilder.
func NewScanContractRenewals(gc *gonveyor.Gonductor) func(context.Context, st.ScanRenewalsInput) (st.ScanRenewalsOutput, error) {
	return func(ctx context.Context, _ st.ScanRenewalsInput) (st.ScanRenewalsOutput, error) {
		due := findDueContracts()
		for _, c := range due {
			manifest, err := clbp.RenewalLauncher.Manifest(st.CheckContractRenewalInput{
				ContractID:  c.ContractID,
				ClientEmail: c.ClientEmail,
			})
			if err != nil {
				return st.ScanRenewalsOutput{}, err
			}
			if err := gc.Launch(ctx, manifest); err != nil {
				return st.ScanRenewalsOutput{}, err
			}
		}
		slog.Info("contract renewal scan complete", "found", len(due))
		return st.ScanRenewalsOutput{Found: len(due)}, nil
	}
}
