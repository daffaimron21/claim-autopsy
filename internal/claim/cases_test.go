package claim

import (
	"os"
	"testing"
)

func loadFixtures(t *testing.T) []Case {
	t.Helper()
	cases, err := LoadCases(os.DirFS("../../testdata/cases"))
	if err != nil {
		t.Fatalf("LoadCases: %v", err)
	}
	return cases
}

func TestLoadCases_ParsesAllFixtures(t *testing.T) {
	cases := loadFixtures(t)
	if len(cases) != 6 {
		t.Fatalf("want 6 cases, got %d", len(cases))
	}
	seen := map[string]bool{}
	for _, c := range cases {
		if c.ID == "" || seen[c.ID] {
			t.Errorf("case id empty or duplicated: %q", c.ID)
		}
		seen[c.ID] = true
		primaries := 0
		for _, d := range c.Diagnoses {
			if d.Primary {
				primaries++
			}
		}
		if primaries != 1 {
			t.Errorf("%s: want exactly 1 primary diagnosis, got %d", c.ID, primaries)
		}
		if len(c.Expected) == 0 {
			t.Errorf("%s: Expected must not be empty", c.ID)
		}
		if c.ResumeText == "" {
			t.Errorf("%s: ResumeText must not be empty", c.ID)
		}
	}
}

func TestFindCase(t *testing.T) {
	cases := loadFixtures(t)
	if _, ok := FindCase(cases, cases[0].ID); !ok {
		t.Error("existing case not found")
	}
	if _, ok := FindCase(cases, "nope"); ok {
		t.Error("unknown case must not be found")
	}
}
