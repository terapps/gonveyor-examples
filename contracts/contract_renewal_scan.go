package contracts

// Its own blueprint, independent of ContractRenewal: one recurring schedule (see
// cmd/publisher schedule-contract-renewal-scan) dispatches ScanContractRenewals, whose
// handler reads live contract data and files one "contract_renewal" launch_request per
// contract found due — so there's nothing to edit here when a renewal date changes, and
// no per-contract schedule to manage.

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/terapps/gonveyor"
	st "github.com/terapps/gonveyor-examples/contracts/stations"
)

var ContractRenewalScan = gonveyor.New("contract_renewal_scan", st.ScanContractRenewals)

var ScanTemplate = gonveyor.NewLaunchTemplate(ContractRenewalScan, func(p st.ScanRenewalsInput) []gonveyor.ManifestOption {
	return []gonveyor.ManifestOption{gonveyor.Seed(st.ScanContractRenewals, p)}
})

// --- Worker ---
//
// The scan handler lives with its blueprint (not in stations/) because it launches
// contract_renewal sub-blueprints via RenewalTemplate — it orchestrates, rather than
// executing a single unit.

// dueContract stubs what a real CRM/contracts query would return: contracts whose end date
// falls within the renewal reminder window. A real implementation would query that database
// directly — scanning live data means there's nothing to keep in sync when a date changes.
type dueContract struct {
	ContractID  string
	ClientEmail string
}

// findDueContracts fakes a variable-size result set (1-50) to exercise the scan's fan-out.
func findDueContracts() []dueContract {
	n := rand.Intn(50) + 1
	due := make([]dueContract, n)
	for i := range due {
		due[i] = dueContract{ContractID: fmt.Sprintf("contract-%d", i), ClientEmail: fmt.Sprintf("client%d@example.com", i)}
	}
	return due
}

// NewScanContractRenewals closes over gc to launch one contract_renewal sub-blueprint per
// due contract — a single LaunchBatch (one transaction), since the scan already holds the
// LaunchTemplate, rather than N launch_requests round-trips.
func NewScanContractRenewals(gc *gonveyor.Gonductor) func(context.Context, st.ScanRenewalsInput) (st.ScanRenewalsOutput, error) {
	return func(ctx context.Context, _ st.ScanRenewalsInput) (st.ScanRenewalsOutput, error) {
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
		due := findDueContracts()
		manifests := make([]gonveyor.BlueprintManifest, len(due))
		for i, c := range due {
			manifest, err := RenewalTemplate.Manifest(st.CheckContractRenewalInput{ContractID: c.ContractID, ClientEmail: c.ClientEmail})
			if err != nil {
				return st.ScanRenewalsOutput{}, err
			}
			manifests[i] = manifest
		}
		if err := gc.LaunchBatch(ctx, manifests); err != nil {
			return st.ScanRenewalsOutput{}, err
		}
		slog.Info("contract renewal scan complete", "found", len(due))
		return st.ScanRenewalsOutput{Found: len(due)}, nil
	}
}
