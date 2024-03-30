package chat

import (
	"context"

	desc "github.com/sparhokm/go-course-ms-chat-server/pkg/chat_v1"
)

func (i *Implementation) Create(ctx context.Context, req *desc.CreateIn) (*desc.CreateOut, error) {
	id, err := i.chatService.Create(ctx, req.GetUserIds())
	if err != nil {
		return nil, err
	}

	return &desc.CreateOut{
		Id: id,
	}, nil
}
