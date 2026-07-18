package stations

import (
	"github.com/terapps/gonveyor"
)

type DocumentInput struct {
	EntityID string // QuoteID or ContractID
	DocType  string // e.g. "proposal", "pricing", "contract", "annex_a"
}
type DocumentOutput struct {
	DocURL   string
	DocType  string
	EntityID string
}

// GenerateQuoteDoc and GenerateContractDoc share DocumentInput/Output and the same handler.
// They are distinct stations because a node key must be unique within a blueprint.
var GenerateQuoteDoc = gonveyor.Define[DocumentInput, DocumentOutput]("generate_quote_doc", gonveyor.WithRoute("tasks.document"))
var GenerateContractDoc = gonveyor.Define[DocumentInput, DocumentOutput]("generate_contract_doc", gonveyor.WithRoute("tasks.document"))
