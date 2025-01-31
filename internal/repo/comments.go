package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/models"
)

func (sq *Sqlite) DeleteComment(commentID int) error {
	stmt := `DELETE FROM comments WHERE id = ?`
	_, err := sq.DB.Exec(stmt, commentID)
	if err != nil {
		return fmt.Errorf("repo.DeleteComment: %w", err)
	}
	return nil
}

func (sq *Sqlite) GetCommentedPostsByUser(userId int) ([]*models.CommentWithPost, error) {
	stmt := `
    SELECT 
        c.id, c.content, c.created_at, 
        p.id AS post_id, p.title AS post_title, p.created AS post_created
    FROM comments c
    INNER JOIN posts p ON c.post_id = p.id
    WHERE c.user_id = ?
    ORDER BY c.created_at DESC
    `
	rows, err := sq.DB.Query(stmt, userId)
	if err != nil {
		return nil, fmt.Errorf("repo.GetCommentedPostsByUser: %w", err)
	}
	defer rows.Close()

	commentsWithPosts := []*models.CommentWithPost{}
	for rows.Next() {
		commentWithPost := &models.CommentWithPost{}
		err := rows.Scan(
			&commentWithPost.Comment.Id,
			&commentWithPost.Comment.Content,
			&commentWithPost.Comment.Created,
			&commentWithPost.Post.Id,
			&commentWithPost.Post.Title,
			&commentWithPost.Post.Created,
		)
		if err != nil {
			return nil, fmt.Errorf("repo.GetCommentedPostsByUser: %w", err)
		}
		commentsWithPosts = append(commentsWithPosts, commentWithPost)
	}
	return commentsWithPosts, nil
}

func (sq *Sqlite) InsertComment(post_id int, user_id int, content string) error {
	stmt := `INSERT INTO comments (post_id, user_id, content)
	VALUES (?, ?, ?);`

	_, err := sq.DB.Exec(stmt, post_id, user_id, content)
	if err != nil {
		return fmt.Errorf("repo.InsertComment: %w", err)
	}
	return nil
}

func (sq *Sqlite) GetAllComments() ([]*models.Comment, error) {
	stmt := `SELECT id FROM comments
    ORDER BY id DESC`

	rows, err := sq.DB.Query(stmt)
	if err != nil {
		return nil, fmt.Errorf("repo.GetAllComments: %w", err)
	}

	defer rows.Close()

	comments := []*models.Comment{}

	for rows.Next() {
		s := &models.Comment{}
		err = rows.Scan(&s.Id)
		if err != nil {
			return nil, fmt.Errorf("repo.GetAllComments: %w", err)
		}
		comments = append(comments, s)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("repo.GetAllComments: %w", err)
	}

	return comments, nil
}

func (sq *Sqlite) GetCommentsForPost(post_id int) ([]models.Comment, error) {
	stmt := `SELECT c.id, c.post_id, c.user_id, u.name, c.content, c.created_at, c.likes, c.dislikes FROM comments c
	JOIN users u ON u.id = c.user_id
	WHERE c.post_id = ?
    ORDER BY c.id`

	rows, err := sq.DB.Query(stmt, post_id)
	if err != nil {
		return nil, fmt.Errorf("repo.GetCommentsForPost: %w", err)
	}

	defer rows.Close()

	comments := []models.Comment{}

	for rows.Next() {
		s := models.Comment{}
		err = rows.Scan(&s.Id, &s.PostId, &s.UserId, &s.Username, &s.Content, &s.Created, &s.Likes, &s.Dislikes)
		if err != nil {
			return nil, fmt.Errorf("repo.GetCommentsForPost: %w", err)
		}
		comments = append(comments, s)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("repo.GetCommentsForPost: %w", err)
	}

	return comments, nil
}

func (sq *Sqlite) UpdateComment(comment_id int, content string) error {
	query := `UPDATE comments SET content = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := sq.DB.Exec(query, content, comment_id)
	return fmt.Errorf("repo.UpdateComment: %w", err)
}

func (sq *Sqlite) GetComment(id int) (*models.Comment, error) {
	stmt := `SELECT id, user_id, post_id, content, likes, dislikes, created_at FROM comments c
	WHERE id = ?`

	row := sq.DB.QueryRow(stmt, id)

	s := &models.Comment{}
	err := row.Scan(&s.Id, &s.UserId, &s.PostId, &s.Content, &s.Likes, &s.Dislikes, &s.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, fmt.Errorf("repo.GetComment: %w", err)
		}
	}

	return s, nil
}

func (sq *Sqlite) GetCommentAuthor(comment_id int) (*models.User, error) {
	stmt := `SELECT u.id from users u
	JOIN comments c ON u.id = c.user_id 
	WHERE c.id = ?`
	row := sq.DB.QueryRow(stmt, comment_id)

	s := &models.User{}
	err := row.Scan(&s.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, fmt.Errorf("repo.GetCommentAuthor: %w", err)
		}
	}
	return s, nil
}
