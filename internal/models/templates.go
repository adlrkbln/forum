package models

type TemplateData struct {
	CurrentYear       int
	Post              *Post
	Posts             []*Post
	Comment           *Comment
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
	LikedPosts        []*Post
	DislikedPosts     []*Post
	CreatedPosts      []*Post
	CommentedPosts    []*CommentWithPost
	Form              any
	IsAuthenticated   bool
	IsModerator       bool
	IsAdmin           bool
}
