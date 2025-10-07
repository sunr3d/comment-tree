package infra

import (
	"context"

	"github.com/sunr3d/comment-tree/models"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.2 --name=Database --output=../../../mocks --filename=mock_database.go --with-expecter
type Database interface {
	Create(ctx context.Context, comment *models.Comment) error
	GetByID(ctx context.Context, id int64) (*models.Comment, error)
	GetByParentID(ctx context.Context, parentID int64, pag *models.PagParam) (*models.CommentsRes, error)
	GetRootComments(ctx context.Context, pag *models.PagParam) (*models.CommentsRes, error)
	Delete(ctx context.Context, id int64) error
}
