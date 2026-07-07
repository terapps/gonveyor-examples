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
	clbp "github.com/terapps/gonveyor-examples/contract-lifecycle/blueprint"
	clh "github.com/terapps/gonveyor-examples/contract-lifecycle/handler"
	clst "github.com/terapps/gonveyor-examples/contract-lifecycle/stations"
	sh "github.com/terapps/gonveyor-examples/simple/handler"
	sst "github.com/terapps/gonveyor-examples/simple/stations"
	tbp "github.com/terapps/gonveyor-examples/transcoding/blueprint"
	th "github.com/terapps/gonveyor-examples/transcoding/handler"
	tst "github.com/terapps/gonveyor-examples/transcoding/stations"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	sbp "github.com/terapps/gonveyor-examples/simple/blueprint"
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
	reg.RegisterLauncher(sbp.Launcher)
	reg.RegisterHandler(sst.SendWelcome, gonveyor.Handle(sst.SendWelcome, sh.SendWelcome))

	// transcoding
	reg.RegisterLauncher(tbp.Launcher)
	reg.RegisterHandler(tst.Download, gonveyor.Handle(tst.Download, th.Download))
	reg.RegisterHandler(tst.Transcode, gonveyor.Handle(tst.Transcode, th.Transcode))
	reg.RegisterHandler(tst.Thumbnail, gonveyor.Handle(tst.Thumbnail, th.Thumbnail))
	reg.RegisterHandler(tst.ExtractAudio, gonveyor.Handle(tst.ExtractAudio, th.ExtractAudio))
	reg.RegisterHandler(tst.Package, gonveyor.Handle(tst.Package, th.Package))

	// contract lifecycle — shared handlers registered once for every station that reuses them,
	// across both phase-1/phase-2 of quote_lifecycle and the independent contract_renewal blueprint
	reg.RegisterLauncher(clbp.QuoteLifecycleLauncher)
	reg.RegisterLauncher(clbp.RenewalLauncher)
	reg.RegisterLauncher(clbp.ScanLauncher)
	reg.RegisterHandler(clst.ScanContractRenewals, gonveyor.Handle(clst.ScanContractRenewals, clh.NewScanContractRenewals(gc)))
	docHandler := gonveyor.HandleFunc(clh.GenerateDocument)
	reg.RegisterHandler(clst.GenerateQuoteDoc, docHandler)
	reg.RegisterHandler(clst.GenerateContractDoc, docHandler)
	emailHandler := gonveyor.HandleFunc(clh.SendEmail)
	reg.RegisterHandler(clst.SendQuoteEmail, emailHandler)
	reg.RegisterHandler(clst.SendContractEmail, emailHandler)
	crmHandler := gonveyor.HandleFunc(clh.SyncCrm)
	reg.RegisterHandler(clst.SyncCrmQuote, crmHandler)
	reg.RegisterHandler(clst.SyncCrmContract, crmHandler)
	reg.RegisterHandler(clst.InitiateSignature, gonveyor.Handle(clst.InitiateSignature, clh.InitiateSignature))
	reg.RegisterHandler(clst.InitiatePayment, gonveyor.Handle(clst.InitiatePayment, clh.InitiatePayment))
	reg.RegisterHandler(clst.BundleContractDocs, gonveyor.Handle(clst.BundleContractDocs, clh.BundleContractDocs))
	reg.RegisterHandler(clst.CreateContract, gonveyor.Handle(clst.CreateContract, clh.CreateContract))
	reg.RegisterHandler(clst.CheckContractRenewal, gonveyor.Handle(clst.CheckContractRenewal, clh.CheckContractRenewal))

	opts := []gonveyor.Option{
		gonveyor.WithRegistry(reg),
		gonveyor.WithBlueprintProducer(),
		gonveyor.WithScheduler(),
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
