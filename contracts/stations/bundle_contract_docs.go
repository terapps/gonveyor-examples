package stations

import (
	"github.com/terapps/gonveyor"
)

type BundleContractDocsInput struct {
	ContractID  string
	ClientEmail string
	DocURLs     []string
}
type BundleContractDocsOutput struct {
	ContractID  string
	ClientEmail string
	DocURLs     []string
}

var BundleContractDocs = gonveyor.Define[BundleContractDocsInput, BundleContractDocsOutput]("bundle_contract_docs")
