package repo

import "forum/internal/models"


func (s *Sqlite) DeleteSessionById(userId int) error {
	stmt := `DELETE FROM sessions WHERE user_id = ?`
	if _, err := s.DB.Exec(stmt, userId); err != nil {
		return err
	}
	return nil
}

func (s *Sqlite) CreateSession(session *models.Session) error {
	stmt := `INSERT INTO sessions(user_id, token, exp_time) VALUES(?, ?, ?)`
	_, err := s.DB.Exec(stmt, session.UserId, session.Token, session.ExpTime)
	if err != nil {
		return err
	}
	return nil
}

func (s *Sqlite) DeleteSessionByToken(token string) error {
	stmt := `DELETE FROM sessions WHERE token = ?`
	if _, err := s.DB.Exec(stmt, token); err != nil {
		return err
	}
	return nil
}

func (s *Sqlite) GetUserIDByToken(token string) (int, error) {
	stmt := `SELECT user_id FROM sessions WHERE token = ?`
	var userID int

	err := s.DB.QueryRow(stmt, token).Scan(&userID)
	if err != nil {
		return -1, err
	}
	return userID, nil
}

func (s *Sqlite) IsSessionValid(token string) bool {
    query := `SELECT 1 FROM sessions WHERE token = ?`
    var exists int
    err := s.DB.QueryRow(query, token).Scan(&exists)
    return err == nil
}
