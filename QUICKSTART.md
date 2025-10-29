# Quick Start Guide - CRM Service

This guide will get you up and running in 5 minutes!

## Prerequisites
- Docker & Docker Compose installed
- Or Go 1.20+ and PostgreSQL

---

## Method 1: Docker (Recommended - 2 minutes) 🐳

### Step 1: Clone and Configure
```bash
git clone https://github.com/tamimorif/CRM-Service.git
cd CRM-Service
cp .env.example .env
```

### Step 2: Start Services
```bash
docker-compose up -d
```

### Step 3: Verify
```bash
# Check health
curl http://localhost:8080/health

# You should see: {"status":"healthy",...}
```

That's it! Your CRM service is running on `http://localhost:8080`

### Accessing Services:
- **API**: http://localhost:8080
- **PgAdmin** (Database UI): http://localhost:5050
  - Email: admin@crm.local
  - Password: admin

### View Logs:
```bash
docker-compose logs -f crm-service
```

### Stop Services:
```bash
docker-compose down
```

---

## Method 2: Local Development (5 minutes) 💻

### Step 1: Setup Database
Install and start PostgreSQL, then create a database:
```bash
psql -U postgres
CREATE DATABASE crm_service;
\q
```

### Step 2: Clone and Configure
```bash
git clone https://github.com/tamimorif/CRM-Service.git
cd CRM-Service
cp .env.example .env
```

### Step 3: Edit .env
```env
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=crm_service
DB_PORT=5432
APP_PORT=8080
AUTH_SERVICE_ADDR=localhost:50051
SKIP_AUTH=true  # For development without auth service
```

### Step 4: Install Dependencies
```bash
go mod download
```

### Step 5: Run the Application
```bash
# Simple run
go run ./cmd/api/main.go

# Or with hot reload (requires air)
make install-tools
make dev
```

### Step 6: Verify
```bash
curl http://localhost:8080/health
```

---

## Testing the API 🧪

### Without Authentication (Development)
Set `SKIP_AUTH=true` in your `.env` file.

### Test Endpoints:

#### 1. Health Check (No auth needed)
```bash
curl http://localhost:8080/health
```

#### 2. Create a Course
```bash
curl -X POST http://localhost:8080/courses \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Web Development Bootcamp",
    "monthly_fee": 1500,
    "duration": 6
  }'
```

#### 3. Get All Courses with Pagination
```bash
curl "http://localhost:8080/courses?page=1&page_size=10&sort=title&order=asc"
```

#### 4. Create a Teacher
```bash
curl -X POST http://localhost:8080/teachers \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John",
    "surname": "Doe",
    "phone": "992900123456",
    "email": "john.doe@example.com"
  }'
```

#### 5. Create a Timetable
```bash
curl -X POST http://localhost:8080/timetables \
  -H "Content-Type: application/json" \
  -d '{
    "classroom": "Room 101",
    "start": "09:00:00",
    "finish": "11:00:00"
  }'
```

#### 6. Create a Group
```bash
curl -X POST http://localhost:8080/groups \
  -H "Content-Type: application/json" \
  -d '{
    "course_id": "YOUR_COURSE_ID",
    "teacher_id": "YOUR_TEACHER_ID",
    "timetable_id": "YOUR_TIMETABLE_ID",
    "title": "Web Dev Group 1",
    "start_date": "2025-11-01"
  }'
```

#### 7. Add a Student to Group
```bash
curl -X POST http://localhost:8080/groups/YOUR_GROUP_ID/students \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane",
    "surname": "Smith",
    "phone": "992900654321",
    "email": "jane.smith@example.com"
  }'
```

#### 8. Search Teachers
```bash
curl "http://localhost:8080/teachers?search=john&page=1&page_size=10"
```

---

## Common Issues & Solutions 🔧

### Issue: "Cannot connect to database"
**Solution**: 
- Check PostgreSQL is running
- Verify credentials in `.env`
- For Docker: `docker-compose ps` to check postgres status

### Issue: "Authentication failed"
**Solution**:
- For development, set `SKIP_AUTH=true` in `.env`
- Or ensure auth service is running at the configured address

### Issue: "Port already in use"
**Solution**:
- Change `APP_PORT` in `.env` to a different port
- Or stop the process using the port: `lsof -ti:8080 | xargs kill`

### Issue: "Docker build failed"
**Solution**:
```bash
docker-compose down -v
docker-compose build --no-cache
docker-compose up -d
```

---

## Useful Commands 📝

### Docker Commands
```bash
# Start services
docker-compose up -d

# Stop services
docker-compose down

# View logs
docker-compose logs -f

# Restart specific service
docker-compose restart crm-service

# Rebuild everything
docker-compose down -v
docker-compose up -d --build
```

### Make Commands
```bash
# Show all commands
make help

# Build project
make build

# Run locally
make run

# Run with hot reload
make dev

# Run tests
make test

# Format code
make fmt

# Start Docker
make docker-up

# Stop Docker
make docker-down
```

### Database Access
```bash
# Via Docker
docker exec -it crm-postgres psql -U crm_user -d crm_service

# Or access PgAdmin at http://localhost:5050
```

---

## API Documentation 📚

For complete API documentation, see:
- **API.md** - Full API reference
- **README.md** - Project overview
- **COMPLETION.md** - What's included

---

## Development Workflow 🔄

### 1. Start Development Environment
```bash
docker-compose up -d postgres pgadmin
make dev
```

### 2. Make Changes
- Edit code
- Changes auto-reload (with `make dev`)

### 3. Test Changes
```bash
make test
curl http://localhost:8080/health
```

### 4. Format & Lint
```bash
make fmt
make lint
```

### 5. Commit
```bash
git add .
git commit -m "Your changes"
git push
```

---

## Production Deployment 🚀

### Quick Production Start
```bash
# Update .env for production
SKIP_AUTH=false
AUTH_SERVICE_ADDR=your-auth-service:50051

# Build and run
docker-compose -f docker-compose.yml up -d
```

### Or build production binary
```bash
make prod-build
./crm-service
```

---

## Need Help? 🆘

1. Check logs: `docker-compose logs -f`
2. Read error messages carefully
3. Verify `.env` configuration
4. Check database connection
5. Ensure ports are not in use
6. Review `API.md` for endpoint details

---

## Next Steps 🎯

1. ✅ Get the service running
2. ✅ Test endpoints with curl or Postman
3. ✅ Read API.md for all endpoints
4. ✅ Explore the code structure
5. ✅ Add your custom features

---

**You're all set! Happy coding! 🎉**

For questions: Open an issue on GitHub
Author: Tamim Orif (@tamimorif)