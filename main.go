package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nerijusdu/esp-tv-api/src/providers"
)

const DISPLAY_WIDTH = 128
const DISPLAY_HEIGHT = 64

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "www/index.html")
	})

	r.Get("/api", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	allProviders := []providers.Provider{
		&providers.PosthogProvider{},
	}
	index := 0
	cursor := ""

	r.Get("/api/tv", func(w http.ResponseWriter, r *http.Request) {
		response := allProviders[index].GetView(cursor)
		view := response.View

		if response.NextCursor == "" {
			cursor = ""
			index++
			if index >= len(allProviders) {
				index = 0
			}
		} else {
			cursor = response.NextCursor
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", fmt.Sprint(len(view.Data)))
		w.Header().Set("X-Refresh-After", fmt.Sprint(view.RefreshAfter))
		w.Write(view.Data)
	})

	fmt.Println("Listening on :8080")
	http.ListenAndServe(":8080", r)
}
