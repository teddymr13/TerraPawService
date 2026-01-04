package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	FullName  string    `json:"fullname"`
	AvatarURL string    `json:"avatar_url"`
	Bio       string    `json:"bio"`
	UserType  string    `json:"user_type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Post struct {
	ID             int         `json:"id"`
	UserID         int         `json:"user_id"`
	Content        string      `json:"content"`
	ImageURL       string      `json:"image_url"`
	User           *User       `json:"user,omitempty"`
	Comments       []Comment   `json:"comments,omitempty"`
	Likes          int         `json:"likes"`
	Media          []PostMedia `json:"media,omitempty"`
	CommentsCount  int         `json:"comments_count"`
	SharesCount    int         `json:"shares_count"`
	BookmarksCount int         `json:"bookmarks_count"`
	IsLiked        bool        `json:"is_liked"`
	IsBookmarked   bool        `json:"is_bookmarked"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

type Comment struct {
	ID         int       `json:"id"`
	PostID     int       `json:"post_id"`
	UserID     int       `json:"user_id"`
	Content    string    `json:"content"`
	User       *User     `json:"user,omitempty"`
	LikesCount int       `json:"likes_count"`
	IsLiked    bool      `json:"is_liked"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type PostMedia struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id"`
	MediaURL  string    `json:"media_url"`
	MediaType string    `json:"media_type"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}

type AnimalMedia struct {
	ID           int       `json:"id"`
	AnimalID     int       `json:"animal_id"`
	MediaURL     string    `json:"media_url"`
	MediaType    string    `json:"media_type"` // 'image' or 'video'
	ThumbnailURL string    `json:"thumbnail_url"`
	SortOrder    int       `json:"sort_order"`
	CreatedAt    time.Time `json:"created_at"`
}

type Animal struct {
	ID          int           `json:"id"`
	SellerID    int           `json:"seller_id"`
	AnimalType  string        `json:"animal_type"`
	Breed       string        `json:"breed"`
	Name        string        `json:"name"`
	Age         int           `json:"age"`
	Description string        `json:"description"`
	Price       float64       `json:"price"`
	ImageURL    string        `json:"image_url"` // Main thumbnail
	Location    string        `json:"location"`
	Rating      float64       `json:"rating"`
	Status      string        `json:"status"`
	Color       string        `json:"color"`
	Gender      string        `json:"gender"`
	Stock       int           `json:"stock"`
	Media       []AnimalMedia `json:"media,omitempty"` // Added
	Seller      *User         `json:"seller,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

type Order struct {
	ID         int       `json:"id"`
	BuyerID    int       `json:"buyer_id"`
	AnimalID   int       `json:"animal_id"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"`
	Quantity   int       `json:"quantity"`
	Animal     *Animal   `json:"animal,omitempty"`
	Buyer      *User     `json:"buyer,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Veterinarian struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	ClinicName     string    `json:"clinic_name"`
	LicenseNumber  string    `json:"license_number"`
	Specialization string    `json:"specialization"`
	Phone          string    `json:"phone"`
	Address        string    `json:"address"`
	Bio            string    `json:"bio"`
	Rating         float64   `json:"rating"`
	User           *User     `json:"user,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Consultation struct {
	ID               int           `json:"id"`
	UserID           int           `json:"user_id"`
	VeterinarianID   int           `json:"veterinarian_id"`
	PetName          string        `json:"pet_name"`
	Symptoms         string        `json:"symptoms"`
	ConsultationType string        `json:"consultation_type"`
	Status           string        `json:"status"`
	ScheduledAt      *time.Time    `json:"scheduled_at"`
	User             *User         `json:"user,omitempty"`
	Veterinarian     *Veterinarian `json:"veterinarian,omitempty"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
}

type Message struct {
	ID         int       `json:"id"`
	SenderID   int       `json:"sender_id"`
	ReceiverID int       `json:"receiver_id"`
	Content    string    `json:"content"`
	IsRead     bool      `json:"is_read"`
	Sender     *User     `json:"sender,omitempty"`
	Receiver   *User     `json:"receiver,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type UserPet struct {
	ID         int       `json:"id"`
	OwnerID    int       `json:"owner_id"`
	Name       string    `json:"name"`
	AnimalType string    `json:"animal_type"`
	Breed      string    `json:"breed"`
	Age        int       `json:"age"`
	ImageURL   string    `json:"image_url"`
	Story      string    `json:"story"`
	CreatedAt  time.Time `json:"created_at"`
}

type MedicalRecord struct {
	ID             int           `json:"id"`
	PetID          int           `json:"pet_id"`
	VeterinarianID int           `json:"veterinarian_id"`
	RecordType     string        `json:"record_type"`
	Description    string        `json:"description"`
	Treatment      string        `json:"treatment"`
	Date           time.Time     `json:"date"`
	Notes          string        `json:"notes"`
	Pet            *UserPet      `json:"pet,omitempty"`
	Veterinarian   *Veterinarian `json:"veterinarian,omitempty"`
	CreatedAt      time.Time     `json:"created_at"`
}

type Notification struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Type      string    `json:"type"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

type SplashEvent struct {
	ID        int       `json:"id"`
	EventName string    `json:"event_name"`
	ImageURL  string    `json:"image_url"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type Wishlist struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	AnimalID  int       `json:"animal_id"`
	Animal    *Animal   `json:"animal,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Review struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	OrderID   *int      `json:"order_id,omitempty"`
	AnimalID  int       `json:"animal_id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	ImageURL  string    `json:"image_url"`
	User      *User     `json:"user,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
