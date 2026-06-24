package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/terapps/gonveyor"
	clbp "github.com/terapps/gonveyor-examples/contract-lifecycle/blueprint"
	clh "github.com/terapps/gonveyor-examples/contract-lifecycle/handler"
	clst "github.com/terapps/gonveyor-examples/contract-lifecycle/stations"
	"github.com/terapps/gonveyor-examples/internal/infra"
	sh "github.com/terapps/gonveyor-examples/simple/handler"
	sst "github.com/terapps/gonveyor-examples/simple/stations"
	tbp "github.com/terapps/gonveyor-examples/transcoding/blueprint"
	th "github.com/terapps/gonveyor-examples/transcoding/handler"
	tst "github.com/terapps/gonveyor-examples/transcoding/stations"

	sbp "github.com/terapps/gonveyor-examples/simple/blueprint"
)

func main() {
	gonveyor.SetLogger(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))

	g, cleanup, err := infra.BuildWorker()
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

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

	// contract lifecycle — shared handlers registered on both phase-1 and phase-2 stations
	g.RegisterBlueprint(clbp.QuoteLifecycle)
	g.RegisterHandler(clst.GenerateQuoteDoc, gonveyor.Handle(clst.GenerateQuoteDoc, clh.GenerateDocument))
	g.RegisterHandler(clst.GenerateContractDoc, gonveyor.Handle(clst.GenerateContractDoc, clh.GenerateDocument))
	g.RegisterHandler(clst.InitiateSignature, gonveyor.Handle(clst.InitiateSignature, clh.InitiateSignature))
	g.RegisterHandler(clst.InitiatePayment, gonveyor.Handle(clst.InitiatePayment, clh.InitiatePayment))
	g.RegisterHandler(clst.SendQuoteEmail, gonveyor.Handle(clst.SendQuoteEmail, clh.SendEmail))
	g.RegisterHandler(clst.BundleContractDocs, gonveyor.Handle(clst.BundleContractDocs, clh.BundleContractDocs))
	g.RegisterHandler(clst.SendContractEmail, gonveyor.Handle(clst.SendContractEmail, clh.SendEmail))
	g.RegisterHandler(clst.SyncCrmQuote, gonveyor.Handle(clst.SyncCrmQuote, clh.SyncCrm))
	g.RegisterHandler(clst.SyncCrmContract, gonveyor.Handle(clst.SyncCrmContract, clh.SyncCrm))
	g.RegisterHandler(clst.CreateContract, gonveyor.Handle(clst.CreateContract, clh.CreateContract))

	log.Println("worker listening...")
	if err := g.Listen(context.Background()); err != nil {
		log.Fatal(err)
	}
}
