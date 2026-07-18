// Package simple is the whole "send a welcome email" example in one file: its station,
// its one-step blueprint, and its handler.
package simple

import (
	"github.com/terapps/gonveyor"
)

// --- Station ---

type WelcomeInput struct {
	UserID string `validate:"required"`
	Email  string `validate:"required,email"`
}

type WelcomeOutput struct {
	SentAt string
}

var SendWelcome = gonveyor.Define[WelcomeInput, WelcomeOutput]("send_welcome")

// --- Blueprint ---

var SimpleDispatch = gonveyor.New("simple_dispatch", SendWelcome)

var Template = gonveyor.NewLaunchTemplate(SimpleDispatch, func(p WelcomeInput) []gonveyor.ManifestOption {
	return []gonveyor.ManifestOption{gonveyor.Seed(SendWelcome, p)}
})

func Manifest(userID, email string) (gonveyor.BlueprintManifest, error) {
	return Template.Manifest(WelcomeInput{UserID: userID, Email: email})
}
