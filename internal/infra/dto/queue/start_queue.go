package queuedto

import (
	"errors"
	"time"

	"github.com/google/uuid"
	processentity "github.com/willjrcom/sales-backend-go/internal/domain/process"
)

var (
	ErrGroupItemIDRequired = errors.New("group item ID is required")
)

type StartQueueInput struct {
	processentity.QueueCommonAttributes
	JoinedAt time.Time `json:"joined_at"`
	IsTest   bool      `json:"is_test"`
}

func (s *StartQueueInput) validate() error {
	if s.GroupItemID == uuid.Nil {
		return ErrGroupItemIDRequired
	}

	if s.IsTest {
		s.JoinedAt = time.Now()
	}

	return nil
}

func (s *StartQueueInput) ToModel() (uuid.UUID, *time.Time, error) {
	if err := s.validate(); err != nil {
		return uuid.Nil, nil, err
	}

	return s.GroupItemID, &s.JoinedAt, nil
}