package blueprint

import (
	"github.com/terapps/gonveyor"
	bp "github.com/terapps/gonveyor/blueprint"
	"github.com/terapps/gonveyor/ledger"
	st "github.com/terapps/gonveyor-examples/simple/stations"
)

var SimpleDispatch = bp.New("simple_dispatch", st.SendWelcome)

func Manifest(userID, email string) (ledger.BlueprintManifest, error) {
	return SimpleDispatch.Manifest(
		gonveyor.Seed(st.SendWelcome, st.WelcomeInput{UserID: userID, Email: email}),
	)
}
