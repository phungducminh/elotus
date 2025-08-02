package db

import "math/rand"


type UserRecord struct {
	ID       int64
	Username string
	Password string
}

type Client interface {
	InsertUser(UserRecord) (int64, error)
}

func NewClient(cfg *Config) Client {
	return &client{}
}

type client struct {
}

func (c *client) InsertUser(r UserRecord) (int64, error) {
	return rand.Int63(), nil
}
