# TerraPaw - Complete Setup & Deployment Guide

## Project Overview

**TerraPaw** is a comprehensive platform for pet lovers, combining three core features:

1. **Community** - Connect with other animal lovers, share posts, and interact
2. **Marketplace** - Buy and sell pets (cats, dogs, etc.)
3. **Consultation** - Book consultations with veterinarians for pet health concerns

The application consists of:
- **Frontend**: React Native for Android
- **Backend**: Go with Gin framework and PostgreSQL

## Directory Structure

```
TerraPaw/
├── TerraPawApp/                    # React Native Frontend
│   ├── src/
│   │   ├── api/                   # API client and endpoints
│   │   ├── context/               # Context (Auth state management)
│   │   ├── screens/               # UI screens
│   │   └── navigation/            # Navigation structure
│   ├── .env                       # Frontend configuration
│   ├── package.json
│   └── App.tsx
│
└── backend/                        # Go Backend Server
    ├── cmd/
    │   └── main.go               # Server entry point
    ├── config/                   # Configuration
    ├── db/                       # Database setup
    ├── models/                   # Data models
    ├── handlers/                 # API handlers
    ├── middleware/               # Auth middleware
    ├── routes/                   # Route definitions
    ├── utils/                    # Utilities
    ├── .env                      # Backend configuration
    ├── go.mod
    ├── go.sum
    └── README.md
```

## Quick Start

### Prerequisites
- Node.js & npm (for frontend)
- Go 1.21+ (for backend)
- PostgreSQL 12+ (for database)
- Android SDK (for mobile development)

### Step 1: Database Setup

```bash
# Create PostgreSQL database
createdb terrapaw
```

### Step 2: Backend Setup

```bash
cd backend

# Copy environment template
cp .env.example .env

# Edit .env with your database credentials
# Set JWT_SECRET to a strong random value
nano .env

# Install dependencies
go mod download

# Run the server
go run cmd/main.go
```

The backend will start on `http://localhost:8080`

### Step 3: Frontend Setup

```bash
cd TerraPawApp

# Install dependencies
npm install

# Edit API configuration
nano .env
# Update API_BASE_URL to match your backend:
# For emulator: http://10.0.2.2:8080/api
# For physical device: http://<YOUR_IP>:8080/api

# Start development server
npm start

# In a new terminal, run on Android
npx react-native run-android
```

## API Documentation

### Base URL
```
http://localhost:8080/api
```

### Authentication Endpoints

**Register**
```
POST /auth/register
Content-Type: application/json

{
  "username": "string",
  "email": "string",
  "password": "string",
  "fullname": "string"
}

Response: { token, user_id, ... }
```

**Login**
```
POST /auth/login
Content-Type: application/json

{
  "email": "string",
  "password": "string"
}

Response: { token, user_id, username, ... }
```

### Community Endpoints

**Create Post**
```
POST /community/posts
Authorization: Bearer <token>
Content-Type: application/json

{
  "content": "string",
  "image_url": "string (optional)"
}
```

**Get Posts**
```
GET /community/posts?page=1&limit=10
```

**Like Post**
```
POST /community/posts/{id}/like
Authorization: Bearer <token>
```

### Marketplace Endpoints

**Get Animals**
```
GET /marketplace/animals?animal_type=cat&page=1&limit=10
```

**Create Animal Listing**
```
POST /marketplace/animals
Authorization: Bearer <token>
Content-Type: application/json

{
  "animal_type": "cat|dog",
  "breed": "string",
  "name": "string",
  "age": number,
  "description": "string",
  "price": number,
  "image_url": "string",
  "location": "string"
}
```

**Create Order**
```
POST /marketplace/orders
Authorization: Bearer <token>
Content-Type: application/json

{
  "animal_id": number
}
```

### Consultation Endpoints

**Get Veterinarians**
```
GET /consultation/veterinarians?page=1&limit=10
```

**Create Consultation**
```
POST /consultation/consultations
Authorization: Bearer <token>
Content-Type: application/json

{
  "veterinarian_id": number,
  "pet_name": "string",
  "symptoms": "string",
  "consultation_type": "online|offline",
  "scheduled_at": "ISO date (optional)"
}
```

## Mobile App Features

### Login/Registration
- User can create account with email
- Secure password storage
- JWT token-based authentication
- Auto-login with saved token

### Community Tab
- View posts from other users
- Create new posts with text content
- Like posts
- Comment on posts
- See post engagement (likes, comments count)

### Marketplace Tab
- Browse available animals
- Filter by animal type (cats, dogs, etc.)
- View detailed animal information
- See seller information and ratings
- Make purchase orders

### Consultation Tab
- View available veterinarians
- See veterinarian ratings and specializations
- Book consultations
- Track consultation status
- Message with veterinarians

## Configuration

### Backend Environment Variables (.env)

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=terrapaw

# Server
PORT=8080

# Security
JWT_SECRET=your-very-strong-secret-key-here
```

### Frontend Environment Variables (.env)

```env
# API Configuration
API_BASE_URL=http://192.168.1.100:8080/api
```

## Running in Production

### Backend Production Deployment

1. **Build binary:**
   ```bash
   go build -o terrapaw cmd/main.go
   ```

2. **Use production database:**
   - Set up PostgreSQL on production server
   - Update .env with production credentials
   - Enable SSL for database connection

3. **Environment variables:**
   - Change JWT_SECRET to strong random value
   - Enable HTTPS (use reverse proxy like Nginx)
   - Configure CORS for production domain

4. **Deploy:**
   ```bash
   export $(cat .env | xargs)
   ./terrapaw
   ```

### Frontend Production Build

1. **Build APK for Android:**
   ```bash
   cd TerraPawApp
   cd android
   ./gradlew assembleRelease
   ```

2. **Update API URL:**
   - Change API_BASE_URL in src/api/client.js to production server
   - Or use environment variables system

3. **Sign APK:**
   - Generate keystore file
   - Sign with: `jarsigner -verbose -sigalg SHA1withRSA -digestalg SHA1 ...`

## Testing the Application

### Test Accounts

After starting the backend, you can register new accounts:

1. **User 1 (Pet Lover)**
   - Username: user1
   - Email: user1@example.com
   - Password: password123

2. **User 2 (Seller)**
   - Username: seller1
   - Email: seller1@example.com
   - Password: password123

3. **User 3 (Veterinarian)**
   - Username: vet1
   - Email: vet1@example.com
   - Password: password123
   - Register as veterinarian after login

### Sample Data

You can add sample data to test:

1. Create posts in Community
2. Create animal listings in Marketplace
3. Register as veterinarian and create profile
4. Book consultations

## Troubleshooting

### Backend Issues

**Port 8080 already in use**
```bash
# Change PORT in .env or kill process:
lsof -i :8080
kill -9 <PID>
```

**Database connection error**
```
# Check PostgreSQL is running and credentials are correct
psql -U postgres -d terrapaw
```

**JWT token errors**
- Clear app cache and re-login
- Check JWT_SECRET in backend .env

### Frontend Issues

**Cannot connect to backend**
- Ensure backend is running on correct port
- Check API_BASE_URL in .env
- For emulator, use 10.0.2.2 instead of localhost
- Check firewall allows connection

**Build errors**
```bash
# Clear caches
rm -rf node_modules
npm install

cd android
./gradlew clean
./gradlew build
```

## Next Steps & Enhancements

### Phase 2 Features
- [ ] Real-time messaging with WebSockets
- [ ] Video call consultation with Agora/Twilio
- [ ] Image upload to cloud storage (AWS S3)
- [ ] Push notifications
- [ ] Payment gateway integration (Stripe/PayPal)
- [ ] Pet profile creation and management
- [ ] Veterinarian appointment calendar
- [ ] User ratings and reviews system

### Performance Optimization
- [ ] Database indexing optimization
- [ ] Implement Redis caching
- [ ] API response pagination
- [ ] CDN for image delivery
- [ ] Database query optimization

### Security Enhancements
- [ ] Rate limiting on API endpoints
- [ ] Email verification for registration
- [ ] Password reset functionality
- [ ] Two-factor authentication
- [ ] Data encryption at rest

## Resources

- [React Native Documentation](https://reactnative.dev)
- [Gin Framework Documentation](https://gin-gonic.com)
- [PostgreSQL Documentation](https://www.postgresql.org/docs)
- [JWT Authentication](https://jwt.io)

## License

MIT License - See LICENSE file for details

## Support

For issues, questions, or suggestions, please open an issue in the repository.

---

**Happy coding! 🐾🐱🐶**
