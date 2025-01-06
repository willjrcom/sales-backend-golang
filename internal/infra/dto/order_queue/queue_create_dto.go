package orderqueuedto

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrGroupItemIDRequired = errors.New("group item ID is required")
)

type QueueCreateDTO struct {
	GroupItemID uuid.UUID `json:"group_item_id"`
	JoinedAt    time.Time `json:"joined_at"`
	IsTest      bool      `json:"is_test"`
}

func (s *QueueCreateDTO) validate() error {
	if s.GroupItemID == uuid.Nil {
		return ErrGroupItemIDRequired
	}

	if s.IsTest {
		s.JoinedAt = time.Now().UTC()
	}

	return nil
}

func (s *QueueCreateDTO) ToDomain() (uuid.UUID, *time.Time, error) {
	if err := s.validate(); err != nil {
		return uuid.Nil, nil, err
	}

	return s.GroupItemID, &s.JoinedAt, nil
}
