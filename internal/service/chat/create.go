package chat

import (
	"context"
)

func (s *serv) Create(ctx context.Context, userNames []string) (int64, error) {
	id, err := s.chatRepository.Create(ctx, userNames)
	if err != nil {
		return 0, err
	}

	return id, nil
}
