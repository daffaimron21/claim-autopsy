package resume

import (
	"context"
	"encoding/json"
	"errors"
	"io/fs"

	"github.com/daffaimron21/claim-autopsy/internal/claim"
)

type fixtureReader struct{ fsys fs.FS }

// NewFixture serves saved findings from <caseID>.json. A missing file means no findings.
func NewFixture(fsys fs.FS) Reader { return fixtureReader{fsys} }

func (f fixtureReader) Read(_ context.Context, c claim.Case) (Result, error) {
	b, err := fs.ReadFile(f.fsys, c.ID+".json")
	if errors.Is(err, fs.ErrNotExist) {
		return Result{}, nil
	}
	if err != nil {
		return Result{}, err
	}
	var fs []claim.Finding
	if err := json.Unmarshal(b, &fs); err != nil {
		return Result{}, err
	}
	return Result{Findings: fs}, nil
}
