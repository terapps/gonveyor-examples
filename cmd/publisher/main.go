package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"

	bp "github.com/terapps/gonveyor-examples/blueprint"
	clbp "github.com/terapps/gonveyor-examples/contract-lifecycle/blueprint"
	"github.com/terapps/gonveyor-examples/internal/infra"
	"github.com/terapps/gonveyor/ledger"
)

const usage = `usage: publisher <command> [flags]

commands:
  simple          submit a simple welcome dispatch
  transcoding     submit a video transcoding workflow
  quote-lifecycle submit a full quote → contract lifecycle workflow
  signal          send a signal to an existing blueprint

flags:
  simple:
    -user-id   string  user ID (default: user-1)
    -email     string  email address (default: user@example.com)

  transcoding:
    -asset-id   string  asset ID (default: asset-1)
    -source-url string  source URL (default: s3://bucket/video.mp4)

  quote-lifecycle:
    -quote-id    string  quote ID (default: quote-1)
    -email       string  client email (default: client@example.com)

  signal:
    -blueprint-id  string  blueprint instance ID (required)
    -key           string  signal key, e.g. await_signature (required)
    -payload       string  JSON payload (default: {})
`

func main() {
	if len(os.Args) < 2 {
		log.Fatal(usage)
	}

	ctx := context.Background()
	cmd := os.Args[1]
	args := os.Args[2:]

	if cmd == "signal" {
		runSignal(ctx, args)
		return
	}

	gc, cleanup, err := infra.BuildGonductor()
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	var manifest ledger.BlueprintManifest

	switch cmd {
	case "simple":
		fs := flag.NewFlagSet("simple", flag.ExitOnError)
		userID := fs.String("user-id", "user-1", "user ID")
		email := fs.String("email", "user@example.com", "email address")
		_ = fs.Parse(args)
		manifest, err = bp.SimpleManifest(*userID, *email)

	case "transcoding":
		fs := flag.NewFlagSet("transcoding", flag.ExitOnError)
		assetID := fs.String("asset-id", "asset-1", "asset ID")
		sourceURL := fs.String("source-url", "s3://bucket/video.mp4", "source URL")
		_ = fs.Parse(args)
		manifest, err = bp.TranscodingManifest(*assetID, *sourceURL)

	case "quote-lifecycle":
		fs := flag.NewFlagSet("quote-lifecycle", flag.ExitOnError)
		quoteID := fs.String("quote-id", "quote-1", "quote ID")
		email := fs.String("email", "client@example.com", "client email")
		_ = fs.Parse(args)
		manifest, err = clbp.Manifest(clbp.Params{
			QuoteID:          *quoteID,
			ClientEmail:      *email,
			QuoteDocTypes:    []string{"proposal", "pricing", "terms"},
			ContractDocTypes: []string{"contract", "annex_a"},
		})

	default:
		log.Fatalf("unknown command %q\n\n%s", cmd, usage)
	}

	if err != nil {
		log.Fatal(err)
	}

	if err := gc.Launch(ctx, manifest); err != nil {
		log.Fatal(err)
	}
	log.Printf("blueprint %s (%s) launched", manifest.Blueprint.ID, manifest.Blueprint.Name)
}

func runSignal(ctx context.Context, args []string) {
	fs := flag.NewFlagSet("signal", flag.ExitOnError)
	blueprintID := fs.String("blueprint-id", "", "blueprint instance ID (required)")
	key := fs.String("key", "", "signal key (required)")
	payloadStr := fs.String("payload", "{}", "JSON payload")
	_ = fs.Parse(args)

	if *blueprintID == "" || *key == "" {
		log.Fatal("signal requires -blueprint-id and -key")
	}

	var payload any
	if err := json.Unmarshal([]byte(*payloadStr), &payload); err != nil {
		log.Fatalf("invalid JSON payload: %v", err)
	}

	gc, cleanup, err := infra.BuildGonductor()
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	if err := gc.SendSignal(ctx, *blueprintID, *key, payload); err != nil {
		log.Fatal(err)
	}
	log.Printf("signal %q sent to blueprint %s", *key, *blueprintID)
}
