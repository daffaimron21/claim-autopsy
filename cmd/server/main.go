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

	var live resume.Reader
	if os.Getenv("LIVE_MODEL") == "true" && os.Getenv("ANTHROPIC_API_KEY") != "" {
		live = resume.Verified(resume.NewAnthropic(os.Getenv("ANTHROPIC_API_KEY"), env("LIVE_MODEL_NAME", "claude-haiku-5-5"), http.DefaultClient))
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

func dirExists(p string) bool { st, err := os.Stat(p); return err == nil && st.IsDir() }
