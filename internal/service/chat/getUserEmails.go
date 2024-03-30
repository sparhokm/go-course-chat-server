package chat

import (
	"context"
	"errors"
)

func (s *serv) CheckAccess(ctx context.Context, chatId int64, userId int64) error {
	ids, err := s.chatRepository.GetUserIds(ctx, chatId)

	if err != nil {
		return errors.New("chat not found")
	}

	found := false
	for _, id := range ids {
		if id == userId {
			found = true
			break
		}
	}

	if !found {
		return errors.New("chat not found")
	}

	return nil
}
