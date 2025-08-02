package storage

import (
	"fmt"
	"math/rand"
)

var (
	ErrNotFound = fmt.Errorf("not found")
)

type Config struct {
	Datasource string
}

type UserRecord struct {
	ID       int64
	Username string
	HashedPassword string
}

type Storage interface {
	InsertUser(*UserRecord) (int64, error)
	GetUserByUserName(username string) (*UserRecord, error)
}

func NewStorage(cfg *Config) Storage {
	return &storage{}
}

type storage struct {
}

func (c *storage) InsertUser(r *UserRecord) (int64, error) {
	return rand.Int63(), nil
}

func (c *storage) GetUserByUserName(username string) (*UserRecord, error) {
	return nil, nil
}
