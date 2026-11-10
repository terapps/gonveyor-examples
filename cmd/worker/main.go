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

	_ "github.com/lib/pq"
	"github.com/terapps/gonveyor"
	"github.com/terapps/gonveyor-examples/cmd/worker/handler"
	"github.com/terapps/gonveyor-examples/contracts"
	clst "github.com/terapps/gonveyor-examples/contracts/stations"
	"github.com/terapps/gonveyor-examples/transcoding"
	tst "github.com/terapps/gonveyor-examples/transcoding/stations"

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

	registerStations := func(reg gonveyor.StationRegistry) {
		reg.RegisterBlueprint(simple.SimpleDispatch,
			simple.SendWelcome.Handle(handler.HandleWelcome),
		)

		reg.RegisterBlueprint(transcoding.Transcoding,
			tst.Download.Handle(handler.HandleDownload),
			tst.Transcode.Handle(handler.HandleTranscode),
			tst.Thumbnail.Handle(handler.HandleThumbnail),
			tst.ExtractAudio.Handle(handler.HandleExtractAudio),
			tst.Package.Handle(handler.HandlePackage),
		)

		reg.RegisterBlueprint(contracts.QuoteLifecycle,
			clst.GenerateQuoteDoc.Handle(handler.HandleDocument),
			clst.SendQuoteEmail.Handle(handler.HandleEmail),
			clst.SyncCrmQuote.Handle(handler.HandleCrm),
			clst.GenerateContractDoc.Handle(handler.HandleDocument),
			clst.SendContractEmail.Handle(handler.HandleEmail),
			clst.SyncCrmContract.Handle(handler.HandleCrm),
			clst.InitiateSignature.Handle(handler.HandleInitiateSignature),
			clst.InitiatePayment.Handle(handler.HandleInitiatePayment),
			clst.BundleContractDocs.Handle(handler.HandleBundleContractDocs),
			clst.CreateContract.Handle(handler.HandleCreateContract),
		)

		reg.RegisterBlueprint(contracts.ContractRenewal,
			clst.GenerateContractDoc.Handle(handler.HandleDocument),
			clst.SendContractEmail.Handle(handler.HandleEmail),
			clst.SyncCrmContract.Handle(handler.HandleCrm),
			clst.CheckContractRenewal.Handle(handler.HandleCheckContractRenewal),
			clst.RenewContract.Handle(handler.HandleRenewContract),
			clst.CloseContract.Handle(handler.HandleCloseContract),
		)

		reg.RegisterBlueprint(contracts.ContractRenewalScan,
			clst.ScanContractRenewals.Handle(handler.NewScanContractRenewals(gc)),
		)
	}

	registerTemplates := func(lr gonveyor.LaunchRegistry) {
		lr.RegisterTemplate(
			simple.Template,
			transcoding.Template,
			contracts.QuoteLifecycleTemplate,
			contracts.RenewalTemplate,
			contracts.ScanTemplate,
		)
	}

	opts := []gonveyor.Option{
		gonveyor.WithStationMaster(registerStations),
		gonveyor.WithBlueprintProducer(registerTemplates),
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

func openDB() *sql.DB {
	db, err := sql.Open("postgres", envOr("POSTGRES_DSN", defaultPostgresDSN))
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	return db
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
