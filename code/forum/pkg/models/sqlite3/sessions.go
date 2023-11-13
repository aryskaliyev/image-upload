package sqlite3

import (
	"database/sql"
	"time"
)

type SessionModel struct {
	DB *sql.DB
}

func (m *SessionModel) Insert(user_id int, token string, expires time.Time) error {
	stmt := `
	INSERT INTO session (user_id, uuid_token, created, expires)
	VALUES(?, ?, CURRENT_TIMESTAMP, ?)`

	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(stmt, user_id, token, expires.Format("2006-01-02 15:04:05"))
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (m *SessionModel) GetUser(token string) int {
	stmt := `
	SELECT user_id
	FROM session
	WHERE uuid_token = ?`

	row := m.DB.QueryRow(stmt, token)

	var user_id int

	err := row.Scan(&user_id)
	if err != nil {
		return 0
	}

	return user_id
}

func (m *SessionModel) GetToken(user_id int) string {
	stmt := `
	SELECT uuid_token
	FROM session
	WHERE user_id = ?`

	row := m.DB.QueryRow(stmt, user_id)

	var token string

	err := row.Scan(&token)
	if err != nil {
		return ""
	}

	return token
}

func (m *SessionModel) Update(user_id int) error {
	stmt := `
	UPDATE session
	SET expires = DATETIME(CURRENT_TIMESTAMP, '+5 minutes')
	WHERE user_id = ?`

	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(stmt, user_id)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (m *SessionModel) Delete(user_id int) error {
	stmt := `
	DELETE FROM session
	WHERE user_id = ?`

	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(stmt, user_id)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
