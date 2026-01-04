package db

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func PatchLargeData() {
	log.Println("Starting LARGE data patching (4000 records)...")

	// Check if we already have enough data
	var count int
	DB.QueryRow("SELECT COUNT(*) FROM animals").Scan(&count)
	if count > 3500 {
		log.Println("Data seems already populated (3500+). Skipping large patch.")
		return
	}

	// Password hash
	password, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	passStr := string(password)

	// --- 1. USERS (500 Users) ---
	log.Println("Creating 500 Users...")
	for i := 0; i < 500; i++ {
		username := fmt.Sprintf("user_bulk_%d", i)
		email := fmt.Sprintf("user_bulk_%d@example.com", i)
		fullname := fmt.Sprintf("User Bulk %d", i)
		_, _ = DB.Exec(`INSERT INTO users (username, email, password, fullname, user_type, avatar_url, bio) 
			VALUES ($1, $2, $3, $4, 'customer', $5, $6) ON CONFLICT DO NOTHING`,
			username, email, passStr, fullname, "https://via.placeholder.com/150", "I love pets!")
	}

	// Get User IDs
	rows, _ := DB.Query("SELECT id FROM users")
	var userIDs []int
	for rows.Next() {
		var id int
		rows.Scan(&id)
		userIDs = append(userIDs, id)
	}
	rows.Close()

	if len(userIDs) == 0 {
		log.Println("No users found, aborting.")
		return
	}

	// --- 2. VETERINARIANS (50 Vets) ---
	log.Println("Creating 50 Vets...")
	for i := 0; i < 50; i++ {
		// uid := userIDs[i%len(userIDs)] // Reuse users as vets (Unused variable removed)
		_, _ = DB.Exec(`INSERT INTO users (username, email, password, fullname, user_type, avatar_url) 
            VALUES ($1, $2, $3, $4, 'veterinarian', $5) 
            ON CONFLICT (username) DO UPDATE SET user_type = 'veterinarian'`,
			fmt.Sprintf("vet_bulk_%d", i), fmt.Sprintf("vet_bulk_%d@example.com", i), passStr, fmt.Sprintf("Dr. Bulk %d", i), "https://via.placeholder.com/150")

		var vetUserID int
		DB.QueryRow("SELECT id FROM users WHERE username = $1", fmt.Sprintf("vet_bulk_%d", i)).Scan(&vetUserID)

		_, _ = DB.Exec(`INSERT INTO veterinarians (user_id, clinic_name, license_number, specialization, phone, address, bio, rating) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT DO NOTHING`,
			vetUserID, fmt.Sprintf("Klinik Hewan %d", i), fmt.Sprintf("LIC-%d", i), "General Vet", "08123456789", "Jl. Hewan No. "+fmt.Sprint(i), "Expert Vet", 4.0+rand.Float64())
	}

	// --- 3. ANIMALS (4000 Items) ---
	log.Println("Creating 4000 Marketplace Items (Live & Food)...")

	// Live Animal Breeds
	breedsMap := map[string][]string{
		"Kucing":   {"Persia", "Anggora", "British Shorthair", "Munchkin", "Domestik"},
		"Anjing":   {"Golden Retriever", "Bulldog", "Poodle", "Husky", "Pomeranian"},
		"Burung":   {"Lovebird", "Kenari", "Murai Batu", "Kakaktua", "Parkit"},
		"Hamster":  {"Syrian", "Winter White", "Roborovski", "Campbell"},
		"Kelinci":  {"Rex", "Anggora", "Flemish Giant", "Netherland Dwarf"},
		"Reptil":   {"Iguana", "Gecko", "Kura-kura", "Ular Corn Snake", "Chameleon"},
		"Serangga": {"Kumbang Tanduk", "Tarantula", "Kalajengking", "Belalang Sembah"},
	}

	// Food Types (Category is "Makanan X")
	foodMap := map[string][]string{
		"Makanan Kucing":   {"Whiskas Tuna", "Royal Canin Kitten", "Me-O Salmon", "Friskies Seafood", "Pro Plan"},
		"Makanan Anjing":   {"Pedigree Chicken", "Royal Canin Puppy", "Alpo Beef", "Science Diet", "Cesar"},
		"Makanan Burung":   {"Pakan Kenari", "Millet Putih", "Voer Burung Juara", "Jangkrik Kering"},
		"Makanan Hamster":  {"Vitakraft Menu", "Biji Bunga Matahari", "Hamster Mix", "Snack Hamster"},
		"Makanan Kelinci":  {"Nova Rabbit Food", "Hay Timothy", "Pelet Kelinci", "Alfafa Hay"},
		"Makanan Reptil":   {"Jangkrik Kering", "Pelet Kura-kura", "Ulat Hongkong Kering", "Calcium Powder"},
		"Makanan Serangga": {"Jelly Pot", "Beetle Jelly", "Protein Mix", "Buah Segar"},
	}

	locations := []string{"Jakarta", "Bandung", "Surabaya", "Medan", "Bali", "Yogyakarta"}

	stmt, _ := DB.Prepare(`INSERT INTO animals (seller_id, animal_type, breed, name, age, description, price, image_url, location, rating, status, color, gender, stock) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, 'available', $11, $12, $13)`)

	// Generate 2000 Live Animals
	baseTypes := []string{"Kucing", "Anjing", "Burung", "Hamster", "Kelinci", "Reptil", "Serangga"}
	for i := 0; i < 2000; i++ {
		sellerID := userIDs[rand.Intn(len(userIDs))]
		aType := baseTypes[rand.Intn(len(baseTypes))]
		breedList := breedsMap[aType]
		breed := breedList[rand.Intn(len(breedList))]

		name := fmt.Sprintf("%s %s", breed, aType)
		desc := "Hewan sehat, vaksin lengkap, siap adopsi."
		price := float64(rand.Intn(5000000) + 500000)
		imgUrl := "https://via.placeholder.com/300"

		// Specific Images based on type
		switch aType {
		case "Kucing":
			imgUrl = "https://images.unsplash.com/photo-1514888286974-6c03e2ca1dba"
		case "Anjing":
			imgUrl = "https://images.unsplash.com/photo-1543466835-00a7907e9de1"
		case "Burung":
			imgUrl = "https://images.unsplash.com/photo-1552728089-57bdde30ebd1"
		case "Hamster":
			imgUrl = "https://images.unsplash.com/photo-1425082661705-1834bfd90591"
		case "Reptil":
			imgUrl = "https://images.unsplash.com/photo-1504450874802-0ed58ffa94df"
		case "Serangga":
			imgUrl = "https://images.unsplash.com/photo-1563207038-7f99849206d9"
		case "Kelinci":
			imgUrl = "https://images.unsplash.com/photo-1585110396000-c9285745b255"
		}

		loc := locations[rand.Intn(len(locations))]
		gender := "Jantan"
		if rand.Intn(2) == 0 {
			gender = "Betina"
		}

		_, _ = stmt.Exec(sellerID, aType, breed, name, rand.Intn(5)+1, desc, price, imgUrl, loc, 4.0+rand.Float64(), "Mixed", gender, rand.Intn(10)+1)
	}

	// Generate 2000 Food Items
	foodTypes := []string{"Makanan Kucing", "Makanan Anjing", "Makanan Burung", "Makanan Hamster", "Makanan Kelinci", "Makanan Reptil", "Makanan Serangga"}
	for i := 0; i < 2000; i++ {
		sellerID := userIDs[rand.Intn(len(userIDs))]
		aType := foodTypes[rand.Intn(len(foodTypes))]
		foodList := foodMap[aType]
		productName := foodList[rand.Intn(len(foodList))]

		name := productName
		breed := "Makanan/Aksesoris"
		desc := "Makanan berkualitas tinggi untuk hewan kesayangan Anda."
		price := float64(rand.Intn(500000) + 15000)

		imgUrl := "https://via.placeholder.com/300?text=Food"
		switch aType {
		case "Makanan Kucing":
			imgUrl = "https://images.unsplash.com/photo-1583337130417-3346a1be7dee"
		case "Makanan Anjing":
			imgUrl = "https://images.unsplash.com/photo-1589924691195-41432c84c161"
		case "Makanan Burung":
			imgUrl = "https://images.unsplash.com/photo-1623366302587-b38b1ddaefd9"
		default:
			imgUrl = "https://images.unsplash.com/photo-1601004890684-d8cbf643f5f2"
		}

		loc := locations[rand.Intn(len(locations))]

		_, _ = stmt.Exec(sellerID, aType, breed, name, 0, desc, price, imgUrl, loc, 4.0+rand.Float64(), "-", "-", rand.Intn(50)+1)
	}
	stmt.Close()

	// --- 4. COMMUNITY POSTS (1000 Posts) ---
	log.Println("Creating 1000 Community Posts...")
	stmtPost, _ := DB.Prepare(`INSERT INTO posts (user_id, content, image_url, likes) VALUES ($1, $2, $3, $4)`)
	for i := 0; i < 1000; i++ {
		uid := userIDs[rand.Intn(len(userIDs))]
		likes := rand.Intn(100)
		_, _ = stmtPost.Exec(uid, fmt.Sprintf("Halo teman-teman! Ini postingan ke-%d saya. #TerraPawCommunity", i), "https://images.unsplash.com/photo-1548199973-03cce0bbc87b", likes)
	}
	stmtPost.Close()

	// --- 5. ORDERS (500 Orders) ---
	log.Println("Creating 500 Orders...")
	var animalIDs []int
	rows, _ = DB.Query("SELECT id FROM animals LIMIT 2000")
	for rows.Next() {
		var id int
		rows.Scan(&id)
		animalIDs = append(animalIDs, id)
	}
	rows.Close()

	stmtOrder, _ := DB.Prepare(`INSERT INTO orders (buyer_id, animal_id, total_price, status, quantity, created_at) VALUES ($1, $2, $3, $4, $5, $6)`)
	statuses := []string{"pending", "completed", "cancelled", "shipped"}

	for i := 0; i < 500; i++ {
		buyerID := userIDs[rand.Intn(len(userIDs))]
		animalID := animalIDs[rand.Intn(len(animalIDs))]
		status := statuses[rand.Intn(len(statuses))]
		date := time.Now().AddDate(0, 0, -rand.Intn(30)) // Past 30 days

		_, _ = stmtOrder.Exec(buyerID, animalID, rand.Float64()*500000, status, 1, date)
	}
	stmtOrder.Close()

	log.Println("LARGE Data patching completed!")
}

func SeedCategories() {
	log.Println("Seeding Categories table...")

	categories := []struct {
		Name string
		Icon string
		Type string
	}{
		{"Kucing", "🐱", "animal"},
		{"Anjing", "🐶", "animal"},
		{"Burung", "🐦", "animal"},
		{"Hamster", "🐹", "animal"},
		{"Kelinci", "🐰", "animal"},
		{"Reptil", "🦎", "animal"},
		{"Serangga", "🦗", "animal"},
		{"Makanan Kucing", "🥫", "food"},
		{"Makanan Anjing", "🦴", "food"},
		{"Makanan Burung", "🌾", "food"},
		{"Makanan Kelinci", "🥕", "food"},
		{"Makanan Hamster", "🌻", "food"},
		{"Makanan Reptil", "🦗", "food"},
	}

	for _, cat := range categories {
		_, err := DB.Exec("INSERT INTO categories (name, icon, type) VALUES ($1, $2, $3) ON CONFLICT (name) DO UPDATE SET icon = $2, type = $3", cat.Name, cat.Icon, cat.Type)
		if err != nil {
			log.Printf("Error seeding category %s: %v", cat.Name, err)
		}
	}
}

func EnsureFoodData() {
	log.Println("Checking for Food Data...")

	var foodCount int
	DB.QueryRow("SELECT COUNT(*) FROM animals WHERE animal_type LIKE 'Makanan %'").Scan(&foodCount)

	if foodCount > 500 {
		log.Println("Food data already exists. Skipping.")
		return
	}

	log.Println("Injecting 1000 Food Items...")

	// Get User IDs to assign as sellers
	rows, _ := DB.Query("SELECT id FROM users LIMIT 100")
	var userIDs []int
	for rows.Next() {
		var id int
		rows.Scan(&id)
		userIDs = append(userIDs, id)
	}
	rows.Close()

	if len(userIDs) == 0 {
		return
	}

	foodMap := map[string][]string{
		"Makanan Kucing":   {"Whiskas Tuna", "Royal Canin Kitten", "Me-O Salmon", "Friskies Seafood", "Pro Plan"},
		"Makanan Anjing":   {"Pedigree Chicken", "Royal Canin Puppy", "Alpo Beef", "Science Diet", "Cesar"},
		"Makanan Burung":   {"Pakan Kenari", "Millet Putih", "Voer Burung Juara", "Jangkrik Kering"},
		"Makanan Hamster":  {"Vitakraft Menu", "Biji Bunga Matahari", "Hamster Mix", "Snack Hamster"},
		"Makanan Kelinci":  {"Nova Rabbit Food", "Hay Timothy", "Pelet Kelinci", "Alfafa Hay"},
		"Makanan Reptil":   {"Jangkrik Kering", "Pelet Kura-kura", "Ulat Hongkong Kering", "Calcium Powder"},
		"Makanan Serangga": {"Jelly Pot", "Beetle Jelly", "Protein Mix", "Buah Segar"},
	}

	locations := []string{"Jakarta", "Bandung", "Surabaya", "Medan", "Bali", "Yogyakarta"}
	foodTypes := []string{"Makanan Kucing", "Makanan Anjing", "Makanan Burung", "Makanan Hamster", "Makanan Kelinci", "Makanan Reptil", "Makanan Serangga"}

	stmt, _ := DB.Prepare(`INSERT INTO animals (seller_id, animal_type, breed, name, age, description, price, image_url, location, rating, status, color, gender, stock) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, 'available', $11, $12, $13)`)

	for i := 0; i < 1000; i++ {
		sellerID := userIDs[rand.Intn(len(userIDs))]
		aType := foodTypes[rand.Intn(len(foodTypes))]
		foodList := foodMap[aType]
		productName := foodList[rand.Intn(len(foodList))]

		name := fmt.Sprintf("%s - %d", productName, i)
		breed := "Makanan/Aksesoris"
		desc := "Makanan berkualitas tinggi untuk hewan kesayangan Anda. Stok selalu baru dan higienis."
		price := float64(rand.Intn(500000) + 15000)

		imgUrl := "https://via.placeholder.com/300?text=Food"
		switch aType {
		case "Makanan Kucing":
			imgUrl = "https://images.unsplash.com/photo-1583337130417-3346a1be7dee"
		case "Makanan Anjing":
			imgUrl = "https://images.unsplash.com/photo-1589924691195-41432c84c161"
		case "Makanan Burung":
			imgUrl = "https://images.unsplash.com/photo-1623366302587-b38b1ddaefd9"
		default:
			imgUrl = "https://images.unsplash.com/photo-1601004890684-d8cbf643f5f2"
		}

		loc := locations[rand.Intn(len(locations))]

		_, _ = stmt.Exec(sellerID, aType, breed, name, 0, desc, price, imgUrl, loc, 4.0+rand.Float64(), "-", "-", rand.Intn(50)+1)
	}
	stmt.Close()
	log.Println("Food Data Injected Successfully!")
}
