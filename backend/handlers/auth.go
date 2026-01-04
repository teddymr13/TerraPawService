package handlers

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
    "strings"

	"github.com/TerraPaw/backend/db"
	"github.com/TerraPaw/backend/models"
	"github.com/TerraPaw/backend/utils"
	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	FullName string `json:"fullname"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err.Error()))
		return
	}
    
    // Normalize email
    req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// Hash password
	hash := sha256.Sum256([]byte(req.Password))
	hashedPassword := hex.EncodeToString(hash[:])

	// Insert user into database
	var userID int
	err := db.DB.QueryRow(
		"INSERT INTO users (username, email, password, fullname) VALUES ($1, $2, $3, $4) RETURNING id",
		req.Username, req.Email, hashedPassword, req.FullName,
	).Scan(&userID)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Registration failed", "Email or username already exists"))
		return
	}

	// Generate token
	token, err := utils.GenerateToken(userID, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Token generation failed", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse("Registration successful", gin.H{
		"user_id":    userID,
		"username":   req.Username,
		"email":      req.Email,
		"fullname":   req.FullName,
		"avatar_url": "", // Default empty or placeholder if set in DB default
		"user_type":  "customer",
		"token":      token,
	}))
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err.Error()))
		return
	}
    
    // Normalize email
    req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// Hash password
	hash := sha256.Sum256([]byte(req.Password))
	hashedPassword := hex.EncodeToString(hash[:])

    fmt.Printf("Login attempt: Email=%s Hash=%s\n", req.Email, hashedPassword)

	// Get user from database (case-insensitive email check for robustness)
	var user models.User
	err := db.DB.QueryRow(
		"SELECT id, username, email, fullname, COALESCE(avatar_url, ''), COALESCE(bio, ''), user_type FROM users WHERE LOWER(email) = LOWER($1) AND password = $2",
		req.Email, hashedPassword,
	).Scan(&user.ID, &user.Username, &user.Email, &user.FullName, &user.AvatarURL, &user.Bio, &user.UserType)

	if err != nil {
        fmt.Println("Login error:", err)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Login failed", "Invalid email or password"))
		} else {
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Login failed", err.Error()))
		}
		return
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Token generation failed", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Login successful", gin.H{
		"user_id":    user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"fullname":   user.FullName,
		"avatar_url": user.AvatarURL,
		"bio":        user.Bio,
		"user_type":  user.UserType,
		"token":      token,
	}))
}

func GetUserProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", "User ID not found in context"))
		return
	}

	var user models.User
	err := db.DB.QueryRow(
		"SELECT id, username, email, fullname, COALESCE(avatar_url, ''), COALESCE(bio, ''), user_type FROM users WHERE id = $1",
		userID,
	).Scan(&user.ID, &user.Username, &user.Email, &user.FullName, &user.AvatarURL, &user.Bio, &user.UserType)

	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("User not found", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("User profile retrieved", user))
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required"`
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err.Error()))
		return
	}

	// Check if user exists
	var email string
	err := db.DB.QueryRow("SELECT email FROM users WHERE email = $1", req.Email).Scan(&email)
	if err != nil {
		// Don't reveal if email exists or not for security, just say email sent if exists
		c.JSON(http.StatusOK, utils.SuccessResponse("If your email is registered, you will receive a password reset link.", nil))
		return
	}

	// Generate reset token (simple random string for demo)
	resetToken, _ := utils.GenerateToken(0, req.Email) // reusing JWT gen for convenience, or simple random string
    // Ideally use crypto/rand for a short code
    
    // For this demo, let's use a simple 6 digit code for easy testing
    resetToken = "123456" // In production, generate random: strconv.Itoa(rand.Intn(999999-100000)+100000)

	// Save token to DB with expiry (15 mins)
	_, err = db.DB.Exec(
		"UPDATE users SET reset_token = $1, reset_token_expiry = NOW() + INTERVAL '15 minutes' WHERE email = $2",
		resetToken, req.Email,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to process request", err.Error()))
		return
	}

	// Mock Send Email
	// In real app: sendEmail(req.Email, resetToken)
    // For demo purposes, we return the token in response so you can test it
	c.JSON(http.StatusOK, utils.SuccessResponse("Password reset code sent (Mock: Code is 123456)", gin.H{
        "mock_code": resetToken,
    }))
}

func ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err.Error()))
		return
	}

	// Verify token and expiry
	var id int
	err := db.DB.QueryRow(
		"SELECT id FROM users WHERE email = $1 AND reset_token = $2 AND reset_token_expiry > NOW()",
		req.Email, req.Token,
	).Scan(&id)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid or expired token", "Please request a new password reset"))
		return
	}

	// Hash new password
	hash := sha256.Sum256([]byte(req.NewPassword))
	hashedPassword := hex.EncodeToString(hash[:])

	// Update password and clear token
	_, err = db.DB.Exec(
		"UPDATE users SET password = $1, reset_token = NULL, reset_token_expiry = NULL WHERE id = $2",
		hashedPassword, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to reset password", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Password has been reset successfully", nil))
}
