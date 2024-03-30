package chat

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/sparhokm/go-course-ms-chat-server/internal/interceptor"
	desc "github.com/sparhokm/go-course-ms-chat-server/pkg/chat_v1"
)

func (i *Implementation) SendMessage(ctx context.Context, req *desc.SendMessageIn) (*emptypb.Empty, error) {
	userId, ok := interceptor.GetUserId(ctx)
	if !ok {
		return nil, status.Errorf(codes.NotFound, "access denied")
	}

	err := i.chatService.CheckAccess(ctx, req.GetChatId(), userId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "chat not found")
	}

	i.mxChannel.RLock()
	chatChan, ok := i.channels[req.GetChatId()]
	i.mxChannel.RUnlock()

	if !ok {
		return nil, status.Errorf(codes.NotFound, "chat not found")
	}

	chatChan <- &desc.Message{
		From:      userId,
		Text:      req.GetText(),
		CreatedAt: timestamppb.New(time.Now()),
	}

	return &emptypb.Empty{}, nil
}
