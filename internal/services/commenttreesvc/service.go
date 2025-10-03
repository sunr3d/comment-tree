package commenttreesvc

import (
	"context"
	"fmt"

	"github.com/sunr3d/comment-tree/internal/interfaces/infra"
	"github.com/sunr3d/comment-tree/internal/interfaces/services"
	"github.com/sunr3d/comment-tree/models"
)

var _ services.CommentTree = (*commentTreeSvc)(nil)

type commentTreeSvc struct {
	repo infra.Database
}

func New(repo infra.Database) *commentTreeSvc {
	return &commentTreeSvc{repo: repo}
}

func (s *commentTreeSvc) WriteComment(ctx context.Context, comment *models.Comment) error {
	if comment.ParentID != nil {
		parent, err := s.repo.GetByID(ctx, *comment.ParentID)
		if err != nil {
			return fmt.Errorf("s.repo.GetByID: %w", err)
		}
		if parent == nil {
			return fmt.Errorf("родительский комментарий с id %d не найден", *comment.ParentID)
		}
		if parent.DeletedAt != nil {
			return fmt.Errorf("родительский комментарий с id %d уже удален", *comment.ParentID)
		}
	}

	return s.repo.Create(ctx, comment)
}

func (s *commentTreeSvc) GetComments(ctx context.Context, parentID int64) ([]models.Comment, error) {
	comment, err := s.repo.GetByID(ctx, parentID)
	if err != nil {
		return nil, fmt.Errorf("s.repo.GetByID: %w", err)
	}
	if comment == nil {
		return nil, fmt.Errorf("комментарий с id %d не найден", parentID)
	}
	if comment.DeletedAt != nil {
		return nil, fmt.Errorf("комментарий с id %d уже удален", parentID)
	}

	return s.repo.GetByParentID(ctx, parentID)
}

func (s *commentTreeSvc) DeleteComment(ctx context.Context, id int64) error {
	comment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("s.repo.GetByID: %w", err)
	}
	if comment == nil {
		return fmt.Errorf("комментарий с id %d не найден", id)
	}

	if comment.DeletedAt != nil {
		return fmt.Errorf("комментарий с id %d уже удален", id)
	}

	return s.repo.Delete(ctx, id)
}
