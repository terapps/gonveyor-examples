package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/terapps/gonveyor"
	bp "github.com/terapps/gonveyor-examples/blueprint"
	clbp "github.com/terapps/gonveyor-examples/contract-lifecycle/blueprint"
	clh "github.com/terapps/gonveyor-examples/contract-lifecycle/handler"
	clst "github.com/terapps/gonveyor-examples/contract-lifecycle/stations"
	"github.com/terapps/gonveyor-examples/handler"
	"github.com/terapps/gonveyor-examples/internal/infra"
)

func main() {
	gonveyor.SetLogger(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))

	g, cleanup, err := infra.BuildWorker()
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	// simple
	g.RegisterBlueprint(bp.SimpleDispatch)
	g.RegisterHandler(bp.SendWelcome, gonveyor.Handle(bp.SendWelcome, handler.SendWelcome))

	// transcoding
	g.RegisterBlueprint(bp.Transcoding)
	g.RegisterHandler(bp.Download, gonveyor.Handle(bp.Download, handler.Download))
	g.RegisterHandler(bp.Transcode, gonveyor.Handle(bp.Transcode, handler.Transcode))
	g.RegisterHandler(bp.Thumbnail, gonveyor.Handle(bp.Thumbnail, handler.Thumbnail))
	g.RegisterHandler(bp.ExtractAudio, gonveyor.Handle(bp.ExtractAudio, handler.ExtractAudio))
	g.RegisterHandler(bp.Package, gonveyor.Handle(bp.Package, handler.Package))

	// contract lifecycle — shared handlers registered on both phase-1 and phase-2 stations
	g.RegisterBlueprint(clbp.QuoteLifecycle)
	g.RegisterHandler(clst.GenerateQuoteDoc, gonveyor.Handle(clst.GenerateQuoteDoc, clh.GenerateDocument))
	g.RegisterHandler(clst.GenerateContractDoc, gonveyor.Handle(clst.GenerateContractDoc, clh.GenerateDocument))
	g.RegisterHandler(clst.InitiateSignature, gonveyor.Handle(clst.InitiateSignature, clh.InitiateSignature))
	g.RegisterHandler(clst.SendQuoteEmail, gonveyor.Handle(clst.SendQuoteEmail, clh.SendEmail))
	g.RegisterHandler(clst.SendContractEmail, gonveyor.Handle(clst.SendContractEmail, clh.SendEmail))
	g.RegisterHandler(clst.SyncCrmQuote, gonveyor.Handle(clst.SyncCrmQuote, clh.SyncCrm))
	g.RegisterHandler(clst.SyncCrmContract, gonveyor.Handle(clst.SyncCrmContract, clh.SyncCrm))
	g.RegisterHandler(clst.CreateContract, gonveyor.Handle(clst.CreateContract, clh.CreateContract))

	log.Println("worker listening...")
	if err := g.Listen(context.Background()); err != nil {
		log.Fatal(err)
	}
}
