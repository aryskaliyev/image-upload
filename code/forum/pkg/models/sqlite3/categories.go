package sqlite3

import (
	"database/sql"

	"lincoln.boris/forum/pkg/models"
)

type CategoryModel struct {
	DB *sql.DB
}

func (m *CategoryModel) All() ([]*models.Category, error) {
	stmt := `SELECT category_id, name FROM category`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	categories := []*models.Category{}

	for rows.Next() {
		c := &models.Category{}

		err = rows.Scan(&c.ID, &c.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (m *CategoryModel) Get(id int) (*models.Category, error) {
	stmt := `SELECT category_id, name FROM category
	WHERE category_id = ?`

	row := m.DB.QueryRow(stmt, id)
	c := &models.Category{}

	err := row.Scan(&c.ID, &c.Name)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return c, nil
}
