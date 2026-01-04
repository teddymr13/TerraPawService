package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/TerraPaw/backend/config"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	cfg := config.LoadConfig()

	// Create connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	DB = db
	log.Println("Database connected successfully")

	// Create tables
	createTables()
}

func createTables() {
	// Users table
	createUserTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) UNIQUE NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		fullname VARCHAR(255),
		avatar_url VARCHAR(500),
		bio TEXT,
		user_type VARCHAR(50) DEFAULT 'customer',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// Posts table (for community)
	createPostsTable := `
	CREATE TABLE IF NOT EXISTS posts (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		content TEXT NOT NULL,
		image_url VARCHAR(500),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// Comments table
	createCommentsTable := `
	CREATE TABLE IF NOT EXISTS comments (
		id SERIAL PRIMARY KEY,
		post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// Comment Likes table
	createCommentLikesTable := `
	CREATE TABLE IF NOT EXISTS comment_likes (
		id SERIAL PRIMARY KEY,
		comment_id INTEGER NOT NULL REFERENCES comments(id) ON DELETE CASCADE,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(comment_id, user_id)
	);`

	// Likes table
	createLikesTable := `
	CREATE TABLE IF NOT EXISTS likes (
		id SERIAL PRIMARY KEY,
		post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(post_id, user_id)
	);`

	// Bookmarks table
	createBookmarksTable := `
    CREATE TABLE IF NOT EXISTS bookmarks (
        id SERIAL PRIMARY KEY,
        post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
        user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        UNIQUE(post_id, user_id)
    );`

	// Post Media table
	createPostMediaTable := `
    CREATE TABLE IF NOT EXISTS post_media (
        id SERIAL PRIMARY KEY,
        post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
        media_url VARCHAR(500) NOT NULL,
        media_type VARCHAR(20) DEFAULT 'image', -- 'image' or 'video'
        sort_order INTEGER DEFAULT 0,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`

	// Post Shares table
	createPostSharesTable := `
	CREATE TABLE IF NOT EXISTS post_shares (
		id SERIAL PRIMARY KEY,
		post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(post_id, user_id)
	);`

	// Animals/Pets table (for marketplace)
	createAnimalsTable := `
	CREATE TABLE IF NOT EXISTS animals (
		id SERIAL PRIMARY KEY,
		seller_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		animal_type VARCHAR(50) NOT NULL,
		breed VARCHAR(255),
		name VARCHAR(255) NOT NULL,
		age INTEGER,
		description TEXT,
		price DECIMAL(10, 2),
		image_url VARCHAR(500),
		location VARCHAR(500),
		rating DECIMAL(3, 2) DEFAULT 0,
		status VARCHAR(50) DEFAULT 'available',
        color VARCHAR(50),
        gender VARCHAR(20),
        stock INTEGER DEFAULT 1,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// Orders table
	createOrdersTable := `
	CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		buyer_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		animal_id INTEGER NOT NULL REFERENCES animals(id) ON DELETE CASCADE,
		total_price DECIMAL(10, 2),
		status VARCHAR(50) DEFAULT 'pending',
        quantity INTEGER DEFAULT 1,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// Veterinarians table
	createVetsTable := `
	CREATE TABLE IF NOT EXISTS veterinarians (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		clinic_name VARCHAR(255),
		license_number VARCHAR(255),
		specialization VARCHAR(255),
		phone VARCHAR(20),
		address VARCHAR(500),
		bio TEXT,
		rating DECIMAL(3, 2) DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// Consultations table
	createConsultationsTable := `
	CREATE TABLE IF NOT EXISTS consultations (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		veterinarian_id INTEGER NOT NULL REFERENCES veterinarians(id) ON DELETE CASCADE,
		pet_name VARCHAR(255),
		symptoms TEXT,
		consultation_type VARCHAR(50) DEFAULT 'online',
		status VARCHAR(50) DEFAULT 'pending',
		scheduled_at TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// Messages table
	createMessagesTable := `
	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		sender_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		receiver_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		content TEXT NOT NULL,
		is_read BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// Reminders table
	createRemindersTable := `
	CREATE TABLE IF NOT EXISTS reminders (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		date TIMESTAMP NOT NULL,
		is_completed BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// User Pets table (Profile -> My Pets)
	createUserPetsTable := `
    CREATE TABLE IF NOT EXISTS user_pets (
        id SERIAL PRIMARY KEY,
        owner_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        name VARCHAR(255) NOT NULL,
        animal_type VARCHAR(50),
        breed VARCHAR(255),
        age INTEGER,
        image_url VARCHAR(500),
        story TEXT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`

	// Medical Records table (Linked to User Pets)
	createMedicalRecordsTable := `
    CREATE TABLE IF NOT EXISTS medical_records (
        id SERIAL PRIMARY KEY,
        pet_id INTEGER NOT NULL REFERENCES user_pets(id) ON DELETE CASCADE,
        veterinarian_id INTEGER REFERENCES veterinarians(id), 
        record_type VARCHAR(50), -- 'vaccine', 'sickness', 'checkup'
        description TEXT,
        treatment VARCHAR(255),
        date TIMESTAMP NOT NULL,
        notes TEXT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`

	// Notifications table
	createNotificationsTable := `
    CREATE TABLE IF NOT EXISTS notifications (
        id SERIAL PRIMARY KEY,
        user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        title VARCHAR(255) NOT NULL,
        message TEXT,
        type VARCHAR(50), -- 'reminder', 'promo', 'order'
        is_read BOOLEAN DEFAULT FALSE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`

	// Wishlist table
	createWishlistsTable := `
    CREATE TABLE IF NOT EXISTS wishlists (
        id SERIAL PRIMARY KEY,
        user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        animal_id INTEGER NOT NULL REFERENCES animals(id) ON DELETE CASCADE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        UNIQUE(user_id, animal_id)
    );`

	// Reviews table
	createReviewsTable := `
    CREATE TABLE IF NOT EXISTS reviews (
        id SERIAL PRIMARY KEY,
        user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        order_id INTEGER REFERENCES orders(id), -- Optional link to order if we want to allow reviews without order strictly (but requirement says "produk yang sudah dibeli")
        animal_id INTEGER NOT NULL REFERENCES animals(id) ON DELETE CASCADE,
        rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
        comment TEXT,
        image_url VARCHAR(500),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`

	// Animal Media table (for multiple photos/videos per animal)
	createAnimalMediaTable := `
    CREATE TABLE IF NOT EXISTS animal_media (
        id SERIAL PRIMARY KEY,
        animal_id INTEGER NOT NULL REFERENCES animals(id) ON DELETE CASCADE,
        media_url VARCHAR(500) NOT NULL,
        media_type VARCHAR(20) DEFAULT 'image', -- 'image' or 'video'
        thumbnail_url VARCHAR(500), -- For videos
        sort_order INTEGER DEFAULT 0,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`

	// Splash Events table
	createSplashEventsTable := `
	CREATE TABLE IF NOT EXISTS splash_events (
		id SERIAL PRIMARY KEY,
		event_name VARCHAR(255) NOT NULL,
		image_url VARCHAR(500) NOT NULL,
		start_date TIMESTAMP NOT NULL,
		end_date TIMESTAMP NOT NULL,
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// Categories table
	createCategoriesTable := `
	CREATE TABLE IF NOT EXISTS categories (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL UNIQUE,
		icon VARCHAR(50),
		type VARCHAR(50) DEFAULT 'animal', -- 'animal' or 'food'
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	tables := []string{
		createUserTable,
		createCategoriesTable,
		createPostsTable,
		createCommentsTable,
		createCommentLikesTable,
		createLikesTable,
		createBookmarksTable,
		createPostMediaTable,
		createPostSharesTable,
		createAnimalsTable,
		createAnimalMediaTable,
		createOrdersTable,
		createVetsTable,
		createConsultationsTable,
		createMessagesTable,
		createRemindersTable,
		createUserPetsTable,
		createMedicalRecordsTable,
		createNotificationsTable,
		createWishlistsTable,
		createReviewsTable,
		createSplashEventsTable,
	}

	for _, tableSQL := range tables {
		if _, err := DB.Exec(tableSQL); err != nil {
			log.Printf("Error creating table: %v", err)
		}
	}

	// Migrations for existing tables
	migrations := []string{
		"ALTER TABLE animals ADD COLUMN IF NOT EXISTS color VARCHAR(50);",
		"ALTER TABLE animals ADD COLUMN IF NOT EXISTS gender VARCHAR(20);",
		"ALTER TABLE animals ADD COLUMN IF NOT EXISTS stock INTEGER DEFAULT 1;",
		"ALTER TABLE orders ADD COLUMN IF NOT EXISTS quantity INTEGER DEFAULT 1;",

		// Indexes for performance
		"CREATE INDEX IF NOT EXISTS idx_animals_type_status ON animals(animal_type, status);",
		"CREATE INDEX IF NOT EXISTS idx_animals_created_at ON animals(created_at DESC);",
		"CREATE INDEX IF NOT EXISTS idx_animals_price ON animals(price);",
		"CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at DESC);",
		"CREATE INDEX IF NOT EXISTS idx_veterinarians_rating ON veterinarians(rating DESC);",
		// Auth migrations
		"ALTER TABLE users ADD COLUMN IF NOT EXISTS reset_token VARCHAR(255);",
		"ALTER TABLE users ADD COLUMN IF NOT EXISTS reset_token_expiry TIMESTAMP;",

		// Community migrations
		"ALTER TABLE posts ADD COLUMN IF NOT EXISTS shares_count INTEGER DEFAULT 0;",
	}

	for _, migration := range migrations {
		if _, err := DB.Exec(migration); err != nil {
			log.Printf("Migration warning (might be already applied): %v", err)
		}
	}

	log.Println("Tables created successfully")
}
