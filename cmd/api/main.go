package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/softclub-go-0-0/crm-service/pkg/config"
	"github.com/softclub-go-0-0/crm-service/pkg/database"
	"github.com/softclub-go-0-0/crm-service/pkg/handlers"
	"github.com/softclub-go-0-0/crm-service/pkg/logger"
	"github.com/softclub-go-0-0/crm-service/pkg/middlewares"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
	"github.com/softclub-go-0-0/crm-service/pkg/services"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           CRM Service API
// @version         1.0
// @description     A robust CRM service for educational institutions.
// @termsOfService  http://swagger.io/terms/

// @contact.name    API Support
// @contact.url     http://www.swagger.io/support
// @contact.email   support@swagger.io

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @host            localhost:8080
// @BasePath        /
// @schemes         http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-Auth-Token
func main() {
	// Initialize logger
	if err := logger.Init(logger.Config{
		Output: "stdout",
		Level:  "debug",
	}); err != nil {
		panic(err)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("cannot load config", err)
	}

	// Initialize database
	db, err := database.DBInit(cfg)
	if err != nil {
		logger.Fatal("cannot connect to database", err)
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

	// Auto-migrate models
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
		logger.Fatal("failed to auto-migrate database models", err)
	}

	// Initialize handlers
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

	// Initialize session handler
	sessionHandler := handlers.NewSessionHandler(sessionService)

	// Initialize router
	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()

	// Register middlewares
	router.Use(middlewares.RequestID())
	router.Use(middlewares.Logging())
	router.Use(middlewares.Recovery())
	router.Use(middlewares.CORS())
	router.Use(middlewares.RateLimit(20, 40)) // 20 req/sec, burst 40

	// Health check endpoints (no auth required)
	router.GET("/health", h.HealthCheck)
	router.GET("/ready", h.ReadinessProbe)
	router.GET("/live", h.LivenessProbe)

	// Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Auth routes (no auth middleware required for login)
	auth := router.Group("/auth")
	{
		auth.POST("/login", sessionHandler.Login)
		// Protected session routes (with auth middleware)
		authProtected := auth.Group("")
		authProtected.Use(middlewares.AuthMiddleware(cfg))
		{
			authProtected.POST("/logout", sessionHandler.Logout)
			authProtected.GET("/sessions", sessionHandler.GetActiveSessions)
			authProtected.DELETE("/sessions/:sessionID", sessionHandler.RevokeSession)
			authProtected.POST("/sessions/revoke-all", sessionHandler.RevokeAllSessions)
		}
	}

	// Auth middleware for all other routes
	router.Use(middlewares.AuthMiddleware(cfg))

	// Routes
	router.GET("/teachers", h.GetAllTeachers)
	router.POST("/teachers", h.CreateTeacher)
	router.GET("/teachers/:teacherID", h.GetOneTeacher)
	router.PUT("/teachers/:teacherID", h.UpdateTeacher)
	router.DELETE("/teachers/:teacherID", h.DeleteTeacher)

	router.GET("/courses", h.GetAllCourses)
	router.POST("/courses", h.CreateCourse)
	router.GET("/courses/:courseID", h.GetOneCourse)
	router.PUT("/courses/:courseID", h.UpdateCourse)
	router.DELETE("/courses/:courseID", h.DeleteCourse)

	router.GET("/timetables", h.GetAllTimetables)
	router.POST("/timetables", h.CreateTimetable)
	router.GET("/timetables/:timetableID", h.GetOneTimetable)
	router.PUT("/timetables/:timetableID", h.UpdateTimetable)
	router.DELETE("/timetables/:timetableID", h.DeleteTimetable)

	// Global students endpoint
	router.GET("/students", h.GetAllStudentsGlobal)

	groups := router.Group("/groups")
	{
		groups.GET("/", h.GetAllGroups)
		groups.POST("/", h.CreateGroup)
		groups.GET("/:groupID", h.GetOneGroup)
		groups.PUT("/:groupID", h.UpdateGroup)
		groups.DELETE("/:groupID", h.DeleteGroup)

		// Attendance routes for group
		groups.POST("/:groupID/attendance", h.MarkAttendance)
		groups.POST("/:groupID/attendance/batch", h.BatchMarkAttendance)
		groups.GET("/:groupID/attendance", h.GetGroupAttendance)

		// Grade routes for group
		groups.POST("/:groupID/grades", h.CreateGrade)
		groups.GET("/:groupID/grades", h.GetGroupGrades)

		students := groups.Group("/:groupID/students")
		{
			students.GET("/", h.GetAllStudents)
			students.POST("/", h.CreateStudent)
			students.GET("/:studentID", h.GetOneStudent)
			students.PUT("/:studentID", h.UpdateStudent)
			students.DELETE("/:studentID", h.DeleteStudent)
		}
	}

	// Student attendance history
	router.GET("/students/:studentID/attendance", h.GetStudentAttendance)
	// Student grade history
	router.GET("/students/:studentID/grades", h.GetStudentGrades)

	// Direct grade access
	router.PUT("/grades/:gradeID", h.UpdateGrade)
	router.DELETE("/grades/:gradeID", h.DeleteGrade)

	// User Management (Admin only)
	users := router.Group("/users")
	users.Use(middlewares.RequireRole(models.RoleAdmin))
	{
		users.POST("/", h.CreateUser)
		users.GET("/:userID", h.GetUser)
		users.PUT("/:userID", h.UpdateUser)
		users.DELETE("/:userID", h.DeleteUser)
		users.PUT("/:userID/password", h.ChangePassword)
	}

	// Payment Management
	payments := router.Group("/payments")
	{
		payments.POST("/", h.CreatePayment)
		payments.GET("/", h.GetAllPayments)
		payments.GET("/:paymentID", h.GetPayment)
		payments.PUT("/:paymentID", h.UpdatePayment)
		payments.DELETE("/:paymentID", h.DeletePayment)
	}

	// Invoice Management
	invoices := router.Group("/invoices")
	{
		invoices.POST("/", h.CreateInvoice)
		invoices.GET("/", h.GetAllInvoices)
		invoices.GET("/:invoiceID", h.GetInvoice)
		invoices.PUT("/:invoiceID", h.UpdateInvoice)
		invoices.DELETE("/:invoiceID", h.DeleteInvoice)
	}

	// Student-specific payment and invoice routes
	router.GET("/students/:studentID/payments", h.GetStudentPayments)
	router.GET("/students/:studentID/invoices", h.GetStudentInvoices)

	// Notification Management
	notifications := router.Group("/notifications")
	{
		notifications.POST("/send", h.SendNotification)
		notifications.POST("/send/bulk", h.SendBulkNotification)
		notifications.GET("/", h.GetAllNotifications)
		notifications.GET("/:notificationID", h.GetNotification)
		notifications.GET("/recipient/:recipient", h.GetNotificationsByRecipient)
		notifications.POST("/:notificationID/retry", h.RetryNotification)
	}

	// Notification Templates
	templates := router.Group("/notification-templates")
	{
		templates.POST("/", h.CreateTemplate)
		templates.GET("/", h.GetAllTemplates)
		templates.GET("/:templateID", h.GetTemplate)
		templates.PUT("/:templateID", h.UpdateTemplate)
		templates.DELETE("/:templateID", h.DeleteTemplate)
	}

	// Document Management
	documents := router.Group("/documents")
	{
		documents.POST("/upload", h.UploadDocument)
		documents.GET("/", h.GetAllDocuments)
		documents.GET("/:documentID", h.GetDocument)
		documents.GET("/:documentID/download", h.DownloadDocument)
		documents.GET("/:entityType/:entityID", h.GetEntityDocuments)
		documents.PUT("/:documentID", h.UpdateDocument)
		documents.POST("/:documentID/approve", h.ApproveDocument)
		documents.DELETE("/:documentID", h.DeleteDocument)
	}

	// Messaging & Communication
	messages := router.Group("/messages")
	{
		messages.POST("/send", h.SendMessage)
		messages.GET("/inbox", h.GetInbox)
		messages.GET("/sent", h.GetSentMessages)
		messages.GET("/announcements", h.GetAnnouncements)
		messages.GET("/:messageID", h.GetMessage)
		messages.POST("/:messageID/read", h.MarkMessageAsRead)
		messages.DELETE("/:messageID", h.DeleteMessage)
	}

	// Analytics & Reporting
	analytics := router.Group("/analytics")
	{
		analytics.GET("/dashboard", h.GetDashboardMetrics)
		analytics.POST("/reports/financial", h.GetFinancialReport)
		analytics.POST("/reports/attendance", h.GetAttendanceReport)
		analytics.GET("/students/:studentID/progress", h.GetStudentProgress)
	}

	// Calendar & Events
	calendar := router.Group("/calendar")
	{
		calendar.POST("/events", h.CreateEvent)
		calendar.GET("/events", h.GetCalendarEvents)
		calendar.GET("/events/:eventID", h.GetEvent)
		calendar.PUT("/events/:eventID", h.UpdateEvent)
		calendar.DELETE("/events/:eventID", h.DeleteEvent)
	}

	// Applications & Enrollment
	applications := router.Group("/applications")
	{
		applications.POST("/", h.CreateApplication)
		applications.GET("/", h.GetAllApplications)
		applications.GET("/:applicationID", h.GetApplication)
		applications.GET("/course/:courseID", h.GetApplicationsByCourse)
		applications.GET("/status/:status", h.GetApplicationsByStatus)
		applications.PUT("/:applicationID", h.UpdateApplication)
		applications.POST("/:applicationID/review", h.ReviewApplication)
		applications.POST("/:applicationID/enroll", h.EnrollApplication)
		applications.DELETE("/:applicationID", h.DeleteApplication)
	}

	// Exams & Results
	exams := router.Group("/exams")
	{
		exams.POST("/", h.CreateExam)
		exams.GET("/", h.GetAllExams)
		exams.GET("/:examID", h.GetExam)
		exams.GET("/course/:courseID", h.GetExamsByCourse)
		exams.GET("/group/:groupID", h.GetExamsByGroup)
		exams.PUT("/:examID", h.UpdateExam)
		exams.POST("/:examID/results", h.SubmitExamResult)
		exams.GET("/:examID/results", h.GetExamResults)
		exams.GET("/student/:studentID/results", h.GetStudentExamResults)
		exams.GET("/:examID/statistics", h.GetExamStatistics)
		exams.DELETE("/:examID", h.DeleteExam)
	}

	// Portals
	portal := router.Group("/portal")
	{
		portal.GET("/student/:studentID", h.GetStudentPortal)
		portal.GET("/teacher/:teacherID", h.GetTeacherPortal)
	}

	// Audit Logs (Admin only)
	router.GET("/audit-logs", middlewares.RequireRole(models.RoleAdmin), h.GetAuditLogs)

	// Start server
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("listen error", err)
		}
	}()

	logger.Infof("Server started on port %s", cfg.Server.Port)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", err)
	}

	logger.Info("Server exiting")
}
