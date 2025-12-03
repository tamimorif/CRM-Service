# **CRM System for Educational Institutions**

A robust and scalable Customer Relationship Management (CRM) system built using **Golang**. This system is specifically designed for educational institutions to manage courses, teachers, students, groups, and timetables efficiently.

---

## **Features**

### Core Features
- ✅ **Complete CRUD Operations** for all entities (Teachers, Students, Courses, Groups, Timetables)
- ✅ **Authentication & Authorization** via gRPC integration with auth service
- ✅ **RESTful API** with Gin framework
- ✅ **Advanced Pagination, Filtering & Search** capabilities
- ✅ **Relationship Management** between entities with GORM preloading
- ✅ **Health Check Endpoints** for monitoring and Kubernetes readiness/liveness probes
- ✅ **Structured Error Handling** with consistent API responses
- ✅ **Docker & Docker Compose** support for easy deployment
- ✅ **PostgreSQL Database** with automatic migrations

### New Features
- ✅ **Parent Portal**: Manage parent profiles and link them to students
- ✅ **Assignment Tracking**: Create and manage homework, projects, and quizzes
- ✅ **Waitlist Management**: Handle course waitlists with priority and notifications
- ✅ **Bulk Operations**: Bulk create students, attendance, and grades
- ✅ **Recurring Invoices**: Automated monthly invoicing system
- ✅ **Advanced Search**: Multi-field filtering and date range support

### Data Models
- **Teachers** - Manage teacher information and their assigned groups
- **Students** - Track student enrollment and group assignments
- **Parents** - Manage parent contact info and student relationships
- **Courses** - Define courses with pricing and duration
- **Groups** - Organize students into groups with assigned teachers and schedules
- **Timetables** - Manage class schedules and classroom assignments
- **Assignments** - Track homework, projects, and grades
- **Waitlists** - Manage prospective students for full courses

---

## **Tech Stack**
- **Backend**: Go 1.24+
- **Framework**: Gin Web Framework
- **Database**: PostgreSQL with GORM
- **Authentication**: gRPC integration with separate auth service
- **Containerization**: Docker & Docker Compose
- **API Protocol**: REST + gRPC
- **Database Tools**: GORM for ORM and migrations

---

## **Quick Start**

### Prerequisites
- Go 1.24 or higher
- PostgreSQL 12+
- Docker & Docker Compose (optional)

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/tamimorif/CRM-Service.git
   cd CRM-Service
   ```

2. **Set up environment variables:**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Install dependencies:**
   ```bash
   make deps
   ```

4. **Run with Docker (Recommended):**
   ```bash
   docker-compose up -d
   ```

   Or **run locally:**
   ```bash
   make run
   ```

The API will be available at `http://localhost:8080`

---

## **API Documentation**

See [API.md](API.md) for comprehensive API documentation including:
- All available endpoints
- Request/response examples
- Authentication details
- Pagination and filtering

### Quick API Examples

**Get all teachers with pagination:**
```bash
curl -H "X-Auth-Token: your_token" \
  "http://localhost:8080/teachers?page=1&page_size=10&search=john"
```

**Create a new course:**
```bash
curl -X POST -H "X-Auth-Token: your_token" \
  -H "Content-Type: application/json" \
  -d '{"title":"Web Development","monthly_fee":1000,"duration":6}' \
  http://localhost:8080/courses
```

**Health check:**
```bash
curl http://localhost:8080/health
```

---

## **Project Structure**
```plaintext
.
├── cmd/
│   ├── api/          # Main API application
│   └── console/      # Console utilities
├── pkg/
│   ├── auth/         # gRPC auth client (generated)
│   ├── config/       # Application configuration
│   ├── database/     # Database configuration
│   ├── dto/          # Data Transfer Objects
│   ├── errors/       # Custom error definitions
│   ├── handlers/     # HTTP request handlers
│   ├── helpers/      # HTTP response helpers & utilities
│   ├── logger/       # Logging configuration
│   ├── middlewares/  # Auth & logging middleware
│   ├── models/       # Database models
│   ├── repository/   # Data access layer
│   └── services/     # Business logic layer
├── tests/            # Integration and unit tests
├── .env.example      # Environment variables template
├── docker-compose.yml
├── Dockerfile
├── Makefile          # Build & development commands
├── API.md            # API documentation
└── README.md
```

---

## **Development**

### Available Make Commands
```bash
make help           # Show all available commands
make build          # Build the application
make run            # Run the application
make dev            # Run with hot reload (requires air)
make test           # Run tests
make test-coverage  # Run tests with coverage report
make docker-up      # Start Docker containers
make docker-down    # Stop Docker containers
make docker-logs    # Show Docker logs
make lint           # Run linter
make fmt            # Format code
```

### Running Tests
```bash
make test
```

### Code Formatting
```bash
make fmt
make lint
```

---

## **Configuration**

### Environment Variables

Create a `.env` file based on `.env.example`:

```env
# Database
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=crm_service
DB_PORT=5432

# Application
APP_PORT=8080

# Auth Service
AUTH_SERVICE_ADDR=localhost:50051

# Development
SKIP_AUTH=false
```

---

## **Deployment**

### Docker Deployment

The project includes a complete Docker setup:

```bash
# Build and start all services
docker-compose up -d

# View logs
docker-compose logs -f crm-service

# Stop services
docker-compose down
```

Services included:
- **crm-service** - Main application
- **postgres** - PostgreSQL database
- **pgadmin** - Database management UI (optional)

### Production Deployment

For production deployment:

1. Update environment variables for production
2. Use the production build:
   ```bash
   make prod-build
   ```
3. Deploy using Docker or directly on servers
4. Configure reverse proxy (nginx/traefik)
5. Set up SSL certificates
6. Configure monitoring and logging

---

## **API Features**

### Pagination
All list endpoints support pagination:
```
GET /teachers?page=1&page_size=10
```

### Search & Filtering
Search across multiple fields:
```
GET /teachers?search=john
GET /courses?search=web
```

### Sorting
Sort by any field:
```
GET /teachers?sort=name&order=asc
GET /courses?sort=created_at&order=desc
```

### Relationship Loading
Related data is automatically loaded for single-item GET requests and can be included in list requests.

---

## **Architecture**

### Microservices Communication
- **CRM Service** (this repository) - Manages educational data
- **Auth Service** (separate) - Handles authentication via gRPC

### Database Schema
- Teachers can have multiple Groups
- Groups belong to one Course, Teacher, and Timetable
- Students belong to one Group
- Soft deletes enabled for all models

---

## **Contributing**

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch:
   ```bash
   git checkout -b feature/amazing-feature
   ```
3. Commit your changes:
   ```bash
   git commit -m "Add amazing feature"
   ```
4. Push to the branch:
   ```bash
   git push origin feature/amazing-feature
   ```
5. Open a Pull Request

### Coding Standards
- Follow Go best practices
- Write tests for new features
- Update documentation
- Run `make fmt` and `make lint` before committing

---

## **Monitoring & Health Checks**

The application provides three health check endpoints:

- **GET /health** - Comprehensive health status with metrics
- **GET /ready** - Kubernetes readiness probe
- **GET /live** - Kubernetes liveness probe

These endpoints are publicly accessible (no authentication required).

---

## **License**

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## **Contact**

**Author**: Tamim Orif  
**GitHub**: [@tamimorif](https://github.com/tamimorif)  
**Project**: [CRM-Service](https://github.com/tamimorif/CRM-Service)

For questions, issues, or collaboration opportunities, feel free to open an issue or reach out!

---

## **Roadmap**

- [ ] Add more comprehensive unit tests
- [ ] Implement GraphQL API
- [ ] Add real-time notifications via WebSocket
- [ ] Implement payment tracking for students
- [ ] Add reporting and analytics dashboard
- [ ] Implement file upload for documents
- [ ] Add email notification system
- [ ] Implement attendance tracking
- [ ] Add grade management system
- [ ] Create admin dashboard UI
