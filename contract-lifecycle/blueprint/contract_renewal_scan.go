package blueprint

// Its own blueprint, independent of ContractRenewal: one recurring schedule (see
// cmd/publisher schedule-contract-renewal-scan) dispatches ScanContractRenewals, whose
// handler reads live contract data and files one "contract_renewal" launch_request per
// contract found due — so there's nothing to edit here when a renewal date changes, and
// no per-contract schedule to manage.

import (
	"github.com/terapps/gonveyor"
	st "github.com/terapps/gonveyor-examples/contract-lifecycle/stations"
)

var ContractRenewalScan = gonveyor.New("contract_renewal_scan", st.ScanContractRenewals)

var ScanLauncher = gonveyor.NewManifestBuilder(ContractRenewalScan, func(p st.ScanRenewalsInput) []gonveyor.ManifestOption {
	return []gonveyor.ManifestOption{gonveyor.Seed(st.ScanContractRenewals, p)}
})
