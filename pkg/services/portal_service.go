package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
	"gorm.io/gorm"
)

// PortalService handles portal operations
type PortalService struct {
	db *gorm.DB
}

// NewPortalService creates a new portal service
func NewPortalService(db *gorm.DB) *PortalService {
	return &PortalService{db: db}
}

// GetStudentDashboard retrieves student portal dashboard data
func (s *PortalService) GetStudentDashboard(ctx context.Context, studentID uuid.UUID) (*dto.StudentPortalDashboard, error) {
	var student models.Student
	if err := s.db.Preload("Group").Preload("Group.Course").First(&student, "id = ?", studentID).Error; err != nil {
		return nil, err
	}

	dashboard := &dto.StudentPortalDashboard{
		Student: dto.StudentSimple{
			ID:   student.ID,
			Name: student.Name,
		},
	}

	// Group and course info
	if student.Group != nil {
		dashboard.Group = &dto.GroupSimple{
			ID:   student.Group.ID,
			Name: student.Group.Name,
		}
		if student.Group.Course != nil {
			dashboard.Course = &dto.CourseSimple{
				ID:    student.Group.Course.ID,
				Title: student.Group.Course.Title,
			}
		}
	}

	// Upcoming classes (next 7 days)
	if student.GroupID != uuid.Nil {
		var timetables []models.Timetable
		s.db.Where("group_id = ?", student.GroupID).Limit(5).Find(&timetables)
		dashboard.UpcomingClasses = make([]dto.TimetableSimple, len(timetables))
		for i, tt := range timetables {
			dashboard.UpcomingClasses[i] = dto.TimetableSimple{
				ID:        tt.ID,
				StartTime: tt.StartTime,
				EndTime:   tt.EndTime,
				Days:      tt.Days,
				Classroom: tt.Classroom,
			}
		}
	}

	// Recent grades
	var grades []models.Grade
	s.db.Where("student_id = ?", studentID).Order("created_at DESC").Limit(5).Find(&grades)
	dashboard.RecentGrades = make([]dto.GradeInfo, len(grades))
	for i, g := range grades {
		dashboard.RecentGrades[i] = dto.GradeInfo{
			ID:        g.ID,
			Score:     float64(g.Value),
			Type:      string(g.Type),
			Notes:     g.Notes,
			CreatedAt: g.CreatedAt,
		}
	}

	// Upcoming exams
	if student.GroupID != uuid.Nil {
		var exams []models.Exam
		s.db.Where("group_id = ? AND start_time > ?", student.GroupID, time.Now()).
			Order("start_time ASC").Limit(5).Find(&exams)
		dashboard.UpcomingExams = make([]dto.ExamSimple, len(exams))
		for i, e := range exams {
			dashboard.UpcomingExams[i] = dto.ExamSimple{
				ID:        e.ID,
				Title:     e.Title,
				Type:      string(e.Type),
				StartTime: e.StartTime,
				Duration:  e.Duration,
				Location:  e.Location,
			}
		}
	}

	// Attendance rate
	var totalAttendance, presentCount int64
	s.db.Model(&models.Attendance{}).Where("student_id = ?", studentID).Count(&totalAttendance)
	s.db.Model(&models.Attendance{}).Where("student_id = ? AND status = ?", studentID, "present").Count(&presentCount)
	if totalAttendance > 0 {
		dashboard.AttendanceRate = float64(presentCount) / float64(totalAttendance) * 100
	}

	// Unread messages
	s.db.Model(&models.Message{}).
		Where("recipient_id = ? AND status != ?", studentID, models.MessageStatusRead).
		Count(&dashboard.UnreadMessages)

	// Pending payments
	s.db.Model(&models.Invoice{}).
		Where("student_id = ? AND status != ?", studentID, "paid").
		Select("COALESCE(SUM(balance), 0)").Scan(&dashboard.PendingPayments)

	// Announcements
	var announcements []models.Message
	s.db.Where("type = ?", models.MessageTypeAnnouncement).
		Order("created_at DESC").Limit(5).Find(&announcements)
	dashboard.Announcements = make([]dto.MessageSimple, len(announcements))
	for i, a := range announcements {
		dashboard.Announcements[i] = dto.MessageSimple{
			ID:        a.ID,
			Subject:   a.Subject,
			Body:      a.Body,
			CreatedAt: a.CreatedAt,
			Priority:  a.Priority,
		}
	}

	return dashboard, nil
}

// GetTeacherDashboard retrieves teacher portal dashboard data
func (s *PortalService) GetTeacherDashboard(ctx context.Context, teacherID uuid.UUID) (*dto.TeacherPortalDashboard, error) {
	var teacher models.Teacher
	if err := s.db.First(&teacher, "id = ?", teacherID).Error; err != nil {
		return nil, err
	}

	dashboard := &dto.TeacherPortalDashboard{
		Teacher: dto.TeacherSimple{
			ID:   teacher.ID,
			Name: teacher.Name,
		},
	}

	// Total students (across all groups taught by teacher)
	var groupIDs []uuid.UUID
	s.db.Model(&models.Group{}).Where("teacher_id = ?", teacherID).Pluck("id", &groupIDs)
	if len(groupIDs) > 0 {
		s.db.Model(&models.Student{}).Where("group_id IN ?", groupIDs).Count(&dashboard.TotalStudents)
		dashboard.TotalGroups = int64(len(groupIDs))
	}

	// Today's classes
	today := time.Now().Weekday().String()
	var timetables []models.Timetable
	if len(groupIDs) > 0 {
		s.db.Where("group_id IN ? AND day_of_week = ?", groupIDs, today).Find(&timetables)
		dashboard.TodayClasses = make([]dto.TimetableSimple, len(timetables))
		for i, tt := range timetables {
			dashboard.TodayClasses[i] = dto.TimetableSimple{
				ID:        tt.ID,
				StartTime: tt.StartTime,
				EndTime:   tt.EndTime,
				Days:      tt.Days,
				Classroom: tt.Classroom,
			}
		}
	}

	// Upcoming exams
	if len(groupIDs) > 0 {
		var exams []models.Exam
		s.db.Where("group_id IN ? AND start_time > ?", groupIDs, time.Now()).
			Order("start_time ASC").Limit(5).Find(&exams)
		dashboard.UpcomingExams = make([]dto.ExamSimple, len(exams))
		for i, e := range exams {
			dashboard.UpcomingExams[i] = dto.ExamSimple{
				ID:        e.ID,
				Title:     e.Title,
				Type:      string(e.Type),
				StartTime: e.StartTime,
				Duration:  e.Duration,
				Location:  e.Location,
			}
		}
	}

	// Pending grading (exams without all results)
	// This is a simplified count - in production you'd want more sophisticated logic
	s.db.Model(&models.Exam{}).
		Where("group_id IN ? AND status = ?", groupIDs, models.ExamStatusCompleted).
		Count(&dashboard.PendingGrading)

	// Unread messages
	s.db.Model(&models.Message{}).
		Where("recipient_id = ? AND status != ?", teacherID, models.MessageStatusRead).
		Count(&dashboard.UnreadMessages)

	// Recent attendance (last 5 sessions)
	if len(groupIDs) > 0 {
		var attendances []models.Attendance
		s.db.Where("group_id IN ?", groupIDs).
			Order("created_at DESC").Limit(5).Find(&attendances)
		dashboard.RecentAttendance = make([]dto.AttendanceSimple, len(attendances))
		for i, a := range attendances {
			dashboard.RecentAttendance[i] = dto.AttendanceSimple{
				ID:        a.ID,
				Date:      a.Date,
				Status:    string(a.Status),
				StudentID: a.StudentID,
				GroupID:   a.GroupID,
			}
		}
	}

	// Announcements
	var announcements []models.Message
	s.db.Where("type = ?", models.MessageTypeAnnouncement).
		Order("created_at DESC").Limit(5).Find(&announcements)
	dashboard.Announcements = make([]dto.MessageSimple, len(announcements))
	for i, a := range announcements {
		dashboard.Announcements[i] = dto.MessageSimple{
			ID:        a.ID,
			Subject:   a.Subject,
			Body:      a.Body,
			CreatedAt: a.CreatedAt,
			Priority:  a.Priority,
		}
	}

	return dashboard, nil
}
