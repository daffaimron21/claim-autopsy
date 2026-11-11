.PHONY: test check dev
test:
	go test ./...
	@if [ -d web/node_modules ]; then cd web && bun run test; fi
check:
	bash scripts/check-kiro-tracked.sh
dev:
	@trap 'kill 0' EXIT; go run ./cmd/server & (cd web && bun run dev) & wait
