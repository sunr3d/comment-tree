package commenttreesvc

import (
	"fmt"
	"strings"

	"github.com/sunr3d/comment-tree/models"
)

func validateComment(comment *models.Comment) error {
	if strings.TrimSpace(comment.Content) == "" {
		return fmt.Errorf("содержимое комментария не может быть пустым")
	}

	if len(comment.Content) > 1000 {
		return fmt.Errorf("содержимое комментария не может превышать 1000 символов")
	}

	if strings.TrimSpace(comment.Author) == "" {
		return fmt.Errorf("автор комментария не может быть пустым")
	}

	if len(comment.Author) > 50 {
		return fmt.Errorf("автор комментария не может превышать 50 символов")
	}

	return nil
}
