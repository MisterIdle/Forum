package logic

type ErrorMessage struct {
	Error string
}

type Data struct {
	Data    interface{}
	Session Session
}

type Session struct {
	Username string
	ID       int
	UUID     string
	LoggedIn bool
}

type Categories struct {
	Categories map[string][]Category
}

type Category struct {
	CategoryID    int
	Name          string
	Description   string
	TotalPosts    int
	TotalComments int
	Global        string
	Posts         []Posts
}

type Posts struct {
	PostID       int
	Title        string
	Content      string
	Username     string
	Timestamp    string
	Likes        int
	Dislikes     int
	Images       []string
	Comments     []Comment
	CategoryName string
	CategoryID   int
}

type Comment struct {
	CommentID int
	PostID    int
	Title     string
	Content   string
	Timestamp string
	Username  string
	Likes     int
	Dislikes  int
}
