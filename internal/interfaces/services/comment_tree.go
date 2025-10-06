package services

import (
	"context"

	"github.com/sunr3d/comment-tree/models"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.2 --name=CommentTree --output=../../../mocks --filename=mock_comment_tree.go --with-expecter
type CommentTree interface {
	WriteComment(ctx context.Context, comment *models.Comment) error
	GetComments(ctx context.Context, parentID int64, pag *models.PagParam) (*models.CommentsRes, error)
	DeleteComment(ctx context.Context, id int64) error
}
