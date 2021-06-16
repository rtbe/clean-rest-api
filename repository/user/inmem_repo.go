package user

import (
	"fmt"
	"sync"

	"github.com/rtbe/clean-rest-api/domain/entity"
)

// InMemRepo is an abstraction layer that manages user entities inside basic in-memory store (map+RWMutex)
type InMemRepo struct {
	store map[string]entity.User
	sync.RWMutex
}

//NewInMemRepo creates a new in-memory repository for User entity
func NewInMemRepo() *InMemRepo {
	m := make(map[string]entity.User)
	return &InMemRepo{
		store: m,
	}
}

//Get gets user from in-memory store
func (r *InMemRepo) Get(id string) (*entity.User, error) {
	u, err := r.store[id]
	if !err {
		return nil, fmt.Errorf("there are no user with id: %s", id)
	}
	return &u, nil
}

//Create creates a new user in in-memory store
func (r *InMemRepo) Create(u *entity.User) (*entity.User, error) {
	r.Lock()
	defer r.Unlock()

	r.store[u.ID] = *u
	return u, nil
}

//Delete deletes user from in-memory store
func (r *InMemRepo) Delete(id string) error {
	r.Lock()
	defer r.Unlock()

	_, err := r.store[id]
	if !err {
		return fmt.Errorf("user with id: %v was not found", id)
	}
	delete(r.store, id)
	return nil
}

//List lists all users from in-memory store
func (r *InMemRepo) List() ([]entity.User, error) {
	if len(r.store) == 0 {
		return nil, fmt.Errorf("there are no users")
	}

	users := make([]entity.User, 0)
	for _, user := range r.store {
		users = append(users, user)
	}
	return users, nil
}
