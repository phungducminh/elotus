package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"elotus.com/hackathon/storage/query"
	_ "github.com/go-sql-driver/mysql"
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
}

func NewStorage(cfg *Config) Storage {
	db, err := sql.Open("mysql", cfg.Datasource)
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Second * time.Duration(cfg.ConnMaxLifeTimeInSeconds))
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	return &storage{db}
}

type storage struct {
	db *sql.DB
}

func (c *storage) InsertUser(r *query.InsertUserParams) (int64, error) {
	queries := query.New(c.db)
	result, err := queries.InsertUser(context.TODO(), *r)
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
