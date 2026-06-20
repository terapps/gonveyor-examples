package handler

import (
	"context"
	"fmt"
	"log/slog"

	bp "github.com/terapps/gonveyor-examples/blueprint"
)

func PrepareContract(_ context.Context, in bp.PrepareContractInput) (bp.PrepareContractOutput, error) {
	slog.Info("preparing contract", "client_id", in.ClientID, "ref", in.ContractRef)
	return bp.PrepareContractOutput{
		ContractID: fmt.Sprintf("ctr-%s-%s", in.ClientID, in.ContractRef),
		Amount:     9900.00,
	}, nil
}

func FinalizeContract(_ context.Context, in bp.FinalizeInput) (bp.FinalizeOutput, error) {
	slog.Info("finalizing contract", "contract_id", in.ContractID, "txn", in.TxnID, "amount", in.Amount)
	return bp.FinalizeOutput{SignedURL: fmt.Sprintf("storage://contracts/%s/signed.pdf", in.ContractID)}, nil
}
