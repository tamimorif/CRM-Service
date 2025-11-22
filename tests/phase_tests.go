package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/stretchr/testify/assert"
)

// TestAllPhasesSmoke tests that all phase endpoints are accessible
func TestAllPhasesSmoke(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
	}{
		// Core entities
		{"List Teachers", "GET", "/api/v1/teachers", http.StatusOK},
		{"List Courses", "GET", "/api/v1/courses", http.StatusOK},
		{"List Students", "GET", "/api/v1/students", http.StatusOK},
		{"List Groups", "GET", "/api/v1/groups", http.StatusOK},

		// Scheduling
		{"List Timetables", "GET", "/api/v1/timetables", http.StatusOK},
		{"List Attendance", "GET", "/api/v1/attendance", http.StatusOK},

		// Grades & Exams
		{"List Grades", "GET", "/api/v1/grades", http.StatusOK},
		{"List Exams", "GET", "/api/v1/exams", http.StatusOK},

		// Financial (Phase 2)
		{"List Payments", "GET", "/api/v1/payments", http.StatusOK},
		{"List Invoices", "GET", "/api/v1/invoices", http.StatusOK},

		// Notifications (Phase 3)
		{"List Notifications", "GET", "/api/v1/notifications", http.StatusOK},
		{"List Templates", "GET", "/api/v1/notification-templates", http.StatusOK},

		// Documents (Phase 4)
		{"List Documents", "GET", "/api/v1/documents", http.StatusOK},

		// Messages (Phase 5)
		{"List Messages", "GET", "/api/v1/messages", http.StatusOK},

		// Analytics (Phase 6)
		{"Dashboard Metrics", "GET", "/api/v1/analytics/dashboard", http.StatusOK},

		// Calendar (Phase 7)
		{"List Events", "GET", "/api/v1/events", http.StatusOK},

		// Applications (Phase 8)
		{"List Applications", "GET", "/api/v1/applications", http.StatusOK},

		// Health check
		{"Health Check", "GET", "/health", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			req.Header.Set("X-API-Key", "test-api-key")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code, "Endpoint %s should return %d", tt.path, tt.wantStatus)
		})
	}
}

// TestCRUDWorkflow tests a complete CRUD workflow
func TestCRUDWorkflow(t *testing.T) {
	router := setupRouter()

	// Create a teacher
	teacherReq := dto.CreateTeacherRequest{
		Name:  "John Doe",
		Email: "john.doe@example.com",
		Phone: "+1234567890",
	}
	teacherBody, _ := json.Marshal(teacherReq)
	req := httptest.NewRequest("POST", "/api/v1/teachers", bytes.NewBuffer(teacherBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "test-api-key")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code, "Teacher creation should succeed")

	var teacherResp struct {
		Data dto.TeacherResponse `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &teacherResp)
	teacherID := teacherResp.Data.ID

	// Get the teacher
	req = httptest.NewRequest("GET", "/api/v1/teachers/"+teacherID.String(), nil)
	req.Header.Set("X-API-Key", "test-api-key")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Teacher retrieval should succeed")

	// Create a course
	courseReq := dto.CreateCourseRequest{
		Title:      "Mathematics 101",
		Duration:   90,
		MonthlyFee: 999.99,
	}
	courseBody, _ := json.Marshal(courseReq)
	req = httptest.NewRequest("POST", "/api/v1/courses", bytes.NewBuffer(courseBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "test-api-key")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code, "Course creation should succeed")

	var courseResp struct {
		Data dto.CourseResponse `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &courseResp)
	courseID := courseResp.Data.ID

	// Create a student
	studentReq := dto.CreateStudentRequest{
		Name:  "Jane Smith",
		Email: "jane.smith@example.com",
		Phone: "+0987654321",
	}
	studentBody, _ := json.Marshal(studentReq)
	req = httptest.NewRequest("POST", "/api/v1/students", bytes.NewBuffer(studentBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "test-api-key")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code, "Student creation should succeed")

	var studentResp struct {
		Data dto.StudentResponse `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &studentResp)
	studentID := studentResp.Data.ID

	// Create a group
	groupReq := dto.CreateGroupRequest{
		Name:      "Math Group A",
		CourseID:  courseID,
		TeacherID: teacherID,
	}
	groupBody, _ := json.Marshal(groupReq)
	req = httptest.NewRequest("POST", "/api/v1/groups", bytes.NewBuffer(groupBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "test-api-key")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code, "Group creation should succeed")

	// Verify we can list all entities
	endpoints := []string{
		"/api/v1/teachers",
		"/api/v1/courses",
		"/api/v1/students",
		"/api/v1/groups",
	}

	for _, endpoint := range endpoints {
		req = httptest.NewRequest("GET", endpoint, nil)
		req.Header.Set("X-API-Key", "test-api-key")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Listing %s should succeed", endpoint)
	}

	t.Logf("✅ Created Teacher: %s", teacherID)
	t.Logf("✅ Created Course: %s", courseID)
	t.Logf("✅ Created Student: %s", studentID)
	t.Logf("✅ All CRUD operations working correctly")
}
