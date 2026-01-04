package db

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func SeedData() {
	log.Println("Starting comprehensive data patching...")

	// Password hash
	password, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	passStr := string(password)

	// --- 1. USERS ---
	// Create base users if they don't exist
	baseUsers := []struct {
		Username string
		Email    string
		Fullname string
		Type     string
	}{
		{"anna", "anna@example.com", "Anna", "customer"},
		{"mia", "mia@example.com", "Mia", "customer"},
		{"budi", "budi@example.com", "Budi", "customer"},
		{"siti", "siti@example.com", "Siti", "customer"},
		{"dr_parlian", "doc@example.com", "Dr. Parlian", "veterinarian"},
		{"dr_sarah", "sarah@example.com", "Dr. Sarah", "veterinarian"},
		{"dr_john", "john@example.com", "Dr. John", "veterinarian"},
	}

	for _, u := range baseUsers {
		_, err := DB.Exec(`INSERT INTO users (username, email, password, fullname, user_type, avatar_url, bio) 
			VALUES ($1, $2, $3, $4, $5, $6, $7) 
			ON CONFLICT (username) DO UPDATE SET fullname = EXCLUDED.fullname`,
			u.Username, u.Email, passStr, u.Fullname, u.Type, "https://via.placeholder.com/150", "Bio for "+u.Fullname)
		if err != nil {
			log.Printf("Error seeding user %s: %v", u.Username, err)
		}
	}

	// Patch up to 20 random users
	for i := 1; i <= 20; i++ {
		username := fmt.Sprintf("user%d", i)
		email := fmt.Sprintf("user%d@example.com", i)
		fullname := fmt.Sprintf("User %d", i)
		_, _ = DB.Exec(`INSERT INTO users (username, email, password, fullname, user_type, avatar_url, bio) 
			VALUES ($1, $2, $3, $4, 'customer', $5, $6) ON CONFLICT DO NOTHING`,
			username, email, passStr, fullname, "https://via.placeholder.com/150", "Auto generated bio")
	}

	// Get all User IDs
	rows, _ := DB.Query("SELECT id, user_type FROM users")
	var userIDs []int
	var vetUserIDs []int
	var customerIDs []int
	for rows.Next() {
		var id int
		var uType string
		rows.Scan(&id, &uType)
		userIDs = append(userIDs, id)
		if uType == "veterinarian" {
			vetUserIDs = append(vetUserIDs, id)
		} else {
			customerIDs = append(customerIDs, id)
		}
	}
	rows.Close()

	// --- 2. VETERINARIANS ---
	// Ensure every vet user has a vet profile
	for _, uid := range vetUserIDs {
		var count int
		DB.QueryRow("SELECT COUNT(*) FROM veterinarians WHERE user_id = $1", uid).Scan(&count)
		if count == 0 {
			_, _ = DB.Exec(`INSERT INTO veterinarians (user_id, clinic_name, license_number, specialization, phone, address, bio, rating) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
				uid, "Klinik Sehat "+fmt.Sprint(uid), fmt.Sprintf("LIC-%d", uid), "General Vet", "08123456789", "Jl. Vet No. "+fmt.Sprint(uid), "Experienced Vet", 4.5+rand.Float64()*0.5)
		}
	}

	// Get Vet IDs
	var vetIDs []int
	rows, _ = DB.Query("SELECT id FROM veterinarians")
	for rows.Next() {
		var id int
		rows.Scan(&id)
		vetIDs = append(vetIDs, id)
	}
	rows.Close()

	// --- 3. ANIMALS (Marketplace) ---
	animalTypes := []string{"Kucing", "Anjing", "Burung", "Hamster", "Kelinci", "Reptil", "Serangga"}
	statuses := []string{"available", "sold", "pending"}
	colors := []string{"White", "Black", "Brown", "Grey", "Mixed", "Golden", "Spotted"}
	genders := []string{"Jantan", "Betina"}

	for i := 0; i < 1000; i++ { // Patch 1000 animals
		sellerID := userIDs[rand.Intn(len(userIDs))]
		aType := animalTypes[rand.Intn(len(animalTypes))]
		name := fmt.Sprintf("%s %d", aType, i)
		color := colors[rand.Intn(len(colors))]
		gender := genders[rand.Intn(len(genders))]
		stock := rand.Intn(10) + 1

		// Check duplicates roughly or just insert
		_, _ = DB.Exec(`INSERT INTO animals (seller_id, animal_type, breed, name, age, description, price, image_url, location, rating, status, color, gender, stock) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`,
			sellerID, aType, "Breed "+fmt.Sprint(i), name, rand.Intn(5)+1, "Description for "+name,
			float64(rand.Intn(5000000)+50000), "https://via.placeholder.com/300", "Jakarta", 4.0+rand.Float64(), statuses[rand.Intn(len(statuses))],
			color, gender, stock)
	}

	// Get Animal IDs
	var animalIDs []int
	rows, _ = DB.Query("SELECT id FROM animals")
	for rows.Next() {
		var id int
		rows.Scan(&id)
		animalIDs = append(animalIDs, id)
	}
	rows.Close()

	// --- 4. POSTS ---
	for i := 0; i < 50; i++ {
		uid := userIDs[rand.Intn(len(userIDs))]
		_, _ = DB.Exec(`INSERT INTO posts (user_id, content, image_url) VALUES ($1, $2, $3)`,
			uid, fmt.Sprintf("Post content number %d by user %d. #TerraPaw", i, uid), "https://via.placeholder.com/400")
	}

	var postIDs []int
	rows, _ = DB.Query("SELECT id FROM posts")
	for rows.Next() {
		var id int
		rows.Scan(&id)
		postIDs = append(postIDs, id)
	}
	rows.Close()

	// --- 5. COMMENTS & 6. LIKES ---
	if len(postIDs) > 0 {
		for i := 0; i < 100; i++ {
			pid := postIDs[rand.Intn(len(postIDs))]
			uid := userIDs[rand.Intn(len(userIDs))]
			_, _ = DB.Exec(`INSERT INTO comments (post_id, user_id, content) VALUES ($1, $2, $3)`,
				pid, uid, fmt.Sprintf("Nice post! Comment #%d", i))

			// Likes (Handle unique constraint)
			_, _ = DB.Exec(`INSERT INTO likes (post_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, pid, uid)
		}
	}

	// --- 7. USER PETS ---
	for i := 0; i < 30; i++ {
		uid := userIDs[rand.Intn(len(userIDs))]
		aType := animalTypes[rand.Intn(len(animalTypes))]
		name := fmt.Sprintf("Pet %d", i)
		_, _ = DB.Exec(`INSERT INTO user_pets (owner_id, name, animal_type, breed, age, image_url, story) 
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			uid, name, aType, "Mix Breed", rand.Intn(10)+1, "https://via.placeholder.com/200", "My lovely pet story.")
	}

	var userPetIDs []int
	rows, _ = DB.Query("SELECT id FROM user_pets")
	for rows.Next() {
		var id int
		rows.Scan(&id)
		userPetIDs = append(userPetIDs, id)
	}
	rows.Close()

	// --- 8. MEDICAL RECORDS ---
	if len(userPetIDs) > 0 && len(vetIDs) > 0 {
		types := []string{"vaccine", "sickness", "checkup"}
		for i := 0; i < 50; i++ {
			petID := userPetIDs[rand.Intn(len(userPetIDs))]
			vetID := vetIDs[rand.Intn(len(vetIDs))]
			rType := types[rand.Intn(len(types))]
			_, _ = DB.Exec(`INSERT INTO medical_records (pet_id, veterinarian_id, record_type, description, treatment, date, notes) 
				VALUES ($1, $2, $3, $4, $5, $6, $7)`,
				petID, vetID, rType, "Medical Record "+rType, "Treatment details...", time.Now().AddDate(0, -rand.Intn(12), 0), "Doctor notes here.")
		}
	}

	// --- 9. ORDERS ---
	if len(animalIDs) > 0 && len(customerIDs) > 0 {
		for i := 0; i < 20; i++ {
			buyerID := customerIDs[rand.Intn(len(customerIDs))]
			animalID := animalIDs[rand.Intn(len(animalIDs))]
			_, _ = DB.Exec(`INSERT INTO orders (buyer_id, animal_id, total_price, status) 
				VALUES ($1, $2, $3, 'completed')`,
				buyerID, animalID, rand.Float64()*1000000)
		}
	}

	// --- 10. CONSULTATIONS ---
	if len(vetIDs) > 0 && len(customerIDs) > 0 {
		statuses := []string{"pending", "scheduled", "completed", "cancelled"}
		for i := 0; i < 30; i++ {
			uid := customerIDs[rand.Intn(len(customerIDs))]
			vid := vetIDs[rand.Intn(len(vetIDs))]
			_, _ = DB.Exec(`INSERT INTO consultations (user_id, veterinarian_id, pet_name, symptoms, status, scheduled_at, consultation_type) 
				VALUES ($1, $2, $3, $4, $5, $6, 'online')`,
				uid, vid, "Pet Name", "Symptoms description...", statuses[rand.Intn(len(statuses))], time.Now().Add(24*time.Hour))
		}
	}

	// --- 11. MESSAGES ---
	if len(userIDs) > 1 {
		for i := 0; i < 50; i++ {
			s := userIDs[rand.Intn(len(userIDs))]
			r := userIDs[rand.Intn(len(userIDs))]
			if s == r {
				continue
			}
			_, _ = DB.Exec(`INSERT INTO messages (sender_id, receiver_id, content, is_read) 
				VALUES ($1, $2, $3, $4)`,
				s, r, fmt.Sprintf("Message content %d", i), rand.Intn(2) == 1)
		}
	}

	// --- 12. REMINDERS ---
	for i := 0; i < 30; i++ {
		uid := userIDs[rand.Intn(len(userIDs))]
		_, _ = DB.Exec(`INSERT INTO reminders (user_id, title, description, date, is_completed) 
			VALUES ($1, $2, $3, $4, $5)`,
			uid, "Reminder Title", "Reminder Desc", time.Now().Add(48*time.Hour), false)
	}

	// --- 13. NOTIFICATIONS ---
	notifTypes := []string{"reminder", "promo", "order", "info"}
	for i := 0; i < 50; i++ {
		uid := userIDs[rand.Intn(len(userIDs))]
		nType := notifTypes[rand.Intn(len(notifTypes))]
		_, _ = DB.Exec(`INSERT INTO notifications (user_id, title, message, type, is_read) 
			VALUES ($1, $2, $3, $4, $5)`,
			uid, "Notification Title "+nType, "This is a notification message.", nType, false)
	}

	// --- 14. REVIEWS ---
	if len(animalIDs) > 0 && len(userIDs) > 0 {
		for i := 0; i < 100; i++ {
			uid := userIDs[rand.Intn(len(userIDs))]
			aid := animalIDs[rand.Intn(len(animalIDs))]
			rating := rand.Intn(5) + 1
			_, _ = DB.Exec(`INSERT INTO reviews (user_id, animal_id, rating, comment, image_url) 
                VALUES ($1, $2, $3, $4, $5)`,
				uid, aid, rating, "Review comment sample...", "")
		}
	}

	log.Println("Data patching completed successfully!")

	// Automatically patch images to be realistic
	PatchImages()
}

func PatchImages() {
	log.Println("Running automatic image patching...")

	// 1. Patch Users (Avatars)
	avatars := []string{
		"https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?ixlib=rb-4.0.3&auto=format&fit=crop&w=150&q=80", // Male
		"https://images.unsplash.com/photo-1494790108377-be9c29b29330?ixlib=rb-4.0.3&auto=format&fit=crop&w=150&q=80", // Female
		"https://images.unsplash.com/photo-1599566150163-29194dcaad36?ixlib=rb-4.0.3&auto=format&fit=crop&w=150&q=80", // Male 2
		"https://images.unsplash.com/photo-1438761681033-6461ffad8d80?ixlib=rb-4.0.3&auto=format&fit=crop&w=150&q=80", // Female 2
	}

	// Only update if placeholder or empty
	_, err := DB.Exec(`UPDATE users SET avatar_url = CASE 
        WHEN id % 4 = 0 THEN $1 
        WHEN id % 4 = 1 THEN $2 
        WHEN id % 4 = 2 THEN $3 
        ELSE $4 END 
        WHERE avatar_url IS NULL OR avatar_url = '' OR avatar_url LIKE '%placeholder%'`,
		avatars[0], avatars[1], avatars[2], avatars[3])
	if err != nil {
		log.Printf("Failed to patch users: %v", err)
	}

	// 2. Patch Animals
	animalImages := map[string][]string{
		"Kucing": {
			"https://images.unsplash.com/photo-1514888286974-6c03e2ca1dba?ixlib=rb-4.0.3&auto=format&fit=crop&w=400&q=80",
			"https://images.unsplash.com/photo-1573865526739-10659fec78a5?ixlib=rb-4.0.3&auto=format&fit=crop&w=400&q=80",
		},
		"Anjing": {
			"https://images.unsplash.com/photo-1583511655857-d19b40a7a54e?ixlib=rb-4.0.3&auto=format&fit=crop&w=400&q=80",
			"https://images.unsplash.com/photo-1543466835-00a7907e9de1?ixlib=rb-4.0.3&auto=format&fit=crop&w=400&q=80",
		},
		"Burung": {
			"https://images.unsplash.com/photo-1552728089-57bdde30ebd1?ixlib=rb-4.0.3&auto=format&fit=crop&w=400&q=80",
			"https://images.unsplash.com/photo-1444464666168-49d633b86797?ixlib=rb-4.0.3&auto=format&fit=crop&w=400&q=80",
		},
		"Hamster": {
			"https://images.unsplash.com/photo-1425082661705-1834bfd09dca?ixlib=rb-4.0.3&auto=format&fit=crop&w=400&q=80",
		},
		"Kelinci": {
			"https://images.unsplash.com/photo-1585110396000-c9285745b504?ixlib=rb-4.0.3&auto=format&fit=crop&w=400&q=80",
			"https://images.unsplash.com/photo-1591382396632-3a4b1c974c81?ixlib=rb-4.0.3&auto=format&fit=crop&w=400&q=80",
		},
		"Reptil": {
			"https://images.unsplash.com/photo-1575535468632-345892291673?ixlib=rb-4.0.3&auto=format&fit=crop&w=400&q=80",
			"https://images.unsplash.com/photo-1533738363-b7f9aef128ce?ixlib=rb-4.0.3&auto=format&fit=crop&w=400&q=80",
		},
		"Serangga": {
			"https://images.unsplash.com/photo-1563404281029-7c85848c2c8f?ixlib=rb-4.0.3&auto=format&fit=crop&w=400&q=80",
			"https://images.unsplash.com/photo-1504198458649-3128b932f49e?ixlib=rb-4.0.3&auto=format&fit=crop&w=400&q=80",
		},
	}

	// Patch Animals
	for _, aType := range []string{"Kucing", "Anjing", "Burung", "Hamster", "Kelinci", "Reptil", "Serangga"} {
		imgs := animalImages[aType]
		img := imgs[0]
		if len(imgs) > 1 {
			// Simple round robin or just pick first for simplicity in SQL
			_, _ = DB.Exec(`UPDATE animals SET image_url = $1 WHERE animal_type = $2 AND (image_url IS NULL OR image_url = '' OR image_url LIKE '%placeholder%')`, img, aType)
		} else {
			_, _ = DB.Exec(`UPDATE animals SET image_url = $1 WHERE animal_type = $2 AND (image_url IS NULL OR image_url = '' OR image_url LIKE '%placeholder%')`, img, aType)
		}
	}

	// 3. Patch User Pets (Similar logic)
	for _, aType := range []string{"Kucing", "Anjing", "Burung", "Hamster"} {
		imgs := animalImages[aType]
		img := imgs[0]
		_, _ = DB.Exec(`UPDATE user_pets SET image_url = $1 WHERE animal_type = $2 AND (image_url IS NULL OR image_url = '' OR image_url LIKE '%placeholder%')`, img, aType)
	}

	// 4. Patch Posts
	postImages := []string{
		"https://images.unsplash.com/photo-1514888286974-6c03e2ca1dba?ixlib=rb-4.0.3&auto=format&fit=crop&w=500&q=80",
		"https://images.unsplash.com/photo-1543466835-00a7907e9de1?ixlib=rb-4.0.3&auto=format&fit=crop&w=500&q=80",
		"https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?ixlib=rb-4.0.3&auto=format&fit=crop&w=500&q=80",
	}

	_, _ = DB.Exec(`UPDATE posts SET image_url = CASE 
        WHEN id % 3 = 0 THEN $1 
        WHEN id % 3 = 1 THEN $2 
        ELSE $3 END 
        WHERE image_url IS NULL OR image_url = '' OR image_url LIKE '%placeholder%'`,
		postImages[0], postImages[1], postImages[2])

	log.Println("Automatic image patching finished.")
}
