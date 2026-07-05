package handler

import (
	"context"
	"encoding/json"
	"log/slog"

	st "github.com/terapps/gonveyor-examples/contract-lifecycle/stations"
	"github.com/terapps/gonveyor/ledger"
	"github.com/terapps/gonveyor/transport/pg"
	"github.com/uptrace/bun"
)

// NewSpawnChildRenewal returns the SpawnChildRenewal handler, closing over db. It files
// a contract_renewal launch_request as a child and returns immediately — it never waits
// on it itself; ChildGate downstream is what suspends until the child is terminal.
func NewSpawnChildRenewal(db *bun.DB) func(context.Context, st.SpawnChildRenewalInput) (st.SpawnChildRenewalOutput, error) {
	return func(ctx context.Context, in st.SpawnChildRenewalInput) (st.SpawnChildRenewalOutput, error) {
		params, err := json.Marshal(st.CheckContractRenewalInput{
			ContractID:  in.ContractID,
			ClientEmail: in.ClientEmail,
		})
		if err != nil {
			return st.SpawnChildRenewalOutput{}, err
		}
		_, err = pg.CreateChildLaunchRequest(ctx, db, "contract_renewal", params,
			ledger.BlueprintIDFromContext(ctx), st.ChildGate.Key())
		return st.SpawnChildRenewalOutput{}, err
	}
}

// AfterChildDone runs once ChildGate has been signaled — the child contract_renewal is
// terminal (succeeded or failed, gonveyor doesn't distinguish here).
func AfterChildDone(_ context.Context, _ struct{}) (struct{}, error) {
	slog.Info("child renewal finished, parent resumed")
	return struct{}{}, nil
}
