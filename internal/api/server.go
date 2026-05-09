package api

import (
	"fmt"
	"gopher-gotchi/internal/brain"
	"gopher-gotchi/internal/window"
	"net/http"
)

func StartServer(p *brain.Pet) {
	http.HandleFunc("/tell", func(w http.ResponseWriter, r *http.Request) {
		cmd := r.URL.Query().Get("cmd")
		if cmd == "" {
			fmt.Fprint(w, "Diana is listening...")
			return
		}

		response := p.HandleCommand(cmd)
		fmt.Fprint(w, response)
	})

	http.HandleFunc("/spritesheet", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/webp")
		w.Write([]byte(window.Spritesheet))
	})

	go http.ListenAndServe(":9090", nil)
}