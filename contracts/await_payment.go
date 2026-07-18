package contracts

import "github.com/terapps/gonveyor"

type PaymentPayload struct {
	TxnID  string
	Amount float64
	PaidAt string
}

// AwaitPayment is a signal — no handler; completed by Gonductor.Send, not a worker.
var AwaitPayment = gonveyor.Signal[PaymentPayload]("await_payment")
