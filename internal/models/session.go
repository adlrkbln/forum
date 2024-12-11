package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	UserId  int
	Token   string
	ExpTime time.Time
}

func NewSession(userId int) *Session {
	return &Session{
		UserId:  userId,
		Token:   uuid.New().String(),
		ExpTime: time.Now().Add(60 * time.Minute),
	}
}
