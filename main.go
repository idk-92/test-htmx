package main

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"main/app/index"
	"main/app/leaderboard"
	"main/template"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	e := echo.New()

	// Other Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// svc := NewDBAppService(pool)

	// Configure the static file directory
	e.Static("/dist", "dist")
	template.NewTemplateRenderer(e)

	// Init routes
	index.InitRoutes(e)
	leaderboard.InitRoutes(e)

	e.Logger.Fatal(e.Start("localhost:1323"))
	return nil

}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
