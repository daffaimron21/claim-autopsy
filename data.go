// Package claimautopsy embeds the synthetic data so the server is a single binary.
package claimautopsy

import "embed"

//go:embed testdata/cases testdata/resume_fixtures
var Data embed.FS
