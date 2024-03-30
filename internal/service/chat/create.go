package chat

import (
	"context"
)

func (s *serv) Create(ctx context.Context, userIds []int64) (int64, error) {
	id, err := s.chatRepository.Create(ctx, userIds)
	if err != nil {
		return 0, err
	}

	return id, nil
}
