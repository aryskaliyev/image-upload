package sqlite3

import (
	"database/sql"
	"strings"

	"lincoln.boris/forum/pkg/models"

	"github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(username, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `
	INSERT INTO useraccount (username, email, hashed_password, created)
	VALUES(?, ?, ?, CURRENT_TIMESTAMP)`

	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(stmt, username, email, string(hashedPassword))
	if err != nil {
		tx.Rollback()
		if err.(sqlite3.Error).Code == sqlite3.ErrConstraint && strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return models.ErrDuplicateEmail
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	stmt := `
	SELECT user_id, hashed_password
	FROM useraccount
	WHERE email = ?`

	var user_id int
	var hashedPassword []byte

	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&user_id, &hashedPassword)
	if err == sql.ErrNoRows {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	return user_id, nil
}
