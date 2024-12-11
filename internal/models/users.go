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
	Created        time.Time
}

type UserSignupForm struct {
	Name               string `form:"name"`
	Email              string `form:"email"`
	Password           string `form:"password"`
	validate.Validator `form:"-"`
}

type UserLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validate.Validator `form:"-"`
}