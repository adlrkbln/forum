package repo

import (
	"database/sql"
	"errors"
	"forum/internal/models"
	"strings"

	"github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func (sq *Sqlite) InsertUser(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created)
    VALUES(?, ?, ?, CURRENT_TIMESTAMP)`

	_, err = sq.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique && strings.Contains(sqliteErr.Error(), "email") {
				return models.ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (sq *Sqlite) AuthenticateUser(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := "SELECT id, hashed_password FROM users WHERE email = ?"

	err := sq.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (sq *Sqlite) Exists(id int) (bool, error) {
	var exists bool
	stmt := `SELECT EXISTS(SELECT 1 FROM users WHERE id = ?);`

	err := sq.DB.QueryRow(stmt, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, err
}

func (sq *Sqlite) GetUserByID(id int) (*models.User, error) {
	var u models.User
	stmt := `SELECT id, name, email, created, role FROM users WHERE id=?`
	err := sq.DB.QueryRow(stmt, id).Scan(&u.Id, &u.Name, &u.Email, &u.Created, &u.Role)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (sq *Sqlite) InsertReportPost(moderator_id int, post_id int, reason string) error {
	_, err := sq.DB.Exec("INSERT INTO reports (post_id, moderator_id, reason) VALUES (?, ?, ?)", post_id, moderator_id, reason)
	return err
}

func (sq *Sqlite) GetAllReports() ([]*models.Report, error) {
	stmt := `SELECT r.id, r.post_id, r.moderator_id, u.name, r.reason, r.status, r.created_at FROM reports r
	JOIN users u ON u.id = r.moderator_id
    ORDER BY r.id DESC`

	rows, err := sq.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	reports := []*models.Report{}

	for rows.Next() {
		s := &models.Report{}
		err = rows.Scan(&s.Id, &s.PostId, &s.ModeratorId, &s.ModeratorName, &s.Reason, &s.Status, &s.Created)
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