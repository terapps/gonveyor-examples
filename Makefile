.PHONY: build worker publisher dev dev-documents load fmt clean

build:
	go build ./...

worker:
	go run ./cmd/worker

publisher:
	go run ./cmd/publisher $(ARGS)

dev:
	air -build.cmd "go build -o ./tmp/worker-default ./cmd/worker" -build.entrypoint "./tmp/worker-default" -build.args_bin "$(ARGS)"

dev-documents:
	air -build.cmd "go build -o ./tmp/worker-documents ./cmd/worker" -build.entrypoint "./tmp/worker-documents" -build.args_bin "-k tasks.document -name documents"

load:
	./scripts/simulate-load.sh

fmt:
	gofmt -w .

clean:
	rm -rf tmp
