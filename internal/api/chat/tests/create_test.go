package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sparhokm/go-course-ms-chat-server/internal/api/chat"
	chatServiceMock "github.com/sparhokm/go-course-ms-chat-server/internal/service/mocks"
	desc "github.com/sparhokm/go-course-ms-chat-server/pkg/chat_v1"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	var (
		ctx = context.Background()

		serviceErr = fmt.Errorf("service error")

		usernames = []string{"Name", "Name2"}

		id        int64 = 100
		createOut       = &desc.CreateOut{
			Id: id,
		}
	)

	tests := []struct {
		name      string
		prepare   func(service *chatServiceMock.MockChatService)
		createIn  *desc.CreateIn
		want      *desc.CreateOut
		expectErr bool
	}{
		{
			name: "success case",
			createIn: &desc.CreateIn{
				Usernames: usernames,
			},
			want:      createOut,
			expectErr: false,
			prepare: func(m *chatServiceMock.MockChatService) {
				m.EXPECT().Create(ctx, usernames).Return(id, nil)
			},
		},
		{
			name: "user service error",
			createIn: &desc.CreateIn{
				Usernames: usernames,
			},
			want:      nil,
			expectErr: true,
			prepare: func(m *chatServiceMock.MockChatService) {
				m.EXPECT().Create(ctx, usernames).Return(0, serviceErr)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			serviceMock := chatServiceMock.NewMockChatService(t)
			tt.prepare(serviceMock)
			api := chat.NewImplementation(serviceMock)

			res, err := api.Create(ctx, tt.createIn)

			require.Equal(t, tt.want, res)
			require.Equal(t, tt.expectErr, err != nil)
		})
	}
}
