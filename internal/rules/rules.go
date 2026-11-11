// Package rules is the deterministic claim rule engine. No model calls here.
// All rule data is a SIMPLIFIED SUBSET for demo; see ruledata/*.json.
package rules

import (
	"embed"
	"encoding/json"
	"strings"
	"time"

	"github.com/daffaimron21/claim-autopsy/internal/claim"
)

type Rule interface {
	ID() string
	Check(c claim.Case) []claim.Finding
}

func All() []Rule {
	return []Rule{datesRule{}, losRule{}, primaryDxRule{}, procSupportRule{}, codeFormatRule{}}
}

func Run(c claim.Case) []claim.Finding {
	var out []claim.Finding
	for _, r := range All() {
		out = append(out, r.Check(c)...)
	}
	return out
}

type group struct {
	Label      string   `json:"label"`
	DxPrefixes []string `json:"dx_prefixes"`
	LosMin     int      `json:"los_min"`
	LosMax     int      `json:"los_max"`
	Source     string   `json:"source"`
	Simplified bool     `json:"simplified"`
}

type procEntry struct {
	Label      string   `json:"label"`
	DxPrefixes []string `json:"dx_prefixes"`
	Source     string   `json:"source"`
	Simplified bool     `json:"simplified"`
}

var data struct {
	CBG struct {
		Groups map[string]group `json:"groups"`
	}
	Proc struct {
		Procedures map[string]procEntry `json:"procedures"`
	}
}

//go:embed ruledata/*.json
var ruleFS embed.FS

func init() {
	for name, dst := range map[string]any{"ruledata/cbg_map.json": &data.CBG, "ruledata/proc_support.json": &data.Proc} {
		b, err := ruleFS.ReadFile(name)
		if err == nil {
			err = json.Unmarshal(b, dst)
		}
		if err != nil {
			panic("rules: bad embedded data " + name + ": " + err.Error())
		}
	}
}

const dateLayout = "2006-01-02"

// losDays returns the length of stay in days; a same-day stay counts as 1.
// ok is false when either date is malformed or discharge is before admit.
func losDays(admit, discharge string) (days int, ok bool) {
	a, err1 := time.Parse(dateLayout, admit)
	d, err2 := time.Parse(dateLayout, discharge)
	if err1 != nil || err2 != nil || d.Before(a) {
		return 0, false
	}
	days = int(d.Sub(a).Hours() / 24)
	if days == 0 {
		days = 1
	}
	return days, true
}

func hasPrefix(code string, prefixes []string) bool {
	code = strings.ToUpper(code)
	for _, p := range prefixes {
		if strings.HasPrefix(code, p) {
			return true
		}
	}
	return false
}

func primaryCode(c claim.Case) string {
	for _, d := range c.Diagnoses {
		if d.Primary {
			return d.Code
		}
	}
	return ""
}
