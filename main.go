package main

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"log"
	"main/app/index"
	"main/app/leaderboard"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/cors"
)

type DBAppService struct {
	pool *pgxpool.Pool
}

// var _ AppService = (*DBAppService)(nil)

func NewDBAppService(pool *pgxpool.Pool) *DBAppService {
	return &DBAppService{pool: pool}
}

func run(ctx context.Context, w io.Writer, args []string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	// pool, err := db.Connect(ctx)
	// if err != nil {
	// 	panic(err)
	// }
	server := &http.Server{Addr: "0.0.0.0:3333", Handler: service()}
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	// Run the server
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	return nil

}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
func service() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		// AllowedOrigins: []string{"https://*sprkl.netlify.app", "http://localhost:5173"},
		AllowedHeaders: []string{"Connect-Protocol-Version", "Content-Type", "X-Custom-Header"},
		// Debug: true,
	})

	r.Use(c.Handler)
	index.InitRoutes(r)
	leaderboard.InitRoutes(r)
	// r.Get("/", templ.Handler(index.IndexPage("asd")).ServeHTTP)
	// r.Get("/leaderboard", templ.Handler(leaderboard.LeaderboardPage()).ServeHTTP)
	// leaderboard.InitRoutes(r)

	r.Get("/slow", func(w http.ResponseWriter, r *http.Request) {
		// Simulates some hard work.
		//
		// We want this handler to complete successfully during a shutdown signal,
		// so consider the work here as some background routine to fetch a long running
		// search query to find as many results as possible, but, instead we cut it short
		// and respond with what we have so far. How a shutdown is handled is entirely
		// up to the developer, as some code blocks are preemptible, and others are not.
		time.Sleep(5 * time.Second)

		w.Write([]byte(fmt.Sprintf("all done.\n")))
	})
	fs := http.FileServer(http.Dir("dist"))

	r.Handle("/dist/*", http.StripPrefix("/dist/", fs))

	return r
}
