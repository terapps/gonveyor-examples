package blueprint

import (
	"github.com/terapps/gonveyor"
	st "github.com/terapps/gonveyor-examples/simple/stations"
	"github.com/terapps/gonveyor/core"
)

var SimpleDispatch = gonveyor.New("simple_dispatch", st.SendWelcome)

var Launcher = gonveyor.NewManifestBuilder(SimpleDispatch, func(p st.WelcomeInput) []gonveyor.ManifestOption {
	return []gonveyor.ManifestOption{gonveyor.Seed(st.SendWelcome, p)}
})

func Manifest(userID, email string) (core.BlueprintManifest, error) {
	return Launcher.Manifest(st.WelcomeInput{UserID: userID, Email: email})
}
