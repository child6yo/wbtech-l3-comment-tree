package controller

import (
	"context"
	"fmt"
	"strconv"

	"github.com/child6yo/wbtech-l3-comment-tree/internal/models"
	"github.com/wb-go/wbf/ginext"
)

type commentsService interface {
	CreateComment(ctx context.Context, comment models.CommentNode) (int, error)
	GetCommentTree(ctx context.Context, id int) ([]*models.CommentNode, error)
	DeleteTreeByParent(ctx context.Context, id int) error
}

// CommentsController http контроллер комментариев.
type CommentsController struct {
	cs commentsService
}

// NewCommentsController создает новый CommentsController.
func NewCommentsController(cs commentsService) *CommentsController {
	return &CommentsController{cs: cs}
}

type createCommentRequest struct {
	Content string `json:"content" binding:"required"`
	ID      int    `json:"id" binding:"omitempty"`
}

// NewComment обрабатывает POST /comments — создание комментария (с указанием родительского).
func (cc *CommentsController) NewComment(c *ginext.Context) {
	var req createCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, ginext.H{"error": "invalid request: " + err.Error()})
		_ = c.Error(fmt.Errorf("validation error: %v", err))
		return
	}

	c.Set("request", req)

	var comment models.CommentNode
	comment.Content = req.Content
	comment.AnswerAt = req.ID

	id, err := cc.cs.CreateComment(c.Request.Context(), comment)
	if err != nil {
		c.JSON(500, ginext.H{"error": "server error: " + err.Error()})
		_ = c.Error(fmt.Errorf("server error: %v", err))
		return
	}

	c.JSON(201, ginext.H{"id": id})
}

// GetComments обрабатывает GET /comments?parent={id} — получение комментария и всех вложенных.
func (cc *CommentsController) GetComments(c *ginext.Context) {
	parent := c.Query("parent")

	var id int
	var err error
	if parent != "" {
		id, err = strconv.Atoi(parent)
		if err != nil {
			c.JSON(400, ginext.H{"error": "invalid request: " + err.Error()})
			_ = c.Error(fmt.Errorf("validation error: %v", err))
			return
		}
	}

	coms, err := cc.cs.GetCommentTree(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, ginext.H{"error": "server error: " + err.Error()})
		_ = c.Error(fmt.Errorf("server error: %v", err))
		return
	}

	c.JSON(200, coms)
}

// DeleteComments обрабатывает DELETE /comments/{id} — удаление комментария и всех вложенных под ним.
func (cc *CommentsController) DeleteComments(c *ginext.Context) {
	parent := c.Param("id")

	id, err := strconv.Atoi(parent)
	if err != nil {
		c.JSON(400, ginext.H{"error": "invalid request: " + err.Error()})
		_ = c.Error(fmt.Errorf("validation error: %v", err))
		return
	}

	err = cc.cs.DeleteTreeByParent(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, ginext.H{"error": "server error: " + err.Error()})
		_ = c.Error(fmt.Errorf("server error: %v", err))
		return
	}

	c.JSON(200, nil)
}
