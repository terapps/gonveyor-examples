package stations

import "github.com/terapps/gonveyor"

type SyncCrmInput struct {
	EntityType string // "quote" or "contract"
	EntityID   string
	Metadata   map[string]string
	DocURLs    []string
}
type SyncCrmOutput struct{}

var SyncCrmQuote = gonveyor.Define[SyncCrmInput, SyncCrmOutput]("sync_crm_quote")
var SyncCrmContract = gonveyor.Define[SyncCrmInput, SyncCrmOutput]("sync_crm_contract")
