# CRM Service Deployment Guide

Complete guide for deploying and running the CRM Service in development and production environments.

## Prerequisites

- **Go**: Version 1.21 or higher
- **PostgreSQL**: Version 14 or higher (for production)
- **SQLite**: For development/testing (included)
- **Git**: For version control

## Environment Setup

### 1. Clone the Repository

```bash
git clone https://github.com/softclub-go-0-0/crm-service.git
cd crm-service
```

### 2. Install Dependencies

```bash
go mod download
go mod tidy
```

### 3. Configure Environment Variables

Create a `.env` file in the project root:

```env
# Server Configuration
SERVER_PORT=8080
SERVER_ENVIRONMENT=development
LOG_LEVEL=debug

# Database Configuration (PostgreSQL)
DB_HOST=localhost
DB_PORT=5432
DB_USER=crm_user
DB_PASSWORD=your_secure_password
DB_NAME=crm_service
DB_SSLMODE=disable

# API Security
API_KEY=your-secret-api-key-here

# File Storage
UPLOAD_DIR=./uploads
MAX_FILE_SIZE=10485760  # 10MB in bytes

# Email Configuration (Optional)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=noreply@yourcompany.com

# SMS Configuration (Optional)
SMS_PROVIDER=twilio
SMS_API_KEY=your-sms-api-key
SMS_FROM=+1234567890
```

### 4. Database Setup

#### PostgreSQL (Production)

```bash
# Create database
createdb crm_service

# Or using psql
psql -U postgres
CREATE DATABASE crm_service;
CREATE USER crm_user WITH PASSWORD 'your_secure_password';
GRANT ALL PRIVILEGES ON DATABASE crm_service TO crm_user;
\q
```

#### SQLite (Development)

SQLite will automatically create the database file on first run.

## Running the Application

### Development Mode

```bash
# Run with hot reload (if using air)
air

# Or run directly
go run cmd/api/main.go
```

The server will start on `http://localhost:8080`

### Production Build

```bash
# Build the binary
go build -o crm-service cmd/api/main.go

# Run the binary
./crm-service
```

### Using Docker (Recommended for Production)

Create a `Dockerfile`:

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o crm-service cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/crm-service .
COPY --from=builder /app/.env .

EXPOSE 8080
CMD ["./crm-service"]
```

Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:14-alpine
    environment:
      POSTGRES_DB: crm_service
      POSTGRES_USER: crm_user
      POSTGRES_PASSWORD: your_secure_password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  crm-service:
    build: .
    ports:
      - "8080:8080"
    depends_on:
      - postgres
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: crm_user
      DB_PASSWORD: your_secure_password
      DB_NAME: crm_service
    volumes:
      - ./uploads:/root/uploads

volumes:
  postgres_data:
```

Run with Docker Compose:

```bash
docker-compose up -d
```

## Database Migrations

The application automatically runs migrations on startup. All models are migrated using GORM's AutoMigrate feature.

To manually verify migrations:

```bash
# Connect to database
psql -U crm_user -d crm_service

# List tables
\dt

# Check specific table
\d students
```

## Testing

### Run All Tests

```bash
go test ./... -v
```

### Run Specific Test

```bash
go test ./tests -v -run TestCRMWorkflow
```

### Run with Coverage

```bash
go test ./... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## API Documentation

### Swagger UI

Access Swagger documentation at:
```
http://localhost:8080/swagger/index.html
```

### Generate/Update Swagger Docs

```bash
# Install swag
go install github.com/swaggo/swag/cmd/swag@latest

# Generate docs
swag init -g cmd/api/main.go -o docs
```

## Monitoring & Logging

### Health Check

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "healthy",
  "timestamp": "2025-11-21T19:00:00Z",
  "database": "connected"
}
```

### Logs

Logs are written to stdout. In production, use a log aggregation service:

```bash
# View logs with Docker
docker-compose logs -f crm-service

# View logs with systemd
journalctl -u crm-service -f
```

## Production Deployment Checklist

- [ ] Set `SERVER_ENVIRONMENT=production`
- [ ] Use strong, unique `API_KEY`
- [ ] Configure PostgreSQL with SSL (`DB_SSLMODE=require`)
- [ ] Set up database backups
- [ ] Configure reverse proxy (Nginx/Caddy)
- [ ] Enable HTTPS/TLS
- [ ] Set up monitoring (Prometheus/Grafana)
- [ ] Configure log aggregation (ELK/Loki)
- [ ] Set up CI/CD pipeline
- [ ] Configure file storage (S3/MinIO for documents)
- [ ] Set resource limits (CPU/Memory)
- [ ] Enable rate limiting
- [ ] Configure CORS properly
- [ ] Set up email/SMS providers
- [ ] Create admin user accounts
- [ ] Test disaster recovery procedures

## Nginx Reverse Proxy Configuration

```nginx
server {
    listen 80;
    server_name api.yourcompany.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## Systemd Service (Linux)

Create `/etc/systemd/system/crm-service.service`:

```ini
[Unit]
Description=CRM Service API
After=network.target postgresql.service

[Service]
Type=simple
User=crm
WorkingDirectory=/opt/crm-service
ExecStart=/opt/crm-service/crm-service
Restart=always
RestartSec=5
Environment="SERVER_ENVIRONMENT=production"
EnvironmentFile=/opt/crm-service/.env

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl enable crm-service
sudo systemctl start crm-service
sudo systemctl status crm-service
```

## Backup & Recovery

### Database Backup

```bash
# Automated daily backup
pg_dump -U crm_user crm_service > backup_$(date +%Y%m%d).sql

# Restore from backup
psql -U crm_user crm_service < backup_20251121.sql
```

### File Storage Backup

```bash
# Backup uploads directory
tar -czf uploads_backup_$(date +%Y%m%d).tar.gz uploads/
```

## Troubleshooting

### Database Connection Issues

```bash
# Check PostgreSQL is running
sudo systemctl status postgresql

# Test connection
psql -U crm_user -d crm_service -h localhost
```

### Port Already in Use

```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>
```

### Migration Errors

```bash
# Drop and recreate database (CAUTION: Data loss!)
dropdb crm_service
createdb crm_service
```

## Performance Optimization

### Database Indexing

Key indexes are automatically created by GORM. Monitor slow queries:

```sql
-- Enable slow query logging in PostgreSQL
ALTER SYSTEM SET log_min_duration_statement = 1000;
SELECT pg_reload_conf();
```

### Connection Pooling

Configure in code or environment:

```go
sqlDB.SetMaxOpenConns(25)
sqlDB.SetMaxIdleConns(5)
sqlDB.SetConnMaxLifetime(5 * time.Minute)
```

## Security Best Practices

1. **Never commit `.env` files** - Add to `.gitignore`
2. **Rotate API keys regularly**
3. **Use HTTPS in production**
4. **Implement rate limiting**
5. **Validate all inputs**
6. **Use prepared statements** (GORM does this automatically)
7. **Regular security audits**
8. **Keep dependencies updated**

## Support & Maintenance

### Update Dependencies

```bash
go get -u ./...
go mod tidy
```

### Monitor Application

```bash
# Check memory usage
ps aux | grep crm-service

# Check disk usage
df -h
du -sh uploads/
```

## Additional Resources

- [API Documentation](./API.md)
- [Go Documentation](https://golang.org/doc/)
- [GORM Documentation](https://gorm.io/docs/)
- [Gin Framework](https://gin-gonic.com/docs/)

---

**Last Updated**: November 21, 2025
**Version**: 1.0.0
