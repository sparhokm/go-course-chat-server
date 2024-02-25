package repository

import (
	"context"
)

type ChatRepository interface {
	Create(ctx context.Context, userNames []string) (int64, error)
	Delete(ctx context.Context, id int64) error
}
