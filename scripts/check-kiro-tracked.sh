#!/usr/bin/env bash
# Kiroween rule: .kiro/ must be at repo root and not gitignored.
cd "$(git rev-parse --show-toplevel)" || exit 1
if [ -z "$(git ls-files .kiro)" ]; then echo "FAIL: .kiro has no tracked files"; exit 1; fi
if [ -n "$(git ls-files -ci --exclude-standard .kiro)" ] || git check-ignore -q .kiro/steering; then echo "FAIL: .kiro is gitignored"; exit 1; fi
echo "OK: .kiro tracked and not ignored"
