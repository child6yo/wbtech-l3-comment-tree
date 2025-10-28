package usecase

import (
	"context"

	"github.com/child6yo/wbtech-l3-comment-tree/internal/models"
)

type commentsRepository interface {
	CreateNewComment(ctx context.Context, comment models.CommentNode) (int, error)
	CreateAnswer(ctx context.Context, comment models.CommentNode) (int, error)
	GetFullCommentTree(ctx context.Context) ([]models.CommentNode, error)
	GetTreeByParent(ctx context.Context, parentID int) ([]models.CommentNode, error)
	DeleteTreeByParent(ctx context.Context, parentID int) error
}

// CommentsService создает, удаляет и возвращает комментарии и их деревья.
type CommentsService struct {
	repo commentsRepository
}

// NewCommentsService создает новый CommentsService.
func NewCommentsService(repo commentsRepository) *CommentsService {
	return &CommentsService{repo: repo}
}

// CreateComment создает новый комментарий. Возвращает его айди.
func (cs *CommentsService) CreateComment(ctx context.Context, comment models.CommentNode) (int, error) {
	if comment.AnswerAt == 0 {
		return cs.repo.CreateNewComment(ctx, comment)
	}
	return cs.repo.CreateAnswer(ctx, comment)
}

// GetCommentTree возвращает дерево комментариев по родительскому комментарию.
// Если id = 0, возвращает полное дерево.
func (cs *CommentsService) GetCommentTree(ctx context.Context, id int) ([]*models.CommentNode, error) {
	var coms []models.CommentNode
	var err error

	if id == 0 {
		coms, err = cs.repo.GetFullCommentTree(ctx)
	} else {
		coms, err = cs.repo.GetTreeByParent(ctx, id)
	}

	if err != nil {
		return nil, err
	}

	return createTree(coms), nil
}

// DeleteTreeByParent удаляет все дочерние комментарии, начиная с родительского.
func (cs *CommentsService) DeleteTreeByParent(ctx context.Context, id int) error {
	return cs.repo.DeleteTreeByParent(ctx, id)
}

func createTree(coms []models.CommentNode) []*models.CommentNode {
	// собираем все узлы в мапу по ID
	nodes := make(map[int]*models.CommentNode)
	for _, c := range coms {
		nodes[c.ID] = &models.CommentNode{
			ID:       c.ID,
			Content:  c.Content,
			AnswerAt: c.AnswerAt,
			Children: nil,
		}
	}

	// собираем корни и строим связи
	var roots []*models.CommentNode
	for _, c := range coms {
		node := nodes[c.ID]
		parentID := c.AnswerAt

		if parentID == 0 || parentID == -1 || nodes[parentID] == nil {
			roots = append(roots, node)
		} else {
			parent := nodes[parentID]
			parent.Children = append(parent.Children, node)
		}
	}

	return roots
}
