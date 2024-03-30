package service

import (
	"context"
)

type ChatService interface {
	Create(ctx context.Context, userIds []int64) (int64, error)
	CheckAccess(ctx context.Context, chatId int64, userId int64) error
	Delete(ctx context.Context, id int64) error
}
