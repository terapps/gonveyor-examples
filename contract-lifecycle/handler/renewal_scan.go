package handler

import (
	"context"
	"encoding/json"
	"log/slog"

	st "github.com/terapps/gonveyor-examples/contract-lifecycle/stations"
	"github.com/terapps/gonveyor/transport/pg"
	"github.com/uptrace/bun"
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

// NewScanContractRenewals returns the ScanContractRenewals handler, closing over db to
// file one contract_renewal launch_request per contract found due.
func NewScanContractRenewals(db *bun.DB) func(context.Context, st.ScanRenewalsInput) (st.ScanRenewalsOutput, error) {
	return func(ctx context.Context, _ st.ScanRenewalsInput) (st.ScanRenewalsOutput, error) {
		due := findDueContracts()
		for _, c := range due {
			params, err := json.Marshal(st.CheckContractRenewalInput{
				ContractID:  c.ContractID,
				ClientEmail: c.ClientEmail,
			})
			if err != nil {
				return st.ScanRenewalsOutput{}, err
			}
			if _, err := pg.CreateLaunchRequest(ctx, db, "contract_renewal", params); err != nil {
				return st.ScanRenewalsOutput{}, err
			}
		}
		slog.Info("contract renewal scan complete", "found", len(due))
		return st.ScanRenewalsOutput{Found: len(due)}, nil
	}
}
