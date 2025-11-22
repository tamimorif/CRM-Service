package dto

import (
	"time"

	"github.com/google/uuid"
)

// StudentPortalDashboard represents student portal dashboard data
type StudentPortalDashboard struct {
	Student         StudentSimple     `json:"student"`
	Group           *GroupSimple      `json:"group,omitempty"`
	Course          *CourseSimple     `json:"course,omitempty"`
	UpcomingClasses []TimetableSimple `json:"upcoming_classes"`
	RecentGrades    []GradeInfo       `json:"recent_grades"`
	UpcomingExams   []ExamSimple      `json:"upcoming_exams"`
	AttendanceRate  float64           `json:"attendance_rate"`
	UnreadMessages  int64             `json:"unread_messages"`
	PendingPayments float64           `json:"pending_payments"`
	Announcements   []MessageSimple   `json:"announcements"`
}

// TeacherPortalDashboard represents teacher portal dashboard data
type TeacherPortalDashboard struct {
	Teacher          TeacherSimple      `json:"teacher"`
	TotalStudents    int64              `json:"total_students"`
	TotalGroups      int64              `json:"total_groups"`
	TodayClasses     []TimetableSimple  `json:"today_classes"`
	UpcomingExams    []ExamSimple       `json:"upcoming_exams"`
	PendingGrading   int64              `json:"pending_grading"`
	UnreadMessages   int64              `json:"unread_messages"`
	RecentAttendance []AttendanceSimple `json:"recent_attendance"`
	Announcements    []MessageSimple    `json:"announcements"`
}

// ExamSimple represents simplified exam info
type ExamSimple struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Type      string    `json:"type"`
	StartTime time.Time `json:"start_time"`
	Duration  int       `json:"duration"`
	Location  string    `json:"location,omitempty"`
}

// MessageSimple represents simplified message info
type MessageSimple struct {
	ID        uuid.UUID `json:"id"`
	Subject   string    `json:"subject,omitempty"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	Priority  int       `json:"priority"`
}

// AttendanceSimple represents simplified attendance info
type AttendanceSimple struct {
	ID        uuid.UUID `json:"id"`
	Date      time.Time `json:"date"`
	Status    string    `json:"status"`
	StudentID uuid.UUID `json:"student_id,omitempty"`
	GroupID   uuid.UUID `json:"group_id,omitempty"`
}

// StudentSchedule represents student's weekly schedule
type StudentSchedule struct {
	StudentID uuid.UUID         `json:"student_id"`
	GroupID   uuid.UUID         `json:"group_id"`
	Schedule  []TimetableSimple `json:"schedule"`
}

// TeacherSchedule represents teacher's weekly schedule
type TeacherSchedule struct {
	TeacherID uuid.UUID         `json:"teacher_id"`
	Schedule  []TimetableSimple `json:"schedule"`
}
