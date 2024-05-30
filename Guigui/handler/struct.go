package handlers

type Category struct {
    CategoryID int
    Name       string
}

type Post struct {
    PostID    int
    Title     string
    Content   string
    Timestamp string
    Username  string
    ImagePaths []string
}

type Comment struct {
    CommentID  int
    Content    string
    Timestamp  string
    Username   string
    Score      int
    ImagePaths []string
}

type PostPreview struct {
    PostID          int
    Title           string
    ContentPreview  string
    FirstImage      string
    CommentCount    int
    PopularityScore float64
}

type WelcomeData struct {
    Categories []Category
    Posts      []PostPreview
}

type PostFormData struct {
    CategoryID int
    UserID     int
    Title      string
    Content    string
    Error      string
}


type CategoryPostsData struct {
    CategoryID int
    Posts      []PostPreview
}
