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

type CreateAnimalRequest struct {
	AnimalType  string  `json:"animal_type" binding:"required"`
	Breed       string  `json:"breed"`
	Name        string  `json:"name" binding:"required"`
	Age         int     `json:"age"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required"`
	ImageURL    string  `json:"image_url"`
	Location    string  `json:"location"`
	Color       string  `json:"color"`
	Gender      string  `json:"gender"`
	Stock       int     `json:"stock"`
}

type CreateOrderRequest struct {
	AnimalID int `json:"animal_id" binding:"required"`
	Quantity int `json:"quantity"`
}

type CreateReviewRequest struct {
	OrderID  *int   `json:"order_id"`
	AnimalID int    `json:"animal_id" binding:"required"`
	Rating   int    `json:"rating" binding:"required,min=1,max=5"`
	Comment  string `json:"comment"`
	ImageURL string `json:"image_url"`
}

type AddWishlistRequest struct {
	AnimalID int `json:"animal_id" binding:"required"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
	Type string `json:"type"`
}

func GetCategories(c *gin.Context) {
	rows, err := db.DB.Query("SELECT id, name, icon, type FROM categories WHERE is_active = TRUE ORDER BY id ASC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch categories", err.Error()))
		return
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var cat Category
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.Icon, &cat.Type); err == nil {
			categories = append(categories, cat)
		}
	}

	// Fallback if empty (should not happen after seed)
	if len(categories) == 0 {
		categories = []Category{}
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Categories retrieved", categories))
}

func GetAnimals(c *gin.Context) {
	animalType := c.Query("animal_type")
	searchQuery := c.Query("search")
	breed := c.Query("breed")
	minPriceStr := c.Query("min_price")
	maxPriceStr := c.Query("max_price")
	sortBy := c.Query("sort") // price_asc, price_desc, newest, oldest

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var rows *sql.Rows
	var err error

	query := `SELECT a.id, a.seller_id, a.animal_type, COALESCE(a.breed, ''), a.name, a.age, COALESCE(a.description, ''), 
	                 a.price, COALESCE(a.image_url, ''), COALESCE(a.location, ''), a.rating, a.status, 
                     COALESCE(a.color, ''), COALESCE(a.gender, ''), COALESCE(a.stock, 0),
                     a.created_at, a.updated_at,
	                 u.id, u.username, u.email, u.fullname, COALESCE(u.avatar_url, ''), COALESCE(u.bio, '')
	          FROM animals a
	          LEFT JOIN users u ON a.seller_id = u.id
	          WHERE a.status = 'available'`

	if animalType != "" {
		query += fmt.Sprintf(" AND a.animal_type ILIKE '%%%s%%'", animalType)
	}

	if searchQuery != "" {
		query += fmt.Sprintf(" AND (a.name ILIKE '%%%s%%' OR a.breed ILIKE '%%%s%%')", searchQuery, searchQuery)
	}

	if breed != "" {
		query += fmt.Sprintf(" AND a.breed ILIKE '%%%s%%'", breed)
	}

	if minPriceStr != "" {
		fmt.Println("Received min_price:", minPriceStr)
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			query += fmt.Sprintf(" AND a.price >= %f", minPrice)
		} else {
			fmt.Println("Error parsing min_price:", err)
		}
	}

	if maxPriceStr != "" {
		fmt.Println("Received max_price:", maxPriceStr)
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			query += fmt.Sprintf(" AND a.price <= %f", maxPrice)
		} else {
			fmt.Println("Error parsing max_price:", err)
		}
	}

	fmt.Println("Executing Query:", query)

	orderBy := "a.created_at DESC"
	switch sortBy {
	case "price_asc":
		orderBy = "a.price ASC"
	case "price_desc":
		orderBy = "a.price DESC"
	case "oldest":
		orderBy = "a.created_at ASC"
	case "newest":
		orderBy = "a.created_at DESC"
	}

	query += fmt.Sprintf(" ORDER BY %s LIMIT %d OFFSET %d", orderBy, limit, offset)

	rows, err = db.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch animals", err.Error()))
		return
	}
	defer rows.Close()

	var animals []models.Animal
	for rows.Next() {
		var animal models.Animal
		var seller models.User
		err := rows.Scan(
			&animal.ID, &animal.SellerID, &animal.AnimalType, &animal.Breed, &animal.Name, &animal.Age,
			&animal.Description, &animal.Price, &animal.ImageURL, &animal.Location, &animal.Rating, &animal.Status,
			&animal.Color, &animal.Gender, &animal.Stock,
			&animal.CreatedAt, &animal.UpdatedAt,
			&seller.ID, &seller.Username, &seller.Email, &seller.FullName, &seller.AvatarURL, &seller.Bio,
		)
		if err == nil {
			animal.Seller = &seller
			animals = append(animals, animal)
		}
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Animals retrieved", animals))
}

func GetAnimal(c *gin.Context) {
	id := c.Param("id")
	animalID, _ := strconv.Atoi(id)

	var animal models.Animal
	var seller models.User

	err := db.DB.QueryRow(
		`SELECT a.id, a.seller_id, a.animal_type, COALESCE(a.breed, ''), a.name, a.age, COALESCE(a.description, ''), 
		        a.price, COALESCE(a.image_url, ''), COALESCE(a.location, ''), a.rating, a.status, 
                COALESCE(a.color, ''), COALESCE(a.gender, ''), COALESCE(a.stock, 0),
                a.created_at, a.updated_at,
		        u.id, u.username, u.email, u.fullname, COALESCE(u.avatar_url, ''), COALESCE(u.bio, '')
		FROM animals a
		LEFT JOIN users u ON a.seller_id = u.id
		WHERE a.id = $1`,
		animalID,
	).Scan(
		&animal.ID, &animal.SellerID, &animal.AnimalType, &animal.Breed, &animal.Name, &animal.Age,
		&animal.Description, &animal.Price, &animal.ImageURL, &animal.Location, &animal.Rating, &animal.Status,
		&animal.Color, &animal.Gender, &animal.Stock,
		&animal.CreatedAt, &animal.UpdatedAt,
		&seller.ID, &seller.Username, &seller.Email, &seller.FullName, &seller.AvatarURL, &seller.Bio,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Animal not found", ""))
		return
	}

	// Fetch Media
	var mediaList []models.AnimalMedia
	mediaRows, err := db.DB.Query("SELECT id, animal_id, media_url, media_type, COALESCE(thumbnail_url, ''), sort_order FROM animal_media WHERE animal_id = $1 ORDER BY sort_order ASC", animal.ID)
	if err == nil {
		defer mediaRows.Close()
		for mediaRows.Next() {
			var m models.AnimalMedia
			if err := mediaRows.Scan(&m.ID, &m.AnimalID, &m.MediaURL, &m.MediaType, &m.ThumbnailURL, &m.SortOrder); err == nil {
				mediaList = append(mediaList, m)
			}
		}
	}
	animal.Media = mediaList

	animal.Seller = &seller
	c.JSON(http.StatusOK, utils.SuccessResponse("Animal details", animal))
}

func CreateAnimal(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req CreateAnimalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err.Error()))
		return
	}

	// Default stock if 0
	if req.Stock <= 0 {
		req.Stock = 1
	}

	var animalID int
	err := db.DB.QueryRow(
		`INSERT INTO animals (seller_id, animal_type, breed, name, age, description, price, image_url, location, status, color, gender, stock)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 'available', $10, $11, $12) RETURNING id`,
		userID, req.AnimalType, req.Breed, req.Name, req.Age, req.Description, req.Price, req.ImageURL, req.Location, req.Color, req.Gender, req.Stock,
	).Scan(&animalID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to create animal listing", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse("Animal listing created", gin.H{"id": animalID}))
}

func CreateOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err.Error()))
		return
	}

	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	// Get animal price and stock
	var price float64
	var stock int
	err := db.DB.QueryRow("SELECT price, stock FROM animals WHERE id = $1", req.AnimalID).Scan(&price, &stock)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Animal not found", ""))
		return
	}

	// Validate stock
	if stock < req.Quantity {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Insufficient stock", ""))
		return
	}

	// Create order
	var orderID int
	err = db.DB.QueryRow(
		"INSERT INTO orders (buyer_id, animal_id, total_price, status, quantity) VALUES ($1, $2, $3, 'pending', $4) RETURNING id",
		userID, req.AnimalID, price*float64(req.Quantity), req.Quantity,
	).Scan(&orderID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to create order", err.Error()))
		return
	}

	// Update animal status and stock
	if stock-req.Quantity == 0 {
		db.DB.Exec("UPDATE animals SET status = 'sold', stock = stock - $1 WHERE id = $2", req.Quantity, req.AnimalID)
	} else {
		db.DB.Exec("UPDATE animals SET stock = stock - $1 WHERE id = $2", req.Quantity, req.AnimalID)
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse("Order created", gin.H{"id": orderID}))
}

func GetOrders(c *gin.Context) {
	userID, _ := c.Get("user_id")

	rows, err := db.DB.Query(
		`SELECT o.id, o.buyer_id, o.animal_id, o.total_price, o.status, o.quantity, o.created_at,
		        a.id, a.seller_id, a.animal_type, COALESCE(a.breed, ''), a.name, a.price, COALESCE(a.image_url, '')
		FROM orders o
		LEFT JOIN animals a ON o.animal_id = a.id
		WHERE o.buyer_id = $1
		ORDER BY o.created_at DESC`,
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch orders", err.Error()))
		return
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		var animal models.Animal
		// Handle quantity scan if it exists in DB, otherwise use default
		var quantity int
		err := rows.Scan(
			&order.ID, &order.BuyerID, &order.AnimalID, &order.TotalPrice, &order.Status, &quantity, &order.CreatedAt,
			&animal.ID, &animal.SellerID, &animal.AnimalType, &animal.Breed, &animal.Name, &animal.Price, &animal.ImageURL,
		)
		if err == nil {
			order.Quantity = quantity
			order.Animal = &animal
			orders = append(orders, order)
		}
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Orders retrieved", orders))
}

// Wishlist Handlers

func AddToWishlist(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req AddWishlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err.Error()))
		return
	}

	_, err := db.DB.Exec("INSERT INTO wishlists (user_id, animal_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", userID, req.AnimalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to add to wishlist", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Added to wishlist", nil))
}

func RemoveFromWishlist(c *gin.Context) {
	userID, _ := c.Get("user_id")
	animalID := c.Param("id")

	_, err := db.DB.Exec("DELETE FROM wishlists WHERE user_id = $1 AND animal_id = $2", userID, animalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to remove from wishlist", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Removed from wishlist", nil))
}

func GetWishlist(c *gin.Context) {
	userID, _ := c.Get("user_id")

	rows, err := db.DB.Query(`
		SELECT w.id, w.animal_id, w.created_at,
		       a.animal_type, a.name, a.price, COALESCE(a.image_url, ''), a.status, a.color, a.gender, a.stock
		FROM wishlists w
		JOIN animals a ON w.animal_id = a.id
		WHERE w.user_id = $1
		ORDER BY w.created_at DESC
	`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch wishlist", err.Error()))
		return
	}
	defer rows.Close()

	var wishlist []models.Wishlist
	for rows.Next() {
		var w models.Wishlist
		var a models.Animal
		err := rows.Scan(&w.ID, &w.AnimalID, &w.CreatedAt,
			&a.AnimalType, &a.Name, &a.Price, &a.ImageURL, &a.Status, &a.Color, &a.Gender, &a.Stock)
		if err == nil {
			w.Animal = &a
			w.UserID = userID.(int)
			wishlist = append(wishlist, w)
		}
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Wishlist retrieved", wishlist))
}

// Review Handlers

func CreateReview(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err.Error()))
		return
	}

	// Check if order exists (optional validation if order_id provided)
	if req.OrderID != nil {
		var status string
		err := db.DB.QueryRow("SELECT status FROM orders WHERE id = $1 AND buyer_id = $2", req.OrderID, userID).Scan(&status)
		if err != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse("Order not found or not owned by user", ""))
			return
		}
	}

	_, err := db.DB.Exec(`
        INSERT INTO reviews (user_id, order_id, animal_id, rating, comment, image_url)
        VALUES ($1, $2, $3, $4, $5, $6)
    `, userID, req.OrderID, req.AnimalID, req.Rating, req.Comment, req.ImageURL)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to submit review", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse("Review submitted", nil))
}

func GetReviews(c *gin.Context) {
	animalID := c.Param("id")

	rows, err := db.DB.Query(`
        SELECT r.id, r.user_id, r.rating, r.comment, COALESCE(r.image_url, ''), r.created_at,
               u.username, u.fullname, COALESCE(u.avatar_url, '')
        FROM reviews r
        JOIN users u ON r.user_id = u.id
        WHERE r.animal_id = $1
        ORDER BY r.created_at DESC
    `, animalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch reviews", err.Error()))
		return
	}
	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var r models.Review
		var u models.User
		err := rows.Scan(&r.ID, &r.UserID, &r.Rating, &r.Comment, &r.ImageURL, &r.CreatedAt,
			&u.Username, &u.FullName, &u.AvatarURL)
		if err == nil {
			r.User = &u
			reviews = append(reviews, r)
		}
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Reviews retrieved", reviews))
}
