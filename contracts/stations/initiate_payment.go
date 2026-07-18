package stations

import (
	"github.com/terapps/gonveyor"
)

type InitiatePaymentInput struct {
	QuoteID     string
	ClientEmail string
	Amount      float64
	DocURLs     []string
}
type InitiatePaymentOutput struct {
	ProcessID  string
	PaymentURL string
}

var InitiatePayment = gonveyor.Define[InitiatePaymentInput, InitiatePaymentOutput]("initiate_payment")
