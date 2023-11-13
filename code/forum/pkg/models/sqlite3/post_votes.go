package sqlite3

import (
	"database/sql"

	"lincoln.boris/forum/pkg/models"
)

type PostVoteModel struct {
	DB *sql.DB
}

func (m *PostVoteModel) Get(post_id, user_id int) (*models.PostVote, error) {
	stmt := `
	SELECT post_id, user_id, vote 
	FROM post_vote
	WHERE post_id = ? AND user_id = ?`

	row := m.DB.QueryRow(stmt, post_id, user_id)

	pv := &models.PostVote{}

	err := row.Scan(&pv.PostID, &pv.UserID, &pv.Vote)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return pv, nil

}

func (m *PostVoteModel) Insert(post_id, user_id, vote int) (int, error) {
	stmt := `
	INSERT INTO post_vote (post_id, user_id, vote)
	VALUES(?, ?, ?)`

	tx, err := m.DB.Begin()
	if err != nil {
		return 0, err
	}

	result, err := tx.Exec(stmt, post_id, user_id, vote)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	post_vote_id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return int(post_vote_id), nil
}

func (m *PostVoteModel) Delete(post_id, user_id int) (int, error) {
	stmt := `
	DELETE FROM post_vote
	WHERE post_id = ? AND user_id = ?`

	tx, err := m.DB.Begin()
	if err != nil {
		return 0, err
	}

	result, err := tx.Exec(stmt, post_id, user_id)
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

func (m *PostVoteModel) Update(post_id, user_id, vote int) (int, error) {
	stmt := `
	UPDATE post_vote
	SET vote = ?
	WHERE post_id = ? AND user_id = ?`

	tx, err := m.DB.Begin()
	if err != nil {
		return 0, err
	}

	result, err := tx.Exec(stmt, post_id, user_id, vote)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	post_vote_id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return int(post_vote_id), nil
}

func (m *PostVoteModel) Sum(post_id int) (int, error) {
	stmt := `
	SELECT COALESCE(SUM(vote), 0)
	FROM post_vote
	WHERE post_id = ?`

	rows, err := m.DB.Query(stmt, post_id)
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
