package stations

import "github.com/terapps/gonveyor"

type ScanRenewalsInput struct{}
type ScanRenewalsOutput struct {
	Found int
}

// ScanContractRenewals is the root of contract_renewal_scan — one recurring schedule that
// files one contract_renewal launch_request per contract found due.
var ScanContractRenewals = gonveyor.Define[ScanRenewalsInput, ScanRenewalsOutput]("scan_contract_renewals")
