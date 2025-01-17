package models

type TemplateData struct {
	CurrentYear       int
	Post              *Post
	Posts             []*Post
	Category          *Category
	Categories        []*Category
	User              *User
	Users             []*User
	Report            *Report
	Reports           []*Report
	ModeratorRequest  *ModeratorRequest
	ModeratorRequests []*ModeratorRequest
	Notification      *Notification
	Notifications     []*Notification
	Form              any
	IsAuthenticated   bool
	IsModerator       bool
	IsAdmin           bool
}
