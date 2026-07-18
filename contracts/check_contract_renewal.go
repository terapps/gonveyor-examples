package contracts

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/terapps/gonveyor"
)

type CheckContractRenewalInput struct {
	ContractID  string `validate:"required"`
	ClientEmail string `validate:"required,email"`
}
type CheckContractRenewalOutput struct {
	ContractID  string
	ClientEmail string
	RenewalURL  string
}

// CheckContractRenewal is the root of the contract_renewal blueprint. SendContractEmail is
// reused as-is from quote_lifecycle — same station, different Wire().
var CheckContractRenewal = gonveyor.Define[CheckContractRenewalInput, CheckContractRenewalOutput]("check_contract_renewal")

// --- Worker ---

func HandleCheckContractRenewal(_ context.Context, in CheckContractRenewalInput) (CheckContractRenewalOutput, error) {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	slog.Info("checking contract renewal", "contract_id", in.ContractID)
	return CheckContractRenewalOutput{ContractID: in.ContractID, ClientEmail: in.ClientEmail, RenewalURL: fmt.Sprintf("https://renew.example.com/contract/%s", in.ContractID)}, nil
}
