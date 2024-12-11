package repo

import (
	"database/sql"
	"log"
)

func (sq *Sqlite) AddLikePost(post_id int, user_id int) error {
	stmt := `UPDATE posts SET likes = likes + 1 WHERE id = ?`

	_, err := sq.DB.Exec(stmt, post_id)
	if err != nil {
		return err
	}

	return nil
}

func (sq *Sqlite) AddDislikePost(post_id int, user_id int) error {
	stmt := `UPDATE posts SET dislikes = dislikes + 1 WHERE id = ?`

	_, err := sq.DB.Exec(stmt, post_id)
	if err != nil {
		return err
	}
	return nil
}

func (sq *Sqlite) CheckUserReactionsPost(post_id int, user_id int) (int, error) {
	stmt := `SELECT reaction FROM user_post_reactions 
             WHERE post_id = ? AND user_id = ?`

	var reaction int
	err := sq.DB.QueryRow(stmt, post_id, user_id).Scan(&reaction)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return reaction, nil
}

func (sq *Sqlite) InsertUserReactionPost(post_id int, user_id int, reaction int) error {
	stmt := `INSERT INTO user_post_reactions (post_id, user_id, reaction)
	VALUES (?, ?, ?);`

	_, err := sq.DB.Exec(stmt, post_id, user_id, reaction)
	if err != nil {
		return err
	}

	return nil
}

func (sq *Sqlite) RemoveUserReactionPost(post_id int, user_id int, reaction int) error {
	stmt := `DELETE FROM user_post_reactions 
	WHERE post_id = ? AND user_id = ?;`

	_, err := sq.DB.Exec(stmt, post_id, user_id)
	if err != nil {
		return err
	}

	if reaction == 1 {
		stmt = `UPDATE posts SET likes = likes - 1 
		WHERE id = ?`
	} else {
		stmt = `UPDATE posts SET dislikes = dislikes - 1 
		WHERE id = ?`
	}

	_, err = sq.DB.Exec(stmt, post_id)
	if err != nil {
		return err
	}
	return nil
}

func (sq *Sqlite) AddLikeComment(comment_id int, user_id int) error {
	stmt := `UPDATE comments SET likes = likes + 1 WHERE id = ?`

	_, err := sq.DB.Exec(stmt, comment_id)
	if err != nil {
		return err
	}
	log.Println("like")

	return nil
}

func (sq *Sqlite) AddDislikeComment(comment_id int, user_id int) error {
	stmt := `UPDATE comments SET dislikes = dislikes + 1 WHERE id = ?`

	_, err := sq.DB.Exec(stmt, comment_id)
	if err != nil {
		return err
	}
	return nil
}

func (sq *Sqlite) CheckUserReactionComment(comment_id int, user_id int) (int, error) {
	stmt := `SELECT reaction FROM comment_reactions 
             WHERE comment_id = ? AND user_id = ?`

	var reaction int
	err := sq.DB.QueryRow(stmt, comment_id, user_id).Scan(&reaction)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return reaction, nil
}

func (sq *Sqlite) InsertUserReactionComment(comment_id int, user_id int, reaction int) error {
	stmt := `INSERT INTO comment_reactions (comment_id, user_id, reaction)
	VALUES (?, ?, ?);`

	_, err := sq.DB.Exec(stmt, comment_id, user_id, reaction)
	if err != nil {
		return err
	}

	return nil
}

func (sq *Sqlite) RemoveUserReactionComment(comment_id int, user_id int, reaction int) error {
	stmt := `DELETE FROM comment_reactions 
	WHERE comment_id = ? AND user_id = ?;`

	_, err := sq.DB.Exec(stmt, comment_id, user_id)
	if err != nil {
		return err
	}

	if reaction == 1 {
		stmt = `UPDATE comments SET likes = likes - 1 
		WHERE id = ?`
	} else {
		stmt = `UPDATE comments SET dislikes = dislikes - 1 
		WHERE id = ?`
	}

	_, err = sq.DB.Exec(stmt, comment_id)
	if err != nil {
		return err
	}
	return nil
}
