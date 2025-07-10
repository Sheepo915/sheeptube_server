package repository

import (
	"fmt"
	"log/slog"
	"sheeptube/internal/db"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	DB     db.DBSession
	logger *slog.Logger

	VideoRepository
	ChannelRepository
}

func NewRepository(db *db.DBSession, logger *slog.Logger) *Repository {
	return &Repository{
		DB:     *db,
		logger: logger,
	}
}

func (r *Repository) WithTx(tx db.DBSession) *Repository {
	return &Repository{
		DB: tx,
	}
}

func (r *Repository) GetTx() (pgx.Tx, error) {
	if r.DB == nil {
		r.logger.Error("repository db is not initialized")
		return nil, fmt.Errorf("db is not initialized")
	}

	tx := r.DB.Tx()
	if tx == nil {
		r.logger.Error("db transaction is not initialized")
		return nil, fmt.Errorf("transaction is not initialized")
	}

	return tx, nil
}
