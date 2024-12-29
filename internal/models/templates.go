package models

type TemplateData struct {
	CurrentYear     int
	Post            *Post
	Posts           []*Post
	Category        *Category
	Categories      []*Category
	User            *User
	Report          *Report
	Reports         []*Report
	Form            any
	IsAuthenticated bool
	IsModerator     bool
}
