package stations

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"strings"
	"time"

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

// --- Worker ---

type emailTemplate struct {
	subject string
	body    string // vars interpolated manually
}

var templates = map[EmailTemplate]emailTemplate{
	TemplateSignatureRequest: {subject: "Votre devis est prêt à signer", body: "Bonjour,\n\nVeuillez signer votre devis :\n{signature_url}\n\nCordialement"},
	TemplateContractSigned:   {subject: "Votre contrat est disponible", body: "Bonjour,\n\nVotre contrat a été finalisé :\n{doc_urls}\n\nCordialement"},
	TemplateContractRenewal:  {subject: "Votre contrat arrive à échéance", body: "Bonjour,\n\nRenouvelez votre contrat :\n{renewal_url}\n\nCordialement"},
}

func HandleEmail(_ context.Context, in SendEmailInput) (SendEmailOutput, error) {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	tmpl, ok := templates[in.Template]
	if !ok {
		return SendEmailOutput{}, fmt.Errorf("unknown email template %q", in.Template)
	}
	vars := in.Vars
	if vars == nil {
		vars = map[string]string{}
	}
	if len(in.DocURLs) > 0 {
		vars["doc_urls"] = strings.Join(in.DocURLs, "\n")
	}
	body := tmpl.body
	for k, v := range vars {
		body = strings.ReplaceAll(body, "{"+k+"}", v)
	}
	slog.Info("sending email", "to", in.To, "template", in.Template, "subject", tmpl.subject)
	slog.Debug("email body", "body", body)
	return SendEmailOutput{}, nil
}
