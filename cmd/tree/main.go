package main

import (
	"github.com/SmAlexAl/tree_server.git/internal/config"
	"github.com/SmAlexAl/tree_server.git/internal/httpServer/tree/add"
	"github.com/SmAlexAl/tree_server.git/internal/httpServer/tree/apply"
	delete2 "github.com/SmAlexAl/tree_server.git/internal/httpServer/tree/delete"
	"github.com/SmAlexAl/tree_server.git/internal/httpServer/tree/get"
	"github.com/SmAlexAl/tree_server.git/internal/httpServer/tree/initDb"
	"github.com/SmAlexAl/tree_server.git/internal/httpServer/tree/reset"
	"github.com/SmAlexAl/tree_server.git/internal/httpServer/tree/update"
	"github.com/SmAlexAl/tree_server.git/internal/lib/viewer/jstree"
	"github.com/SmAlexAl/tree_server.git/internal/storage/cache"
	"github.com/SmAlexAl/tree_server.git/internal/storage/fixtures"
	"github.com/SmAlexAl/tree_server.git/internal/storage/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	cfg := config.MustLoad()

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		os.Exit(1)
	}

	fixtures := fixtures.New()
	err = storage.SaveLeafs(fixtures.GetCollection())

	if err != nil {
		os.Exit(1)
	}

	cacheStorage, err := cache.New()

	if err != nil {
		os.Exit(1)
	}

	viewer := jstree.New()

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	//midleware

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/tree", add.New(cacheStorage, viewer))
	router.Post("/tree/init", initDb.New(cacheStorage, storage, viewer))
	router.Put("/tree/reset", reset.New(cacheStorage, storage, fixtures, viewer))

	router.Get("/tree", get.New(cacheStorage, storage, viewer))
	router.Put("/tree/update", update.New(cacheStorage, viewer))
	router.Put("/tree/delete", delete2.New(cacheStorage, viewer))
	router.Post("/tree/apply", apply.New(cacheStorage, storage, viewer))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		os.Exit(1)
	}
}
