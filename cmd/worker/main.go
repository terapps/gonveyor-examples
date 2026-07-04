package main

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"os"

	"github.com/terapps/gonveyor"
	clbp "github.com/terapps/gonveyor-examples/contract-lifecycle/blueprint"
	clh "github.com/terapps/gonveyor-examples/contract-lifecycle/handler"
	clst "github.com/terapps/gonveyor-examples/contract-lifecycle/stations"
	sh "github.com/terapps/gonveyor-examples/simple/handler"
	sst "github.com/terapps/gonveyor-examples/simple/stations"
	tbp "github.com/terapps/gonveyor-examples/transcoding/blueprint"
	th "github.com/terapps/gonveyor-examples/transcoding/handler"
	tst "github.com/terapps/gonveyor-examples/transcoding/stations"
	bunledger "github.com/terapps/gonveyor/ledger/bun"
	"github.com/terapps/gonveyor/transport/pg"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	sbp "github.com/terapps/gonveyor-examples/simple/blueprint"
)

const defaultPostgresDSN = "postgres://gonveyor:gonveyor@localhost:5432/gonveyor?sslmode=disable"

func main() {
	gonveyor.SetLogger(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))

	db := openDB()
	defer func() { _ = db.Close() }()

	worker := pg.NewWorker(db)
	g := gonveyor.NewGonveyor(bunledger.New(db), worker)

	// simple
	g.RegisterBlueprint(sbp.SimpleDispatch)
	g.RegisterHandler(sst.SendWelcome, gonveyor.Handle(sst.SendWelcome, sh.SendWelcome))

	// transcoding
	g.RegisterBlueprint(tbp.Transcoding)
	g.RegisterHandler(tst.Download, gonveyor.Handle(tst.Download, th.Download))
	g.RegisterHandler(tst.Transcode, gonveyor.Handle(tst.Transcode, th.Transcode))
	g.RegisterHandler(tst.Thumbnail, gonveyor.Handle(tst.Thumbnail, th.Thumbnail))
	g.RegisterHandler(tst.ExtractAudio, gonveyor.Handle(tst.ExtractAudio, th.ExtractAudio))
	g.RegisterHandler(tst.Package, gonveyor.Handle(tst.Package, th.Package))

	// contract lifecycle — shared handlers registered once for every station that reuses them,
	// across both phase-1/phase-2 of quote_lifecycle and the independent contract_renewal blueprint
	g.RegisterBlueprint(clbp.QuoteLifecycle)
	g.RegisterBlueprint(clbp.ContractRenewal)
	g.RegisterHandlers(gonveyor.HandleFunc(clh.GenerateDocument), clst.GenerateQuoteDoc, clst.GenerateContractDoc)
	g.RegisterHandlers(gonveyor.HandleFunc(clh.SendEmail), clst.SendQuoteEmail, clst.SendContractEmail)
	g.RegisterHandlers(gonveyor.HandleFunc(clh.SyncCrm), clst.SyncCrmQuote, clst.SyncCrmContract)
	g.RegisterHandler(clst.InitiateSignature, gonveyor.Handle(clst.InitiateSignature, clh.InitiateSignature))
	g.RegisterHandler(clst.InitiatePayment, gonveyor.Handle(clst.InitiatePayment, clh.InitiatePayment))
	g.RegisterHandler(clst.BundleContractDocs, gonveyor.Handle(clst.BundleContractDocs, clh.BundleContractDocs))
	g.RegisterHandler(clst.CreateContract, gonveyor.Handle(clst.CreateContract, clh.CreateContract))
	g.RegisterHandler(clst.CheckContractRenewal, gonveyor.Handle(clst.CheckContractRenewal, clh.CheckContractRenewal))

	log.Println("worker listening...")
	if err := g.Listen(context.Background()); err != nil {
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
