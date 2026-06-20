package blueprint

import (
	"github.com/terapps/gonveyor"
	"github.com/terapps/gonveyor/blueprint"
	"github.com/terapps/gonveyor/ledger"
)

type WelcomeInput struct {
	UserID string
	Email  string
}

type WelcomeOutput struct {
	SentAt string
}

var SendWelcome = blueprint.Define[WelcomeInput, WelcomeOutput]("send_welcome")

var SimpleDispatch = blueprint.New("simple_dispatch", SendWelcome)

func SimpleManifest(userID, email string) (ledger.BlueprintManifest, error) {
	return SimpleDispatch.Manifest(
		gonveyor.Seed(SendWelcome, WelcomeInput{UserID: userID, Email: email}),
	)
}
