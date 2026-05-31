package api

import (
	"bytes"
	_ "embed"
	"net/http"
	"time"
)

//go:embed spritesheet.webp
var Spritesheet []byte

func StartServer() {
	http.HandleFunc("/spritesheet", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/webp")
		http.ServeContent(w, r, "spritesheet.webp", time.Now(), bytes.NewReader(Spritesheet))
	})

	go http.ListenAndServe(":9090", nil)
}
