package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"elotus.com/hackathon/storage/query"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

var (
	ErrNotFound = fmt.Errorf("not found")
)

type Config struct {
	Datasource               string
	ConnMaxLifeTimeInSeconds int
	MaxOpenConns             int
	MaxIdleConns             int
}

type UserRecord struct {
	ID             int64
	Username       string
	HashedPassword string
}

type Storage interface {
	InsertUser(*query.InsertUserParams) (int64, error)
	GetUserByUserName(username string) (*query.User, error)
	InsertFile(*query.InsertFileParams) (int64, error)
}

func NewStorage(lg *zap.Logger, cfg *Config) Storage {
	db, err := sql.Open("mysql", cfg.Datasource)
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Second * time.Duration(cfg.ConnMaxLifeTimeInSeconds))
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	return &storage{
		db: db,
		lg: lg,
	}
}

type storage struct {
	db *sql.DB
	lg *zap.Logger
}

func (c *storage) InsertUser(params *query.InsertUserParams) (int64, error) {
	queries := query.New(c.db)
	result, err := queries.InsertUser(context.TODO(), *params)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (c *storage) GetUserByUserName(username string) (*query.User, error) {
	queries := query.New(c.db)
	user, err := queries.GetUserByUsername(context.TODO(), username)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *storage) InsertFile(params *query.InsertFileParams) (int64, error) {
	queries := query.New(c.db)
	result, err := queries.InsertFile(context.TODO(), *params)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
