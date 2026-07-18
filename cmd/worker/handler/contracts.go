package handler

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"strings"
	"time"

	"github.com/terapps/gonveyor"
	"github.com/terapps/gonveyor-examples/contracts"
	clst "github.com/terapps/gonveyor-examples/contracts/stations"
)

func HandleDocument(_ context.Context, in clst.DocumentInput) (clst.DocumentOutput, error) {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	slog.Info("generating document", "entity_id", in.EntityID, "doc_type", in.DocType)
	return clst.DocumentOutput{DocURL: fmt.Sprintf("storage://%s/%s.pdf", in.EntityID, in.DocType), DocType: in.DocType, EntityID: in.EntityID}, nil
}

func HandleCrm(_ context.Context, in clst.SyncCrmInput) (clst.SyncCrmOutput, error) {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	slog.Info("syncing to CRM", "entity_type", in.EntityType, "entity_id", in.EntityID, "metadata", in.Metadata, "docs", len(in.DocURLs))
	return clst.SyncCrmOutput{}, nil
}

type emailTemplate struct {
	subject string
	body    string // vars interpolated manually
}

var templates = map[clst.EmailTemplate]emailTemplate{
	clst.TemplateSignatureRequest: {subject: "Votre devis est prêt à signer", body: "Bonjour,\n\nVeuillez signer votre devis :\n{signature_url}\n\nCordialement"},
	clst.TemplateContractSigned:   {subject: "Votre contrat est disponible", body: "Bonjour,\n\nVotre contrat a été finalisé :\n{doc_urls}\n\nCordialement"},
	clst.TemplateContractRenewal:  {subject: "Votre contrat arrive à échéance", body: "Bonjour,\n\nRenouvelez votre contrat :\n{renewal_url}\n\nCordialement"},
}

func HandleEmail(_ context.Context, in clst.SendEmailInput) (clst.SendEmailOutput, error) {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	tmpl, ok := templates[in.Template]
	if !ok {
		return clst.SendEmailOutput{}, fmt.Errorf("unknown email template %q", in.Template)
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
	return clst.SendEmailOutput{}, nil
}

func HandleInitiateSignature(_ context.Context, in clst.InitiateSignatureInput) (clst.InitiateSignatureOutput, error) {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	slog.Info("initiating signature process", "quote_id", in.QuoteID, "docs", len(in.DocURLs))
	processID := fmt.Sprintf("sig-%s", in.QuoteID)
	return clst.InitiateSignatureOutput{ProcessID: processID, SignatureURL: fmt.Sprintf("https://sign.example.com/process/%s", processID)}, nil
}

func HandleInitiatePayment(_ context.Context, in clst.InitiatePaymentInput) (clst.InitiatePaymentOutput, error) {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	slog.Info("initiating payment process", "quote_id", in.QuoteID, "client_email", in.ClientEmail, "amount", in.Amount)
	processID := fmt.Sprintf("pay-%s", in.QuoteID)
	return clst.InitiatePaymentOutput{ProcessID: processID, PaymentURL: fmt.Sprintf("https://pay.example.com/process/%s", processID)}, nil
}

func HandleCreateContract(_ context.Context, in clst.CreateContractInput) (clst.CreateContractOutput, error) {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	slog.Info("creating contract", "quote_id", in.QuoteID, "signature_id", in.SignatureID, "txn_id", in.TxnID, "amount", in.Amount)
	return clst.CreateContractOutput{ContractID: fmt.Sprintf("ctr-%s", in.QuoteID)}, nil
}

func HandleBundleContractDocs(_ context.Context, in clst.BundleContractDocsInput) (clst.BundleContractDocsOutput, error) {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	slog.Info("bundling contract docs", "contract_id", in.ContractID, "count", len(in.DocURLs))
	return clst.BundleContractDocsOutput{ContractID: in.ContractID, ClientEmail: in.ClientEmail, DocURLs: in.DocURLs}, nil
}

func HandleCheckContractRenewal(_ context.Context, in clst.CheckContractRenewalInput) (clst.CheckContractRenewalOutput, error) {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	slog.Info("checking contract renewal", "contract_id", in.ContractID)
	return clst.CheckContractRenewalOutput{ContractID: in.ContractID, ClientEmail: in.ClientEmail, RenewalURL: fmt.Sprintf("https://renew.example.com/contract/%s", in.ContractID)}, nil
}

// --- Scan launcher ---

// dueContract stubs what a real CRM/contracts query would return: contracts whose end date
// falls within the renewal reminder window.
type dueContract struct {
	ContractID  string
	ClientEmail string
}

func findDueContracts() []dueContract {
	n := rand.Intn(50) + 1
	due := make([]dueContract, n)
	for i := range due {
		due[i] = dueContract{ContractID: fmt.Sprintf("contract-%d", i), ClientEmail: fmt.Sprintf("client%d@example.com", i)}
	}
	return due
}

// NewScanContractRenewals closes over gc to launch one contract_renewal sub-blueprint per
// due contract — a single LaunchBatch, since the scan already holds the LaunchTemplate.
func NewScanContractRenewals(gc *gonveyor.Gonductor) func(context.Context, clst.ScanRenewalsInput) (clst.ScanRenewalsOutput, error) {
	return func(ctx context.Context, _ clst.ScanRenewalsInput) (clst.ScanRenewalsOutput, error) {
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
		due := findDueContracts()
		manifests := make([]gonveyor.BlueprintManifest, len(due))
		for i, c := range due {
			manifest, err := contracts.RenewalTemplate.Manifest(clst.CheckContractRenewalInput{ContractID: c.ContractID, ClientEmail: c.ClientEmail})
			if err != nil {
				return clst.ScanRenewalsOutput{}, err
			}
			manifests[i] = manifest
		}
		if err := gc.LaunchBatch(ctx, manifests); err != nil {
			return clst.ScanRenewalsOutput{}, err
		}
		slog.Info("contract renewal scan complete", "found", len(due))
		return clst.ScanRenewalsOutput{Found: len(due)}, nil
	}
}
