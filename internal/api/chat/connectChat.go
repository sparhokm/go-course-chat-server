package chat

import (
	"github.com/sparhokm/go-course-ms-chat-server/internal/interceptor"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	desc "github.com/sparhokm/go-course-ms-chat-server/pkg/chat_v1"
)

func (i *Implementation) ConnectChat(req *desc.ConnectChatIn, stream desc.ChatV1_ConnectChatServer) error {
	userId, ok := interceptor.GetUserId(stream.Context())
	if !ok {
		return status.Errorf(codes.NotFound, "access denied")
	}

	err := i.chatService.CheckAccess(stream.Context(), req.GetChatId(), userId)
	if err != nil {
		return status.Errorf(codes.NotFound, "chat not found")
	}

	i.mxChannel.RLock()
	chatChan, ok := i.channels[req.GetChatId()]
	i.mxChannel.RUnlock()
	if !ok {
		i.mxChannel.Lock()
		chatChan = make(chan *desc.Message, 100)
		i.channels[req.GetChatId()] = chatChan
		i.mxChannel.Unlock()
	}

	i.mxChat.Lock()
	if _, okChat := i.chats[req.GetChatId()]; !okChat {
		i.chats[req.GetChatId()] = &Chat{
			streams: make(map[int64]desc.ChatV1_ConnectChatServer),
		}
	}
	i.mxChat.Unlock()

	i.chats[req.GetChatId()].m.Lock()
	i.chats[req.GetChatId()].streams[userId] = stream
	i.chats[req.GetChatId()].m.Unlock()

	for {
		select {
		case msg, okCh := <-chatChan:
			if !okCh {
				return nil
			}

			for _, st := range i.chats[req.GetChatId()].streams {
				if err := st.Send(msg); err != nil {
					return err
				}
			}

		case <-stream.Context().Done():
			i.chats[req.GetChatId()].m.Lock()
			delete(i.chats[req.GetChatId()].streams, userId)
			i.chats[req.GetChatId()].m.Unlock()
			return nil
		}
	}
}
