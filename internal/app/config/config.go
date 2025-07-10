package config

import (
	"flag"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
)

type Config struct {
	Port   string
	Env    string
	DBPath string
}

func (c *Config) ParseFlag() {
	flag.StringVar(&c.Port, "port", os.Getenv("PORT"), "Server address port")
	flag.StringVar(&c.Env, "env", os.Getenv("ENV"), "Environment config (development|staging|production)")
	flag.StringVar(&c.DBPath, "db_path", os.Getenv("DB_PATH"), "Database path")

	flag.Parse()
}

func InitializeRepository(db pgx.Conn, logger *slog.Logger) {

}
