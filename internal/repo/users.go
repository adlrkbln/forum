package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/models"
	"strings"

	"github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func (sq *Sqlite) InsertUser(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return fmt.Errorf("repo.InsertUser: %w", err)
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
		return fmt.Errorf("repo.InsertUser: %w", err)
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
			return 0, fmt.Errorf("repo.AuthenticateUser: %w", err)
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, fmt.Errorf("repo.AuthenticateUser: %w", err)
		}
	}

	return id, nil
}

func (sq *Sqlite) Exists(id int) (bool, error) {
	var exists bool
	stmt := `SELECT EXISTS(SELECT 1 FROM users WHERE id = ?);`

	err := sq.DB.QueryRow(stmt, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("repo.Exists: %w", err)
	}
	return exists, err
}

func (sq *Sqlite) GetUserByID(id int) (*models.User, error) {
	var u models.User
	stmt := `SELECT id, name, email, created, role FROM users WHERE id=?`
	err := sq.DB.QueryRow(stmt, id).Scan(&u.Id, &u.Name, &u.Email, &u.Created, &u.Role)
	if err != nil {
		return nil, fmt.Errorf("repo.GetUserByID: %w", err)
	}
	return &u, nil
}

func (sq *Sqlite) GetAllUsers() ([]*models.User, error) {
	stmt := `SELECT id, name, email, role FROM users 
	WHERE role = 'User' OR role = 'Moderator'
	ORDER BY role, name `
	rows, err := sq.DB.Query(stmt)
	if err != nil {
		return nil, fmt.Errorf("repo.GetAllUsers: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		u := &models.User{}
		if err := rows.Scan(&u.Id, &u.Name, &u.Email, &u.Role); err != nil {
			return nil, fmt.Errorf("repo.GetAllUsers: %w", err)
		}
		users = append(users, u)
	}
	return users, nil
}

func (sq *Sqlite) InsertReportPost(moderator_id int, post_id int, reason string) error {
	_, err := sq.DB.Exec("INSERT INTO reports (post_id, moderator_id, reason) VALUES (?, ?, ?)", post_id, moderator_id, reason)
	if err != nil {
		return fmt.Errorf("repo.InsertReportPost: %w", err)
	}
	return nil
}

func (sq *Sqlite) GetAllReports() ([]*models.Report, error) {
	stmt := `SELECT r.id, r.post_id, r.moderator_id, u.name, r.reason, r.status, r.created_at FROM reports r
	JOIN users u ON u.id = r.moderator_id
    ORDER BY r.id DESC`

	rows, err := sq.DB.Query(stmt)
	if err != nil {
		return nil, fmt.Errorf("repo.GetAllReports: %w", err)
	}

	defer rows.Close()

	reports := []*models.Report{}

	for rows.Next() {
		s := &models.Report{}
		err = rows.Scan(&s.Id, &s.PostId, &s.ModeratorId, &s.ModeratorName, &s.Reason, &s.Status, &s.Created)
		if err != nil {
			return nil, fmt.Errorf("repo.GetAllReports: %w", err)
		}
		reports = append(reports, s)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("repo.GetAllReports: %w", err)
	}

	return reports, nil
}

func (sq *Sqlite) RequestModeratorRole(user_id int) error {
	stmt := "INSERT INTO moderator_requests (user_id, status) VALUES (?, 'Pending')"
	_, err := sq.DB.Exec(stmt, user_id)
	if err != nil {
		return fmt.Errorf("repo.RequestModeratorRole: %w", err)
	}
	return nil
}

func (sq *Sqlite) GetAllRequests() ([]*models.ModeratorRequest, error) {
	stmt := `
        SELECT r.id, r.user_id, u.name AS username, r.status, r.requested_at
        FROM moderator_requests r
        JOIN users u ON r.user_id = u.id
		ORDER BY r.id DESC
    `
	rows, err := sq.DB.Query(stmt)
	if err != nil {
		return nil, fmt.Errorf("repo.GetAllRequests: %w", err)
	}
	defer rows.Close()

	var requests []*models.ModeratorRequest
	for rows.Next() {
		s := &models.ModeratorRequest{}
		if err := rows.Scan(&s.Id, &s.UserId, &s.Username, &s.Status, &s.RequestedAt); err != nil {
			return nil, fmt.Errorf("repo.GetAllRequests: %w", err)
		}
		requests = append(requests, s)
	}
	return requests, nil
}

func (sq *Sqlite) PromoteUserToModerator(request_id int) error {
	tx, err := sq.DB.Begin()
	if err != nil {
		return fmt.Errorf("repo.PromoteUserToModerator: %w", err)
	}

	var user_id int
	err = tx.QueryRow("SELECT user_id FROM moderator_requests WHERE id = ?", request_id).Scan(&user_id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("repo.PromoteUserToModerator: %w", err)
	}

	_, err = tx.Exec("UPDATE users SET role = 'Moderator' WHERE id = ?", user_id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("repo.PromoteUserToModerator: %w", err)
	}

	_, err = tx.Exec("UPDATE moderator_requests SET status = 'Approved' WHERE user_id = ?", user_id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("repo.PromoteUserToModerator: %w", err)
	}

	return tx.Commit()
}

func (sq *Sqlite) DenyModeratorRequest(request_id int) error {
	query := "UPDATE moderator_requests SET status = 'Denied' WHERE id = ?"
	_, err := sq.DB.Exec(query, request_id)
	if err != nil {
		return fmt.Errorf("repo.DenyModeratorRequest: %w", err)
	}
	return nil
}

func (sq *Sqlite) GetUserModeratorRequests(user_id int) ([]*models.ModeratorRequest, error) {
	stmt := `SELECT r.id, r.user_id, u.name AS username, r.status, r.requested_at
        FROM moderator_requests r
        JOIN users u ON r.user_id = u.id
        WHERE r.user_id = ?`

	rows, err := sq.DB.Query(stmt, user_id)
	if err != nil {
		return nil, fmt.Errorf("repo.GetUserModeratorRequests: %w", err)
	}
	defer rows.Close()

	var requests []*models.ModeratorRequest
	for rows.Next() {
		s := &models.ModeratorRequest{}
		if err := rows.Scan(&s.Id, &s.UserId, &s.Username, &s.Status, &s.RequestedAt); err != nil {
			return nil, fmt.Errorf("repo.GetUserModeratorRequests: %w", err)
		}
		requests = append(requests, s)
	}
	return requests, nil
}

func (sq *Sqlite) GetModeratorReports(user_id int) ([]*models.Report, error) {
	stmt := `SELECT r.id, r.post_id, r.moderator_id, u.name, r.reason, r.status, r.created_at FROM reports r
	JOIN users u ON u.id = r.moderator_id
	WHERE r.moderator_id = ?
    ORDER BY r.id DESC`

	rows, err := sq.DB.Query(stmt, user_id)
	if err != nil {
		return nil, fmt.Errorf("repo.GetModeratorReports: %w", err)
	}

	defer rows.Close()

	reports := []*models.Report{}

	for rows.Next() {
		s := &models.Report{}
		err = rows.Scan(&s.Id, &s.PostId, &s.ModeratorId, &s.ModeratorName, &s.Reason, &s.Status, &s.Created)
		if err != nil {
			return nil, fmt.Errorf("repo.GetModeratorReports: %w", err)
		}
		reports = append(reports, s)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("repo.GetModeratorReports: %w", err)
	}

	return reports, nil
}

func (sq *Sqlite) DemoteModerator(userID int) error {
	query := "UPDATE users SET role = 'User' WHERE id = ? AND role = 'Moderator'"
	_, err := sq.DB.Exec(query, userID)
	return fmt.Errorf("repo.DemoteModerator: %w", err)
}

func (sq *Sqlite) GetUserByEmail(email string) (*models.User, error) {
	var u models.User
	stmt := `SELECT id, name, email, created, role FROM users WHERE email=?`
	err := sq.DB.QueryRow(stmt, email).Scan(&u.Id, &u.Name, &u.Email, &u.Created, &u.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("repo.GetUserByEmail: %w", err)
	}
	return &u, nil
}
