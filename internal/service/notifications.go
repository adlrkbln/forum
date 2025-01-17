package service

import "forum/internal/models"

func (s *service) NotifyUser(userId int, postId int, notifType, message string) error {
    notification := &models.Notification{
        UserId:  userId,
        PostId:  postId,
        Type:    notifType,
        Message: message,
        Read:    false,
    }
    return s.repo.CreateNotification(notification)
}

func (s *service) GetUnreadNotifications(userId int) ([]*models.Notification, error) {
    return s.repo.GetUnreadNotifications(userId)
}

func (s *service) MarkNotificationAsRead(notificationId int) error {
    return s.repo.MarkNotificationAsRead(notificationId)
}

func (s *service) GetNotifications() ([]*models.Notification, error) {
	return s.repo.GetNotifications()
}