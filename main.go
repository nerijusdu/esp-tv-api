package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/nerijusdu/esp-tv-api/src/providers"
	"github.com/nerijusdu/esp-tv-api/src/util"
)

var providerMap = map[string]providers.Provider{
	"bsky":    &providers.BskyProvider{},
	"video":   &providers.VideoProvider{},
	"image":   &providers.ImageProvider{},
	"time":    &providers.TimeProvider{},
	"posthog": &providers.PosthogProvider{},
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	config, err := util.LoadConfig()
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "www/index.html")
	})

	r.Get("/api", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	allProviders := []providers.Provider{}
	for name := range config.Providers {
		provider, found := providerMap[name]
		if !found {
			fmt.Printf("Provider %s not found\n", name)
			continue
		}

		err := provider.Init(config.Providers[name])
		if err != nil {
			fmt.Printf("Provider %s init failed: %s\n", name, err)
			continue
		}

		allProviders = append(allProviders, provider)
	}

	index := 0
	cursor := ""

	r.Get("/api/tv", func(w http.ResponseWriter, r *http.Request) {
		response, error := allProviders[index].GetView(cursor)
		if error != nil {
			fmt.Printf("Provider %d failed: %s\n", index, error)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(error.Error()))
			return
		}

		view := response.View

		if response.NextCursor == "" {
			cursor = ""
			index++
			if index >= len(allProviders) {
				index = 0
			}
			if config.ViewDelay > 0 {
				view.RefreshAfter = config.ViewDelay
			}
		} else {
			cursor = response.NextCursor
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", fmt.Sprint(len(view.Data)))
		w.Header().Set("X-Refresh-After", fmt.Sprint(view.RefreshAfter))
		w.Write(view.Data)
	})

	port := 8080
	if config.Server.Port > 0 {
		port = config.Server.Port
	}

	fmt.Printf("Listening on :%d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
