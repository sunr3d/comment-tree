package httphandlers

import (
	"net/http"
	"strings"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"

	"github.com/sunr3d/comment-tree/models"
)

func (h *Handler) getRootComments(c *ginext.Context, req *getCommentsReq) {
	pag := h.buildPagination(c, req)

	result, err := h.svc.GetRootComments(c.Request.Context(), pag)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("svc.GetRootComments")
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "внутренняя ошибка сервера"})
		return
	}

	h.sendCommentsResp(c, result)
}

func (h *Handler) getCommentsByParent(c *ginext.Context, req *getCommentsReq) {
	pag := h.buildPagination(c, req)

	result, err := h.svc.GetComments(c.Request.Context(), *req.ParentID, pag)
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

	h.sendCommentsResp(c, result)
}

func (h *Handler) buildPagination(c *ginext.Context, req *getCommentsReq) *models.PagParam {
	if c.Query("page") == "" && c.Query("limit") == "" && c.Query("sort") == "" && c.Query("search") == "" {
		return nil
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Sort == "" {
		req.Sort = "created_at_asc"
	}

	return &models.PagParam{
		Page:   req.Page,
		Limit:  req.Limit,
		Sort:   req.Sort,
		Search: req.Search,
	}
}

func (h *Handler) sendCommentsResp(c *ginext.Context, result *models.CommentsRes) {
	commentsDTO := make([]comment, len(result.Comments))
	for i, c := range result.Comments {
		commentsDTO[i] = comment{
			ID:        c.ID,
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
