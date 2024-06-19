package logic

// Data struct contains the main data and session information.
type Data struct {
	Data    interface{} // Main data to be passed to the template
	Session Session     // Session information to be passed to the template
}

// Session struct contains information about the user's session.
type Session struct {
	Username string // Username of the logged-in user
	ID       int    // User ID
	UUID     string // User UUID
	Rank     string // User rank (e.g., admin, moderator, user)
	LoggedIn bool   // Indicates if the user is logged in
	HasError bool   // Indicates if there is an error in the session
	Message  string // Error or information message
}

// Error struct is used to encapsulate error information.
type Error struct {
	HasError bool   // Indicates if there is an error
	Message  string // Error message
}

// DashBoard struct is used to represent data for the dashboard page.
type DashBoard struct {
	Users      []string // List of usernames
	Categories []string // List of category names
	Profile    Profile  // Profile data for the user
}

// Profile struct contains detailed information about a user's profile.
type Profile struct {
	Username      string     // Username
	UUID          string     // User UUID
	Picture       string     // Profile picture URL or file name
	Rank          string     // User rank
	Timestamp     string     // Account creation timestamp
	TotalPosts    int        // Total number of posts made by the user
	TotalComments int        // Total number of comments made by the user
	TotalLikes    int        // Total number of likes received by the user
	TotalDislikes int        // Total number of dislikes received by the user
	Posts         []Posts    // List of posts made by the user
	Comments      []Comments // List of comments made by the user
}

// Categories struct contains data related to categories.
type Categories struct {
	Globals       map[string][]Category // Map of global categories
	AllCategories []string              // List of all category names
	AllGlobals    []string              // List of all global category names
}

// Category struct contains detailed information about a category.
type Category struct {
	CategoryID    int     // Category ID
	Name          string  // Category name
	Description   string  // Category description
	TotalPosts    int     // Total number of posts in the category
	TotalComments int     // Total number of comments in the category
	Global        string  // Global category name
	Posts         []Posts // List of posts in the category
}

// Posts struct contains detailed information about a post.
type Posts struct {
	CategoryName string     // Name of the category the post belongs to
	CategoryID   int        // ID of the category the post belongs to
	PostID       int        // Post ID
	Title        string     // Post title
	Content      string     // Post content
	Username     string     // Username of the author
	Timestamp    string     // Timestamp when the post was created
	LikesPost    int        // Number of likes on the post
	DislikesPost int        // Number of dislikes on the post
	Images       []string   // List of image URLs or file names associated with the post
	Comments     []Comments // List of comments on the post
}

// Comments struct contains detailed information about a comment.
type Comments struct {
	CommentID       int     // Comment ID
	PostID          int     // ID of the post the comment belongs to
	Title           string  // Title of the comment
	Content         string  // Content of the comment
	Timestamp       string  // Timestamp when the comment was created
	Username        string  // Username of the commenter
	LikesComment    int     // Number of likes on the comment
	DislikesComment int     // Number of dislikes on the comment
	Sessions        Session // Session information (could be used for context in templates)
}
