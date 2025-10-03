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

func TestWriteComment_ValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		comment *models.Comment
		wantErr string
	}{
		{
			name: "пустое содержимое",
			comment: &models.Comment{
				Content: "",
				Author:  "Тестер",
			},
			wantErr: "содержимое комментария не может быть пустым",
		},
		{
			name: "слишком длинное содержимое",
			comment: &models.Comment{
				Content: string(make([]byte, 1001)), // 1001 символ
				Author:  "Тестер",
			},
			wantErr: "содержимое комментария не может превышать 1000 символов",
		},
		{
			name: "без автора",
			comment: &models.Comment{
				Content: "Комментарий",
				Author:  "",
			},
			wantErr: "автор комментария не может быть пустым",
		},
		{
			name: "слишком длинный автор",
			comment: &models.Comment{
				Content: "Комментарий",
				Author:  string(make([]byte, 51)), // 51 символ
			},
			wantErr: "автор комментария не может превышать 50 символов",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewDatabase(t)
			svc := New(repo)
			ctx := context.Background()

			err := svc.WriteComment(ctx, tt.comment)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

// GetComments tests.
func TestGetComments_OK(t *testing.T) {
	repo := mocks.NewDatabase(t)
	svc := New(repo)

	ctx := context.Background()
	parentID := int64(1)

	expectedComments := []models.Comment{
		{
			ID:        1,
			ParentID:  &parentID,
			Content:   "Комментарий 1",
			Author:    "Автор 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
			Level:     1,
		},
		{
			ID:        2,
			ParentID:  &parentID,
			Content:   "Комментарий 2",
			Author:    "Автор 2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
			Level:     1,
		},
	}

	repo.EXPECT().GetByParentID(ctx, parentID).Return(expectedComments, nil)

	comments, err := svc.GetComments(ctx, parentID)

	assert.NoError(t, err)
	assert.Equal(t, expectedComments, comments)
}

func TestGetComments_InvalidParentID(t *testing.T) {
	repo := mocks.NewDatabase(t)
	svc := New(repo)
	ctx := context.Background()

	comments, err := svc.GetComments(ctx, 0)

	assert.Error(t, err)
	assert.Nil(t, comments)
	assert.Contains(t, err.Error(), "id родительского комментария должен быть больше 0")
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
