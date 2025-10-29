# CRM Service API Documentation

Base URL: `http://localhost:8080`

## Authentication

All endpoints (except health checks) require authentication via token in the header:
```
X-Auth-Token: your_token_here
```
or
```
Authorization: Bearer your_token_here
```

---

## Health Check Endpoints

### GET /health
Get comprehensive health status of the application

**Response:**
```json
{
  "status": "healthy",
  "version": "1.0.0",
  "timestamp": "2025-10-29T10:00:00Z",
  "uptime": "1h30m20s",
  "database": {
    "status": "healthy",
    "response_time_ms": 5
  },
  "system": {
    "go_version": "go1.20",
    "num_goroutine": 10,
    "num_cpu": 8
  },
  "services": {
    "auth_service": "healthy"
  }
}
```

### GET /ready
Readiness probe for Kubernetes

### GET /live
Liveness probe for Kubernetes

---

## Teachers

### GET /teachers
Get all teachers with pagination, search, and sorting

**Query Parameters:**
- `page` (int, default: 1) - Page number
- `page_size` (int, default: 10, max: 100) - Items per page
- `sort` (string, default: "created_at") - Sort field
- `order` (string, default: "asc") - Sort order (asc/desc)
- `search` (string) - Search in name, surname, phone, email

**Example:** `/teachers?page=1&page_size=10&search=john&sort=name&order=asc`

**Response:**
```json
{
  "success": true,
  "message": "Teachers retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "name": "John",
      "surname": "Doe",
      "phone": "992900123456",
      "email": "john@example.com",
      "created_at": "2025-10-29T10:00:00Z",
      "updated_at": "2025-10-29T10:00:00Z",
      "groups": []
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 10,
    "total_pages": 5,
    "total_count": 47,
    "has_next": true,
    "has_prev": false
  }
}
```

### POST /teachers
Create a new teacher

**Request Body:**
```json
{
  "name": "John",
  "surname": "Doe",
  "phone": "992900123456",
  "email": "john@example.com"
}
```

### GET /teachers/:teacherID
Get a specific teacher by ID

### PUT /teachers/:teacherID
Update a teacher

**Request Body:**
```json
{
  "name": "John",
  "surname": "Doe",
  "phone": "992900123456",
  "email": "john@example.com"
}
```

### DELETE /teachers/:teacherID
Delete a teacher

---

## Courses

### GET /courses
Get all courses with pagination and search

**Query Parameters:**
- `page`, `page_size`, `sort`, `order` - Same as teachers
- `search` - Search in course title

### POST /courses
Create a new course

**Request Body:**
```json
{
  "title": "Web Development",
  "monthly_fee": 1000,
  "duration": 6
}
```

### GET /courses/:courseID
Get a specific course

### PUT /courses/:courseID
Update a course

### DELETE /courses/:courseID
Delete a course

---

## Timetables

### GET /timetables
Get all timetables

### POST /timetables
Create a new timetable

**Request Body:**
```json
{
  "classroom": "Room 101",
  "start": "09:00:00",
  "finish": "11:00:00"
}
```

### GET /timetables/:timetableID
Get a specific timetable

### PUT /timetables/:timetableID
Update a timetable

### DELETE /timetables/:timetableID
Delete a timetable

---

## Groups

### GET /groups
Get all groups with related data (course, teacher, timetable, students)

### POST /groups
Create a new group

**Request Body:**
```json
{
  "course_id": "uuid",
  "teacher_id": "uuid",
  "timetable_id": "uuid",
  "title": "Web Dev Group 1",
  "start_date": "2025-11-01"
}
```

### GET /groups/:groupID
Get a specific group with all relations

### PUT /groups/:groupID
Update a group

### DELETE /groups/:groupID
Delete a group

---

## Students

### GET /students
Get all students globally

### GET /groups/:groupID/students
Get all students in a specific group

### POST /groups/:groupID/students
Create a new student in a group

**Request Body:**
```json
{
  "name": "Jane",
  "surname": "Smith",
  "phone": "992900654321",
  "email": "jane@example.com"
}
```

### GET /groups/:groupID/students/:studentID
Get a specific student

### PUT /groups/:groupID/students/:studentID
Update a student

**Request Body:**
```json
{
  "group_id": "uuid",
  "name": "Jane",
  "surname": "Smith",
  "phone": "992900654321",
  "email": "jane@example.com"
}
```

### DELETE /groups/:groupID/students/:studentID
Delete a student

---

## Error Responses

All endpoints return structured error responses:

```json
{
  "success": false,
  "message": "Error message",
  "errors": "Detailed error information",
  "timestamp": "2025-10-29T10:00:00Z"
}
```

**Common HTTP Status Codes:**
- 200 OK - Request successful
- 201 Created - Resource created successfully
- 400 Bad Request - Invalid request
- 401 Unauthorized - Authentication required
- 403 Forbidden - Insufficient permissions
- 404 Not Found - Resource not found
- 409 Conflict - Resource already exists
- 422 Unprocessable Entity - Validation error
- 500 Internal Server Error - Server error

---

## Development

### Running with Docker
```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Running locally
```bash
# Install dependencies
make deps

# Run the application
make run

# Run with hot reload
make dev

# Run tests
make test
```

### Environment Variables
See `.env.example` for all required environment variables.