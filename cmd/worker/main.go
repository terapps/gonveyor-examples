package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"log/slog"
	"os"
	ossignal "os/signal"
	"strings"
	"syscall"

	"github.com/terapps/gonveyor"
	"github.com/terapps/gonveyor-examples/contracts"
	"github.com/terapps/gonveyor-examples/transcoding"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	"github.com/terapps/gonveyor-examples/simple"
)

const defaultPostgresDSN = "postgres://gonveyor:gonveyor@localhost:5432/gonveyor?sslmode=disable"

// routingKeysFlag accumulates across repeated occurrences (-k a -k b) and also splits
// each occurrence on commas (-k a,b), so both styles work and can be mixed.
type routingKeysFlag []string

func (f *routingKeysFlag) String() string { return strings.Join(*f, ",") }
func (f *routingKeysFlag) Set(v string) error {
	*f = append(*f, strings.Split(v, ",")...)
	return nil
}

func main() {
	var routingKeys routingKeysFlag
	flag.Var(&routingKeys, "routing-keys", "routing keys to poll, repeatable and/or comma-separated (default: gonveyor.default only)")
	flag.Var(&routingKeys, "k", "shorthand for -routing-keys")
	name := flag.String("name", "", "worker name recorded in worker_instances")
	flag.Parse()

	db := openDB()
	defer func() { _ = db.Close() }()

	gc := gonveyor.NewGonductor(db)
	reg := gonveyor.NewStationRegistry()

	// simple
	reg.RegisterBlueprint(simple.SimpleDispatch, gonveyor.Handlers{
		simple.SendWelcome: gonveyor.Handle(simple.SendWelcome, simple.HandleWelcome),
	})

	// transcoding
	reg.RegisterBlueprint(transcoding.Transcoding, gonveyor.Handlers{
		transcoding.Download:     gonveyor.Handle(transcoding.Download, transcoding.HandleDownload),
		transcoding.Transcode:    gonveyor.Handle(transcoding.Transcode, transcoding.HandleTranscode),
		transcoding.Thumbnail:    gonveyor.Handle(transcoding.Thumbnail, transcoding.HandleThumbnail),
		transcoding.ExtractAudio: gonveyor.Handle(transcoding.ExtractAudio, transcoding.HandleExtractAudio),
		transcoding.Package:      gonveyor.Handle(transcoding.Package, transcoding.HandlePackage),
	})

	// contract lifecycle — shared handlers
	docHandler := gonveyor.HandleFunc(contracts.HandleDocument)
	emailHandler := gonveyor.HandleFunc(contracts.HandleEmail)
	crmHandler := gonveyor.HandleFunc(contracts.HandleCrm)

	reg.RegisterBlueprint(contracts.QuoteLifecycle, gonveyor.Handlers{
		contracts.GenerateQuoteDoc:    docHandler,
		contracts.SendQuoteEmail:      emailHandler,
		contracts.SyncCrmQuote:        crmHandler,
		contracts.GenerateContractDoc: docHandler,
		contracts.SendContractEmail:   emailHandler,
		contracts.SyncCrmContract:     crmHandler,
		contracts.InitiateSignature:   gonveyor.Handle(contracts.InitiateSignature, contracts.HandleInitiateSignature),
		contracts.InitiatePayment:     gonveyor.Handle(contracts.InitiatePayment, contracts.HandleInitiatePayment),
		contracts.BundleContractDocs:  gonveyor.Handle(contracts.BundleContractDocs, contracts.HandleBundleContractDocs),
		contracts.CreateContract:      gonveyor.Handle(contracts.CreateContract, contracts.HandleCreateContract),
	})

	reg.RegisterBlueprint(contracts.ContractRenewal, gonveyor.Handlers{
		contracts.GenerateContractDoc:  docHandler,
		contracts.SendContractEmail:    emailHandler,
		contracts.SyncCrmContract:      crmHandler,
		contracts.CheckContractRenewal: gonveyor.Handle(contracts.CheckContractRenewal, contracts.HandleCheckContractRenewal),
	})

	reg.RegisterBlueprint(contracts.ContractRenewalScan, gonveyor.Handlers{
		contracts.ScanContractRenewals: gonveyor.Handle(contracts.ScanContractRenewals, contracts.NewScanContractRenewals(gc)),
	})

	templates := []gonveyor.AnyLaunchTemplate{
		simple.Template,
		transcoding.Template,
		contracts.QuoteLifecycleTemplate,
		contracts.RenewalTemplate,
		contracts.ScanTemplate,
	}

	opts := []gonveyor.Option{
		gonveyor.WithRegistry(reg),
		gonveyor.WithBlueprintProducer(templates),
		gonveyor.WithScheduler(),
		gonveyor.WithDiscovery(),
		gonveyor.WithLogger(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))),
	}
	if len(routingKeys) > 0 {
		opts = append(opts, gonveyor.WithRoutingKeys(routingKeys...))
	}
	if *name != "" {
		opts = append(opts, gonveyor.WithName(*name))
	}
	g := gonveyor.NewGonveyor(db, opts...)

	ctx, stop := ossignal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	log.Println("worker listening...")
	if err := g.Listen(ctx); err != nil {
		log.Fatal(err)
	}
}

func openDB() *bun.DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(envOr("POSTGRES_DSN", defaultPostgresDSN))))
	return bun.NewDB(sqldb, pgdialect.New())
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
