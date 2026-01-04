package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/TerraPaw/backend/db"
	"github.com/TerraPaw/backend/models"
	"github.com/TerraPaw/backend/utils"
	"github.com/gin-gonic/gin"
)

// GetCurrentSplash returns the active splash screen based on current date
func GetCurrentSplash(c *gin.Context) {
	var splash models.SplashEvent
	
	// Query to find active splash event where current date is between start and end date
	// Using CURRENT_TIMESTAMP from DB or passing Go time
	query := `
		SELECT id, event_name, image_url, start_date, end_date, is_active 
		FROM splash_events 
		WHERE is_active = TRUE 
		AND start_date <= NOW() 
		AND end_date >= NOW() 
		ORDER BY start_date DESC 
		LIMIT 1`

	err := db.DB.QueryRow(query).Scan(
		&splash.ID,
		&splash.EventName,
		&splash.ImageURL,
		&splash.StartDate,
		&splash.EndDate,
		&splash.IsActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// No active event, return default null or specific code
			c.JSON(http.StatusOK, utils.SuccessResponse("No active event splash", nil))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch splash event", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Active splash event found", splash))
}

// CreateSplashEvent (Admin only - simplified for now)
func CreateSplashEvent(c *gin.Context) {
	var input struct {
		EventName string `json:"event_name" binding:"required"`
		ImageURL  string `json:"image_url" binding:"required"`
		StartDate string `json:"start_date" binding:"required"` // Format: YYYY-MM-DD
		EndDate   string `json:"end_date" binding:"required"`   // Format: YYYY-MM-DD
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err.Error()))
		return
	}

	// Parse dates
	start, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid start_date format (YYYY-MM-DD)", err.Error()))
		return
	}
	end, err := time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid end_date format (YYYY-MM-DD)", err.Error()))
		return
	}
	
	// Set end date to end of day
	end = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	query := `
		INSERT INTO splash_events (event_name, image_url, start_date, end_date) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id`

	var id int
	err = db.DB.QueryRow(query, input.EventName, input.ImageURL, start, end).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to create splash event", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse("Splash event created", gin.H{"id": id}))
}
