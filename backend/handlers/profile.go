package handlers

import (
	"net/http"

	"github.com/TerraPaw/backend/db"
	"github.com/TerraPaw/backend/models"
	"github.com/gin-gonic/gin"
)

// GetMyPets returns all pets owned by the current user
func GetMyPets(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	rows, err := db.DB.Query("SELECT id, owner_id, name, animal_type, COALESCE(breed, ''), age, COALESCE(image_url, ''), COALESCE(story, ''), created_at FROM user_pets WHERE owner_id = $1 ORDER BY created_at DESC", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pets"})
		return
	}
	defer rows.Close()

	var pets []models.UserPet
	for rows.Next() {
		var p models.UserPet
		if err := rows.Scan(&p.ID, &p.OwnerID, &p.Name, &p.AnimalType, &p.Breed, &p.Age, &p.ImageURL, &p.Story, &p.CreatedAt); err != nil {
			continue
		}
		pets = append(pets, p)
	}

	if pets == nil {
		pets = []models.UserPet{}
	}

	c.JSON(http.StatusOK, gin.H{"data": pets})
}

// GetMedicalRecords returns all medical records for the current user's pets
func GetMedicalRecords(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	query := `
		SELECT mr.id, mr.pet_id, mr.veterinarian_id, mr.record_type, COALESCE(mr.description, ''), COALESCE(mr.treatment, ''), mr.date, COALESCE(mr.notes, ''), mr.created_at,
		       p.name, p.animal_type,
		       COALESCE(v.clinic_name, ''), u.fullname
		FROM medical_records mr
		JOIN user_pets p ON mr.pet_id = p.id
		LEFT JOIN veterinarians v ON mr.veterinarian_id = v.id
		LEFT JOIN users u ON v.user_id = u.id
		WHERE p.owner_id = $1
		ORDER BY mr.date DESC
	`

	rows, err := db.DB.Query(query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch medical records"})
		return
	}
	defer rows.Close()

	var records []gin.H
	for rows.Next() {
		var mr models.MedicalRecord
		var petName, petType, clinicName, vetName string
		var vetID *int // Handle null vet

		if err := rows.Scan(
			&mr.ID, &mr.PetID, &vetID, &mr.RecordType, &mr.Description, &mr.Treatment, &mr.Date, &mr.Notes, &mr.CreatedAt,
			&petName, &petType,
			&clinicName, &vetName,
		); err != nil {
			continue
		}

		// Construct response object
		record := gin.H{
			"id":          mr.ID,
			"pet_id":      mr.PetID,
			"pet_name":    petName,
			"pet_type":    petType,
			"record_type": mr.RecordType,
			"description": mr.Description,
			"treatment":   mr.Treatment,
			"date":        mr.Date,
			"notes":       mr.Notes,
			"doctor_name": vetName,
			"clinic_name": clinicName,
		}
		records = append(records, record)
	}

	if records == nil {
		records = []gin.H{}
	}

	c.JSON(http.StatusOK, gin.H{"data": records})
}

// GetNotifications returns notifications for the current user
func GetNotifications(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	rows, err := db.DB.Query("SELECT id, user_id, title, message, type, is_read, created_at FROM notifications WHERE user_id = $1 ORDER BY created_at DESC", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var n models.Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Message, &n.Type, &n.IsRead, &n.CreatedAt); err != nil {
			continue
		}
		notifications = append(notifications, n)
	}

	if notifications == nil {
		notifications = []models.Notification{}
	}

	c.JSON(http.StatusOK, gin.H{"data": notifications})
}

// CreateUserPet adds a new pet for the user
func CreateUserPet(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		Name       string `json:"name" binding:"required"`
		AnimalType string `json:"animal_type" binding:"required"`
		Breed      string `json:"breed"`
		Age        int    `json:"age"`
		ImageURL   string `json:"image_url"`
		Story      string `json:"story"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var petID int
	err := db.DB.QueryRow(`
		INSERT INTO user_pets (owner_id, name, animal_type, breed, age, image_url, story)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		userID, input.Name, input.AnimalType, input.Breed, input.Age, input.ImageURL, input.Story,
	).Scan(&petID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create pet"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Pet created successfully", "id": petID})
}

// GetUserStats returns statistics for the user profile (Pets, Orders, Vouchers)
func GetUserStats(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var petCount, orderCount int

	// Count Pets
	if err := db.DB.QueryRow("SELECT COUNT(*) FROM user_pets WHERE owner_id = $1", userID).Scan(&petCount); err != nil {
		petCount = 0
	}

	// Count Orders
	if err := db.DB.QueryRow("SELECT COUNT(*) FROM orders WHERE buyer_id = $1", userID).Scan(&orderCount); err != nil {
		orderCount = 0
	}

	// Voucher Count (Mock for now, as voucher system is not fully implemented)
	// In a real app, this would be SELECT COUNT(*) FROM user_vouchers WHERE user_id = ...
	voucherCount := 5

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"pets":     petCount,
			"orders":   orderCount,
			"vouchers": voucherCount,
		},
	})
}
