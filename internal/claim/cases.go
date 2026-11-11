package claim

import (
	"encoding/json"
	"fmt"
	"io/fs"
)

// LoadCases reads every *.json file at the root of fsys, in filename order.
func LoadCases(fsys fs.FS) ([]Case, error) {
	names, err := fs.Glob(fsys, "*.json")
	if err != nil {
		return nil, err
	}
	cases := make([]Case, 0, len(names))
	for _, name := range names {
		b, err := fs.ReadFile(fsys, name)
		if err != nil {
			return nil, err
		}
		var c Case
		if err := json.Unmarshal(b, &c); err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		cases = append(cases, c)
	}
	return cases, nil
}

func FindCase(cases []Case, id string) (Case, bool) {
	for _, c := range cases {
		if c.ID == id {
			return c, true
		}
	}
	return Case{}, false
}
