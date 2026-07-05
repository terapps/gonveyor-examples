package blueprint

// Workflow: Démo parent/enfant
//
// SpawnChildRenewal lance un contract_renewal comme enfant sans jamais attendre dessus
// lui-même — le blocage se fait de façon asynchrone via ChildGate (un Signal), complété
// automatiquement par le poller "awaits" du worker une fois l'enfant terminé (succès ou
// échec, peu importe). AfterChildDone ne se dispatch qu'une fois ce signal reçu.
//
//   SpawnChildRenewal ──> ChildGate (signal) ──> AfterChildDone

import (
	"github.com/terapps/gonveyor"
	st "github.com/terapps/gonveyor-examples/contract-lifecycle/stations"
)

var ParentChildDemo = gonveyor.New("parent_child_demo",
	st.SpawnChildRenewal, // root — dispatched via Seed at manifest time
	gonveyor.Wire(st.ChildGate, gonveyor.After[struct{}](st.SpawnChildRenewal)),
	gonveyor.Wire(st.AfterChildDone, gonveyor.After[struct{}](st.ChildGate)),
)

var ParentChildDemoLauncher = gonveyor.NewManifestBuilder(ParentChildDemo, func(p st.SpawnChildRenewalInput) []gonveyor.ManifestOption {
	return []gonveyor.ManifestOption{gonveyor.Seed(st.SpawnChildRenewal, p)}
})
