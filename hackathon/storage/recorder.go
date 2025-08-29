package storage

import (
	"math/rand"
	"sync"

	"elotus.com/hackathon/storage/query"
)

// ensure compile-time check StorageRecorder imply Storage interface
var _ Storage = NewRecorder()

// StorageRecorder: used for internal testing
type StorageRecorder struct {
	mu    sync.RWMutex
	users map[string]*query.User
}

func NewRecorder() *StorageRecorder {
	return &StorageRecorder{
		mu:    sync.RWMutex{},
		users: make(map[string]*query.User),
	}
}

func (c *StorageRecorder) InsertUser(r *query.InsertUserParams) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	id := rand.Int63()
	user := &query.User{
		ID:             id,
		Username:       r.Username,
		HashedPassword: r.HashedPassword,
	}
	c.users[r.Username] = user
	return id, nil
}

func (c *StorageRecorder) GetUserByUserName(username string) (*query.User, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	usr, ok := c.users[username]
	if !ok {
		return nil, ErrNotFound
	}
	return usr, nil
}

func (c *StorageRecorder) InsertFile(*query.InsertFileParams) (int64, error) {
	return rand.Int63(), nil
}
