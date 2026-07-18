package stations

import "github.com/terapps/gonveyor"

type PaymentPayload struct {
	TxnID  string
	Amount float64
	PaidAt string
}

var AwaitPayment = gonveyor.Signal[PaymentPayload]("await_payment")
