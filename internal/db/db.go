package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBConfig struct {
	DatabaseURL string

	// Connection Pool config
	PoolMaxConn          int
	PoolMinConn          int
	PoolMaxConnLifetime  int
	PoolMaxConnIdleTime  int
	PoolHeathCheckPeriod int
}

type Postgres struct {
	pool      *pgxpool.Pool
	tx        pgx.Tx
	TxOptions *pgx.TxOptions
	ctx       context.Context
}

type DBSession interface {
	Begin(ctx context.Context) (DBSession, error)
	Transaction(ctx context.Context, f func(txCtx context.Context) error) error
	Rollback() error
	Commit() error
	Context() context.Context
	Tx() pgx.Tx
	Close()
}

type DBKey struct{}

func (p Postgres) Begin(ctx context.Context) (DBSession, error) {
	tx, err := p.pool.BeginTx(ctx, *p.TxOptions)
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, DBKey{}, tx)

	return &Postgres{
		pool:      p.pool,
		tx:        tx,
		TxOptions: p.TxOptions,
		ctx:       ctx,
	}, nil
}

func (p *Postgres) Transaction(ctx context.Context, f func(context.Context) error) error {
	tx, err := p.pool.BeginTx(ctx, *p.TxOptions)
	if err != nil {
		return err
	}
	ctx = context.WithValue(ctx, DBKey{}, tx)
	p.ctx = ctx
	p.tx = tx
	if err := f(ctx); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func (p Postgres) Rollback() error {
	if p.tx == nil {
		return fmt.Errorf("no transaction to commit")
	}
	return p.tx.Rollback(context.WithValue(p.ctx, DBKey{}, p.tx))
}

func (p *Postgres) Commit() error {
	if p.tx == nil {
		return fmt.Errorf("no transaction to commit")
	}
	return p.tx.Commit(context.WithValue(p.ctx, DBKey{}, p.tx))
}

func (p *Postgres) Context() context.Context {
	return p.ctx
}

func (p Postgres) Tx() pgx.Tx {
	return p.tx
}

func (p *Postgres) Close() {
	p.pool.Close()
}

func ConnectPool(ctx context.Context, cfg *DBConfig, txOptions *pgx.TxOptions) (DBSession, error) {
	var pgPool *pgxpool.Pool

	config, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	// TODO: Create env
	config.MaxConns = int32(cfg.PoolMaxConn)
	config.MinConns = int32(cfg.PoolMinConn)
	config.MaxConnLifetime = time.Duration(cfg.PoolMaxConnLifetime) * time.Minute
	config.MaxConnIdleTime = time.Duration(cfg.PoolMaxConnIdleTime) * time.Minute
	config.HealthCheckPeriod = time.Duration(cfg.PoolHeathCheckPeriod) * time.Minute

	pgPool, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err = pgPool.Ping(ctx); err != nil {
		return nil, err
	}

	return &Postgres{
		pool:      pgPool,
		TxOptions: txOptions,
		ctx:       ctx,
	}, nil
}

func DB(ctx context.Context, fallback *pgxpool.Pool) *pgxpool.Pool {
	db := ctx.Value(DBKey{})
	if db == nil {
		return fallback
	}
	return db.(*pgxpool.Pool)
}
