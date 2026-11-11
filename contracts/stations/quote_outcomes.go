package stations

import "github.com/terapps/gonveyor"

type QuoteOutcomeInput struct {
	QuoteID     string `validate:"required"`
	ClientEmail string `validate:"required,email"`
	Reason      string
}

type QuoteOutcomeOutput struct{}

// MarkQuoteDeclined runs only if the client refused to sign (When on AwaitSignature).
var MarkQuoteDeclined = gonveyor.Define[QuoteOutcomeInput, QuoteOutcomeOutput]("mark_quote_declined")

// NotifyPaymentFailed runs only if the payment failed (When on AwaitPayment).
var NotifyPaymentFailed = gonveyor.Define[QuoteOutcomeInput, QuoteOutcomeOutput]("notify_payment_failed")
