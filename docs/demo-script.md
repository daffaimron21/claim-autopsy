# Demo script (target 2:50, English or English subtitles, no copyrighted music)

Open the live app on `case-06` before recording. Run once beforehand so the page is warm.

| Time | Show | Say |
|---|---|---|
| 0:00 | Title slide or app header | "When BPJS returns a hospital claim, you get one line. 'Pending: length of stay.' A coder then guesses what is really wrong. This is Claim Autopsy." |
| 0:20 | `case-06`: case file, return reason, then **Run autopsy** | "BPJS says the stay is too long. Watch the report." |
| 0:40 | Report: primary cause (model) and other cause (rule) | "The rule engine confirms nine days against a one-to-five-day group. But the primary cause is different: the discharge summary says the patient went home on July 4th. The real problem is a wrong discharge date." |
| 1:10 | Click the quote, then **Why this finding?** on both | "Every finding shows its evidence and the rule behind it. The model's quote is checked against the summary word for word. If it is not in the text, we drop it." |
| 1:40 | `case-05`: run autopsy | "Here the coded data is clean. The only cause is in free text: kidney injury and dialysis that were never coded. A rule engine alone cannot see this." |
| 2:05 | `case-04`: run, then toggle EN/ID | "Reversed dates and a malformed code. Deterministic rules, same input same result. The interface is bilingual." |
| 2:25 | README or `.kiro/` folder | TODO (owner): say only what you really did with Kiro. |
| 2:45 | Header with the synthetic-data line | "All data is synthetic and the rules are a simplified subset for demo. Claim Autopsy." |
