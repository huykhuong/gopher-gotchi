package api

import (
	"fmt"
	"gopher-gotchi/internal/brain"

	"net/http"
)

func StartServer(p *brain.Pet, spritesheet []byte) {
	http.HandleFunc("/tell", func(w http.ResponseWriter, r *http.Request) {
		cmd := r.URL.Query().Get("cmd")
		if cmd == "" {
			fmt.Fprint(w, "Diana is listening...")
			return
		}

		p.HandleCommand(cmd)
	})

	http.HandleFunc("/spritesheet", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/webp")
		w.Write(spritesheet)
	})

	go http.ListenAndServe(":9090", nil)
}
