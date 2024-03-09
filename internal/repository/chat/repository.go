package user

import (
	"context"
	"fmt"

	"github.com/sparhokm/go-course-ms-auth/pkg/client/db"

	sq "github.com/Masterminds/squirrel"
)

const (
	tableName = "chats"

	idColumn        = "id"
	userNamesColumn = "user_names"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) *repo {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, userNames []string) (int64, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(userNamesColumn).
		Values(userNames).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}
	var id int64
	err = r.db.DB().ScanOneContext(ctx, &id, db.Query{Name: "chat.create", QueryRaw: query}, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to create chat: %v", err)
	}

	return id, nil
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	builder := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.DB().ExecContext(ctx, db.Query{Name: "user.delete", QueryRaw: query}, args...)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	return nil
}
