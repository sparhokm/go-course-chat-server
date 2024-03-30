package user

import (
	"context"
	"fmt"

	"github.com/sparhokm/go-course-ms-auth/pkg/client/db"

	sq "github.com/Masterminds/squirrel"
)

const (
	tableName = "chats"

	idColumn      = "id"
	userIdsColumn = "user_ids"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) *repo {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, userIds []int64) (int64, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(userIdsColumn).
		Values(userIds).
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

	tag, err := r.db.DB().ExecContext(ctx, db.Query{Name: "user.delete", QueryRaw: query}, args...)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *repo) GetUserIds(ctx context.Context, id int64) ([]int64, error) {
	builder := sq.Select(userIdsColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{idColumn: id}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	row := r.db.DB().QueryRowContext(ctx, db.Query{Name: "user.get", QueryRaw: query}, args...)

	var ids []int64
	err = row.Scan(&ids)
	if err != nil {
		return nil, err
	}

	return ids, nil
}
