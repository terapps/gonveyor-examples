package contracts

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/terapps/gonveyor"
)

type ScanRenewalsInput struct{}
type ScanRenewalsOutput struct {
	Found int
}

// ScanContractRenewals is the root of contract_renewal_scan — one recurring schedule that
// files one contract_renewal launch_request per contract found due.
var ScanContractRenewals = gonveyor.Define[ScanRenewalsInput, ScanRenewalsOutput]("scan_contract_renewals")

// --- Worker ---

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
func NewScanContractRenewals(gc *gonveyor.Gonductor) func(context.Context, ScanRenewalsInput) (ScanRenewalsOutput, error) {
	return func(ctx context.Context, _ ScanRenewalsInput) (ScanRenewalsOutput, error) {
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
		due := findDueContracts()
		manifests := make([]gonveyor.BlueprintManifest, len(due))
		for i, c := range due {
			manifest, err := RenewalTemplate.Manifest(CheckContractRenewalInput{ContractID: c.ContractID, ClientEmail: c.ClientEmail})
			if err != nil {
				return ScanRenewalsOutput{}, err
			}
			manifests[i] = manifest
		}
		if err := gc.LaunchBatch(ctx, manifests); err != nil {
			return ScanRenewalsOutput{}, err
		}
		slog.Info("contract renewal scan complete", "found", len(due))
		return ScanRenewalsOutput{Found: len(due)}, nil
	}
}
