# Backend CORS & Configuration Guide

## Current Backend Setup (main.go)

Your Go Fiber backend is already configured with:

```go
// CORS (Flutter Web SAFE)
app.Use(cors.New(cors.Config{
    AllowOrigins: "*", // DEV ONLY
    AllowMethods: "GET,POST,PUT,DELETE,OPTIONS,PATCH",
    AllowHeaders: "Origin, Content-Type, Accept, Authorization",
    MaxAge:       300,
}))
```

✅ **CORS is already enabled** for all origins (development mode)

---

## Production CORS Configuration

For production, update `main.go` CORS settings:

```go
app.Use(cors.New(cors.Config{
    AllowOrigins: "https://yourdomain.com,https://www.yourdomain.com",
    AllowMethods: "GET,POST,PUT,DELETE,OPTIONS,PATCH",
    AllowHeaders: "Origin, Content-Type, Accept, Authorization",
    MaxAge:       300,
}))
```

---

## Backend Health Check

The backend has a built-in health check endpoint:

```bash
curl http://localhost:3000/health
```

Response:
```json
{
  "status": "ok",
  "service": "attendance-backend",
  "port": "3000"
}
```

This is used by the Flutter app to verify backend connectivity on startup.

---

## Starting the Backend

### Prerequisites
- Go 1.21+
- PostgreSQL (via Neon)
- Environment variables configured in `.env`

### Commands

```bash
# Run the server
go run main.go

# Build production binary
go build -o attendance-backend

# Run built binary
./attendance-backend
```

### Expected Output
```
🚀 Server running on port 3000
Health: http://localhost:3000/health
```

---

## .env Configuration

Current `.env` file:
```env
DATABASE_URL=postgresql://neondb_owner:npg_1FqAb4McvrIU@ep-cold-mountain-a12ximyd-pooler.ap-southeast-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require
APP_PORT=3000
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=rhizalhose@gmail.com
SMTP_PASS=fkjhnxvdjqdpvtvi
```

⚠️ **Security Note**: Change SMTP credentials and database URL for production!

---

## Database Seeding

The backend automatically seeds a superadmin user on startup:

```go
// In main.go
seeder.SeedSuperAdmin()
```

**Default Superadmin Credentials** (from seed.go):
- Student ID: `ADMIN001`
- Email: `admin@attendance.local`
- Password: Check `seeder/seed.go` for default password

---

## Rate Limiting

Rate limiting is enabled by default:

```go
app.Use(middleware.RateLimit)
```

Configuration is in `middleware/rate_limit.go`

---

## Request Logging

All requests are logged:

```go
app.Use(middleware.RequestLogger())
```

Logs are stored in the `logs/` directory.

---

## Security Headers

Security headers are automatically added:

```go
app.Use(middleware.SecurityHeaders)
```

Headers include:
- X-Content-Type-Options: nosniff
- X-Frame-Options: DENY
- X-XSS-Protection: 1; mode=block

---

## Flutter App Connection Flow

```
1. App starts → main.dart
   ↓
2. UserSession.initialize() - Load saved token
   ↓
3. ApiService.healthCheck() - Verify backend
   ↓
4. If token exists → Navigate to dashboard
   If no token → Navigate to login
   ↓
5. Login → AuthService.login()
   ↓
6. Backend returns JWT token
   ↓
7. Token saved to secure storage
   ↓
8. Navigate to dashboard
   ↓
9. All subsequent API calls include JWT token
```

---

## Troubleshooting Connection Issues

### Backend not responding

**Error**: `Connection refused`

**Solution**:
```bash
# 1. Check if backend is running
ps aux | grep "go run"

# 2. Check if port 3000 is in use
netstat -an | grep 3000

# 3. Start backend
cd Attendancebackend
go run main.go

# 4. Test health endpoint
curl http://localhost:3000/health
```

### CORS errors in Flutter

**Error**: `CORS policy error`

**Solution**:
- Ensure CORS is enabled in backend (it is)
- Check if backend is running
- Verify `AllowOrigins` includes your client origin

### Token not working

**Error**: `401 Unauthorized`

**Solution**:
1. Verify token is saved to secure storage:
   ```dart
   final token = await UserSession.getToken();
   print('Token: $token');
   ```

2. Verify token is included in request:
   - Check `api_service.dart` `_getHeaders()` method
   - Token should be in `Authorization: Bearer {token}` header

3. Check token expiration:
   - Backend may have token TTL
   - Implement refresh token mechanism

---

## API Response Format

All backend endpoints follow this response format:

### Success (2xx)
```json
{
  "status": "success",
  "message": "Operation successful",
  "data": { /* response data */ }
}
```

### Error (4xx/5xx)
```json
{
  "status": "error",
  "error": "Error message",
  "code": "ERROR_CODE"
}
```

---

## Default API Responses

### Login Success
```json
{
  "user_id": "uuid",
  "student_id": "2024-001",
  "email": "user@example.com",
  "role": "student|admin|superadmin",
  "access_token": "jwt_token",
  "refresh_token": "refresh_token",
  "message": "Login successful",
  "status": "success"
}
```

### Login Failure (401)
```json
{
  "error": "Invalid credentials",
  "status": "error"
}
```

### Registration Success
```json
{
  "student_id": "2024-001",
  "message": "Registration successful. Please check your email to verify your account.",
  "status": "success"
}
```

---

## Environment-Specific Configuration

### Development
```env
APP_PORT=3000
DATABASE_URL=postgresql://localhost:5432/attendance_dev
CORS_ORIGINS=*
```

### Staging
```env
APP_PORT=3000
DATABASE_URL=postgresql://staging-db.example.com:5432/attendance
CORS_ORIGINS=https://staging.yourdomain.com
```

### Production
```env
APP_PORT=8080
DATABASE_URL=postgresql://prod-db.example.com:5432/attendance
CORS_ORIGINS=https://yourdomain.com
```

---

## Monitoring & Logging

### Logs Location
- Path: `Attendancebackend/logs/`
- Format: JSON logs with timestamps

### Monitoring Endpoints
```
GET /health - Server health status
GET /         - Server info
```

---

## Deployment Checklist

- [ ] Update `.env` with production credentials
- [ ] Change `CORS_ORIGINS` from `*` to specific domain
- [ ] Set up HTTPS/SSL certificate
- [ ] Configure database backups
- [ ] Set up monitoring and alerting
- [ ] Configure log aggregation
- [ ] Update Flutter base URL to production domain
- [ ] Test all endpoints before going live
- [ ] Set up CI/CD pipeline

