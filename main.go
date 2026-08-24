// confidential-secrets-demo proves possession of a keyserver-released secret
// without ever disclosing it: DEMO_SECRET arrives as an env var inside the
// enclave, and the only thing this server emits is an HMAC over a
// caller-chosen challenge.
package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func main() {
	secret := os.Getenv("DEMO_SECRET")

	http.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("GET /", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]bool{"secret_present": secret != ""})
	})

	// GET /prove?challenge=<text> returns hex(HMAC-SHA256(DEMO_SECRET, text)).
	// A verifier who knows the secret can recompute it; nobody else learns
	// anything about the value.
	http.HandleFunc("GET /prove", func(w http.ResponseWriter, r *http.Request) {
		challenge := r.URL.Query().Get("challenge")
		if secret == "" || challenge == "" {
			http.Error(w, "need a provisioned secret and a ?challenge=", http.StatusBadRequest)
			return
		}
		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write([]byte(challenge))
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"challenge": challenge,
			"hmac":      hex.EncodeToString(mac.Sum(nil)),
		})
	})

	log.Println("confidential-secrets-demo listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
