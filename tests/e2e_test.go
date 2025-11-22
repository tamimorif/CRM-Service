package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite" // SQLite for testing
	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/handlers"
	"github.com/softclub-go-0-0/crm-service/pkg/middlewares"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
	"github.com/softclub-go-0-0/crm-service/pkg/services"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupRouter() *gin.Engine {
	// Set up test config
	os.Setenv("SERVER_ENVIRONMENT", "test")
	os.Setenv("LOG_LEVEL", "debug")

	// cfg := &config.Config{ ... } ignored as we use SQLite directly

	// Initialize SQLite In-Memory DB
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic("failed to connect to database: " + err.Error())
	}

	// Auto Migrate all models
	err = db.AutoMigrate(
		&models.Teacher{},
		&models.Course{},
		&models.Student{},
		&models.Group{},
		&models.Timetable{},
		&models.Attendance{},
		&models.Grade{},
		&models.User{},
		&models.AuditLog{},
		&models.Session{},
		&models.Payment{},
		&models.Invoice{},
		&models.Discount{},
		&models.Scholarship{},
		&models.Notification{},
		&models.NotificationTemplate{},
		&models.Document{},
		&models.Message{},
		&models.Event{},
		&models.Application{},
		&models.Exam{},
		&models.ExamResult{},
	)
	if err != nil {
		panic("failed to migrate database: " + err.Error())
	}

	// Initialize services
	teacherService := services.NewTeacherService(db)
	courseService := services.NewCourseService(db)
	studentService := services.NewStudentService(db)
	groupService := services.NewGroupService(db)
	timetableService := services.NewTimetableService(db)
	attendanceService := services.NewAttendanceService(db)
	gradeService := services.NewGradeService(db)
	healthService := services.NewHealthService(db)
	userService := services.NewUserService(db)
	auditService := services.NewAuditService(db)
	sessionService := services.NewSessionService(db)
	paymentService := services.NewPaymentService(db)
	invoiceService := services.NewInvoiceService(db)
	notificationService := services.NewNotificationService(db)
	templateService := services.NewTemplateService(db)
	documentService := services.NewDocumentService(db)
	messageService := services.NewMessageService(db)
	analyticsService := services.NewAnalyticsService(db)
	calendarService := services.NewCalendarService(db)
	applicationService := services.NewApplicationService(db)
	examService := services.NewExamService(db)
	portalService := services.NewPortalService(db)

	h := handlers.NewHandler(
		teacherService,
		courseService,
		studentService,
		groupService,
		timetableService,
		attendanceService,
		gradeService,
		healthService,
		userService,
		auditService,
		sessionService,
		paymentService,
		invoiceService,
		notificationService,
		templateService,
		documentService,
		messageService,
		analyticsService,
		calendarService,
		applicationService,
		examService,
		portalService,
	)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middlewares.RequestID())
	router.Use(middlewares.Logging())
	router.Use(middlewares.Recovery())

	// Skip Auth for testing simplicity
	// router.Use(middlewares.AuthMiddleware(cfg))

	router.POST("/teachers", h.CreateTeacher)
	router.POST("/courses", h.CreateCourse)
	router.POST("/timetables", h.CreateTimetable)

	groups := router.Group("/groups")
	{
		groups.POST("/", h.CreateGroup)
		groups.POST("/:groupID/attendance", h.MarkAttendance)
		groups.POST("/:groupID/grades", h.CreateGrade)

		students := groups.Group("/:groupID/students")
		{
			students.POST("/", h.CreateStudent)
		}
	}

	return router
}

func TestCRMWorkflow(t *testing.T) {
	router := setupRouter()

	// 1. Create Teacher
	teacherReq := dto.CreateTeacherRequest{
		Name:    "John",
		Surname: "Doe",
		Phone:   "992900000001",
		Email:   "john.doe@example.com",
	}
	teacherID := performRequest(t, router, "POST", "/teachers", teacherReq)
	assert.NotEmpty(t, teacherID)

	// 2. Create Course
	courseReq := dto.CreateCourseRequest{
		Title:      "Go Programming",
		MonthlyFee: 100.0,
		Duration:   3,
	}
	courseID := performRequest(t, router, "POST", "/courses", courseReq)
	assert.NotEmpty(t, courseID)

	// 3. Create Timetable
	timetableReq := dto.CreateTimetableRequest{
		StartTime: "09:00",
		EndTime:   "11:00",
		Days:      "Mon,Wed,Fri",
		Classroom: "Room 101",
	}
	timetableID := performRequest(t, router, "POST", "/timetables", timetableReq)
	assert.NotEmpty(t, timetableID)

	// 4. Create Group
	groupReq := dto.CreateGroupRequest{
		Name:        "Go-101",
		StartDate:   time.Now(),
		CourseID:    parseUUID(courseID),
		TeacherID:   parseUUID(teacherID),
		TimetableID: parseUUID(timetableID),
		Capacity:    20,
	}
	groupID := performRequest(t, router, "POST", "/groups/", groupReq)
	assert.NotEmpty(t, groupID)

	// 5. Create Student
	studentReq := dto.CreateStudentRequest{
		Name:    "Alice",
		Surname: "Smith",
		Phone:   "992900000002",
		Email:   "alice@example.com",
	}
	studentID := performRequest(t, router, "POST", fmt.Sprintf("/groups/%s/students/", groupID), studentReq)
	assert.NotEmpty(t, studentID)

	// 6. Mark Attendance
	attendanceReq := dto.CreateAttendanceRequest{
		StudentID: parseUUID(studentID),
		Date:      time.Now().Format("2006-01-02"),
		Status:    "present",
		Notes:     "On time",
	}
	performRequest(t, router, "POST", fmt.Sprintf("/groups/%s/attendance", groupID), attendanceReq)

	// 7. Add Grade
	gradeReq := dto.CreateGradeRequest{
		StudentID: parseUUID(studentID),
		Value:     95,
		Type:      "exam",
		Date:      time.Now().Format("2006-01-02"),
		Notes:     "Excellent work",
	}
	performRequest(t, router, "POST", fmt.Sprintf("/groups/%s/grades", groupID), gradeReq)
}

func performRequest(t *testing.T, router *gin.Engine, method, path string, body interface{}) string {
	jsonValue, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check for success status codes (200 or 201)
	if w.Code != http.StatusOK && w.Code != http.StatusCreated {
		t.Logf("Request failed with status %d: %s", w.Code, w.Body.String())
	}
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusCreated)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	if id, ok := response["id"].(string); ok {
		return id
	}
	// Check if data.id exists (wrapped response)
	if data, ok := response["data"].(map[string]interface{}); ok {
		if id, ok := data["id"].(string); ok {
			return id
		}
	}

	return ""
}

func parseUUID(s string) uuid.UUID {
	id, _ := uuid.Parse(s)
	return id
}
