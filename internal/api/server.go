package api

import (
	"net/http"
)

func StartServer(spritesheet []byte) {
	http.HandleFunc("/spritesheet", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/webp")
		w.Write(spritesheet)
	})

	go http.ListenAndServe(":9090", nil)
}
