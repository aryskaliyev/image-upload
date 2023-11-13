package sqlite3

import (
	"database/sql"

	"lincoln.boris/forum/pkg/models"
)

type CommentModel struct {
	DB *sql.DB
}

func (m *CommentModel) Insert(body string, post_id, user_id int) (int, error) {
	stmt := `
	INSERT INTO comment (body, post_id, user_id, created)
	VALUES(?, ?, ?, CURRENT_TIMESTAMP)`

	tx, err := m.DB.Begin()
	if err != nil {
		return 0, err
	}

	result, err := tx.Exec(stmt, body, post_id, user_id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	comment_id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return int(comment_id), nil
}

func (m *CommentModel) Get(post_id int) ([]*models.Comment, error) {
	stmt := `
	SELECT comment.comment_id, comment.post_id, comment.body, comment.user_id, comment.created, useraccount.username
	FROM comment
	INNER JOIN useraccount ON comment.user_id = useraccount.user_id
	WHERE comment.post_id = ?
	ORDER BY comment.created DESC`

	rows, err := m.DB.Query(stmt, post_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	comments := []*models.Comment{}

	for rows.Next() {
		c := &models.Comment{}

		err := rows.Scan(&c.ID, &c.PostID, &c.Body, &c.UserID, &c.Created, &c.Username)
		if err == sql.ErrNoRows {
			return nil, models.ErrNoRecord
		} else if err != nil {
			return nil, err
		}

		comments = append(comments, c)
	}

	return comments, nil
}
