package users

import (
	"errors"
	"sync"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}
type userrepository interface {
	create(user User) (User, error)
	getbyid(id int) (User, error)
	getAll() ([]User, error)
	update(users User) (User, error)
	delete(id int) error
}

type mockuserrepository struct {
	mu     sync.Mutex
	users  map[int]User
	nextid int
}

func Newmockuserrepository() userrepository {
	return &mockuserrepository{
		users:  make(map[int]User),
		nextid: 1,
	}

}
func (r *mockuserrepository) Create(user User) (User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	user.ID = r.nextid
	r.nextid++
	r.users[user.ID] = user
	return user, nil
}
func (r *mockuserrepository) Getbyid(id int) (User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	user, exists := r.users[id]
	if !exists {
		return user, errors.New("user not found")
	}
	return user, nil
}

func (r *mockuserrepository) GetAll() ([]User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []User
	for _, user := range r.users {
		result = append(result, user)
	}
	return result, nil
}

func (r *mockuserrepository) Update(users User) (user error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, exists := r.users[user.ID]
	if !exists {
		return errors.New("user not found")
	}
	r.users[users.ID] = users
	return nil

}
func (r *mockuserrepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, exists := r.users[id]
	if !exists {
		return errors.New("user not found")
	}
	delete(r.users, id)
	return nil

}
