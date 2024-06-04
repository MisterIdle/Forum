package logic

type Category struct {
	CategoryID int
	Name       string
}

type Post struct {
	PostID     int
	Title      string
	Content    string
	Timestamp  string
	Username   string
	UserID     int
	ImagePaths []string
}

type PostData struct {
	Post     Post
	Comments []Comment
}

type Comment struct {
	CommentID  int
	Content    string
	Timestamp  string
	Username   string
	UserID     int
	Score      int
	ImagePaths []string
}

type WelcomeData struct {
	Categories      []Category
	Posts           []PostPreview
	PostCount       int
	CommentCount    int
	UniqueUserCount int
	LatestMember    string
}

type PostPreview struct {
	PostID          int
	Title           string
	ContentPreview  string
	FirstImage      string
	CommentCount    int
	PopularityScore float64
}

type PostFormData struct {
    CategoryID  int
    UserID      int
    Title       string
    Content     string
    ImagePaths  []string
    Error       string
    Categories  []Category
}

type CategoryPostsData struct {
	CategoryID int
	Posts      []PostPreview
}

type SearchResult struct {
	Posts []PostResult `json:"posts"`
	Users []UserResult `json:"users"`
}

type PostResult struct {
	PostID       int    `json:"post_id"`
	Title        string `json:"title"`
	CategoryName string `json:"category_name"`
	UniqueUsers  int    `json:"unique_users"`
}

type UserResult struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
}
