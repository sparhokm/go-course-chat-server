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

func TestDelete(t *testing.T) {
	t.Parallel()

	var (
		ctx = context.Background()

		serviceErr = fmt.Errorf("service error")

		id       int64 = 100
		deleteIn       = &desc.DeleteIn{
			Id: id,
		}
	)

	tests := []struct {
		name      string
		prepare   func(service *chatServiceMock.MockChatService)
		request   *desc.DeleteIn
		want      *desc.CreateOut
		expectErr bool
	}{
		{
			name:      "success case",
			request:   deleteIn,
			expectErr: false,
			prepare: func(m *chatServiceMock.MockChatService) {
				m.EXPECT().Delete(ctx, id).Return(nil)
			},
		},
		{
			name:      "chat service error",
			request:   deleteIn,
			expectErr: true,
			prepare: func(m *chatServiceMock.MockChatService) {
				m.EXPECT().Delete(ctx, id).Return(serviceErr)
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

			_, err := api.Delete(ctx, tt.request)

			require.Equal(t, tt.expectErr, err != nil)
		})
	}
}
