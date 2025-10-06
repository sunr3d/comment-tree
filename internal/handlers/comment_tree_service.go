package httphandlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"

	"github.com/sunr3d/comment-tree/models"
)

func (h *Handler) writeComment(c *ginext.Context) {
	var req createCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "некорректный JSON"})
		return
	}

	if req.Content == "" {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "комментарий не может быть пустым"})
		return
	}

	if req.Author == "" {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "автор не может быть пустым"})
		return
	}

	if len(req.Content) > 1000 {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "комментарий не может быть длиннее 1000 символов"})
		return
	}

	if len(req.Author) > 50 {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "автор не может быть длиннее 50 символов"})
		return
	}

	comment := &models.Comment{
		ParentID: req.ParentID,
		Content:  req.Content,
		Author:   req.Author,
	}

	if err := h.svc.WriteComment(c.Request.Context(), comment); err != nil {
		if strings.Contains(err.Error(), "не найден") || strings.Contains(err.Error(), "уже удален") {
			zlog.Logger.Error().Err(err).Msg("svc.WriteComment")
			c.JSON(http.StatusNotFound, ginext.H{"error": "комментарий не найден"})
			return
		}
		zlog.Logger.Error().Err(err).Msg("svc.WriteComment")
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "внутренняя ошибка сервера"})
		return
	}

	c.JSON(http.StatusOK, ginext.H{"message": "комментарий успешно создан"})
}

func (h *Handler) getComments(c *ginext.Context) {
	var req getCommentsReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "некорректный запрос"})
		return
	}
	if req.ParentID < 1 {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "id родительского комментария должен быть больше 0"})
		return
	}

	var pag *models.PagParam
	if c.Query("page") != "" || c.Query("limit") != "" || c.Query("sort") != "" {
		if req.Page <= 0 {
			req.Page = 1
		}
		if req.Limit <= 0 {
			req.Limit = 20
		}
		if req.Sort == "" {
			req.Sort = "created_at_asc"
		}
		pag = &models.PagParam{
			Page:  req.Page,
			Limit: req.Limit,
			Sort:  req.Sort,
		}
	}

	result, err := h.svc.GetComments(c.Request.Context(), req.ParentID, pag)
	if err != nil {
		if strings.Contains(err.Error(), "не найден") || strings.Contains(err.Error(), "уже удален") {
			zlog.Logger.Error().Err(err).Msg("svc.GetComments")
			c.JSON(http.StatusNotFound, ginext.H{"error": "комментарий не найден"})
			return
		}
		zlog.Logger.Error().Err(err).Msg("svc.GetComments")
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "внутренняя ошибка сервера"})
		return
	}

	commentsDTO := make([]comment, len(result.Comments))
	for i, c := range result.Comments {
		commentsDTO[i] = comment{
			Content:   c.Content,
			Author:    c.Author,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			DeletedAt: c.DeletedAt,
			Level:     c.Level,
		}
	}

	out := getCommentsResp{
		Comments: commentsDTO,
		Total:    result.Total,
		Page:     result.Page,
		Limit:    result.Limit,
		Pages:    result.Pages,
	}

	c.JSON(http.StatusOK, out)

}

func (h *Handler) deleteComment(c *ginext.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "некорректный id комментария"})
		return
	}
	if id < 1 {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "id комментария должен быть больше 0"})
		return
	}

	if err := h.svc.DeleteComment(c.Request.Context(), id); err != nil {
		if strings.Contains(err.Error(), "не найден") {
			zlog.Logger.Error().Err(err).Msg("svc.DeleteComment")
			c.JSON(http.StatusNotFound, ginext.H{"error": "комментарий не найден"})
			return
		}
		zlog.Logger.Error().Err(err).Msg("svc.DeleteComment")
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "внутренняя ошибка сервера"})
		return
	}

	c.JSON(http.StatusOK, ginext.H{"message": "комментарий успешно удален"})
}
