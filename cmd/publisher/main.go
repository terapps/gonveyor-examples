package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/terapps/gonveyor"
	clbp "github.com/terapps/gonveyor-examples/contract-lifecycle/blueprint"
	clst "github.com/terapps/gonveyor-examples/contract-lifecycle/stations"
	sbp "github.com/terapps/gonveyor-examples/simple/blueprint"
	tbp "github.com/terapps/gonveyor-examples/transcoding/blueprint"
	"github.com/terapps/gonveyor/ledger"
	bunledger "github.com/terapps/gonveyor/ledger/bun"
	"github.com/terapps/gonveyor/transport/pg"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

const defaultPostgresDSN = "postgres://gonveyor:gonveyor@localhost:5432/gonveyor?sslmode=disable"

func buildGonductor() (*gonveyor.Gonductor, func(), error) {
	db := openDB()
	return gonveyor.NewGonductor(bunledger.New(db)), func() { _ = db.Close() }, nil
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

const usage = `usage: publisher <command> [flags]

commands:
  simple                         submit a simple welcome dispatch
  transcoding                     submit a video transcoding workflow
  quote-lifecycle                  submit a full quote → contract lifecycle workflow
  contract-renewal                 submit a standalone contract renewal reminder
  schedule-contract-renewal-scan   register the recurring contract renewal scan
  parent-child-demo                spawn a contract_renewal as a child and wait on it
  signal                           send a signal to an existing blueprint

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
    -amount      float   quote amount (default: 100)

  contract-renewal:
    -contract-id string  contract ID (default: contract-1)
    -email       string  client email (default: client@example.com)

  schedule-contract-renewal-scan:
    -cron  string  cron expression, standard 5-field or "@every 1h30m" (default: "0 9 * * *")

  parent-child-demo:
    -contract-id string  contract ID for the child renewal (default: contract-1)
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
	if cmd == "schedule-contract-renewal-scan" {
		runScheduleContractRenewalScan(ctx, args)
		return
	}

	gc, cleanup, err := buildGonductor()
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
		manifest, err = sbp.Manifest(*userID, *email)

	case "transcoding":
		fs := flag.NewFlagSet("transcoding", flag.ExitOnError)
		assetID := fs.String("asset-id", "asset-1", "asset ID")
		sourceURL := fs.String("source-url", "s3://bucket/video.mp4", "source URL")
		_ = fs.Parse(args)
		manifest, err = tbp.Manifest(*assetID, *sourceURL)

	case "quote-lifecycle":
		fs := flag.NewFlagSet("quote-lifecycle", flag.ExitOnError)
		quoteID := fs.String("quote-id", "quote-1", "quote ID")
		email := fs.String("email", "client@example.com", "client email")
		amount := fs.Float64("amount", 100, "quote amount")
		_ = fs.Parse(args)
		manifest, err = clbp.Manifest(clbp.Params{
			QuoteID:          *quoteID,
			ClientEmail:      *email,
			Amount:           *amount,
			QuoteDocTypes:    []string{"proposal", "pricing", "terms"},
			ContractDocTypes: []string{"contract", "annex_a"},
		})

	case "contract-renewal":
		fs := flag.NewFlagSet("contract-renewal", flag.ExitOnError)
		contractID := fs.String("contract-id", "contract-1", "contract ID")
		email := fs.String("email", "client@example.com", "client email")
		_ = fs.Parse(args)
		manifest, err = clbp.RenewalManifest(clst.CheckContractRenewalInput{
			ContractID:  *contractID,
			ClientEmail: *email,
		})

	case "parent-child-demo":
		fs := flag.NewFlagSet("parent-child-demo", flag.ExitOnError)
		contractID := fs.String("contract-id", "contract-1", "contract ID")
		email := fs.String("email", "client@example.com", "client email")
		_ = fs.Parse(args)
		manifest, err = clbp.ParentChildDemoLauncher.Manifest(clst.SpawnChildRenewalInput{
			ContractID:  *contractID,
			ClientEmail: *email,
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

	gc, cleanup, err := buildGonductor()
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	if err := gc.SendSignal(ctx, *blueprintID, *key, payload); err != nil {
		log.Fatal(err)
	}
	log.Printf("signal %q sent to blueprint %s", *key, *blueprintID)
}

// runScheduleContractRenewalScan registers the recurring contract_renewal_scan launch —
// one schedule total, not one per contract: the scan reads live contract data every run
// and files a contract_renewal launch_request per contract found due.
func runScheduleContractRenewalScan(ctx context.Context, args []string) {
	fs := flag.NewFlagSet("schedule-contract-renewal-scan", flag.ExitOnError)
	cronExpr := fs.String("cron", "0 9 * * *", `cron expression, standard 5-field or "@every 1h30m"`)
	_ = fs.Parse(args)

	db := openDB()
	defer func() { _ = db.Close() }()

	id, err := pg.CreateSchedule(ctx, db, "contract_renewal_scan", *cronExpr, []byte(`{}`))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("schedule %s registered for contract_renewal_scan (%s)", id, *cronExpr)
}
