package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/models"
)

func (sq *Sqlite) InsertPost(user_id int, title string, content string, image_path string) (int, error) {
	stmt := `INSERT INTO posts (user_id, title, content, image_path, created)
	VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP);`

	result, err := sq.DB.Exec(stmt, user_id, title, content, image_path)
	if err != nil {
		return 0, fmt.Errorf("repo.InsertPost: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("repo.InsertPost: %w", err)
	}

	return int(id), nil
}

func (sq *Sqlite) GetPost(id int) (*models.Post, error) {
	stmt := `SELECT p.id, p.user_id, p.title, p.content, p.image_path, p.likes, p.dislikes, p.created, u.name FROM posts p
	JOIN users u ON p.user_id = u.id
	WHERE p.id = ?`

	row := sq.DB.QueryRow(stmt, id)

	s := &models.Post{}
	err := row.Scan(&s.Id, &s.UserId, &s.Title, &s.Content, &s.ImagePath, &s.Likes, &s.Dislikes, &s.Created, &s.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, fmt.Errorf("repo.GetPost: %w", err)
		}
	}

	return s, nil
}

func (sq *Sqlite) GetAllPosts() ([]*models.Post, error) {
	stmt := `SELECT id, user_id, title, content, likes, dislikes, created FROM posts
    ORDER BY id DESC`

	rows, err := sq.DB.Query(stmt)
	if err != nil {
		return nil, fmt.Errorf("repo.GetAllPosts: %w", err)
	}

	defer rows.Close()

	posts := []*models.Post{}

	for rows.Next() {
		s := &models.Post{}
		err = rows.Scan(&s.Id, &s.UserId, &s.Title, &s.Content, &s.Likes, &s.Dislikes, &s.Created)
		if err != nil {
			return nil, fmt.Errorf("repo.GetAllPosts: %w", err)
		}
		posts = append(posts, s)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("repo.GetAllPosts: %w", err)
	}

	return posts, nil
}

func (sq *Sqlite) GetCreatedPosts(user_id int) ([]*models.Post, error) {
	stmt := `SELECT id, user_id, title, content, likes, dislikes, created FROM posts
	WHERE user_id = ?
    ORDER BY id DESC`

	rows, err := sq.DB.Query(stmt, user_id)
	if err != nil {
		return nil, fmt.Errorf("repo.GetCreatedPosts: %w", err)
	}

	defer rows.Close()

	posts := []*models.Post{}

	for rows.Next() {
		s := &models.Post{}
		err = rows.Scan(&s.Id, &s.UserId, &s.Title, &s.Content, &s.Likes, &s.Dislikes, &s.Created)
		if err != nil {
			return nil, fmt.Errorf("repo.GetCreatedPosts: %w", err)
		}
		posts = append(posts, s)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("repo.GetCreatedPosts: %w", err)
	}

	return posts, nil
}

func (sq *Sqlite) GetLikedPosts(user_id int) ([]*models.Post, error) {
	stmt := `SELECT p.id, p.user_id, p.title, p.content, p.likes, p.dislikes, p.created FROM posts p
	JOIN user_post_reactions upr ON p.id = upr.post_id
	WHERE upr.user_id = ? AND upr.reaction = 1
    ORDER BY p.id DESC`

	rows, err := sq.DB.Query(stmt, user_id)
	if err != nil {
		return nil, fmt.Errorf("repo.GetLikedPosts: %w", err)
	}

	defer rows.Close()

	posts := []*models.Post{}

	for rows.Next() {
		s := &models.Post{}
		err = rows.Scan(&s.Id, &s.UserId, &s.Title, &s.Content, &s.Likes, &s.Dislikes, &s.Created)
		if err != nil {
			return nil, fmt.Errorf("repo.GetLikedPosts: %w", err)
		}
		posts = append(posts, s)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("repo.GetLikedPosts: %w", err)
	}

	return posts, nil
}

func (sq *Sqlite) GetDislikedPosts(user_id int) ([]*models.Post, error) {
	stmt := `SELECT p.id, p.user_id, p.title, p.content, p.likes, p.dislikes, p.created FROM posts p
	JOIN user_post_reactions upr ON p.id = upr.post_id
	WHERE upr.user_id = ? AND upr.reaction = -1
    ORDER BY p.id DESC`

	rows, err := sq.DB.Query(stmt, user_id)
	if err != nil {
		return nil, fmt.Errorf("repo.GetDislikedPosts: %w", err)
	}

	defer rows.Close()

	posts := []*models.Post{}

	for rows.Next() {
		s := &models.Post{}
		err = rows.Scan(&s.Id, &s.UserId, &s.Title, &s.Content, &s.Likes, &s.Dislikes, &s.Created)
		if err != nil {
			return nil, fmt.Errorf("repo.GetDislikedPosts: %w", err)
		}
		posts = append(posts, s)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("repo.GetDislikedPosts: %w", err)
	}

	return posts, nil
}

func (sq *Sqlite) FindReportsForPost(post_id int) ([]*models.Report, error) {
	stmt := `SELECT id, post_id, moderator_id, reason, status, created_at FROM reports WHERE post_id = ? AND status = 'Pending';`

	rows, err := sq.DB.Query(stmt, post_id)
	if err != nil {
		return nil, fmt.Errorf("repo.FindReportsForPost: %w", err)
	}

	defer rows.Close()

	reports := []*models.Report{}

	for rows.Next() {
		s := &models.Report{}
		err = rows.Scan(&s.Id, &s.PostId, &s.ModeratorId, &s.Reason, &s.Status, &s.Created)
		if err != nil {
			return nil, fmt.Errorf("repo.FindReportsForPost: %w", err)
		}
		reports = append(reports, s)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("repo.FindReportsForPost: %w", err)
	}
	return reports, nil
}

func (sq *Sqlite) ChangeReportStatus(report_id int) error {
	stmt := `UPDATE reports SET status = 'Resolved' 
	WHERE id = ? AND status = 'Pending';`

	_, err := sq.DB.Exec(stmt, report_id)
	if err != nil {
		return fmt.Errorf("repo.ChangeReportStatus: %w", err)
	}
	return nil
}

func (sq *Sqlite) DeletePost(post_id int) error {
	stmt := `DELETE FROM posts WHERE id = ?`

	_, err := sq.DB.Exec(stmt, post_id)
	if err != nil {
		return fmt.Errorf("repo.DeletePost: %w", err)
	}
	return nil
}

func (sq *Sqlite) GetPostAuthor(post_id int) (*models.User, error) {
	stmt := `SELECT u.id from users u
	JOIN posts p ON u.id = p.user_id 
	WHERE p.id = ?`
	row := sq.DB.QueryRow(stmt, post_id)

	s := &models.User{}
	err := row.Scan(&s.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, fmt.Errorf("repo.GetPostAuthor: %w", err)
		}
	}
	return s, err
}

func (sq *Sqlite) UpdatePost(post_id int, title, content string) error {
	query := `UPDATE posts SET title = ?, content = ?, updated = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := sq.DB.Exec(query, title, content, post_id)
	return err
}
