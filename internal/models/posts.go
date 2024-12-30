package models

import (
	"forum/internal/validate"
	"time"
)

type Post struct {
	Id         int
	UserId     int
	Username   string
	Title      string
	Content    string
	Categories []Category
	Comments   []Comment
	Likes      int
	Dislikes   int
	Created    time.Time
}

type PostCreateForm struct {
	Title              string
	Content            string
	CategoryIds        []int
	CategoryName       []string
	validate.Validator `form:"-"`
}

type Category struct {
	Id   int
	Name string
}

type CategoryCreateForm struct {
	Name              string
	validate.Validator `form:"-"`
}

type Comment struct {
	Id       int
	PostId   int
	UserId   int
	Username string
	Content  string
	Likes    int
	Dislikes int
	Created  time.Time
}

type CommentCreateForm struct {
	Content            string
	validate.Validator `form:"-"`
}
