package repo

import (
    "fmt"
    "forum/internal/models"
)

func (sq *Sqlite) CreateNotification(notification *models.Notification) error {
    stmt := `
    INSERT INTO notifications (user_id, post_id, type, message, read)
    VALUES (?, ?, ?, ?, ?)
    `
    _, err := sq.DB.Exec(stmt, notification.UserId, notification.PostId, notification.Type, notification.Message, notification.Read)
    if err != nil {
        return fmt.Errorf("repo.CreateNotification %w", err)
    }
    return nil
}

func (sq *Sqlite) GetUnreadNotifications(userId int) ([]*models.Notification, error) {
    stmt := `
    SELECT id, user_id, post_id, type, message, created_at, read
    FROM notifications
    WHERE user_id = ? AND read = 0
    ORDER BY created_at DESC
    `
    rows, err := sq.DB.Query(stmt, userId)
    if err != nil {
        return nil, fmt.Errorf("repo.GetUnreadNotifications %w", err)
    }
    defer rows.Close()

    notifications := []*models.Notification{}
    for rows.Next() {
        n := &models.Notification{}
        err := rows.Scan(&n.Id, &n.UserId, &n.PostId, &n.Type, &n.Message, &n.CreatedAt, &n.Read)
        if err != nil {
            return nil, fmt.Errorf("repo.GetUnreadNotifications %w", err)
        }
        notifications = append(notifications, n)
    }
    return notifications, nil
}

func (sq *Sqlite) MarkNotificationAsRead(notificationId int) error {
    stmt := `
    UPDATE notifications
    SET read = 1
    WHERE id = ?
    `
    _, err := sq.DB.Exec(stmt, notificationId)
    if err != nil {
        return fmt.Errorf("repo.MarkNotificationAsRead %w", err)
    }
    return nil
}

func (sq *Sqlite) GetNotifications() ([]*models.Notification, error)  {
	var notifications []*models.Notification

	rows, err := sq.DB.Query("SELECT id FROM notifications")
	if err != nil {
		return nil, fmt.Errorf("repo.GetNotifications %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		notification := &models.Notification{}
		if err := rows.Scan(&notification.Id); err != nil {
			return nil, fmt.Errorf("repo.GetNotifications %w", err)
		}
		notifications = append(notifications, notification)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repo.GetNotifications %w", err)
	}
	return notifications, nil
}