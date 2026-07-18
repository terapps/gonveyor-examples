package stations

import (
	"github.com/terapps/gonveyor"
)

type EmailTemplate string

const (
	TemplateSignatureRequest EmailTemplate = "signature_request"
	TemplateContractSigned   EmailTemplate = "contract_signed"
	TemplateContractRenewal  EmailTemplate = "contract_renewal"
)

type SendEmailInput struct {
	To       string
	Template EmailTemplate
	Vars     map[string]string // template variables
	DocURLs  []string          // attachments / links
}
type SendEmailOutput struct{}

// SendQuoteEmail and SendContractEmail share SendEmailInput and the same handler.
var SendQuoteEmail = gonveyor.Define[SendEmailInput, SendEmailOutput]("send_quote_email")
var SendContractEmail = gonveyor.Define[SendEmailInput, SendEmailOutput]("send_contract_email")
