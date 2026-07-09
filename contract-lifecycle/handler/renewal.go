package handler

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	st "github.com/terapps/gonveyor-examples/contract-lifecycle/stations"
)

// CheckContractRenewal simulates detecting that a contract is nearing its expiry date
// and needs a renewal reminder. A real implementation would query contract end dates.
func CheckContractRenewal(_ context.Context, in st.CheckContractRenewalInput) (st.CheckContractRenewalOutput, error) {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	slog.Info("checking contract renewal", "contract_id", in.ContractID)
	return st.CheckContractRenewalOutput{
		ContractID:  in.ContractID,
		ClientEmail: in.ClientEmail,
		RenewalURL:  fmt.Sprintf("https://renew.example.com/contract/%s", in.ContractID),
	}, nil
}
