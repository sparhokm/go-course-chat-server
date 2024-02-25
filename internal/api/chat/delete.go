package chat

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/sparhokm/go-course-ms-chat-server/pkg/chat_v1"
)

func (i *Implementation) Delete(ctx context.Context, req *desc.DeleteIn) (*emptypb.Empty, error) {
	err := i.chatService.Delete(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
