package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/TerraPaw/backend/db"
	"github.com/TerraPaw/backend/models"
	"github.com/TerraPaw/backend/utils"
	"github.com/gin-gonic/gin"
)

type CreatePostRequest struct {
	Content string `json:"content" binding:"required"`
	Media   []struct {
		MediaURL  string `json:"media_url" binding:"required"`
		MediaType string `json:"media_type" binding:"required"`
	} `json:"media"`
}

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

func CreatePost(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", "User ID not found"))
		return
	}

	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err.Error()))
		return
	}

	var postID int
	// Start transaction
	tx, err := db.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Database error", err.Error()))
		return
	}

	err = tx.QueryRow(
		"INSERT INTO posts (user_id, content) VALUES ($1, $2) RETURNING id",
		userID, req.Content,
	).Scan(&postID)

	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to create post", err.Error()))
		return
	}

	// Insert Media
	for i, m := range req.Media {
		_, err = tx.Exec(
			"INSERT INTO post_media (post_id, media_url, media_type, sort_order) VALUES ($1, $2, $3, $4)",
			postID, m.MediaURL, m.MediaType, i,
		)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to save media", err.Error()))
			return
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to commit transaction", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse("Post created successfully", gin.H{"id": postID}))
}

func GetPosts(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	// Get current user ID for is_liked check
	var currentUserID int
	if val, exists := c.Get("user_id"); exists {
		if id, ok := val.(int); ok {
			currentUserID = id
		}
	}

	pageNum, _ := strconv.Atoi(page)
	limitNum, _ := strconv.Atoi(limit)
	offset := (pageNum - 1) * limitNum

	rows, err := db.DB.Query(
		`SELECT p.id, p.user_id, p.content, COALESCE(p.image_url, ''), p.created_at, p.updated_at,
		        COALESCE(u.id, 0), COALESCE(u.username, 'Unknown'), COALESCE(u.email, ''), COALESCE(u.fullname, ''), COALESCE(u.avatar_url, ''), COALESCE(u.bio, ''),
		        COUNT(DISTINCT l.user_id) as like_count,
                (SELECT COUNT(*) FROM comments WHERE post_id = p.id) as comment_count,
                COUNT(DISTINCT b.user_id) as bookmark_count,
                COUNT(DISTINCT s.user_id) as shares_count,
                EXISTS(SELECT 1 FROM likes WHERE post_id = p.id AND user_id = $3) as is_liked,
                EXISTS(SELECT 1 FROM bookmarks WHERE post_id = p.id AND user_id = $3) as is_bookmarked
		FROM posts p
		LEFT JOIN users u ON p.user_id = u.id
		LEFT JOIN likes l ON p.id = l.post_id
        LEFT JOIN bookmarks b ON p.id = b.post_id
        LEFT JOIN post_shares s ON p.id = s.post_id
		GROUP BY p.id, u.id
		ORDER BY p.created_at DESC
		LIMIT $1 OFFSET $2`,
		limitNum, offset, currentUserID,
	)

	if err != nil {
		fmt.Printf("Error querying posts: %v\n", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch posts", err.Error()))
		return
	}
	defer rows.Close()

	var posts []models.Post
	var postIDs []int

	for rows.Next() {
		var post models.Post
		var user models.User
		var likeCount int
		var commentCount int
		var bookmarkCount int
		var sharesCount int
		var isLiked bool
		var isBookmarked bool

		err := rows.Scan(
			&post.ID, &post.UserID, &post.Content, &post.ImageURL, &post.CreatedAt, &post.UpdatedAt,
			&user.ID, &user.Username, &user.Email, &user.FullName, &user.AvatarURL, &user.Bio,
			&likeCount, &commentCount, &bookmarkCount, &sharesCount, &isLiked, &isBookmarked,
		)

		if err != nil {
			fmt.Printf("Error scanning post row: %v\n", err)
			continue
		}

		post.User = &user
		post.Likes = likeCount
		post.CommentsCount = commentCount
		post.BookmarksCount = bookmarkCount
		post.SharesCount = sharesCount
		post.IsLiked = isLiked
		post.IsBookmarked = isBookmarked
		post.Media = []models.PostMedia{} // Initialize empty slice

		posts = append(posts, post)
		postIDs = append(postIDs, post.ID)
	}

	// Fetch Media for these posts if any
	if len(postIDs) > 0 {
		// Simple N+1 for now is fine or a single query with IN
		// For simplicity and to avoid dynamic SQL construction for IN clause in Go sql driver without helper:
		// We will just query all media for these posts.
		// Or simpler: iterate and query. For 10 posts, 10 queries is fast enough for now.
		// Better: Use a loop.

		for i := range posts {
			mediaRows, err := db.DB.Query(
				"SELECT id, post_id, media_url, media_type, sort_order FROM post_media WHERE post_id = $1 ORDER BY sort_order ASC",
				posts[i].ID,
			)
			if err == nil {
				for mediaRows.Next() {
					var pm models.PostMedia
					if err := mediaRows.Scan(&pm.ID, &pm.PostID, &pm.MediaURL, &pm.MediaType, &pm.SortOrder); err == nil {
						posts[i].Media = append(posts[i].Media, pm)
					}
				}
				mediaRows.Close()
			}
		}
	}

	fmt.Printf("Returning %d posts\n", len(posts))
	if len(posts) > 0 {
		fmt.Printf("First post ID: %d, User: %s, Shares: %d\n", posts[0].ID, posts[0].User.Username, posts[0].SharesCount)
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Posts retrieved", posts))
}

func GetPost(c *gin.Context) {
	postID := c.Param("id")

	// Get current user ID for is_liked check
	var currentUserID int
	if val, exists := c.Get("user_id"); exists {
		if id, ok := val.(int); ok {
			currentUserID = id
		}
	}

	var post models.Post
	var user models.User
	var likeCount int
	var commentCount int
	var bookmarkCount int
	var sharesCount int
	var isLiked bool
	var isBookmarked bool

	err := db.DB.QueryRow(
		`SELECT p.id, p.user_id, p.content, COALESCE(p.image_url, ''), p.created_at, p.updated_at,
		        COALESCE(u.id, 0), COALESCE(u.username, 'Unknown'), COALESCE(u.email, ''), COALESCE(u.fullname, ''), COALESCE(u.avatar_url, ''), COALESCE(u.bio, ''),
		        COUNT(DISTINCT l.user_id) as like_count,
                (SELECT COUNT(*) FROM comments WHERE post_id = p.id) as comment_count,
                COUNT(DISTINCT b.user_id) as bookmark_count,
                COUNT(DISTINCT s.user_id) as shares_count,
                EXISTS(SELECT 1 FROM likes WHERE post_id = p.id AND user_id = $2) as is_liked,
                EXISTS(SELECT 1 FROM bookmarks WHERE post_id = p.id AND user_id = $2) as is_bookmarked
		FROM posts p
		LEFT JOIN users u ON p.user_id = u.id
		LEFT JOIN likes l ON p.id = l.post_id
        LEFT JOIN bookmarks b ON p.id = b.post_id
        LEFT JOIN post_shares s ON p.id = s.post_id
		WHERE p.id = $1
		GROUP BY p.id, u.id`,
		postID, currentUserID,
	).Scan(
		&post.ID, &post.UserID, &post.Content, &post.ImageURL, &post.CreatedAt, &post.UpdatedAt,
		&user.ID, &user.Username, &user.Email, &user.FullName, &user.AvatarURL, &user.Bio,
		&likeCount, &commentCount, &bookmarkCount, &sharesCount, &isLiked, &isBookmarked,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, utils.ErrorResponse("Post not found", ""))
		} else {
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch post", err.Error()))
		}
		return
	}

	post.User = &user
	post.Likes = likeCount
	post.CommentsCount = commentCount
	post.BookmarksCount = bookmarkCount
	post.SharesCount = sharesCount
	post.IsLiked = isLiked
	post.IsBookmarked = isBookmarked
	post.Media = []models.PostMedia{}

	// Fetch Media
	mediaRows, err := db.DB.Query(
		"SELECT id, post_id, media_url, media_type, sort_order FROM post_media WHERE post_id = $1 ORDER BY sort_order ASC",
		post.ID,
	)
	if err == nil {
		for mediaRows.Next() {
			var pm models.PostMedia
			if err := mediaRows.Scan(&pm.ID, &pm.PostID, &pm.MediaURL, &pm.MediaType, &pm.SortOrder); err == nil {
				post.Media = append(post.Media, pm)
			}
		}
		mediaRows.Close()
	}

	// Get comments (detailed list)
	commentRows, err := db.DB.Query(
		`SELECT c.id, c.post_id, c.user_id, c.content, c.created_at,
		        u.id, u.username, u.email, COALESCE(u.fullname, ''), COALESCE(u.avatar_url, ''), COALESCE(u.bio, ''),
                (SELECT COUNT(*) FROM comment_likes WHERE comment_id = c.id) as likes_count,
                EXISTS(SELECT 1 FROM comment_likes WHERE comment_id = c.id AND user_id = $2) as is_liked
		FROM comments c
		INNER JOIN users u ON c.user_id = u.id
		WHERE c.post_id = $1
		ORDER BY c.created_at DESC`,
		post.ID, currentUserID,
	)

	if err == nil {
		defer commentRows.Close()
		for commentRows.Next() {
			var comment models.Comment
			var commentUser models.User
			var likesCount int
			var isLiked bool

			err := commentRows.Scan(
				&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt,
				&commentUser.ID, &commentUser.Username, &commentUser.Email, &commentUser.FullName,
				&commentUser.AvatarURL, &commentUser.Bio,
				&likesCount, &isLiked,
			)
			if err == nil {
				comment.User = &commentUser
				comment.LikesCount = likesCount
				comment.IsLiked = isLiked
				post.Comments = append(post.Comments, comment)
			}
		}
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Post retrieved", post))
}

func LikePost(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", "User ID not found"))
		return
	}

	postID := c.Param("id")

	_, err := db.DB.Exec(
		"INSERT INTO likes (post_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
		postID, userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to like post", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Post liked", nil))
}

func UnlikePost(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", "User ID not found"))
		return
	}

	postID := c.Param("id")

	_, err := db.DB.Exec(
		"DELETE FROM likes WHERE post_id = $1 AND user_id = $2",
		postID, userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to unlike post", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Post unliked", nil))
}

func CreateComment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", "User ID not found"))
		return
	}

	postID := c.Param("id")

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err.Error()))
		return
	}

	var commentID int
	err := db.DB.QueryRow(
		"INSERT INTO comments (post_id, user_id, content) VALUES ($1, $2, $3) RETURNING id",
		postID, userID, req.Content,
	).Scan(&commentID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to create comment", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse("Comment created", gin.H{"id": commentID}))
}

func BookmarkPost(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", "User ID not found"))
		return
	}

	postID := c.Param("id")

	_, err := db.DB.Exec(
		"INSERT INTO bookmarks (post_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
		postID, userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to bookmark post", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Post bookmarked", nil))
}

func UnbookmarkPost(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", "User ID not found"))
		return
	}

	postID := c.Param("id")

	_, err := db.DB.Exec(
		"DELETE FROM bookmarks WHERE post_id = $1 AND user_id = $2",
		postID, userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to unbookmark post", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Post unbookmarked", nil))
}

func SharePost(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", "User ID not found"))
		return
	}

	postID := c.Param("id")

	_, err := db.DB.Exec(
		"INSERT INTO post_shares (post_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
		postID, userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to record share", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Share recorded", nil))
}

func LikeComment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", "User ID not found"))
		return
	}

	commentID := c.Param("id")

	_, err := db.DB.Exec(
		"INSERT INTO comment_likes (comment_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
		commentID, userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to like comment", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Comment liked", nil))
}

func UnlikeComment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", "User ID not found"))
		return
	}

	commentID := c.Param("id")

	_, err := db.DB.Exec(
		"DELETE FROM comment_likes WHERE comment_id = $1 AND user_id = $2",
		commentID, userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to unlike comment", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Comment unliked", nil))
}
