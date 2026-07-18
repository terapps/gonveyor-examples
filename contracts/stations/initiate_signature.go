package stations

import "github.com/terapps/gonveyor"

type InitiateSignatureInput struct {
	QuoteID     string
	ClientEmail string
	DocURLs     []string
}
type InitiateSignatureOutput struct {
	ProcessID    string
	SignatureURL string
}

var InitiateSignature = gonveyor.Define[InitiateSignatureInput, InitiateSignatureOutput]("initiate_signature")
