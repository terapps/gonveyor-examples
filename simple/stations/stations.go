package stations

import "github.com/terapps/gonveyor"

type WelcomeInput struct {
	UserID string `validate:"required"`
	Email  string `validate:"required,email"`
}

type WelcomeOutput struct {
	SentAt string
}

var SendWelcome = gonveyor.Define[WelcomeInput, WelcomeOutput]("send_welcome")
