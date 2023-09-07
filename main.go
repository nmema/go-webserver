package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	router := chi.NewRouter()

	apiCfg := apiConfig{
		fileserverHits: 0,
	}

	fshandler := apiCfg.middlewareMetricsInc(
		http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))),
	)

	router.Handle("/app", fshandler)
	router.Handle("/app/*", fshandler)
	router.Get("/healthz", handlerReadiness)
	router.Get("/metrics", apiCfg.handlerMetrics)
	router.Get("/reset", apiCfg.handlerReset)

	corsMux := middlewareCors(router)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
