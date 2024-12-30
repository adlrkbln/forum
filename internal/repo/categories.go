package repo

import (
	"fmt"
	"forum/internal/models"
)

func (sq *Sqlite) GetCategories() ([]*models.Category, error) {
	var categories []*models.Category

	rows, err := sq.DB.Query("SELECT id, name FROM category")
	if err != nil {
		return nil, fmt.Errorf("repo.GetCategories %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		category := &models.Category{}
		if err := rows.Scan(&category.Id, &category.Name); err != nil {
			return nil, fmt.Errorf("repo.GetCategories %w", err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repo.GetCategories %w", err)
	}
	return categories, nil
}

func (sq *Sqlite) GetPostByCategory(id int) ([]*models.Post, error) {
	stmt := `SELECT p.id, p.user_id, p.title, p.content, p.created FROM posts p
	JOIN post_category pc ON p.id = pc.post_id 
	JOIN category c ON pc.category_id = c.id
	WHERE c.id = ?
    ORDER BY p.id`

	rows, err := sq.DB.Query(stmt, id)
	if err != nil {
		return nil, fmt.Errorf("repo.GetPostByCategory: %w", err)
	}

	defer rows.Close()

	posts := []*models.Post{}

	for rows.Next() {
		s := &models.Post{}
		err = rows.Scan(&s.Id, &s.UserId, &s.Title, &s.Content, &s.Created)
		if err != nil {
			return nil, fmt.Errorf("repo.GetPostByCategory: %w", err)
		}
		posts = append(posts, s)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("repo.GetPostByCategory: %w", err)
	}

	return posts, nil
}

func (sq *Sqlite) PostCategoryPost(post_id int, category_id int) error {
	stmt := `INSERT INTO post_category (category_id, post_id)
	VALUES (?, ?);`

	_, err := sq.DB.Exec(stmt, category_id, post_id)
	if err != nil {
		return fmt.Errorf("repo.PostCategoryPost: %w", err)
	}

	return nil
}

func (sq *Sqlite) CreateCategory(name string) error {
	stmt := `INSERT INTO category (name) VALUES (?);`
	_, err := sq.DB.Exec(stmt, name)
	if err != nil {
		return fmt.Errorf("repo.CreateCategory: %w", err)
	}
	return nil
}

func (sq *Sqlite) DeleteCategory(id int) error {
	stmt := `DELETE FROM category WHERE id = ?;`
	_, err := sq.DB.Exec(stmt, id)
	if err != nil {
		return fmt.Errorf("repo.DeleteCategory: %w", err)
	}
	return nil
}
