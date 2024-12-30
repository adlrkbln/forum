package models

import (
	"forum/internal/validate"
	"time"
)

type User struct {
	Id             int
	Name           string
	Email          string
	HashedPassword []byte
	Role           string
	Created        time.Time
}

type UserSignupForm struct {
	Name               string `form:"name"`
	Email              string `form:"email"`
	Password           string `form:"password"`
	validate.Validator `form:"-"`
}

type UserLoginForm struct {
	Email              string `form:"email"`
	Password           string `form:"password"`
	validate.Validator `form:"-"`
}

type Report struct {
	Id            int
	PostId        int
	ModeratorId   int
	ModeratorName string
	Reason        string
	Status        string // "Pending", "Resolved"
	Created       time.Time
}

type ModeratorRequest struct {
	Id          int
	UserId      int
	Username    string
	Status      string
	RequestedAt string
}
