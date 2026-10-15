package main

import (
	"crypto/rand"
	mathrand "math/rand"
	"fmt"
	"log"
	"net/http"
)

var companies = []string{
	"google",
	"anthropic",
	"youtube",
	"tesla",
	"facebook",
	"bbc",
	"apple",
	"amazon",
	"microsoft",
	"openai",
	"pinterest",
	"samsung",
	"oracle",
	"medium",
	"stackoverflow",
	"canva",
	"figma",
}

// Generate a UUID v4 without external dependencies.
func uuid4() (string, error) {
	var b [16]byte

	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}

	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf(
		"%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16],
	), nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fmt.Fprintln(w, "<!DOCTYPE html>")
	fmt.Fprintln(w, "<html><body>")
	fmt.Fprintf(w, "<h1>Mock Crawler Page: %s</h1>\n", r.URL.Path)
	fmt.Fprintln(w, "<h2>Discovered URLs</h2><ul>")

	for i := 0; i < 20; i++ {
		company := companies[mathrand.Intn(len(companies))]

		id, err := uuid4()
		if err != nil {
			http.Error(w, "Failed to generate UUID", http.StatusInternalServerError)
			return
		}

		url := fmt.Sprintf("http://localhost:8000/%s/%s", company, id)

		fmt.Fprintf(
			w,
			`<li><a href="%s">%s</a></li>`+"\n",
			url,
			url,
		)
	}

	fmt.Fprintln(w, "</ul></body></html>")
}

func main() {
	http.HandleFunc("/", handler)

	log.Println("Mock web server running at http://localhost:8000")
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}