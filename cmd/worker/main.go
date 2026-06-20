package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/terapps/gonveyor"
	bp "github.com/terapps/gonveyor-examples/blueprint"
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

	// blueprints
	g.RegisterBlueprint(bp.SimpleDispatch)
	g.RegisterBlueprint(bp.Transcoding)
	g.RegisterBlueprint(bp.ContractFlow)

	// simple
	g.RegisterHandler(bp.SendWelcome, gonveyor.Handle(bp.SendWelcome, handler.SendWelcome))

	// transcoding
	g.RegisterHandler(bp.Download, gonveyor.Handle(bp.Download, handler.Download))
	g.RegisterHandler(bp.Transcode, gonveyor.Handle(bp.Transcode, handler.Transcode))
	g.RegisterHandler(bp.Thumbnail, gonveyor.Handle(bp.Thumbnail, handler.Thumbnail))
	g.RegisterHandler(bp.ExtractAudio, gonveyor.Handle(bp.ExtractAudio, handler.ExtractAudio))
	g.RegisterHandler(bp.Package, gonveyor.Handle(bp.Package, handler.Package))

	// contract
	g.RegisterHandler(bp.PrepareContract, gonveyor.Handle(bp.PrepareContract, handler.PrepareContract))
	g.RegisterHandler(bp.FinalizeContract, gonveyor.Handle(bp.FinalizeContract, handler.FinalizeContract))

	log.Println("worker listening...")
	if err := g.Listen(context.Background()); err != nil {
		log.Fatal(err)
	}
}
