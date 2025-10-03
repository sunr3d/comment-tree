package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"

	"github.com/sunr3d/comment-tree/internal/config"
	"github.com/sunr3d/comment-tree/internal/interfaces/infra"
	"github.com/sunr3d/comment-tree/models"
)

const (
	qCreate        = `INSERT INTO comments (parent_id, content, author) VALUES ($1, $2, $3)`
	qGetByID       = `SELECT id, parent_id, content, author, created_at, updated_at, deleted_at FROM comments WHERE id = $1 AND deleted_at IS NULL`
	qGetByParentID = `
    WITH RECURSIVE comment_tree AS (
        SELECT id, parent_id, content, author, created_at, updated_at, deleted_at, 0 as level
        FROM comments 
        WHERE id = $1 AND deleted_at IS NULL
        
        UNION ALL
        
        SELECT c.id, c.parent_id, c.content, c.author, c.created_at, c.updated_at, c.deleted_at, ct.level + 1
        FROM comments c
        INNER JOIN comment_tree ct ON c.parent_id = ct.id
        WHERE c.deleted_at IS NULL
    )
    SELECT id, parent_id, content, author, created_at, updated_at, deleted_at, level 
    FROM comment_tree 
    ORDER BY level, created_at`
	qDelete = `UPDATE comments SET deleted_at = NOW() WHERE id = $1`
)

var _ infra.Database = (*postgresRepo)(nil)

type postgresRepo struct {
	db *dbpg.DB
}

func New(ctx context.Context, cfg config.DBConfig) (infra.Database, error) {
	db, err := dbpg.New(cfg.DSN, nil, &dbpg.Options{})
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("dbpg.New")
		return nil, fmt.Errorf("не удалось создать подключение к БД: %w", err)
	}

	pCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := db.Master.PingContext(pCtx); err != nil {
		_ = db.Master.Close()
		zlog.Logger.Error().Err(err).Msg("db.Master.PingContext")
		return nil, fmt.Errorf("таймаут пинг к БД: %w", err)
	}

	return &postgresRepo{db: db}, nil
}

func (r *postgresRepo) Create(ctx context.Context, comment *models.Comment) error {
	_, err := r.db.ExecWithRetry(
		ctx,
		retry.Strategy{Attempts: 3},
		qCreate,
		comment.ParentID,
		comment.Content,
		comment.Author,
	)

	return err
}

func (r *postgresRepo) GetByID(ctx context.Context, id int64) (*models.Comment, error) {
	row, err := r.db.QueryRowWithRetry(
		ctx,
		retry.Strategy{Attempts: 3},
		qGetByID,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("r.db.QueryRowWithRetry: %w", err)
	}

	var out models.Comment
	if err := row.Scan(
		&out.ID,
		&out.ParentID,
		&out.Content,
		&out.Author,
		&out.CreatedAt,
		&out.UpdatedAt,
		&out.DeletedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("row.Scan: %w", err)
	}

	return &out, nil
}

func (r *postgresRepo) GetByParentID(ctx context.Context, parentID int64) ([]models.Comment, error) {
	rows, err := r.db.QueryWithRetry(
		ctx,
		retry.Strategy{Attempts: 3},
		qGetByParentID,
		parentID,
	)
	if err != nil {
		_ = rows.Close()
		return nil, fmt.Errorf("r.db.QueryWithRetry: %w", err)
	}
	defer rows.Close()

	var out []models.Comment

	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(
			&comment.ID,
			&comment.ParentID,
			&comment.Content,
			&comment.Author,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.DeletedAt,
			&comment.Level,
		); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		out = append(out, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}

	return out, nil
}

func (r *postgresRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecWithRetry(
		ctx,
		retry.Strategy{Attempts: 3},
		qDelete,
		id,
	)

	return err
}
