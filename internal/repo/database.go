package repo

import (
	"database/sql"
	"fmt"
)

type Sqlite struct {
	DB *sql.DB
}

func OpenDB(dsn string) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	queries := []string{
		`CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			image_path TEXT DEFAULT '/static/img/default.png',
			likes INTEGER DEFAULT 0,
			dislikes INTEGER DEFAULT 0,
			created TIMESTAMP,
			updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			token TEXT NOT NULL,
			exp_time TIMESTAMP NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS user_post_reactions (
			user_id INTEGER,
			post_id INTEGER,
			reaction INTEGER, -- 1 for like, -1 for dislike
			PRIMARY KEY (user_id, post_id)
		);`,
		`CREATE TABLE IF NOT EXISTS comment_reactions (
            comment_id INTEGER,
            user_id INTEGER,
            reaction INTEGER NOT NULL, -- 1 for like, -1 for dislike
            PRIMARY KEY (comment_id, user_id),
            FOREIGN KEY (comment_id) REFERENCES comments(id),
            FOREIGN KEY (user_id) REFERENCES users(id)
        );`,
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL,
			hashed_password CHAR(60) NOT NULL,
			created DATETIME NOT NULL,
			role TEXT NOT NULL DEFAULT 'User',
			UNIQUE(email)
		);`,
		`CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			likes INTEGER DEFAULT 0,
			dislikes INTEGER DEFAULT 0,
			created_at DATETIME,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (post_id) REFERENCES posts(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		);`,
		`CREATE TABLE IF NOT EXISTS reports (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			moderator_id INTEGER NOT NULL,
			reason TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'Pending', -- Status: Pending, Resolved
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(post_id) REFERENCES posts(id),
			FOREIGN KEY(moderator_id) REFERENCES users(id)
		);`,
		`CREATE TABLE IF NOT EXISTS moderator_requests (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			status TEXT NOT NULL DEFAULT 'Pending', -- Pending, Approved, Denied
			requested_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS category (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS post_category (
			category_id INTEGER,
			post_id INTEGER, 
			PRIMARY KEY (category_id, post_id),
			FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY (category_id) REFERENCES category(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS notifications (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			post_id INTEGER,
			type TEXT NOT NULL,
			message TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			read BOOLEAN DEFAULT 0,
			FOREIGN KEY (user_id) REFERENCES users (id),
			FOREIGN KEY (post_id) REFERENCES posts (id)
		);`,
	}
	for _, query := range queries {
		stmt, err := db.Prepare(query)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", "repo.OpenDB", err)
		}
		_, err = stmt.Exec()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", "repo.OpenDB", err)
		}
		stmt.Close()
	}
	return &Sqlite{DB: db}, nil
}

func NewDB(dsn string) (*Sqlite, error) {
	return OpenDB(dsn)
}
