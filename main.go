package main

import (
	"net/http"
	"log"
	"os"
	"time"
	"context"

	"github.com/google/easypki/pkg/easypki"
	"github.com/google/easypki/pkg/store"
	"github.com/boltdb/bolt"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"

	"easypki-ui/settings"
	"easypki-ui/config"
	"easypki-ui/api"
	"os/signal"
	"syscall"
)

func main() {
	sp := settings.PkiSettings{}
	sp.Create()

	if sp.DbPath == "" {
		log.Fatal("Arg db_path must be set.")
	}
	if sp.BundleName == "" && sp.ConfigPath == "" {
		log.Fatal("One of bundle_name or config_path must be set.")
	}

	db, err := bolt.Open(sp.DbPath, 0600, nil)
	if err != nil {
		log.Fatalf("Failed opening bolt database %v: %v", sp.DbPath, err)
	}
	defer db.Close()

	cfg := config.Config{
		Store:   &config.Yaml{Path: sp.ConfigPath},
		EasyPKI: &easypki.EasyPKI{Store: &store.Bolt{DB: db}},
	}
	cfg.Init()

	ws := settings.WebServerSettings{}
	ws.Create()

	r := mux.NewRouter()


	a := api.API{}
	a.Setup(&cfg, r.PathPrefix("/api").Subrouter())

	r.Use(mux.CORSMethodMiddleware(r))

	corsOpts := handlers.AllowedOrigins([]string{"http://localhost:8080"})

	srv := &http.Server{
		Addr: ws.Address,

		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,

		Handler: handlers.CORS(corsOpts)(handlers.LoggingHandler(os.Stdout, r)),
	}

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
			// Lets exit if we hit this error
			c <- syscall.SIGTERM
		}
	}()

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), ws.GracefulTimeout)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
