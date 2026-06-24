package infra

import (
	"database/sql"
	"os"

	"github.com/rabbitmq/amqp091-go"
	"github.com/terapps/gonveyor"
	evamqp "github.com/terapps/gonveyor/events/amqp"
	bunledger "github.com/terapps/gonveyor/ledger/bun"
	"github.com/terapps/gonveyor/transport/amqp"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

const (
	defaultAMQPURL     = "amqp://gonveyor:gonveyor@localhost:5672/"
	defaultPostgresDSN = "postgres://gonveyor:gonveyor@localhost:5432/gonveyor?sslmode=disable"
	QueueName          = "gonveyor"
	EventsExchange     = "gonveyor.events"
)

func defaultQueue() (*amqp.Queue, error) {
	return amqp.NewQueue(QueueName, amqp.WithDeadLetter("gonveyor.dlx"))
}

func BuildWorker() (*gonveyor.Gonveyor, func(), error) {
	db := openDB()
	queue, err := defaultQueue()
	if err != nil {
		_ = db.Close()
		return nil, nil, err
	}

	conn, err := amqp.Dial(envOr("AMQP_URL", defaultAMQPURL))
	if err != nil {
		_ = db.Close()
		return nil, nil, err
	}

	worker, err := conn.NewWorker(queue)
	if err != nil {
		_ = conn.Close()
		_ = db.Close()
		return nil, nil, err
	}

	dispatcher, err := conn.NewDispatcher(queue)
	if err != nil {
		_ = conn.Close()
		_ = db.Close()
		return nil, nil, err
	}

	pub, rawConn, err := buildEventPublisher()
	if err != nil {
		_ = worker.Close()
		_ = dispatcher.Close()
		_ = conn.Close()
		_ = db.Close()
		return nil, nil, err
	}

	cleanup := func() {
		_ = pub.Close()
		_ = rawConn.Close()
		_ = worker.Close()
		_ = dispatcher.Close()
		_ = conn.Close()
		_ = db.Close()
	}

	g := gonveyor.NewGonveyor(bunledger.New(db), dispatcher, worker, gonveyor.WithEventPublisher(pub))
	return g, cleanup, nil
}

func BuildGonductor() (*gonveyor.Gonductor, func(), error) {
	db := openDB()
	queue, err := defaultQueue()
	if err != nil {
		_ = db.Close()
		return nil, nil, err
	}

	conn, err := amqp.Dial(envOr("AMQP_URL", defaultAMQPURL))
	if err != nil {
		_ = db.Close()
		return nil, nil, err
	}

	dispatcher, err := conn.NewDispatcher(queue)
	if err != nil {
		_ = conn.Close()
		_ = db.Close()
		return nil, nil, err
	}

	pub, rawConn, err := buildEventPublisher()
	if err != nil {
		_ = dispatcher.Close()
		_ = conn.Close()
		_ = db.Close()
		return nil, nil, err
	}

	cleanup := func() {
		_ = pub.Close()
		_ = rawConn.Close()
		_ = dispatcher.Close()
		_ = conn.Close()
		_ = db.Close()
	}

	gc := gonveyor.NewGonductor(bunledger.New(db), dispatcher, gonveyor.WithEventPublisher(pub))
	return gc, cleanup, nil
}

func buildEventPublisher() (*evamqp.Publisher, *amqp091.Connection, error) {
	rawConn, err := amqp091.Dial(envOr("AMQP_URL", defaultAMQPURL))
	if err != nil {
		return nil, nil, err
	}
	pub, err := evamqp.New(rawConn, EventsExchange)
	if err != nil {
		_ = rawConn.Close()
		return nil, nil, err
	}
	return pub, rawConn, nil
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
