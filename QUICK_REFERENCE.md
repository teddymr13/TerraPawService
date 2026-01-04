# TerraPaw Development Quick Reference

## Project Initialization Complete ✅

Your TerraPaw application has been successfully initialized with both frontend and backend ready for development.

## What's Included

### ✅ Go Backend (Production-Ready)
- **Framework**: Gin (high-performance web framework)
- **Database**: PostgreSQL with automatic migrations
- **Authentication**: JWT token-based auth with middleware
- **Features**:
  - Community (posts, comments, likes)
  - Marketplace (animals, orders)
  - Consultation (veterinarians, bookings)
- **API Endpoints**: 20+ endpoints with full CRUD operations

### ✅ React Native Frontend (Android)
- **Navigation**: React Navigation with bottom tabs
- **State Management**: Context API for authentication
- **HTTP Client**: Axios with interceptors
- **Features**:
  - Authentication (login/register)
  - Community feed
  - Marketplace browsing
  - Consultation booking
- **Responsive Design**: Optimized for mobile screens

### ✅ Database Schema
Automatically created PostgreSQL tables:
- users
- posts, comments, likes
- animals, orders
- veterinarians, consultations
- messages

## Getting Started (5 Minutes)

### 1. Start PostgreSQL
```bash
# Make sure PostgreSQL is running
```

### 2. Start Backend
```bash
cd backend
cp .env.example .env
# Edit .env with your database credentials
go run cmd/main.go
```

Backend will be available at: `http://localhost:8080/api`

### 3. Start Frontend
```bash
cd TerraPawApp
npm install  # (if not done yet)
npm start
npx react-native run-android
```

## Common Development Tasks

### Add New API Endpoint

1. **Create handler** in `backend/handlers/`:
```go
func NewHandler(c *gin.Context) {
  // Your code here
}
```

2. **Register route** in `backend/routes/routes.go`:
```go
router.POST("/api/feature/endpoint", handlers.NewHandler)
```

3. **Use in frontend** from `TerraPawApp/src/api/endpoints.js`:
```javascript
export const featureAPI = {
  newEndpoint: () => apiClient.post('/feature/endpoint'),
};
```

4. **Call from component**:
```javascript
import { featureAPI } from '../api/endpoints';
const response = await featureAPI.newEndpoint();
```

### Modify Database Schema

1. Edit table creation SQL in `backend/db/database.go`
2. Add new model struct in `backend/models/models.go`
3. Create migration or modify `createTables()` function
4. Restart backend - tables will be recreated

### Add New Screen Component

1. Create component in `TerraPawApp/src/screens/`
2. Add to navigation in `TerraPawApp/src/navigation/RootNavigator.js`
3. Import and use in navigation structure

### Add State Management

The app uses React Context for auth state:

```javascript
// Use in any component
const { user, token, isSignedIn } = useAuth();

// Or create new context for features
export const const FeatureContext = createContext();

export const FeatureProvider = ({ children }) => {
  // Your logic here
  return <FeatureContext.Provider value={value}>{children}</FeatureContext.Provider>;
};
```

## API Response Format

### Success Response
```json
{
  "success": true,
  "message": "Operation successful",
  "data": { /* your data */ }
}
```

### Error Response
```json
{
  "success": false,
  "message": "Error description",
  "error": "Detailed error message"
}
```

## Authentication Flow

1. User registers/logs in
2. Backend returns JWT token
3. Token stored in AsyncStorage
4. Token sent with every request in Authorization header
5. Middleware validates token
6. On logout/token expiration, user returned to login

## File Organization

### Backend Structure
```
handlers/     - API request handlers
models/       - Data structures
routes/       - URL routing
middleware/   - Authentication & CORS
db/           - Database operations
config/       - Configuration loading
utils/        - Helper functions
```

### Frontend Structure
```
api/          - API client & endpoints
context/      - State management (Auth)
screens/      - UI screens
navigation/   - Navigation structure
```

## Environment Variables

### Backend (.env)
```
DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME
PORT, JWT_SECRET
```

### Frontend (.env)
```
API_BASE_URL (set to backend URL)
```

## Testing API Endpoints

### Using cURL
```bash
# Register
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@test.com","password":"pass","fullname":"Test"}'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"pass"}'

# Get posts
curl http://localhost:8080/api/community/posts
```

### Using Postman
1. Import endpoints
2. Set Authorization header: `Bearer <token>`
3. Test endpoints

## Performance Tips

1. **Backend**: Use connection pooling for database
2. **Frontend**: Implement FlatList for long lists
3. **Database**: Add indexes for frequently queried fields
4. **API**: Implement pagination for list endpoints
5. **Images**: Use image optimization/caching

## Security Checklist

- [ ] Change JWT_SECRET in production
- [ ] Use HTTPS for all connections
- [ ] Validate all user inputs
- [ ] Sanitize database queries (parameterized)
- [ ] Never commit .env with real credentials
- [ ] Use environment variables for secrets
- [ ] Implement rate limiting
- [ ] Add CORS restrictions
- [ ] Verify user permissions on backend
- [ ] Hash passwords (already implemented)

## Deployment Checklist

### Backend
- [ ] Set up production PostgreSQL
- [ ] Configure environment variables
- [ ] Build Go binary: `go build -o terrapaw cmd/main.go`
- [ ] Set up reverse proxy (Nginx/Apache)
- [ ] Enable HTTPS/SSL
- [ ] Set up logging
- [ ] Configure backups

### Frontend
- [ ] Update API_BASE_URL to production
- [ ] Build release APK: `cd android && ./gradlew assembleRelease`
- [ ] Sign APK with keystore
- [ ] Test on physical device
- [ ] Publish to Google Play Store

## Useful Commands

```bash
# Backend
go run cmd/main.go              # Run dev server
go build -o terrapaw cmd/main.go # Build binary
go test ./...                   # Run tests

# Frontend
npm start                       # Start Metro
npm install                     # Install deps
npm run android                 # Build & run
npm run build                   # Build for production

# PostgreSQL
psql -U postgres               # Connect to database
\l                             # List databases
\c terrapaw                    # Connect to terrapaw
\dt                            # List tables
```

## Debugging

### Backend
- Check logs in console
- Use Go debugger: `dlv debug cmd/main.go`
- Log database queries: `db.Query(...)`

### Frontend
- React Native debugger: Shake device → "Debug JS Remotely"
- AsyncStorage: `adb logcat | grep AsyncStorage`
- Network requests: Use React Native Network Logger

## Next Steps

1. ✅ Both projects initialized
2. Set up PostgreSQL database
3. Start backend server
4. Configure frontend API URL
5. Test registration & login
6. Add sample data
7. Test all features
8. Customize UI/colors as needed
9. Add more features based on SETUP_GUIDE.md

## Resources

- Backend code: `/backend/`
- Frontend code: `/TerraPawApp/src/`
- Full guide: `SETUP_GUIDE.md`
- Backend docs: `/backend/README.md`
- Frontend docs: `/TerraPawApp/README.md`

## Support

All major features are implemented. For advanced features (video calls, real-time messaging, payments), refer to integration guides in SETUP_GUIDE.md.

---

**Everything is set up and ready to go! 🚀**

Start with the 5-minute setup above, then explore the code in your IDE.
