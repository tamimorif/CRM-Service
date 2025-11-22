package dto

import (
	"time"

	"github.com/google/uuid"
)

// DashboardMetrics represents overall dashboard statistics
type DashboardMetrics struct {
	TotalStudents      int64   `json:"total_students"`
	ActiveStudents     int64   `json:"active_students"`
	TotalTeachers      int64   `json:"total_teachers"`
	TotalCourses       int64   `json:"total_courses"`
	ActiveGroups       int64   `json:"active_groups"`
	TotalRevenue       float64 `json:"total_revenue"`
	PendingPayments    float64 `json:"pending_payments"`
	AttendanceRate     float64 `json:"attendance_rate"`
	UnreadMessages     int64   `json:"unread_messages"`
	PendingDocuments   int64   `json:"pending_documents"`
	ThisMonthRevenue   float64 `json:"this_month_revenue"`
	ThisMonthEnrolment int64   `json:"this_month_enrolment"`
}

// StudentProgressReport represents a student's progress
type StudentProgressReport struct {
	StudentID       uuid.UUID   `json:"student_id"`
	StudentName     string      `json:"student_name"`
	CourseName      string      `json:"course_name"`
	GroupName       string      `json:"group_name"`
	AttendanceRate  float64     `json:"attendance_rate"`
	AverageGrade    float64     `json:"average_grade"`
	TotalClasses    int         `json:"total_classes"`
	AttendedClasses int         `json:"attended_classes"`
	LatestGrades    []GradeInfo `json:"latest_grades"`
	PaymentStatus   string      `json:"payment_status"`
	OutstandingFees float64     `json:"outstanding_fees"`
	LastActive      time.Time   `json:"last_active"`
}

// GradeInfo represents grade information
type GradeInfo struct {
	ID        uuid.UUID `json:"id"`
	Score     float64   `json:"score"`
	Type      string    `json:"type"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
}

// FinancialReport represents financial summary
type FinancialReport struct {
	Period          string          `json:"period"` // daily, weekly, monthly, yearly
	TotalRevenue    float64         `json:"total_revenue"`
	TotalPaid       float64         `json:"total_paid"`
	TotalPending    float64         `json:"total_pending"`
	TotalOverdue    float64         `json:"total_overdue"`
	TotalInvoices   int64           `json:"total_invoices"`
	PaidInvoices    int64           `json:"paid_invoices"`
	PendingInvoices int64           `json:"pending_invoices"`
	OverdueInvoices int64           `json:"overdue_invoices"`
	TopCourses      []CourseRevenue `json:"top_courses"`
}

// CourseRevenue represents revenue by course
type CourseRevenue struct {
	CourseID   uuid.UUID `json:"course_id"`
	CourseName string    `json:"course_name"`
	Revenue    float64   `json:"revenue"`
	Students   int64     `json:"students"`
}

// AttendanceReport represents attendance analytics
type AttendanceReport struct {
	Period            string            `json:"period"`
	AverageAttendance float64           `json:"average_attendance"`
	TotalSessions     int64             `json:"total_sessions"`
	TotalPresent      int64             `json:"total_present"`
	TotalAbsent       int64             `json:"total_absent"`
	AttendanceByGroup []GroupAttendance `json:"attendance_by_group"`
}

// GroupAttendance represents attendance for a group
type GroupAttendance struct {
	GroupID        uuid.UUID `json:"group_id"`
	GroupName      string    `json:"group_name"`
	AttendanceRate float64   `json:"attendance_rate"`
	TotalStudents  int       `json:"total_students"`
}

// GradeDistribution represents grade analytics
type GradeDistribution struct {
	CourseID     uuid.UUID        `json:"course_id"`
	CourseName   string           `json:"course_name"`
	AverageGrade float64          `json:"average_grade"`
	GradeRanges  map[string]int64 `json:"grade_ranges"` // A, B, C, D, F counts
	TotalGrades  int64            `json:"total_grades"`
}

// ReportRequest represents a report request with filters
type ReportRequest struct {
	ReportType string     `json:"report_type" binding:"required"` // financial, attendance, student_progress, etc
	StartDate  *time.Time `json:"start_date,omitempty"`
	EndDate    *time.Time `json:"end_date,omitempty"`
	Period     string     `json:"period,omitempty"` // daily, weekly, monthly, yearly
	StudentID  *uuid.UUID `json:"student_id,omitempty"`
	CourseID   *uuid.UUID `json:"course_id,omitempty"`
	GroupID    *uuid.UUID `json:"group_id,omitempty"`
}
