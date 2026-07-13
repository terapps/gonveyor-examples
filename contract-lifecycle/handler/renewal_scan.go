package handler

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

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

// findDueContracts fakes a variable-size result set (1-50 contracts) to exercise the
// scan's fan-out under load, instead of always spawning a single sub-blueprint.
func findDueContracts() []dueContract {
	n := rand.Intn(50) + 1
	due := make([]dueContract, n)
	for i := range due {
		due[i] = dueContract{
			ContractID:  fmt.Sprintf("contract-%d", i),
			ClientEmail: fmt.Sprintf("client%d@example.com", i),
		}
	}
	return due
}

// NewScanContractRenewals returns the ScanContractRenewals handler, closing over gc to
// launch one contract_renewal sub-blueprint per contract found due — a direct Launch,
// not the launch_requests mailbox, since the scan already has the LaunchTemplate.
func NewScanContractRenewals(gc *gonveyor.Gonductor) func(context.Context, st.ScanRenewalsInput) (st.ScanRenewalsOutput, error) {
	return func(ctx context.Context, _ st.ScanRenewalsInput) (st.ScanRenewalsOutput, error) {
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
		due := findDueContracts()
		for _, c := range due {
			manifest, err := clbp.RenewalTemplate.Manifest(st.CheckContractRenewalInput{
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
