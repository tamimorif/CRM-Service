package handlers

import (
	"github.com/softclub-go-0-0/crm-service/pkg/services"
)

// Handler holds all the service dependencies for HTTP handlers
type Handler struct {
	teacherService          services.TeacherService
	courseService           services.CourseService
	studentService          services.StudentService
	groupService            services.GroupService
	timetableService        services.TimetableService
	attendanceService       services.AttendanceService
	gradeService            services.GradeService
	healthService           services.HealthService
	userService             services.UserService
	auditService            services.AuditService
	sessionService          services.SessionService
	paymentService          services.PaymentService
	invoiceService          services.InvoiceService
	notificationService     *services.NotificationService
	templateService         *services.TemplateService
	documentService         *services.DocumentService
	messageService          *services.MessageService
	analyticsService        *services.AnalyticsService
	calendarService         *services.CalendarService
	applicationService      *services.ApplicationService
	examService             *services.ExamService
	portalService           *services.PortalService
	parentService           *services.ParentService
	assignmentService       *services.AssignmentService
	waitlistService         *services.WaitlistService
	bulkService             *services.BulkService
	recurringInvoiceService *services.RecurringInvoiceService
	advancedSearchService   *services.AdvancedSearchService
}

// NewHandler creates a new Handler instance
func NewHandler(
	teacherService services.TeacherService,
	courseService services.CourseService,
	studentService services.StudentService,
	groupService services.GroupService,
	timetableService services.TimetableService,
	attendanceService services.AttendanceService,
	gradeService services.GradeService,
	healthService services.HealthService,
	userService services.UserService,
	auditService services.AuditService,
	sessionService services.SessionService,
	paymentService services.PaymentService,
	invoiceService services.InvoiceService,
	notificationService *services.NotificationService,
	templateService *services.TemplateService,
	documentService *services.DocumentService,
	messageService *services.MessageService,
	analyticsService *services.AnalyticsService,
	calendarService *services.CalendarService,
	applicationService *services.ApplicationService,
	examService *services.ExamService,
	portalService *services.PortalService,
	parentService *services.ParentService,
	assignmentService *services.AssignmentService,
	waitlistService *services.WaitlistService,
	bulkService *services.BulkService,
	recurringInvoiceService *services.RecurringInvoiceService,
	advancedSearchService *services.AdvancedSearchService,
) *Handler {
	return &Handler{
		teacherService:          teacherService,
		courseService:           courseService,
		studentService:          studentService,
		groupService:            groupService,
		timetableService:        timetableService,
		attendanceService:       attendanceService,
		gradeService:            gradeService,
		healthService:           healthService,
		userService:             userService,
		auditService:            auditService,
		sessionService:          sessionService,
		paymentService:          paymentService,
		invoiceService:          invoiceService,
		notificationService:     notificationService,
		templateService:         templateService,
		documentService:         documentService,
		messageService:          messageService,
		analyticsService:        analyticsService,
		calendarService:         calendarService,
		applicationService:      applicationService,
		examService:             examService,
		portalService:           portalService,
		parentService:           parentService,
		assignmentService:       assignmentService,
		waitlistService:         waitlistService,
		bulkService:             bulkService,
		recurringInvoiceService: recurringInvoiceService,
		advancedSearchService:   advancedSearchService,
	}
}
