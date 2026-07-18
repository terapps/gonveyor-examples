package blueprint

import (
	"github.com/terapps/gonveyor"
	st "github.com/terapps/gonveyor-examples/simple/stations"
)

var SimpleDispatch = gonveyor.New("simple_dispatch", st.SendWelcome)

var Template = gonveyor.NewLaunchTemplate(SimpleDispatch, func(p st.WelcomeInput) []gonveyor.ManifestOption {
	return []gonveyor.ManifestOption{gonveyor.Seed(st.SendWelcome, p)}
})

func Manifest(userID, email string) (gonveyor.BlueprintManifest, error) {
	return Template.Manifest(st.WelcomeInput{UserID: userID, Email: email})
}
