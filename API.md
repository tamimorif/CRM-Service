# CRM Service API Documentation

Complete API reference for the CRM Service with 67+ endpoints across 10 feature phases.

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication

All endpoints require API Key authentication via header:
```
X-API-Key: your-api-key-here
```

Some endpoints also require role-based access control (RBAC).

---

## 📚 Core Entities

### Teachers
- `POST /teachers` - Create teacher
- `GET /teachers` - List all teachers (paginated)
- `GET /teachers/:id` - Get teacher details
- `PUT /teachers/:id` - Update teacher
- `DELETE /teachers/:id` - Delete teacher

### Courses
- `POST /courses` - Create course
- `GET /courses` - List all courses (paginated)
- `GET /courses/:id` - Get course details
- `PUT /courses/:id` - Update course
- `DELETE /courses/:id` - Delete course

### Students
- `POST /students` - Create student
- `GET /students` - List all students (paginated)
- `GET /students/:id` - Get student details
- `PUT /students/:id` - Update student
- `DELETE /students/:id` - Delete student

### Groups
- `POST /groups` - Create group
- `GET /groups` - List all groups (paginated)
- `GET /groups/:id` - Get group details
- `PUT /groups/:id` - Update group
- `DELETE /groups/:id` - Delete group

---

## 📅 Scheduling & Attendance

### Timetables
- `POST /timetables` - Create timetable
- `GET /timetables` - List all timetables
- `GET /timetables/:id` - Get timetable details
- `PUT /timetables/:id` - Update timetable
- `DELETE /timetables/:id` - Delete timetable

### Attendance
- `POST /attendance` - Mark attendance
- `GET /attendance` - List attendance records
- `GET /attendance/:id` - Get attendance details
- `GET /attendance/student/:studentID` - Get student attendance
- `GET /attendance/group/:groupID` - Get group attendance
- `PUT /attendance/:id` - Update attendance
- `DELETE /attendance/:id` - Delete attendance

---

## 📊 Grades & Exams

### Grades
- `POST /grades` - Create grade
- `GET /grades` - List all grades
- `GET /grades/:id` - Get grade details
- `GET /grades/student/:studentID` - Get student grades
- `PUT /grades/:id` - Update grade
- `DELETE /grades/:id` - Delete grade

### Exams
- `POST /exams` - Create exam
- `GET /exams` - List all exams (paginated)
- `GET /exams/:examID` - Get exam details
- `GET /exams/course/:courseID` - Get exams by course
- `GET /exams/group/:groupID` - Get exams by group
- `PUT /exams/:examID` - Update exam
- `POST /exams/:examID/results` - Submit exam result
- `GET /exams/:examID/results` - Get all exam results
- `GET /exams/student/:studentID/results` - Get student results
- `GET /exams/:examID/statistics` - Get exam statistics
- `DELETE /exams/:examID` - Delete exam

---

## 💰 Financial Management

### Payments
- `POST /payments` - Create payment
- `GET /payments` - List all payments (paginated)
- `GET /payments/:id` - Get payment details
- `GET /payments/student/:studentID` - Get student payments
- `PUT /payments/:id` - Update payment
- `DELETE /payments/:id` - Delete payment

### Invoices
- `POST /invoices` - Create invoice
- `GET /invoices` - List all invoices (paginated)
- `GET /invoices/:id` - Get invoice details
- `GET /invoices/student/:studentID` - Get student invoices
- `PUT /invoices/:id` - Update invoice
- `DELETE /invoices/:id` - Delete invoice

---

## 🔔 Notifications

### Notifications
- `POST /notifications` - Send notification
- `GET /notifications` - List all notifications (paginated)
- `GET /notifications/:id` - Get notification details
- `GET /notifications/user/:userID` - Get user notifications
- `PUT /notifications/:id/read` - Mark as read
- `DELETE /notifications/:id` - Delete notification

### Notification Templates
- `POST /notification-templates` - Create template
- `GET /notification-templates` - List all templates
- `GET /notification-templates/:id` - Get template details
- `PUT /notification-templates/:id` - Update template
- `DELETE /notification-templates/:id` - Delete template

---

## 📄 Document Management

### Documents
- `POST /documents` - Upload document
- `GET /documents` - List all documents (paginated)
- `GET /documents/:id` - Get document details
- `GET /documents/:id/download` - Download document
- `GET /documents/entity/:entityType/:entityID` - Get entity documents
- `PUT /documents/:id` - Update document metadata
- `PUT /documents/:id/approve` - Approve document
- `PUT /documents/:id/reject` - Reject document
- `DELETE /documents/:id` - Delete document

---

## 💬 Communication Hub

### Messages
- `POST /messages` - Send message
- `GET /messages` - List all messages (paginated)
- `GET /messages/:id` - Get message details
- `GET /messages/inbox` - Get inbox messages
- `GET /messages/sent` - Get sent messages
- `PUT /messages/:id/read` - Mark as read
- `DELETE /messages/:id` - Delete message

---

## 📈 Analytics & Reporting

### Analytics
- `GET /analytics/dashboard` - Get dashboard metrics
- `GET /analytics/financial` - Get financial analytics
- `GET /analytics/student-progress/:studentID` - Get student progress
- `GET /analytics/attendance/:groupID` - Get attendance analytics

---

## 📅 Calendar & Events

### Events
- `POST /events` - Create event
- `GET /events` - List all events (paginated)
- `GET /events/:id` - Get event details
- `GET /events/calendar` - Get calendar events (date range)
- `PUT /events/:id` - Update event
- `DELETE /events/:id` - Delete event

---

## 🎓 Enrollment & Admission

### Applications
- `POST /applications` - Submit application
- `GET /applications` - List all applications (paginated)
- `GET /applications/:applicationID` - Get application details
- `GET /applications/course/:courseID` - Get applications by course
- `GET /applications/status/:status` - Get applications by status
- `PUT /applications/:applicationID` - Update application
- `POST /applications/:applicationID/review` - Review application
- `POST /applications/:applicationID/enroll` - Enroll student
- `DELETE /applications/:applicationID` - Delete application

---

## 🏠 Student & Teacher Portals

### Portals
- `GET /portal/student/:studentID` - Get student dashboard
- `GET /portal/teacher/:teacherID` - Get teacher dashboard

---

## 👤 User Management

### Users
- `POST /users` - Create user
- `GET /users` - List all users (paginated)
- `GET /users/:id` - Get user details
- `PUT /users/:id` - Update user
- `DELETE /users/:id` - Delete user

### Sessions
- `POST /auth/login` - Login
- `POST /auth/logout` - Logout
- `GET /auth/sessions` - List user sessions
- `DELETE /auth/sessions/:id` - Revoke session

---

## 🔍 Audit Logs

### Audit
- `GET /audit-logs` - Get audit logs (Admin only)

---

## Example Requests

### Create a Student
```bash
curl -X POST http://localhost:8080/api/v1/students \
  -H "X-API-Key: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "+1234567890"
  }'
```

### Submit Exam Result
```bash
curl -X POST http://localhost:8080/api/v1/exams/{examID}/results \
  -H "X-API-Key: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": "uuid-here",
    "marks_obtained": 85,
    "total_marks": 100,
    "remarks": "Excellent performance"
  }'
```

### Get Student Dashboard
```bash
curl -X GET http://localhost:8080/api/v1/portal/student/{studentID} \
  -H "X-API-Key: your-api-key"
```

---

## Error Responses

All endpoints return standard error responses:

```json
{
  "success": false,
  "message": "Error description",
  "data": null
}
```

### HTTP Status Codes
- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `500` - Internal Server Error

---

## Rate Limiting

Currently no rate limiting is implemented. Consider adding in production.

## Pagination

List endpoints support pagination via query parameters:
- `page` - Page number (default: 1)
- `page_size` - Items per page (default: 10)
- `search` - Search term (optional)

Response format:
```json
{
  "data": [...],
  "total": 100,
  "page": 1,
  "page_size": 10,
  "total_pages": 10
}
```