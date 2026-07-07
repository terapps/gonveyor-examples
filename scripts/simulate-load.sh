#!/usr/bin/env bash
# Continuously launches blueprints in rotation to simulate steady load.
# SLEEP=seconds between launches (default 1).
set -euo pipefail
cd "$(dirname "$0")/.."

go build -o ./tmp/publisher ./cmd/publisher

commands=(
	simple
	transcoding
	quote-lifecycle
	contract-renewal
)

i=0
while true; do
	cmd="${commands[$((i % ${#commands[@]}))]}"
	./tmp/publisher "$cmd" || true
	i=$((i + 1))
	sleep "${SLEEP:-1}"
done
