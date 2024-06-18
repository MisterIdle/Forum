package logic

type ErrorMessage struct {
	Error string
}

type Data struct {
	Data         interface{}
	Session      Session
	ErrorMessage ErrorMessage
}

type Session struct {
	Username string
	ID       int
	UUID     string
	Rank     string
	LoggedIn bool
}

type DashBoard struct {
	Users      []string
	Categories []string
	Profile    Profile
}

type Profile struct {
	Username      string
	UUID          string
	Picture       string
	Rank          string
	Timestamp     string
	TotalPosts    int
	TotalComments int
	TotalLikes    int
	TotalDislikes int
	Posts         []Posts
	Comments      []Comments
}

type Categories struct {
	Globals       map[string][]Category
	AllCategories []string
	AllGlobals    []string
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
	CategoryName string
	CategoryID   int
	PostID       int
	Title        string
	Content      string
	Username     string
	Timestamp    string
	LikesPost    int
	DislikesPost int
	Images       []string
	Comments     []Comments
}

type Comments struct {
	CommentID       int
	PostID          int
	Title           string
	Content         string
	Timestamp       string
	Username        string
	LikesComment    int
	DislikesComment int
	Sessions        Session
}
