package app

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"

	"github.com/sparhokm/go-course-ms-auth/pkg/client/db"
	"github.com/sparhokm/go-course-ms-auth/pkg/client/db/pg"
	"github.com/sparhokm/go-course-ms-auth/pkg/client/db/transaction"

	"github.com/sparhokm/go-course-ms-chat-server/internal/api/chat"
	"github.com/sparhokm/go-course-ms-chat-server/internal/closer"
	"github.com/sparhokm/go-course-ms-chat-server/internal/config"
	"github.com/sparhokm/go-course-ms-chat-server/internal/repository"
	chatRepository "github.com/sparhokm/go-course-ms-chat-server/internal/repository/chat"
	"github.com/sparhokm/go-course-ms-chat-server/internal/service"
	chatService "github.com/sparhokm/go-course-ms-chat-server/internal/service/chat"
)

type serviceProvider struct {
	config *config.Config

	dbClient  db.Client
	txManager db.TxManager

	chatRepository repository.ChatRepository

	chatService service.ChatService

	chatImpl *chat.Implementation
}

func newServiceProvider(config *config.Config) *serviceProvider {
	return &serviceProvider{config: config}
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		dbc, err := pgxpool.New(ctx, s.config.PGConfig.DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		cl := pg.New(dbc)
		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) ChatRepository(ctx context.Context) repository.ChatRepository {
	if s.chatRepository == nil {
		s.chatRepository = chatRepository.NewRepository(s.DBClient(ctx))
	}

	return s.chatRepository
}

func (s *serviceProvider) ChatService(ctx context.Context) service.ChatService {
	if s.chatService == nil {
		s.chatService = chatService.NewService(
			s.ChatRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.chatService
}

func (s *serviceProvider) ChatImpl(ctx context.Context) *chat.Implementation {
	if s.chatImpl == nil {
		s.chatImpl = chat.NewImplementation(s.ChatService(ctx))
	}

	return s.chatImpl
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	return s.config.GRPCConfig
}
