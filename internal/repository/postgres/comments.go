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

// CreateNewComment создает новый комментарий.
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
	query := fmt.Sprintf(`SELECT id, content, answer_at FROM %s ORDER BY id`, tableComments)
	rows, err := cr.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("repository: %v", err)
	}
	defer func() {
		_ = rows.Close()
	}()
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: %v", err)
	}

	var rowsData []models.CommentNode

	err = nodesListByRows(rows, &rowsData)

	return rowsData, err
}

// GetTreeByParent возвращает только дочерние комментарии по родительскому айди.
func (cr *CommentsRepository) GetTreeByParent(ctx context.Context, parentID int) ([]models.CommentNode, error) {
	var rowsData []models.CommentNode

	query := fmt.Sprintf(`SELECT id, content FROM %s WHERE id=$1`, tableComments)
	row := cr.db.QueryRowContext(ctx, query, parentID)
	if row.Err() != nil {
		return nil, fmt.Errorf("repository: %v", row.Err())
	}
	var id int
	var content string
	err := row.Scan(&id, &content)
	if err != nil {
		return nil, fmt.Errorf("repository: %v", err)
	}
	rowsData = append(rowsData, models.CommentNode{ID: id, Content: content})

	query = fmt.Sprintf(`WITH RECURSIVE children AS (
		SELECT id, content, answer_at
		FROM %s
		WHERE answer_at=$1 

		UNION ALL

		SELECT c.id, c.content, c.answer_at
		FROM comments c
		INNER JOIN children ch ON c.answer_at = ch.id
	)
	SELECT id, content, answer_at
	FROM children
	ORDER BY id;`, tableComments)
	rows, err := cr.db.QueryContext(ctx, query, parentID)
	if err != nil {
		return nil, fmt.Errorf("repository: %v", err)
	}
	defer func() {
		_ = rows.Close()
	}()
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: %v", err)
	}

	err = nodesListByRows(rows, &rowsData)

	return rowsData, err
}

// DeleteTreeByParent рекурсивно удаляет все дочерние комментарии, начиная с родительского.
func (cr *CommentsRepository) DeleteTreeByParent(ctx context.Context, parentID int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, tableComments)
	_, err := cr.db.ExecContext(ctx, query, parentID)
	if err != nil {
		return fmt.Errorf("repository: %v", err)
	}
	return nil
}

func nodesListByRows(rows *sql.Rows, rowsData *[]models.CommentNode) error {
	for rows.Next() {
		var id int
		var content string
		var answerAt sql.NullInt64
		if err := rows.Scan(&id, &content, &answerAt); err != nil {
			return fmt.Errorf("repository: %v", err)
		}
		var a int
		if answerAt.Valid {
			a = int(answerAt.Int64)
		}
		*rowsData = append(*rowsData, models.CommentNode{ID: id, Content: content, AnswerAt: a})
	}

	return nil
}
