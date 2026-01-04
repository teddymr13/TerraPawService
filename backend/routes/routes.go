package routes

import (
	h "github.com/TerraPaw/backend/handlers"
	"github.com/TerraPaw/backend/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Auth routes (public)
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/forgot-password", h.ForgotPassword)
		auth.POST("/reset-password", h.ResetPassword)
		auth.GET("/profile", middleware.AuthMiddleware(), h.GetUserProfile)
	}

	// Profile routes (User Personal Data)
	profile := router.Group("/api/profile")
	profile.Use(middleware.AuthMiddleware())
	{
		profile.GET("/pets", h.GetMyPets)
		profile.POST("/pets", h.CreateUserPet)
		profile.GET("/medical-records", h.GetMedicalRecords)
		profile.GET("/notifications", h.GetNotifications)
		profile.GET("/stats", h.GetUserStats) // New endpoint for profile stats
	}

	// Community routes
	community := router.Group("/api/community")
	community.Use(middleware.AuthMiddleware())
	{
		community.POST("/posts", h.CreatePost)
		community.GET("/posts", h.GetPosts)
		community.GET("/posts/:id", h.GetPost)
		community.POST("/posts/:id/like", h.LikePost)
		community.DELETE("/posts/:id/like", h.UnlikePost)
		community.POST("/posts/:id/bookmark", h.BookmarkPost)
		community.DELETE("/posts/:id/bookmark", h.UnbookmarkPost)
		community.POST("/posts/:id/share", h.SharePost)
		community.POST("/posts/:id/comments", h.CreateComment)
		community.POST("/comments/:id/like", h.LikeComment)
		community.DELETE("/comments/:id/like", h.UnlikeComment)
	}

	// Marketplace routes
	marketplace := router.Group("/api/marketplace")
	{
		marketplace.GET("/animals", h.GetAnimals)
		marketplace.GET("/animals/:id", h.GetAnimal)
		marketplace.GET("/animals/:id/reviews", h.GetReviews)
		marketplace.GET("/categories", h.GetCategories)
	}

	marketplaceProtected := router.Group("/api/marketplace")
	marketplaceProtected.Use(middleware.AuthMiddleware())
	{
		marketplaceProtected.POST("/animals", h.CreateAnimal)
		marketplaceProtected.POST("/orders", h.CreateOrder)
		marketplaceProtected.GET("/orders", h.GetOrders)

		// Wishlist
		marketplaceProtected.POST("/wishlist", h.AddToWishlist)
		marketplaceProtected.DELETE("/wishlist/:id", h.RemoveFromWishlist)
		marketplaceProtected.GET("/wishlist", h.GetWishlist)

		// Reviews
		marketplaceProtected.POST("/reviews", h.CreateReview)
	}

	// Consultation routes
	consultation := router.Group("/api/consultation")
	{
		consultation.GET("/veterinarians", h.GetVeterinarians)
		consultation.GET("/veterinarians/:id", h.GetVeterinarian)
	}

	consultationProtected := router.Group("/api/consultation")
	consultationProtected.Use(middleware.AuthMiddleware())
	{
		consultationProtected.POST("/veterinarians/register", h.RegisterVeterinarian)
		consultationProtected.POST("/consultations", h.CreateConsultation)
		consultationProtected.GET("/consultations", h.GetConsultations)
		consultationProtected.GET("/consultations/:id", h.GetConsultation)
		consultationProtected.PUT("/consultations/:id/status", h.UpdateConsultationStatus)
	}

	// Chat routes
	chat := router.Group("/api/chat")
	chat.Use(middleware.AuthMiddleware())
	{
		chat.POST("/messages", h.SendMessage)
		chat.GET("/messages", h.GetMessages)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Config routes (Public)
	config := router.Group("/api/config")
	{
		config.GET("/splash", h.GetCurrentSplash)
		// config.POST("/splash", h.CreateSplashEvent) // Admin only, uncomment if needed
	}
}
