package stations

import "github.com/terapps/gonveyor"

type SignaturePayload struct {
	SignatureID string
	SignedAt    string
}

// AwaitSignature is a signal — no handler; completed by Gonductor.Send, not a worker.
var AwaitSignature = gonveyor.Signal[SignaturePayload]("await_signature")
