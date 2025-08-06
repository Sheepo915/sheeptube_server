package main

import (
	// "context"
	// "os"
	"sheeptube/internal/app"
	// "sheeptube/internal/app/config"
	"sheeptube/internal/app/handler"
	"sheeptube/internal/app/service"
	// "sheeptube/internal/db"
	// "sheeptube/internal/minio"

	// "github.com/jackc/pgx/v5"
)

func main() {
	// ctx := context.Background()

	// cfg := config.Config{}
	// cfg.ParseFlag()

	// conn, err := pgx.Connect(ctx, "user=admin password=p@ssw0rd dbname=sheep_tube")
	// if err != nil {
	// 	os.Exit(1)
	// }
	// defer conn.Close(ctx)

	// queries := db.New(conn)

	// m, err := minio.NewClient(cfg.MinioConfig)
	// if err != nil {
	// 	os.Exit(1)
	// }

	// s := service.NewService(queries)
	s := service.NewService()
	h := handler.NewHandler(s)

	app := app.NewApp(h)

	app.Run()
}
