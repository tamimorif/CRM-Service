package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
	"gorm.io/gorm"
)

// AnalyticsService handles analytics and reporting
type AnalyticsService struct {
	db *gorm.DB
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(db *gorm.DB) *AnalyticsService {
	return &AnalyticsService{db: db}
}

// GetDashboardMetrics retrieves overall dashboard metrics
func (s *AnalyticsService) GetDashboardMetrics(ctx context.Context) (*dto.DashboardMetrics, error) {
	metrics := &dto.DashboardMetrics{}

	// Count students
	s.db.Model(&models.Student{}).Count(&metrics.TotalStudents)
	s.db.Model(&models.Student{}).Where("status = ?", "active").Count(&metrics.ActiveStudents)

	// Count teachers and courses
	s.db.Model(&models.Teacher{}).Where("status = ?", "active").Count(&metrics.TotalTeachers)
	s.db.Model(&models.Course{}).Count(&metrics.TotalCourses)
	s.db.Model(&models.Group{}).Count(&metrics.ActiveGroups)

	// Financial metrics
	var totalRevenue, pendingPayments float64
	s.db.Model(&models.Payment{}).Where("status = ?", models.PaymentCompleted).Select("COALESCE(SUM(amount), 0)").Scan(&totalRevenue)
	s.db.Model(&models.Invoice{}).Where("status = ?", "pending").Select("COALESCE(SUM(balance), 0)").Scan(&pendingPayments)
	metrics.TotalRevenue = totalRevenue
	metrics.PendingPayments = pendingPayments

	// This month metrics
	startOfMonth := time.Now().AddDate(0, 0, -time.Now().Day()+1)
	var monthRevenue float64
	s.db.Model(&models.Payment{}).
		Where("status = ? AND created_at >= ?", models.PaymentCompleted, startOfMonth).
		Select("COALESCE(SUM(amount), 0)").Scan(&monthRevenue)
	metrics.ThisMonthRevenue = monthRevenue

	s.db.Model(&models.Student{}).Where("created_at >= ?", startOfMonth).Count(&metrics.ThisMonthEnrolment)

	// Attendance rate
	var totalAttendance, presentCount int64
	s.db.Model(&models.Attendance{}).Count(&totalAttendance)
	s.db.Model(&models.Attendance{}).Where("status = ?", "present").Count(&presentCount)
	if totalAttendance > 0 {
		metrics.AttendanceRate = float64(presentCount) / float64(totalAttendance) * 100
	}

	// Unread messages
	s.db.Model(&models.Message{}).Where("status != ?", models.MessageStatusRead).Count(&metrics.UnreadMessages)

	// Pending documents
	s.db.Model(&models.Document{}).Where("status = ?", models.DocumentStatusPending).Count(&metrics.PendingDocuments)

	return metrics, nil
}

// GetFinancialReport generates financial report
func (s *AnalyticsService) GetFinancialReport(ctx context.Context, req dto.ReportRequest) (*dto.FinancialReport, error) {
	report := &dto.FinancialReport{
		Period: req.Period,
	}

	query := s.db.Model(&models.Invoice{})
	if req.StartDate != nil {
		query = query.Where("created_at >= ?", req.StartDate)
	}
	if req.EndDate != nil {
		query = query.Where("created_at <= ?", req.EndDate)
	}

	// Total invoices
	query.Count(&report.TotalInvoices)

	// Revenue by status
	query.Where("status = ?", models.InvoicePaid).Select("COALESCE(SUM(total_amount), 0)").Scan(&report.TotalPaid)
	query.Where("status = ?", "pending").Select("COALESCE(SUM(balance), 0)").Scan(&report.TotalPending)
	query.Where("status = ?", models.InvoiceOverdue).Select("COALESCE(SUM(balance), 0)").Scan(&report.TotalOverdue)

	report.TotalRevenue = report.TotalPaid + report.TotalPending + report.TotalOverdue

	// Count by status
	query.Where("status = ?", models.InvoicePaid).Count(&report.PaidInvoices)
	query.Where("status = ?", "pending").Count(&report.PendingInvoices)
	query.Where("status = ?", models.InvoiceOverdue).Count(&report.OverdueInvoices)

	// Top courses by revenue
	type courseResult struct {
		CourseID uuid.UUID
		Revenue  float64
		Students int64
	}
	var topCourses []courseResult
	s.db.Model(&models.Invoice{}).
		Select("course_id, SUM(total_amount) as revenue, COUNT(DISTINCT student_id) as students").
		Where("course_id IS NOT NULL").
		Group("course_id").
		Order("revenue DESC").
		Limit(5).
		Scan(&topCourses)

	report.TopCourses = make([]dto.CourseRevenue, 0)
	for _, tc := range topCourses {
		var course models.Course
		if err := s.db.First(&course, "id = ?", tc.CourseID).Error; err == nil {
			report.TopCourses = append(report.TopCourses, dto.CourseRevenue{
				CourseID:   tc.CourseID,
				CourseName: course.Title,
				Revenue:    tc.Revenue,
				Students:   tc.Students,
			})
		}
	}

	return report, nil
}

// GetStudentProgress gets progress report for a student
func (s *AnalyticsService) GetStudentProgress(ctx context.Context, studentID uuid.UUID) (*dto.StudentProgressReport, error) {
	var student models.Student
	if err := s.db.Preload("Group").Preload("Group.Course").First(&student, "id = ?", studentID).Error; err != nil {
		return nil, err
	}

	report := &dto.StudentProgressReport{
		StudentID:   student.ID,
		StudentName: student.Name + " " + "",
	}

	if student.Group != nil {
		report.GroupName = student.Group.Name
		if student.Group.Course != nil {
			report.CourseName = student.Group.Course.Title
		}
	}

	// Attendance statistics
	var totalClasses, attendedClasses int64
	s.db.Model(&models.Attendance{}).Where("student_id = ?", studentID).Count(&totalClasses)
	s.db.Model(&models.Attendance{}).Where("student_id = ? AND status = ?", studentID, "present").Count(&attendedClasses)

	report.TotalClasses = int(totalClasses)
	report.AttendedClasses = int(attendedClasses)
	if totalClasses > 0 {
		report.AttendanceRate = float64(attendedClasses) / float64(totalClasses) * 100
	}

	// Grade average
	var avgGrade float64
	s.db.Model(&models.Grade{}).Where("student_id = ?", studentID).Select("COALESCE(AVG(grade), 0)").Scan(&avgGrade)
	report.AverageGrade = avgGrade

	// Latest grades
	var grades []models.Grade
	s.db.Where("student_id = ?", studentID).Order("created_at DESC").Limit(5).Find(&grades)
	report.LatestGrades = make([]dto.GradeInfo, 0)
	for _, g := range grades {
		report.LatestGrades = append(report.LatestGrades, dto.GradeInfo{
			ID:        g.ID,
			Score:     float64(g.Value),
			Type:      string(g.Type),
			Notes:     g.Notes,
			CreatedAt: g.CreatedAt,
		})
	}

	// Outstanding fees
	var outstandingFees float64
	s.db.Model(&models.Invoice{}).
		Where("student_id = ? AND status != ?", studentID, models.InvoicePaid).
		Select("COALESCE(SUM(balance), 0)").Scan(&outstandingFees)
	report.OutstandingFees = outstandingFees

	if outstandingFees > 0 {
		report.PaymentStatus = "outstanding"
	} else {
		report.PaymentStatus = "paid"
	}

	report.LastActive = student.UpdatedAt

	return report, nil
}

// GetAttendanceReport generates attendance report
func (s *AnalyticsService) GetAttendanceReport(ctx context.Context, req dto.ReportRequest) (*dto.AttendanceReport, error) {
	report := &dto.AttendanceReport{
		Period: req.Period,
	}

	query := s.db.Model(&models.Attendance{})
	if req.StartDate != nil {
		query = query.Where("created_at >= ?", req.StartDate)
	}
	if req.EndDate != nil {
		query = query.Where("created_at <= ?", req.EndDate)
	}
	if req.GroupID != nil {
		query = query.Where("group_id = ?", req.GroupID)
	}

	query.Count(&report.TotalSessions)
	query.Where("status = ?", "present").Count(&report.TotalPresent)
	query.Where("status = ?", "absent").Count(&report.TotalAbsent)

	if report.TotalSessions > 0 {
		report.AverageAttendance = float64(report.TotalPresent) / float64(report.TotalSessions) * 100
	}

	// Attendance by group
	type groupResult struct {
		GroupID         uuid.UUID
		TotalSessions   int64
		PresentSessions int64
	}
	var groupResults []groupResult
	s.db.Model(&models.Attendance{}).
		Select("group_id, COUNT(*) as total_sessions, SUM(CASE WHEN status = ? THEN 1 ELSE 0 END) as present_sessions", "present").
		Group("group_id").
		Scan(&groupResults)

	report.AttendanceByGroup = make([]dto.GroupAttendance, 0)
	for _, gr := range groupResults {
		var group models.Group
		if err := s.db.First(&group, "id = ?", gr.GroupID).Error; err == nil {
			rate := float64(0)
			if gr.TotalSessions > 0 {
				rate = float64(gr.PresentSessions) / float64(gr.TotalSessions) * 100
			}
			report.AttendanceByGroup = append(report.AttendanceByGroup, dto.GroupAttendance{
				GroupID:        gr.GroupID,
				GroupName:      group.Name,
				AttendanceRate: rate,
				TotalStudents:  len(group.Students),
			})
		}
	}

	return report, nil
}
