package config

import (
	"flag"
	"os"
	"sheeptube/internal/minio"
)

type Config struct {
	Port   string
	Env    string
	DBPath string

	*minio.MinioConfig
}

func (c *Config) ParseFlag() {
	flag.StringVar(&c.Port, "port", os.Getenv("PORT"), "Server address port")
	flag.StringVar(&c.Env, "env", os.Getenv("ENV"), "Environment config (development|staging|production)")
	flag.StringVar(&c.DBPath, "db_path", os.Getenv("DB_PATH"), "Database path")

	flag.Parse()
}
