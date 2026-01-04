# TerraPaw Backend API

A production-ready Go backend server for the TerraPaw application, built with Gin framework and PostgreSQL database.

## Features

### Authentication
- User registration and login
- JWT token-based authentication
- Secure password hashing
- User profile management

### Community
- Create and share posts
- Like posts
- Comment on posts
- User interactions

### Marketplace
- List animals (cats, dogs, etc.) for sale
- Browse and search animals
- Filter by animal type
- Purchase orders

### Consultation
- Veterinarian registration and profiles
- Book consultations
- Manage consultation status
- Real-time messaging support

## Project Structure

```
backend/
├── cmd/
│   └── main.go              # Application entry point
├── config/
│   └── config.go            # Configuration management
├── db/
│   └── database.go          # Database initialization & migrations
├── models/
│   └── models.go            # Data models
├── handlers/
│   ├── auth.go              # Authentication endpoints
│   ├── community.go         # Community feature endpoints
│   ├── marketplace.go       # Marketplace endpoints
│   └── consultation.go      # Consultation endpoints
├── middleware/
│   └── auth.go              # Authentication middleware
├── routes/
│   └── routes.go            # Route definitions
├── utils/
│   ├── jwt.go               # JWT utilities
│   └── response.go          # Response utilities
├── .env.example             # Environment variables template
├── go.mod                   # Go module definition
├── go.sum                   # Dependency checksums
└── README.md
```

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Git

## Installation

### 1. Install Go Dependencies

From the `backend` directory, run:

```bash
go mod download
```

### 2. Set Up PostgreSQL Database

Create a new database for TerraPaw:

```bash
# Using PostgreSQL CLI
createdb terrapaw
```

Or use a GUI tool like pgAdmin.

### 3. Configure Environment Variables

Copy `.env.example` to `.env` and update the values:

```bash
cp .env.example .env
```

Edit `.env` with your database credentials and configuration:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=terrapaw
PORT=8080
JWT_SECRET=your-super-secret-key-change-in-production
```

### 4. Run the Application

```bash
# From the backend directory
go run cmd/main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### Authentication

```
POST /api/auth/register
POST /api/auth/login
GET /api/auth/profile (requires token)
```

### Community

```
GET /api/community/posts
GET /api/community/posts/:id
POST /api/community/posts (requires token)
POST /api/community/posts/:id/like (requires token)
DELETE /api/community/posts/:id/like (requires token)
POST /api/community/posts/:id/comments (requires token)
```

### Marketplace

```
GET /api/marketplace/animals
GET /api/marketplace/animals/:id
POST /api/marketplace/animals (requires token)
POST /api/marketplace/orders (requires token)
GET /api/marketplace/orders (requires token)
```

### Consultation

```
GET /api/consultation/veterinarians
GET /api/consultation/veterinarians/:id
POST /api/consultation/veterinarians/register (requires token)
POST /api/consultation/consultations (requires token)
GET /api/consultation/consultations (requires token)
GET /api/consultation/consultations/:id (requires token)
PUT /api/consultation/consultations/:id/status (requires token)
```

## Authentication

The API uses JWT (JSON Web Tokens) for authentication:

1. User registers or logs in
2. Server returns a JWT token
3. Client includes token in `Authorization` header: `Bearer <token>`
4. Token is valid for 7 days

## Database Schema

The application automatically creates the following tables on startup:

- `users` - User accounts and profiles
- `posts` - Community posts
- `comments` - Post comments
- `likes` - Post likes
- `animals` - Pet listings
- `orders` - Purchase orders
- `veterinarians` - Veterinarian profiles
- `consultations` - Consultation records
- `messages` - Direct messages

## Error Handling

The API returns consistent error responses:

```json
{
  "success": false,
  "message": "Error description",
  "error": "Detailed error message"
}
```

## CORS

The API has CORS enabled to allow requests from any origin. Modify the CORS middleware in `cmd/main.go` to restrict origins in production.

## Security Considerations

1. **Change JWT_SECRET** - Set a strong, random secret in production
2. **Use HTTPS** - Always use HTTPS in production
3. **Database credentials** - Use environment variables, never hardcode
4. **Rate limiting** - Consider adding rate limiting middleware
5. **Input validation** - All inputs are validated before processing

## Building for Production

```bash
# Build binary
go build -o terrapaw cmd/main.go

# Or build for specific OS
GOOS=linux GOARCH=amd64 go build -o terrapaw cmd/main.go
```

## Deployment

### Docker

Create a Dockerfile:

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o terrapaw cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/terrapaw .
EXPOSE 8080
CMD ["./terrapaw"]
```

Build and run:

```bash
docker build -t terrapaw .
docker run -p 8080:8080 --env-file .env terrapaw
```

### Cloud Platforms

The application can be deployed to:
- AWS (EC2, Lambda with API Gateway)
- Google Cloud (Cloud Run, Compute Engine)
- Azure (App Service, Container Instances)
- Heroku
- DigitalOcean

## Troubleshooting

### Database Connection Error

```
failed to ping database
```

Check:
- PostgreSQL is running
- Database credentials in `.env` are correct
- Database exists

### Port Already in Use

```
listen tcp :8080: bind: address already in use
```

Change the port in `.env` or kill the process using port 8080.

### CORS Errors

Ensure the frontend URL is allowed in the CORS middleware or update to allow all origins during development.

## Performance Optimization

- Database queries are optimized with proper indexing
- Connection pooling is configured in the database driver
- Response caching can be implemented for read-heavy operations
- Consider adding pagination to list endpoints

## Testing

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./handlers

# Verbose output
go test -v ./...
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## License

MIT License - See LICENSE file for details

## Support

For issues and questions, please open an issue on the GitHub repository.
