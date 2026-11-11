# Claim Autopsy

Find out why a BPJS claim came back. Kiroween 2026 entry, category **Frankenstein**.

A returned claim usually comes with a one-line reason. Claim Autopsy dissects it and reports the primary cause,
other causes, the evidence, a suggested fix, and whether the claim can be resubmitted.

Three parts stitched together:

1. **Rule engine** (`internal/rules`): deterministic Go code over the coded data. Every finding names its rule.
2. **Resume reader** (`internal/resume`): a language model reads the free-text discharge summary and may only report
   mismatches from a closed list. Each finding must quote the summary verbatim. A quote that is not in the text is dropped.
3. **Report assembler** (`internal/report`): ranks findings and decides the verdict. The model never ranks.

## Important notes

- **All data is synthetic.** No real patient data, no BPJS or SATUSEHAT credentials.
- **Rules are a simplified subset for demo.** `internal/rules/ruledata/*.json` is designed for this demo and is not an official tariff table or BPJS rule set.
- This is a demo, not a claims system. It does not submit anything anywhere.
- The saved results in `testdata/resume_fixtures/` were **recorded on 2026-11-11 from `claude-haiku-5-5`** with `go run ./cmd/record-fixtures`. The app still labels them "saved results, not a live model call" unless a live call is made.

- The eval test (`internal/report/eval_test.go`) compares the recorded model output and the rules with hand-written expectations for 6 synthetic cases. Six cases written by the same people who wrote the rules is a regression test, **not a measurement of model accuracy**.
- Only two of the three verdicts can currently occur, because every finding is fixable.

## Run it

```bash
make test                    # Go tests and web tests
make dev                     # API on :8080, web on :5173
```

Production-like (single container):

```bash
docker build -t claim-autopsy .
docker run -p 8080:8080 claim-autopsy
```

| Env var | Meaning |
|---|---|
| `PORT` | Listen port, default `8080` |
| `LIVE_MODEL` | `true` turns on live model calls. Default off: saved results only |
| `ANTHROPIC_API_KEY` | Required when `LIVE_MODEL=true`. Never commit it |
| `LIVE_MODEL_NAME` | Default `claude-haiku-5-5` |
| `WEB_DIR` | Built web app, default `web/dist` |

If the live model fails, times out, has no key, or the rate limit is hit, the report is still returned from saved results and marked as degraded. A successful live result is cached per case, so the paid model is called at most once per case per process.

## License

MIT
