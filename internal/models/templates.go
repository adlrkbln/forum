package models

type TemplateData struct {
	CurrentYear     int
	Post            *Post
	Posts           []*Post
	Category        *Category
	Categories      []*Category
	User            *User
	Form            any
	IsAuthenticated bool
}
