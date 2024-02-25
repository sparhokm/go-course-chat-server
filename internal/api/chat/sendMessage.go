package chat

import (
	"context"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/sparhokm/go-course-ms-chat-server/pkg/chat_v1"
)

func (i *Implementation) SendMessage(ctx context.Context, req *desc.SendMessageIn) (*emptypb.Empty, error) {
	_ = ctx
	log.Printf("Send message %+v", req)

	return &emptypb.Empty{}, nil
}
