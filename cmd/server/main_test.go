package main

import (
	"context"
	"net/http"
	"testing"

	"github.com/daffaimron21/claim-autopsy/internal/claim"
)

func fakeEnv(m map[string]string) func(string) string { return func(k string) string { return m[k] } }

func TestNewLive_Off_IsNil(t *testing.T) {
	if r := newLive(fakeEnv(map[string]string{"ANTHROPIC_API_KEY": "k"}), http.DefaultClient); r != nil {
		t.Fatal("live must be nil unless LIVE_MODEL=true")
	}
}

func TestNewLive_OnWithoutKey_StillLiveSoItDegradesVisibly(t *testing.T) {
	r := newLive(fakeEnv(map[string]string{"LIVE_MODEL": "true"}), http.DefaultClient)
	if r == nil {
		t.Fatal("LIVE_MODEL=true with no key must still build the live reader so reports are marked degraded")
	}
	if _, err := r.Read(context.Background(), claim.Case{ID: "x"}); err == nil {
		t.Fatal("empty key must make the live read fail")
	}
}
