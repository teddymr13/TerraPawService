package handlers

import (
	"net/http"
	"strconv"

	"github.com/TerraPaw/backend/db"
	"github.com/TerraPaw/backend/models"
	"github.com/TerraPaw/backend/utils"
	"github.com/gin-gonic/gin"
)

type SendMessageRequest struct {
	ReceiverID int    `json:"receiver_id" binding:"required"`
	Content    string `json:"content" binding:"required"`
}

// SendMessage sends a message from current user to another user
func SendMessage(c *gin.Context) {
	senderID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", "User ID not found"))
		return
	}

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err.Error()))
		return
	}

	var messageID int
	err := db.DB.QueryRow(
		`INSERT INTO messages (sender_id, receiver_id, content) VALUES ($1, $2, $3) RETURNING id`,
		senderID, req.ReceiverID, req.Content,
	).Scan(&messageID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to send message", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse("Message sent", gin.H{"id": messageID}))
}

// GetMessages retrieves messages between current user and another user (or all messages if no partner specified, but usually we filter by partner)
// For simplicity, let's allow getting messages for a specific partner_id query param
func GetMessages(c *gin.Context) {
	currentUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", "User ID not found"))
		return
	}

	partnerIDStr := c.Query("partner_id")
	if partnerIDStr == "" {
		// If no partner specified, maybe return list of conversations?
		// For now, let's require partner_id to fetch chat history
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("partner_id required", ""))
		return
	}

	partnerID, _ := strconv.Atoi(partnerIDStr)

	rows, err := db.DB.Query(
		`SELECT m.id, m.sender_id, m.receiver_id, m.content, m.is_read, m.created_at,
		        s.fullname as sender_name, COALESCE(s.avatar_url, '') as sender_avatar,
		        r.fullname as receiver_name, COALESCE(r.avatar_url, '') as receiver_avatar
		 FROM messages m
		 JOIN users s ON m.sender_id = s.id
		 JOIN users r ON m.receiver_id = r.id
		 WHERE (m.sender_id = $1 AND m.receiver_id = $2) 
		    OR (m.sender_id = $2 AND m.receiver_id = $1)
		 ORDER BY m.created_at ASC`,
		currentUserID, partnerID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch messages", err.Error()))
		return
	}
	defer rows.Close()

	var messages []gin.H
	for rows.Next() {
		var m models.Message
		var senderName, senderAvatar, receiverName, receiverAvatar string

		err := rows.Scan(
			&m.ID, &m.SenderID, &m.ReceiverID, &m.Content, &m.IsRead, &m.CreatedAt,
			&senderName, &senderAvatar, &receiverName, &receiverAvatar,
		)

		if err == nil {
			msg := gin.H{
				"id":          m.ID,
				"sender_id":   m.SenderID,
				"receiver_id": m.ReceiverID,
				"content":     m.Content,
				"is_read":     m.IsRead,
				"created_at":  m.CreatedAt,
				"sender": gin.H{
					"name":   senderName,
					"avatar": senderAvatar,
				},
			}
			messages = append(messages, msg)
		}
	}

	if messages == nil {
		messages = []gin.H{}
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Messages retrieved", messages))
}
