package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/child6yo/wbtech-l3-comment-tree/internal/models"
	"github.com/wb-go/wbf/dbpg"
)

// CommentsRepository отвечает за область хранения комментариев.
type CommentsRepository struct {
	db *dbpg.DB
}

// NewCommentsRepository создает новый экземпляр CommentsRepository.
func NewCommentsRepository(db *dbpg.DB) *CommentsRepository {
	return &CommentsRepository{db: db}
}

// CreateComment создает новый комментарий.
func (cr *CommentsRepository) CreateNewComment(ctx context.Context, comment models.CommentNode) (int, error) {
	query := fmt.Sprintf(`INSERT INTO %s (content) 
                                values ($1) RETURNING id`, tableComments)
	res := cr.db.QueryRowContext(ctx, query, comment.Content)
	if res.Err() != nil {
		return 0, fmt.Errorf("repository: %v", res.Err())
	}

	var id int
	err := res.Scan(&id)

	return int(id), err
}

// CreateAnswer создает ответ на существующий комментарий.
func (cr *CommentsRepository) CreateAnswer(ctx context.Context, comment models.CommentNode) (int, error) {
	query := fmt.Sprintf(`INSERT INTO %s (content, answer_at) 
                                values ($1, $2) RETURNING id`, tableComments)
	res := cr.db.QueryRowContext(ctx, query, comment.Content, comment.AnswerAt)
	if res.Err() != nil {
		return 0, fmt.Errorf("repository: %v", res.Err())
	}

	var id int
	err := res.Scan(&id)

	return int(id), err
}

// GetFullCommentTree возвращает полное дерево комментариев.
func (cr *CommentsRepository) GetFullCommentTree(ctx context.Context) ([]models.CommentNode, error) {
	query := `SELECT id, content, answer_at FROM comments ORDER BY id`
	rows, err := cr.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("repository: %v", err)
	}
	defer rows.Close()

	var rowsData []models.CommentNode

	for rows.Next() {
		var id int
		var content string
		var answerAt sql.NullInt64
		if err := rows.Scan(&id, &content, &answerAt); err != nil {
			return nil, fmt.Errorf("repository: %v", err)
		}
		var a int
		if answerAt.Valid {
			a = int(answerAt.Int64)
		}
		rowsData = append(rowsData, models.CommentNode{ID: id, Content: content, AnswerAt: a})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: %v", err)
	}

	return rowsData, nil
}
