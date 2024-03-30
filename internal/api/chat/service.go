package chat

import (
	"sync"

	"github.com/sparhokm/go-course-ms-chat-server/internal/service"
	desc "github.com/sparhokm/go-course-ms-chat-server/pkg/chat_v1"
)

type Chat struct {
	streams map[int64]desc.ChatV1_ConnectChatServer
	m       sync.RWMutex
}

type Implementation struct {
	desc.UnimplementedChatV1Server
	chatService service.ChatService

	chats  map[int64]*Chat
	mxChat sync.RWMutex

	channels  map[int64]chan *desc.Message
	mxChannel sync.RWMutex
}

func NewImplementation(chatService service.ChatService) *Implementation {
	return &Implementation{
		chatService: chatService,
		channels:    make(map[int64]chan *desc.Message),
		chats:       make(map[int64]*Chat),
	}
}
