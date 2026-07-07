.PHONY: build worker publisher dev dev-documents load fmt clean

build:
	go build ./...

worker:
	go run ./cmd/worker

publisher:
	go run ./cmd/publisher $(ARGS)

dev:
	air -build.args_bin "$(ARGS)"

dev-documents:
	air -build.args_bin "-k tasks.document -name documents"

load:
	./scripts/simulate-load.sh

fmt:
	gofmt -w .

clean:
	rm -rf tmp
