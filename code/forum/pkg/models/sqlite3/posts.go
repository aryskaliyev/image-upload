package sqlite3

import (
	"database/sql"

	"lincoln.boris/forum/pkg/models"
)

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) Insert(user_id int, title, body string, categories []int, img string) (int, error) {
	post_stmt := `
	INSERT INTO post (user_id, title, body, image, created)
	VALUES(?, ?, ?, ?, CURRENT_TIMESTAMP)`

	post_category_stmt := `
	INSERT INTO post_category (post_id, category_id)
	VALUES(?, ?)`

	tx, err := m.DB.Begin()
	if err != nil {
		return 0, err
	}

	result, err := tx.Exec(post_stmt, user_id, title, body, img)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	post_id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	for _, category := range categories {
		_, err = tx.Exec(post_category_stmt, post_id, category)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return int(post_id), nil
}

func (m *PostModel) Get(post_id int) (*models.Post, error) {
	post_stmt := `
	SELECT post.post_id, post.title, post.body, post.image, post.created, useraccount.username FROM post
	INNER JOIN useraccount ON post.user_id = useraccount.user_id
	WHERE post_id = ?`

	p := &models.Post{}

	row := m.DB.QueryRow(post_stmt, post_id)
	err := row.Scan(&p.ID, &p.Title, &p.Body, &p.Image, &p.Created, &p.Author)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return p, nil
}

func (m *PostModel) AllPosts(user_id, upvoted int) ([]*models.Post, error) {
	var stmt string
	var args []interface{}

	if user_id > 0 && upvoted == 0 {
		stmt = `
		SELECT post.post_id, post.title, post.created, useraccount.username
		FROM post
		INNER JOIN useraccount ON post.user_id = useraccount.user_id
		WHERE post.user_id = ?
		ORDER BY post.created DESC`

		args = append(args, user_id)
	} else if user_id > 0 && upvoted > 0 {
		stmt = `
		SELECT post.post_id, post.title, post.created, useraccount.username
		FROM post
		INNER JOIN post_vote ON post.post_id = post_vote.post_id
		INNER JOIN useraccount ON post.user_id = useraccount.user_id
		WHERE post_vote.user_id = ? AND post_vote.vote = ?`

		args = append(args, user_id)
		args = append(args, upvoted)
	} else {
		stmt = `
		SELECT post.post_id, post.title, post.created, useraccount.username
		FROM post
		INNER JOIN useraccount ON post.user_id = useraccount.user_id
		ORDER BY post.created DESC`
	}

	rows, err := m.DB.Query(stmt, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	posts := []*models.Post{}

	for rows.Next() {
		p := &models.Post{}

		err = rows.Scan(&p.ID, &p.Title, &p.Created, &p.Author)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (m *PostModel) ByCategory(category_id int) ([]*models.Post, error) {
	stmt := `
	SELECT post.post_id, post.title, post.created, useraccount.username FROM post
	INNER JOIN post_category ON post_category.post_id = post.post_id
	INNER JOIN category ON category.category_id = post_category.category_id
	INNER JOIN useraccount ON post.user_id = useraccount.user_id
	WHERE category.category_id = ?
	ORDER BY post.created DESC`

	rows, err := m.DB.Query(stmt, category_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	posts := []*models.Post{}

	for rows.Next() {
		p := &models.Post{}

		err = rows.Scan(&p.ID, &p.Title, &p.Created, &p.Author)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
