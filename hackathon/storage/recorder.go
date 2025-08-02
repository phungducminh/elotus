package storage

import (
	"math/rand"
	"sync"
)

// StorageRecorder: used for internal testing
type StorageRecorder struct {
	mu    sync.RWMutex
	users map[string]*UserRecord
}

func NewRecorder() *StorageRecorder {
	return &StorageRecorder{
		mu:    sync.RWMutex{},
		users: make(map[string]*UserRecord),
	}
}

func (c *StorageRecorder) InsertUser(r *UserRecord) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	id := rand.Int63()
	r.ID = id
	c.users[r.Username] = r
	return id, nil
}

func (c *StorageRecorder) GetUserByUserName(username string) (*UserRecord, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	usr, ok := c.users[username]
	if !ok {
		return nil, ErrNotFound
	}
	return usr, nil
}
