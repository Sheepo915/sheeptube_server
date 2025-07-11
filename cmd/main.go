package main

import (
	"context"
	"os"
	"sheeptube/internal/app"
	"sheeptube/internal/app/handler"
	"sheeptube/internal/app/service"
	"sheeptube/internal/db"

	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, "user=admin password=p@ssw0rd dbname=sheep_tube")
	if err != nil {
		os.Exit(1)
	}
	defer conn.Close(ctx)

	queries := db.New(conn)

	s := service.NewService(queries)
	h := handler.NewHandler(s)

	app := app.NewApp(h)

	app.Run()
}
