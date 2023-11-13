package sqlite3

import (
	"database/sql"

	"lincoln.boris/forum/pkg/models"
)

type PostCategoryModel struct {
	DB *sql.DB
}

func (m *PostCategoryModel) Get(post_id int) ([]*models.PostCategory, error) {
	stmt := `
	SELECT post_category.category_id, category.name
	FROM post_category
	INNER JOIN category ON post_category.category_id = category.category_id
	WHERE post_category.post_id = ?`

	rows, err := m.DB.Query(stmt, post_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	post_categories := []*models.PostCategory{}

	for rows.Next() {
		pc := &models.PostCategory{}

		err := rows.Scan(&pc.CategoryID, &pc.Name)
		if err == sql.ErrNoRows {
			return nil, models.ErrNoRecord
		} else if err != nil {
			return nil, err
		}

		post_categories = append(post_categories, pc)
	}

	return post_categories, nil
}
