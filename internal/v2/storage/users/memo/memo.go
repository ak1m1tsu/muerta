package memo

import (
	"context"
	"sync"

	"github.com/romankravchuk/muerta/internal/v2/data"
	"github.com/romankravchuk/muerta/internal/v2/lib/errors"
	"github.com/romankravchuk/muerta/internal/v2/storage/users"
)

type Storage struct {
	users   map[string]*data.User
	usersMu sync.Mutex
}

func New() *Storage {
	return &Storage{
		users: make(map[string]*data.User),
	}
}

func (s *Storage) Create(_ context.Context, user *data.User) error {
	const op = "storage.users.memo.Storage.Create"

	s.usersMu.Lock()
	defer s.usersMu.Unlock()

	if _, ok := s.users[user.Email]; ok {
		return errors.WithOp(op, users.ErrExists)
	}

	s.users[user.Email] = user
	return nil
}

func (s *Storage) Update(_ context.Context, user *data.User) error {
	const op = "storage.users.memo.Storage.Update"

	s.usersMu.Lock()
	defer s.usersMu.Unlock()

	if _, ok := s.users[user.Email]; !ok {
		return errors.WithOp(op, users.ErrNotFound)
	}

	s.users[user.Email] = user
	return nil
}

func (s *Storage) Delete(_ context.Context, email string) error {
	const op = "storage.users.memo.Storage.Delete"

	s.usersMu.Lock()
	defer s.usersMu.Unlock()

	if _, ok := s.users[email]; !ok {
		return errors.WithOp(op, users.ErrNotFound)
	}

	delete(s.users, email)
	return nil
}

func (s *Storage) FindByEmail(_ context.Context, email string) (*data.User, error) {
	const op = "storage.users.memo.Storage.FindByEmail"

	s.usersMu.Lock()
	defer s.usersMu.Unlock()

	if user, ok := s.users[email]; ok {
		return &data.User{
			ID:                user.ID,
			FirstName:         user.FirstName,
			LastName:          user.LastName,
			Email:             user.Email,
			EncryptedPassword: user.EncryptedPassword,
			CreatedAt:         user.CreatedAt,
			UpdatedAt:         user.UpdatedAt,
			DeletedAt:         user.DeletedAt,
		}, nil
	}

	return nil, errors.WithOp(op, users.ErrNotFound)
}

func (s *Storage) FindMany(_ context.Context, _ data.UserFilter) ([]data.User, error) {
	const op = "storage.users.memo.Storage.FindMany"

	s.usersMu.Lock()
	defer s.usersMu.Unlock()

	if len(s.users) == 0 {
		return nil, errors.WithOp(op, users.ErrNotFound)
	}

	users := make([]data.User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, *user)
	}

	return users, nil
}
