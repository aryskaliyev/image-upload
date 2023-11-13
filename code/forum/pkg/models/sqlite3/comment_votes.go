package sqlite3

import (
	"database/sql"

	"lincoln.boris/forum/pkg/models"
)

type CommentVoteModel struct {
	DB *sql.DB
}

func (m *CommentVoteModel) Get(comment_id, user_id int) (*models.CommentVote, error) {
	stmt := `SELECT comment_id, user_id, vote FROM comment_vote
	WHERE comment_id = ? AND user_id = ?`

	row := m.DB.QueryRow(stmt, comment_id, user_id)

	cv := &models.CommentVote{}

	err := row.Scan(&cv.UserID, &cv.CommentID, &cv.Vote)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return cv, nil

}

func (m *CommentVoteModel) Insert(comment_id, user_id, vote int) (int, error) {
	stmt := `INSERT INTO comment_vote (comment_id, user_id, vote)
	VALUES(?, ?, ?)`

	tx, err := m.DB.Begin()
	if err != nil {
		return 0, err
	}

	result, err := tx.Exec(stmt, comment_id, user_id, vote)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	comment_vote_id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return int(comment_vote_id), nil
}

func (m *CommentVoteModel) Delete(comment_id, user_id int) (int, error) {
	stmt := `DELETE FROM comment_vote
	WHERE comment_id = ? AND user_id = ?`

	tx, err := m.DB.Begin()
	if err != nil {
		return 0, err
	}

	result, err := tx.Exec(stmt, comment_id, user_id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	cnt, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return int(cnt), nil
}

func (m *CommentVoteModel) Update(comment_id, user_id, vote int) (int, error) {
	stmt := `UPDATE comment_vote
	SET vote = ?
	WHERE comment_id = ? AND user_id = ?`

	tx, err := m.DB.Begin()
	if err != nil {
		return 0, err
	}

	result, err := tx.Exec(stmt, comment_id, user_id, vote)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	comment_vote_id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return int(comment_vote_id), nil
}

func (m *CommentVoteModel) Sum(comment_id int) (int, error) {
	stmt := `
	SELECT COALESCE(SUM(vote), 0)
	FROM comment_vote
	WHERE comment_id = ?`

	rows, err := m.DB.Query(stmt, comment_id)
	if err != nil {
		return 0, err
	}

	var sum int64

	for rows.Next() {
		err = rows.Scan(&sum)
		if err != nil {
			return 0, err
		}
	}

	return int(sum), nil
}
