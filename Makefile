.PHONY: test check
test:
	go test ./...
	@if [ -d web/node_modules ]; then cd web && bun run test; fi
check:
	bash scripts/check-kiro-tracked.sh
