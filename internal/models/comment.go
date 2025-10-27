package models

// CommentNode модель узла дерева комментариев.
type CommentNode struct {
	ID       int           `json:"id"`
	Content  string        `json:"content"`
	AnswerAt int           `json:"answer_at"`
	Children []*CommentNode `json:"children"`
}
