package repo

import (
	"database/sql"
	"errors"
	"forum/internal/models"
)

func (sq *Sqlite) InsertPost(user_id int, title string, content string) (int, error) {
	stmt := `INSERT INTO posts (user_id, title, content, created)
	VALUES (?, ?, ?, CURRENT_TIMESTAMP);`

	result, err := sq.DB.Exec(stmt, user_id, title, content)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (sq *Sqlite) GetPost(id int) (*models.Post, error) {
	stmt := `SELECT p.id, p.user_id, p.title, p.content, p.likes, p.dislikes, p.created, u.name FROM posts p
	JOIN users u ON p.user_id = u.id
	WHERE p.id = ?`

	row := sq.DB.QueryRow(stmt, id)

	s := &models.Post{}
	err := row.Scan(&s.Id, &s.UserId, &s.Title, &s.Content, &s.Likes, &s.Dislikes, &s.Created, &s.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (sq *Sqlite) GetAllPosts() ([]*models.Post, error) {
	stmt := `SELECT id, user_id, title, content, likes, dislikes, created FROM posts
    ORDER BY id DESC`

	rows, err := sq.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	posts := []*models.Post{}

	for rows.Next() {
		s := &models.Post{}
		err = rows.Scan(&s.Id, &s.UserId, &s.Title, &s.Content, &s.Likes, &s.Dislikes, &s.Created)
		if err != nil {
			return nil, err
		}
		posts = append(posts, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (sq *Sqlite) GetAllComments() ([]*models.Comment, error) {
	stmt := `SELECT id FROM comments
    ORDER BY id DESC`

	rows, err := sq.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	comments := []*models.Comment{}

	for rows.Next() {
		s := &models.Comment{}
		err = rows.Scan(&s.Id)
		if err != nil {
			return nil, err
		}
		comments = append(comments, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
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
		return nil, err
	}

	defer rows.Close()

	comments := []models.Comment{}

	for rows.Next() {
		s := models.Comment{}
		err = rows.Scan(&s.Id, &s.PostId, &s.UserId, &s.Username, &s.Content, &s.Created, &s.Likes, &s.Dislikes)
		if err != nil {
			return nil, err
		}
		comments = append(comments, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (sq *Sqlite) InsertComment(post_id int, user_id int, content string) error {
	stmt := `INSERT INTO comments (post_id, user_id, content)
	VALUES (?, ?, ?);`

	_, err := sq.DB.Exec(stmt, post_id, user_id, content)
	if err != nil {
		return err
	}
	return nil
}

func (sq *Sqlite) GetCreatedPosts(user_id int) ([]*models.Post, error) {
	stmt := `SELECT id, user_id, title, content, likes, dislikes, created FROM posts
	WHERE user_id = ?
    ORDER BY id DESC`

	rows, err := sq.DB.Query(stmt, user_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	posts := []*models.Post{}

	for rows.Next() {
		s := &models.Post{}
		err = rows.Scan(&s.Id, &s.UserId, &s.Title, &s.Content, &s.Likes, &s.Dislikes, &s.Created)
		if err != nil {
			return nil, err
		}
		posts = append(posts, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
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
		return nil, err
	}

	defer rows.Close()

	posts := []*models.Post{}

	for rows.Next() {
		s := &models.Post{}
		err = rows.Scan(&s.Id, &s.UserId, &s.Title, &s.Content, &s.Likes, &s.Dislikes, &s.Created)
		if err != nil {
			return nil, err
		}
		posts = append(posts, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (sq *Sqlite) FindReportsForPost(post_id int) ([]*models.Report, error) {
	stmt := `SELECT id, post_id, moderator_id, reason, status, created_at FROM reports WHERE post_id = ? AND status = 'Pending';`

	rows, err := sq.DB.Query(stmt, post_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	reports := []*models.Report{}

	for rows.Next() {
		s := &models.Report{}
		err = rows.Scan(&s.Id, &s.PostId, &s.ModeratorId, &s.Reason, &s.Status, &s.Created)
		if err != nil {
			return nil, err
		}
		reports = append(reports, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return reports, nil
}

func (sq *Sqlite) ChangeReportStatus(report_id int) error {
	stmt := `UPDATE reports SET status = 'Resolved' 
	WHERE id = ? AND status = 'Pending';`

	_, err := sq.DB.Exec(stmt, report_id)
	if err != nil {
		return err
	}
	return nil
}

func (sq *Sqlite) DeletePost(post_id int) error {
	stmt := `DELETE FROM posts WHERE id = ?`

	_, err := sq.DB.Exec(stmt, post_id)
	if err != nil {
		return err
	}
	return nil
}

func (sq *Sqlite) DeleteComment(commentID int) error {
	stmt := `DELETE FROM comments WHERE id = ?`
	_, err := sq.DB.Exec(stmt, commentID)
	if err != nil {
		return err
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
			return nil, err
		}
	}
	return s, err
}
