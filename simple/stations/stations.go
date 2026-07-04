package stations

import "github.com/terapps/gonveyor"

type WelcomeInput struct {
	UserID string
	Email  string
}

type WelcomeOutput struct {
	SentAt string
}

var SendWelcome = gonveyor.Define[WelcomeInput, WelcomeOutput]("send_welcome")
