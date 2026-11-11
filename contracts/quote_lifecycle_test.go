package contracts

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	st "github.com/terapps/gonveyor-examples/contracts/stations"
)

// The two signals pick the outcome: a signed + paid quote becomes a contract, a refusal or a
// failed payment runs its own branch and skips CreateContract (and everything after it).
func TestQuoteLifecycle_Outcomes(t *testing.T) {
	tests := []struct {
		name      string
		signature st.SignaturePayload
		payment   st.PaymentPayload
		want      map[string]bool // station key -> runs
	}{
		{"signed and paid", st.SignaturePayload{SignatureID: "s1"}, st.PaymentPayload{TxnID: "t1"},
			map[string]bool{"create_contract": true, "mark_quote_declined": false, "notify_payment_failed": false}},
		{"signature refused", st.SignaturePayload{Declined: true, Reason: "too expensive"}, st.PaymentPayload{TxnID: "t1"},
			map[string]bool{"create_contract": false, "mark_quote_declined": true, "notify_payment_failed": false}},
		{"payment failed", st.SignaturePayload{SignatureID: "s1"}, st.PaymentPayload{Failed: true, Reason: "card declined"},
			map[string]bool{"create_contract": false, "mark_quote_declined": false, "notify_payment_failed": true}},
		{"both", st.SignaturePayload{Declined: true}, st.PaymentPayload{Failed: true},
			map[string]bool{"create_contract": false, "mark_quote_declined": true, "notify_payment_failed": true}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sig, err := json.Marshal(tc.signature)
			require.NoError(t, err)
			pay, err := json.Marshal(tc.payment)
			require.NoError(t, err)
			outputs := map[string][]json.RawMessage{
				st.AwaitSignature.Key(): {sig},
				st.AwaitPayment.Key():   {pay},
			}
			for key, want := range tc.want {
				run, err := QuoteLifecycle.Node(key).ShouldRun(outputs)
				require.NoError(t, err)
				assert.Equal(t, want, run, key)
			}
		})
	}
}

func TestQuoteLifecycleTemplate_Builds(t *testing.T) {
	m, err := QuoteLifecycleTemplate.Manifest(Params{
		QuoteID: "q1", ClientEmail: "a@b.com", Amount: 100,
		QuoteDocTypes: []string{"proposal"}, ContractDocTypes: []string{"contract"},
	})
	require.NoError(t, err)
	keys := map[string]bool{}
	for _, u := range m.Units {
		keys[u.Key] = true
	}
	assert.True(t, keys["mark_quote_declined"])
	assert.True(t, keys["notify_payment_failed"])
}
