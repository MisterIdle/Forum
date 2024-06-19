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
		errorPage(w, r)
		return
	}

	newComment(id, content, getUsernameByUUID(getSessionUUID(r)))

	reloadPageWithoutError(w, r)
}

// DeleteCommentHandler handles the deletion of a comment
func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	commentID := r.FormValue("comment_id")

	id, _ := strconv.Atoi(commentID)

	deleteComment(id)
	reloadPageWithoutError(w, r)
}

// LikeCommentHandler handles liking a comment
func LikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	commentID := r.FormValue("comment_id")

	id, err := strconv.Atoi(commentID)
	if err != nil {
		errorPage(w, r)
		return
	}

	if !isUserLoggedIn(r) {
		logginPage(w, r)
		return
	}

	userID := getIDByUUID(getSessionUUID(r))

	if hasUserDislikedComment(id, userID) {
		removeDislikeComment(id, userID)
	}

	if !hasUserLikedComment(id, userID) {
		newLikeComment(id, userID)
	} else {
		removeLikeComment(id, userID)
	}

	reloadPageWithoutError(w, r)
}

// DislikeCommentHandler handles disliking a comment
func DislikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	commentID := r.FormValue("comment_id")

	id, err := strconv.Atoi(commentID)
	if err != nil {
		errorPage(w, r)
		return
	}

	if !isUserLoggedIn(r) {
		logginPage(w, r)
		return
	}

	userID := getIDByUUID(getSessionUUID(r))

	if hasUserLikedComment(id, userID) {
		removeLikeComment(id, userID)
	}

	if !hasUserDislikedComment(id, userID) {
		newDislikeComment(id, userID)
	} else {
		removeDislikeComment(id, userID)
	}

	reloadPageWithoutError(w, r)
}
