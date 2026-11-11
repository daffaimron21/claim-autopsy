package main

import (
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"

	claimautopsy "github.com/daffaimron21/claim-autopsy"
	"github.com/daffaimron21/claim-autopsy/internal/api"
	"github.com/daffaimron21/claim-autopsy/internal/claim"
	"github.com/daffaimron21/claim-autopsy/internal/resume"
)

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func main() {
	casesFS, _ := fs.Sub(claimautopsy.Data, "testdata/cases")
	fixFS, _ := fs.Sub(claimautopsy.Data, "testdata/resume_fixtures")
	cases, err := claim.LoadCases(casesFS)
	if err != nil {
		log.Fatal(err)
	}
	saved := resume.Verified(resume.NewFixture(fixFS))

	live := newLive(os.Getenv, http.DefaultClient)
	if live != nil {
		log.Println("live model: on")
	} else {
		log.Println("live model: off (saved results only)")
	}

	mux := http.NewServeMux()
	mux.Handle("/api/", api.NewHandler(cases, live, saved, api.NewLimiter(6*time.Second, 3)))
	if dir := env("WEB_DIR", "web/dist"); dirExists(dir) {
		mux.Handle("/", http.FileServer(http.Dir(dir)))
	}
	srv := &http.Server{Addr: ":" + env("PORT", "8080"), Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	log.Println("listening on", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

// newLive returns the live resume reader, or nil unless LIVE_MODEL=true. With LIVE_MODEL=true and no key
// it is still built: every read fails, so reports are served from saved results and marked degraded.
func newLive(getenv func(string) string, client *http.Client) resume.Reader {
	if getenv("LIVE_MODEL") != "true" {
		return nil
	}
	model := getenv("LIVE_MODEL_NAME")
	if model == "" {
		model = "claude-haiku-5-5"
	}
	return resume.Verified(resume.NewAnthropic(getenv("ANTHROPIC_API_KEY"), model, client))
}

func dirExists(p string) bool { st, err := os.Stat(p); return err == nil && st.IsDir() }
