// Code by MisterIdle
// Try Pistol Hand on Itch.io: https://misteridle.itch.io/pistol-hand :D

package logic

import (
	"net/http"
	"strconv"
)

// getCommentsDataByPostID retrieves the comment data by post ID
func getCommentsDataByPostID(postID int, session Session) ([]Comments, error) {
	comments, err := fetchCommentsByPostID(postID)
	if err != nil {
		return nil, err
	}

	var commentsData []Comments
	for _, comment := range comments {
		commentData := Comments{
			CommentID:       comment.CommentID,
			PostID:          comment.PostID,
			Title:           comment.Title,
			Content:         comment.Content,
			Timestamp:       comment.Timestamp,
			Username:        getUsernameByCommentID(comment.CommentID),
			LikesComment:    getLikesByCommentID(comment.CommentID),
			DislikesComment: getDislikesByCommentID(comment.CommentID),
			Sessions:        session,
		}
		commentsData = append(commentsData, commentData)
	}

	return commentsData, nil
}

// CreateCommentHandler handles the creation of a new comment
func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.FormValue("post_id")
	content := r.FormValue("content")

	id, err := strconv.Atoi(postID)
	if err != nil {
		errorPage(w, r) // Display error page if post ID is invalid
		return
	}

	if containsAllHtmlTags(content) {
		reloadPageWithError(w, r, "HTML tags are not allowed") // Display error if content contains HTML tags
		return
	}

	newComment(id, content, getUsernameByUUID(getSessionUUID(r))) // Create new comment

	createLogs("Comment created") // Log the creation
	reloadPageWithoutError(w, r)  // Reload the page without error
}

// DeleteCommentHandler handles the deletion of a comment
func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	commentID := r.FormValue("comment_id")

	id, _ := strconv.Atoi(commentID)

	deleteComment(id)             // Delete the comment
	createLogs("Comment deleted") // Log the deletion
	reloadPageWithoutError(w, r)  // Reload the page without error
}

// LikeCommentHandler handles liking a comment
func LikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	commentID := r.FormValue("comment_id")

	id, err := strconv.Atoi(commentID)
	if err != nil {
		errorPage(w, r) // Display error page if comment ID is invalid
		return
	}

	if !isUserLoggedIn(r) {
		logginPage(w, r) // Redirect to login page if user is not logged in
		return
	}

	userID := getIDByUUID(getSessionUUID(r)) // Get user ID from session

	if hasUserDislikedComment(id, userID) {
		removeDislikeComment(id, userID) // Remove dislike if user has disliked the comment
	}

	if !hasUserLikedComment(id, userID) {
		newLikeComment(id, userID) // Add like if user hasn't liked the comment
	} else {
		removeLikeComment(id, userID) // Remove like if user has already liked the comment
	}

	createLogs("User " + getUsernameByUUID(getSessionUUID(r)) + " liked comment " + commentID)
	reloadPageWithoutError(w, r) // Reload the page without error
}

// DislikeCommentHandler handles disliking a comment
func DislikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	commentID := r.FormValue("comment_id")

	id, err := strconv.Atoi(commentID)
	if err != nil {
		errorPage(w, r) // Display error page if comment ID is invalid
		return
	}

	if !isUserLoggedIn(r) {
		logginPage(w, r) // Redirect to login page if user is not logged in
		return
	}

	userID := getIDByUUID(getSessionUUID(r)) // Get user ID from session

	if hasUserLikedComment(id, userID) {
		removeLikeComment(id, userID) // Remove like if user has liked the comment
	}

	if !hasUserDislikedComment(id, userID) {
		newDislikeComment(id, userID) // Add dislike if user hasn't disliked the comment
	} else {
		removeDislikeComment(id, userID) // Remove dislike if user has already disliked the comment
	}

	createLogs("User " + getUsernameByUUID(getSessionUUID(r)) + " disliked comment " + commentID)
	reloadPageWithoutError(w, r) // Reload the page without error
}
