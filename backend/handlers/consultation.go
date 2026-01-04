package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/TerraPaw/backend/db"
	"github.com/TerraPaw/backend/models"
	"github.com/TerraPaw/backend/utils"
	"github.com/gin-gonic/gin"
)

type CreateConsultationRequest struct {
	VeterinarianID   int        `json:"veterinarian_id" binding:"required"`
	PetName          string     `json:"pet_name" binding:"required"`
	Symptoms         string     `json:"symptoms" binding:"required"`
	ConsultationType string     `json:"consultation_type"`
	ScheduledAt      *time.Time `json:"scheduled_at"`
}

type VeterinarianRegistrationRequest struct {
	ClinicName     string `json:"clinic_name" binding:"required"`
	LicenseNumber  string `json:"license_number" binding:"required"`
	Specialization string `json:"specialization"`
	Phone          string `json:"phone"`
	Address        string `json:"address"`
	Bio            string `json:"bio"`
}

func RegisterVeterinarian(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", "User ID not found"))
		return
	}

	var req VeterinarianRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err.Error()))
		return
	}

	var vetID int
	err := db.DB.QueryRow(
		`INSERT INTO veterinarians (user_id, clinic_name, license_number, specialization, phone, address, bio)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		userID, req.ClinicName, req.LicenseNumber, req.Specialization, req.Phone, req.Address, req.Bio,
	).Scan(&vetID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to register veterinarian", err.Error()))
		return
	}

	// Update user type
	db.DB.Exec("UPDATE users SET user_type = 'veterinarian' WHERE id = $1", userID)

	c.JSON(http.StatusCreated, utils.SuccessResponse("Veterinarian registered", gin.H{"id": vetID}))
}

func GetVeterinarians(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	pageNum, _ := strconv.Atoi(page)
	limitNum, _ := strconv.Atoi(limit)
	offset := (pageNum - 1) * limitNum

	rows, err := db.DB.Query(
		`SELECT v.id, v.user_id, v.clinic_name, v.license_number, v.specialization, v.phone, 
		        v.address, v.bio, v.rating, v.created_at, v.updated_at,
		        u.id, u.username, u.email, u.fullname, u.avatar_url, u.bio
		FROM veterinarians v
		LEFT JOIN users u ON v.user_id = u.id
		ORDER BY v.rating DESC
		LIMIT $1 OFFSET $2`,
		limitNum, offset,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch veterinarians", err.Error()))
		return
	}
	defer rows.Close()

	var vets []models.Veterinarian
	for rows.Next() {
		var vet models.Veterinarian
		var user models.User

		err := rows.Scan(
			&vet.ID, &vet.UserID, &vet.ClinicName, &vet.LicenseNumber, &vet.Specialization, &vet.Phone,
			&vet.Address, &vet.Bio, &vet.Rating, &vet.CreatedAt, &vet.UpdatedAt,
			&user.ID, &user.Username, &user.Email, &user.FullName, &user.AvatarURL, &user.Bio,
		)

		if err == nil {
			vet.User = &user
			vets = append(vets, vet)
		}
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Veterinarians retrieved", vets))
}

func GetVeterinarian(c *gin.Context) {
	vetID := c.Param("id")

	var vet models.Veterinarian
	var user models.User

	err := db.DB.QueryRow(
		`SELECT v.id, v.user_id, v.clinic_name, v.license_number, v.specialization, v.phone,
		        v.address, v.bio, v.rating, v.created_at, v.updated_at,
		        u.id, u.username, u.email, u.fullname, u.avatar_url, u.bio
		FROM veterinarians v
		LEFT JOIN users u ON v.user_id = u.id
		WHERE v.id = $1`,
		vetID,
	).Scan(
		&vet.ID, &vet.UserID, &vet.ClinicName, &vet.LicenseNumber, &vet.Specialization, &vet.Phone,
		&vet.Address, &vet.Bio, &vet.Rating, &vet.CreatedAt, &vet.UpdatedAt,
		&user.ID, &user.Username, &user.Email, &user.FullName, &user.AvatarURL, &user.Bio,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, utils.ErrorResponse("Veterinarian not found", ""))
		} else {
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch veterinarian", err.Error()))
		}
		return
	}

	vet.User = &user
	c.JSON(http.StatusOK, utils.SuccessResponse("Veterinarian retrieved", vet))
}

func CreateConsultation(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", "User ID not found"))
		return
	}

	var req CreateConsultationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err.Error()))
		return
	}

	if req.ConsultationType == "" {
		req.ConsultationType = "online"
	}

	var consultationID int
	err := db.DB.QueryRow(
		`INSERT INTO consultations (user_id, veterinarian_id, pet_name, symptoms, consultation_type, status, scheduled_at)
		VALUES ($1, $2, $3, $4, $5, 'pending', $6) RETURNING id`,
		userID, req.VeterinarianID, req.PetName, req.Symptoms, req.ConsultationType, req.ScheduledAt,
	).Scan(&consultationID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to create consultation", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse("Consultation created", gin.H{"id": consultationID}))
}

func GetConsultations(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", "User ID not found"))
		return
	}

	rows, err := db.DB.Query(
		`SELECT co.id, co.user_id, co.veterinarian_id, co.pet_name, co.symptoms, co.consultation_type,
		        co.status, co.scheduled_at, co.created_at, co.updated_at,
		        v.id, v.clinic_name, v.specialization, v.phone
		FROM consultations co
		LEFT JOIN veterinarians v ON co.veterinarian_id = v.id
		WHERE co.user_id = $1
		ORDER BY co.created_at DESC`,
		userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch consultations", err.Error()))
		return
	}
	defer rows.Close()

	var consultations []models.Consultation
	for rows.Next() {
		var consultation models.Consultation
		var vet models.Veterinarian

		err := rows.Scan(
			&consultation.ID, &consultation.UserID, &consultation.VeterinarianID, &consultation.PetName,
			&consultation.Symptoms, &consultation.ConsultationType, &consultation.Status, &consultation.ScheduledAt,
			&consultation.CreatedAt, &consultation.UpdatedAt,
			&vet.ID, &vet.ClinicName, &vet.Specialization, &vet.Phone,
		)

		if err == nil {
			consultation.Veterinarian = &vet
			consultations = append(consultations, consultation)
		}
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Consultations retrieved", consultations))
}

func GetConsultation(c *gin.Context) {
	consultationID := c.Param("id")

	var consultation models.Consultation
	var vet models.Veterinarian
	var user models.User

	err := db.DB.QueryRow(
		`SELECT co.id, co.user_id, co.veterinarian_id, co.pet_name, co.symptoms, co.consultation_type,
		        co.status, co.scheduled_at, co.created_at, co.updated_at,
		        v.id, v.clinic_name, v.specialization, v.phone,
		        u.id, u.username, u.email, u.fullname
		FROM consultations co
		LEFT JOIN veterinarians v ON co.veterinarian_id = v.id
		LEFT JOIN users u ON co.user_id = u.id
		WHERE co.id = $1`,
		consultationID,
	).Scan(
		&consultation.ID, &consultation.UserID, &consultation.VeterinarianID, &consultation.PetName,
		&consultation.Symptoms, &consultation.ConsultationType, &consultation.Status, &consultation.ScheduledAt,
		&consultation.CreatedAt, &consultation.UpdatedAt,
		&vet.ID, &vet.ClinicName, &vet.Specialization, &vet.Phone,
		&user.ID, &user.Username, &user.Email, &user.FullName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, utils.ErrorResponse("Consultation not found", ""))
		} else {
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch consultation", err.Error()))
		}
		return
	}

	consultation.Veterinarian = &vet
	consultation.User = &user
	c.JSON(http.StatusOK, utils.SuccessResponse("Consultation retrieved", consultation))
}

func UpdateConsultationStatus(c *gin.Context) {
	consultationID := c.Param("id")

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err.Error()))
		return
	}

	_, err := db.DB.Exec(
		"UPDATE consultations SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2",
		req.Status, consultationID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to update consultation", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Consultation updated", nil))
}
