package stations

import "github.com/terapps/gonveyor"

type SignaturePayload struct {
	SignatureID string
	SignedAt    string
}

var AwaitSignature = gonveyor.Signal[SignaturePayload]("await_signature")
