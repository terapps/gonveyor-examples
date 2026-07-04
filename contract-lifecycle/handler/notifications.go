package handler

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	st "github.com/terapps/gonveyor-examples/contract-lifecycle/stations"
)

type emailTemplate struct {
	subject string
	body    string // Go template syntax, vars interpolated manually
}

var templates = map[st.EmailTemplate]emailTemplate{
	st.TemplateSignatureRequest: {
		subject: "Votre devis est prêt à signer",
		body:    "Bonjour,\n\nVeuillez signer votre devis en cliquant sur le lien suivant :\n{signature_url}\n\nCordialement",
	},
	st.TemplateContractSigned: {
		subject: "Votre contrat est disponible",
		body:    "Bonjour,\n\nVotre contrat a été finalisé. Vous trouverez vos documents ci-joints :\n{doc_urls}\n\nCordialement",
	},
	st.TemplateContractRenewal: {
		subject: "Votre contrat arrive à échéance",
		body:    "Bonjour,\n\nVotre contrat arrive bientôt à échéance. Renouvelez-le en cliquant sur le lien suivant :\n{renewal_url}\n\nCordialement",
	},
}

func SendEmail(_ context.Context, in st.SendEmailInput) (st.SendEmailOutput, error) {
	tmpl, ok := templates[in.Template]
	if !ok {
		return st.SendEmailOutput{}, fmt.Errorf("unknown email template %q", in.Template)
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

	slog.Info("sending email",
		"to", in.To,
		"template", in.Template,
		"subject", tmpl.subject,
	)
	slog.Debug("email body", "body", body)

	return st.SendEmailOutput{}, nil
}
