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
	"github.com/terapps/gonveyor-examples/cmd/worker/handler"
	"github.com/terapps/gonveyor-examples/contracts"
	clst "github.com/terapps/gonveyor-examples/contracts/stations"
	"github.com/terapps/gonveyor-examples/transcoding"
	tst "github.com/terapps/gonveyor-examples/transcoding/stations"
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
		simple.SendWelcome: gonveyor.Handle(simple.SendWelcome, handler.HandleWelcome),
	})

	// transcoding
	reg.RegisterBlueprint(transcoding.Transcoding, gonveyor.Handlers{
		tst.Download:     gonveyor.Handle(tst.Download, handler.HandleDownload),
		tst.Transcode:    gonveyor.Handle(tst.Transcode, handler.HandleTranscode),
		tst.Thumbnail:    gonveyor.Handle(tst.Thumbnail, handler.HandleThumbnail),
		tst.ExtractAudio: gonveyor.Handle(tst.ExtractAudio, handler.HandleExtractAudio),
		tst.Package:      gonveyor.Handle(tst.Package, handler.HandlePackage),
	})

	// contract lifecycle — shared handlers
	docHandler := gonveyor.HandleFunc(handler.HandleDocument)
	emailHandler := gonveyor.HandleFunc(handler.HandleEmail)
	crmHandler := gonveyor.HandleFunc(handler.HandleCrm)

	reg.RegisterBlueprint(contracts.QuoteLifecycle, gonveyor.Handlers{
		clst.GenerateQuoteDoc:    docHandler,
		clst.SendQuoteEmail:      emailHandler,
		clst.SyncCrmQuote:        crmHandler,
		clst.GenerateContractDoc: docHandler,
		clst.SendContractEmail:   emailHandler,
		clst.SyncCrmContract:     crmHandler,
		clst.InitiateSignature:   gonveyor.Handle(clst.InitiateSignature, handler.HandleInitiateSignature),
		clst.InitiatePayment:     gonveyor.Handle(clst.InitiatePayment, handler.HandleInitiatePayment),
		clst.BundleContractDocs:  gonveyor.Handle(clst.BundleContractDocs, handler.HandleBundleContractDocs),
		clst.CreateContract:      gonveyor.Handle(clst.CreateContract, handler.HandleCreateContract),
	})

	reg.RegisterBlueprint(contracts.ContractRenewal, gonveyor.Handlers{
		clst.GenerateContractDoc:  docHandler,
		clst.SendContractEmail:    emailHandler,
		clst.SyncCrmContract:      crmHandler,
		clst.CheckContractRenewal: gonveyor.Handle(clst.CheckContractRenewal, handler.HandleCheckContractRenewal),
	})

	reg.RegisterBlueprint(contracts.ContractRenewalScan, gonveyor.Handlers{
		clst.ScanContractRenewals: gonveyor.Handle(clst.ScanContractRenewals, handler.NewScanContractRenewals(gc)),
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
