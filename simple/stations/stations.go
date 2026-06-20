package stations

import "github.com/terapps/gonveyor/blueprint"

type WelcomeInput struct {
	UserID string
	Email  string
}

type WelcomeOutput struct {
	SentAt string
}

var SendWelcome = blueprint.Define[WelcomeInput, WelcomeOutput]("send_welcome")
