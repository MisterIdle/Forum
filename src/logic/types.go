package logic

type ErrorMessage struct {
	Error string
}

type Data struct {
	Data     interface{}
	LoggedIn bool
}

type Session struct {
	Username string
	LoggedIn bool
}

type Categories struct {
	Categories map[string][]Category
}

type Category struct {
	CategoryID int
	Name       string
	Global     string
	Posts      []Posts
}

type Posts struct {
	PostID    int
	Title     string
	Content   string
	Timestamp string
	Comments  []Comment
}

type Comment struct {
	CommentID int
	Content   string
	Timestamp string
}
