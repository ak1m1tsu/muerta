package memo

import (
	"context"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/romankravchuk/muerta/internal/v2/data"
	"github.com/romankravchuk/muerta/internal/v2/lib/errors"
	"github.com/romankravchuk/muerta/internal/v2/storage/sessions"
)

type Storage struct {
	sessions   map[string]string
	sessionsMu sync.Mutex
}

func New() *Storage {
	return &Storage{
		sessions: make(map[string]string),
	}
}

func (s *Storage) Get(_ context.Context, key string) (*data.TokenDetails, error) {
	const op = "storage.sessions.memo.Storage.Get"

	s.sessionsMu.Lock()
	defer s.sessionsMu.Unlock()

	session, ok := s.sessions[key]
	if !ok {
		return nil, errors.WithOp(op, sessions.ErrTokenNotFound)
	}

	var details *data.TokenDetails
	if err := sonic.UnmarshalString(session, &details); err != nil {
		return nil, err
	}

	return details, nil
}

func (s *Storage) Set(_ context.Context, details *data.TokenDetails) error {
	const op = "storage.sessions.memo.Storage.Set"

	s.sessionsMu.Lock()
	defer s.sessionsMu.Unlock()

	session, err := sonic.MarshalString(details)
	if err != nil {
		return errors.WithOp(op, err)
	}

	s.sessions[details.Payload.Email] = session

	return nil
}
