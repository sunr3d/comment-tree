package commenttreesvc

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/sunr3d/comment-tree/mocks"
	"github.com/sunr3d/comment-tree/models"
)

// WriteComment tests.
func TestWriteComment_OK(t *testing.T) {
	repo := mocks.NewDatabase(t)
	svc := New(repo)

	ctx := context.Background()
	comment := &models.Comment{
		ParentID: nil,
		Content:  "Тестовый комментарий",
		Author:   "Тестер",
	}

	repo.EXPECT().
		Create(ctx, comment).
		Return(nil)

	err := svc.WriteComment(ctx, comment)

	assert.NoError(t, err)
}

func TestWriteComment_WithParentID_OK(t *testing.T) {
	repo := mocks.NewDatabase(t)
	svc := New(repo)

	ctx := context.Background()
	parentID := int64(1)
	comment := &models.Comment{
		ParentID: &parentID,
		Content:  "Ответ на комментарий",
		Author:   "Тестер",
	}

	parentComment := &models.Comment{
		ID:        parentID,
		ParentID:  nil,
		Content:   "Тестовый родительский комментарий",
		Author:    "Тестер родительский",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
		Level:     0,
	}

	repo.EXPECT().
		GetByID(ctx, parentID).
		Return(parentComment, nil)

	repo.EXPECT().
		Create(ctx, comment).
		Return(nil)

	err := svc.WriteComment(ctx, comment)

	assert.NoError(t, err)
}

func TestWriteComment_WithParentID_NotFound(t *testing.T) {
	repo := mocks.NewDatabase(t)
	svc := New(repo)

	ctx := context.Background()
	parentID := int64(42)
	comment := &models.Comment{
		ParentID: &parentID,
		Content:  "Ответ на несуществующий комментарий",
		Author:   "Тестер",
	}

	repo.EXPECT().
		GetByID(ctx, parentID).
		Return(nil, nil)

	err := svc.WriteComment(ctx, comment)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "родительский комментарий с id 42 не найден")
}

func TestWriteComment_WithParentID_Deleted(t *testing.T) {
	repo := mocks.NewDatabase(t)
	svc := New(repo)

	ctx := context.Background()
	parentID := int64(1)
	comment := &models.Comment{
		ParentID: &parentID,
		Content:  "Ответ на удаленный комментарий",
		Author:   "Тестер",
	}

	now := time.Now()
	parentComment := &models.Comment{
		ID:        parentID,
		ParentID:  nil,
		Content:   "Тестовый родительский комментарий",
		Author:    "Тестер родительский",
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: &now,
		Level:     0,
	}

	repo.EXPECT().
		GetByID(ctx, parentID).
		Return(parentComment, nil)

	err := svc.WriteComment(ctx, comment)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "родительский комментарий с id 1 уже удален")
}

// GetComments tests.
func TestGetComments_OK(t *testing.T) {
	repo := mocks.NewDatabase(t)
	svc := New(repo)

	ctx := context.Background()
	parentID := int64(1)
	pag := &models.PagParam{
		Page:   1,
		Limit:  20,
		Sort:   "created_at_desc",
		Search: "",
	}

	parentComment := &models.Comment{
		ID:        parentID,
		ParentID:  nil,
		Content:   "Комментарий 1",
		Author:    "Автор 1",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
		Level:     0,
	}

	expectedResult := &models.CommentsRes{
		Comments: []models.Comment{
			{
				ID:        2,
				ParentID:  &parentID,
				Content:   "Ответ 1",
				Author:    "Автор 2",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				DeletedAt: nil,
				Level:     1,
			},
		},
		Total: 1,
		Page:  1,
		Limit: 20,
		Pages: 1,
	}

	repo.EXPECT().GetByID(ctx, parentID).Return(parentComment, nil)
	repo.EXPECT().GetByParentID(ctx, parentID, pag).Return(expectedResult, nil)

	result, err := svc.GetComments(ctx, parentID, pag)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestGetComments_WithNilPagination(t *testing.T) {
	repo := mocks.NewDatabase(t)
	svc := New(repo)

	ctx := context.Background()
	parentID := int64(1)

	parentComment := &models.Comment{
		ID:        parentID,
		ParentID:  nil,
		Content:   "Родительский комментарий",
		Author:    "Автор",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
		Level:     0,
	}

	expectedResult := &models.CommentsRes{
		Comments: []models.Comment{},
		Total:    0,
		Page:     1,
		Limit:    20,
		Pages:    1,
	}

	expectedPag := &models.PagParam{
		Page:   1,
		Limit:  20,
		Sort:   "created_at_asc",
		Search: "",
	}

	repo.EXPECT().GetByID(ctx, parentID).Return(parentComment, nil)
	repo.EXPECT().GetByParentID(ctx, parentID, expectedPag).Return(expectedResult, nil)

	result, err := svc.GetComments(ctx, parentID, nil)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestGetComments_ParentDeleted(t *testing.T) {
	repo := mocks.NewDatabase(t)
	svc := New(repo)

	ctx := context.Background()
	parentID := int64(1)
	pag := &models.PagParam{
		Page:   1,
		Limit:  20,
		Sort:   "created_at_desc",
		Search: "",
	}

	now := time.Now()
	parentComment := &models.Comment{
		ID:        parentID,
		ParentID:  nil,
		Content:   "Удаленный комментарий",
		Author:    "Автор",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: &now,
		Level:     0,
	}

	expectedResult := &models.CommentsRes{
		Comments: []models.Comment{},
		Total:    0,
		Page:     1,
		Limit:    20,
		Pages:    1,
	}

	repo.EXPECT().GetByID(ctx, parentID).Return(parentComment, nil)
	repo.EXPECT().GetByParentID(ctx, parentID, pag).Return(expectedResult, nil)

	result, err := svc.GetComments(ctx, parentID, pag)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}

// DeleteComment tests.
func TestDeleteComment_OK(t *testing.T) {
	repo := mocks.NewDatabase(t)
	svc := New(repo)

	ctx := context.Background()
	commentID := int64(1)

	comment := &models.Comment{
		ID:        commentID,
		ParentID:  nil,
		Content:   "Комментарий для удаления",
		Author:    "Автор",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
		Level:     0,
	}

	repo.EXPECT().GetByID(ctx, commentID).Return(comment, nil)
	repo.EXPECT().Delete(ctx, commentID).Return(nil)

	err := svc.DeleteComment(ctx, commentID)

	assert.NoError(t, err)
}

func TestDeleteComment_NotFound(t *testing.T) {
	repo := mocks.NewDatabase(t)
	svc := New(repo)

	ctx := context.Background()
	commentID := int64(42)

	repo.EXPECT().GetByID(ctx, commentID).Return(nil, nil)

	err := svc.DeleteComment(ctx, commentID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "комментарий с id 42 не найден")
}

func TestDeleteComment_AlreadyDeleted(t *testing.T) {
	repo := mocks.NewDatabase(t)
	svc := New(repo)

	ctx := context.Background()
	commentID := int64(1)
	now := time.Now()

	comment := &models.Comment{
		ID:        commentID,
		ParentID:  nil,
		Content:   "Уже удаленный комментарий",
		Author:    "Автор",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: &now,
		Level:     0,
	}

	repo.EXPECT().GetByID(ctx, commentID).Return(comment, nil)

	err := svc.DeleteComment(ctx, commentID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "комментарий с id 1 уже удален")
}
