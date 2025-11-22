package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateTeacherRequest represents a request to create a teacher
type CreateTeacherRequest struct {
	Name    string `json:"name" binding:"required,min=2,max=100"`
	Surname string `json:"surname" binding:"required,min=2,max=100"`
	Phone   string `json:"phone" binding:"required,len=12"`
	Email   string `json:"email" binding:"omitempty,email"`
}

// UpdateTeacherRequest represents a request to update a teacher
type UpdateTeacherRequest struct {
	Name    string `json:"name" binding:"required,min=2,max=100"`
	Surname string `json:"surname" binding:"required,min=2,max=100"`
	Phone   string `json:"phone" binding:"required,len=12"`
	Email   string `json:"email" binding:"omitempty,email"`
}

// TeacherResponse represents a teacher response
type TeacherResponse struct {
	ID        uuid.UUID     `json:"id"`
	Name      string        `json:"name"`
	Surname   string        `json:"surname"`
	Phone     string        `json:"phone"`
	Email     string        `json:"email"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	Groups    []GroupSimple `json:"groups,omitempty"`
}

// TeacherSimple represents a simplified teacher (for nested responses)
type TeacherSimple struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Surname string    `json:"surname"`
	Email   string    `json:"email"`
}

// CreateCourseRequest represents a request to create a course
type CreateCourseRequest struct {
	Title      string  `json:"title" binding:"required,min=2,max=200"`
	MonthlyFee float64 `json:"monthly_fee" binding:"required,min=0"`
	Duration   int     `json:"duration" binding:"required,min=1,max=60"`
}

// UpdateCourseRequest represents a request to update a course
type UpdateCourseRequest struct {
	Title      string  `json:"title" binding:"required,min=2,max=200"`
	MonthlyFee float64 `json:"monthly_fee" binding:"required,min=0"`
	Duration   int     `json:"duration" binding:"required,min=1,max=60"`
}

// CourseResponse represents a course response
type CourseResponse struct {
	ID         uuid.UUID     `json:"id"`
	Title      string        `json:"title"`
	MonthlyFee float64       `json:"monthly_fee"`
	Duration   int           `json:"duration"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
	Groups     []GroupSimple `json:"groups,omitempty"`
}

// CourseSimple represents a simplified course
type CourseSimple struct {
	ID         uuid.UUID `json:"id"`
	Title      string    `json:"title"`
	MonthlyFee float64   `json:"monthly_fee"`
	Duration   int       `json:"duration"`
}

// CreateStudentRequest represents a request to create a student
type CreateStudentRequest struct {
	Name    string `json:"name" binding:"required,min=2,max=100"`
	Surname string `json:"surname" binding:"required,min=2,max=100"`
	Phone   string `json:"phone" binding:"required,len=12"`
	Email   string `json:"email" binding:"omitempty,email"`
}

// UpdateStudentRequest represents a request to update a student
type UpdateStudentRequest struct {
	Name    string `json:"name" binding:"required,min=2,max=100"`
	Surname string `json:"surname" binding:"required,min=2,max=100"`
	Phone   string `json:"phone" binding:"required,len=12"`
	Email   string `json:"email" binding:"omitempty,email"`
}

// StudentResponse represents a student response
type StudentResponse struct {
	ID        uuid.UUID   `json:"id"`
	Name      string      `json:"name"`
	Surname   string      `json:"surname"`
	Phone     string      `json:"phone"`
	Email     string      `json:"email"`
	GroupID   uuid.UUID   `json:"group_id"`
	Group     GroupSimple `json:"group,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// StudentSimple represents a simplified student
type StudentSimple struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Surname string    `json:"surname"`
	Phone   string    `json:"phone"`
}

// CreateGroupRequest represents a request to create a group
type CreateGroupRequest struct {
	Name        string    `json:"name" binding:"required,min=2,max=100"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	CourseID    uuid.UUID `json:"course_id" binding:"required"`
	TeacherID   uuid.UUID `json:"teacher_id" binding:"required"`
	TimetableID uuid.UUID `json:"timetable_id" binding:"required"`
	Capacity    int       `json:"capacity" binding:"required,min=1,max=100"`
}

// UpdateGroupRequest represents a request to update a group
type UpdateGroupRequest struct {
	Name        string    `json:"name" binding:"required,min=2,max=100"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	CourseID    uuid.UUID `json:"course_id" binding:"required"`
	TeacherID   uuid.UUID `json:"teacher_id" binding:"required"`
	TimetableID uuid.UUID `json:"timetable_id" binding:"required"`
	Capacity    int       `json:"capacity" binding:"required,min=1,max=100"`
}

// GroupResponse represents a group response
type GroupResponse struct {
	ID           uuid.UUID       `json:"id"`
	Name         string          `json:"name"`
	StartDate    time.Time       `json:"start_date"`
	CourseID     uuid.UUID       `json:"course_id"`
	TeacherID    uuid.UUID       `json:"teacher_id"`
	TimetableID  uuid.UUID       `json:"timetable_id"`
	Capacity     int             `json:"capacity"`
	StudentCount int             `json:"student_count"`
	Course       CourseSimple    `json:"course,omitempty"`
	Teacher      TeacherSimple   `json:"teacher,omitempty"`
	Timetable    TimetableSimple `json:"timetable,omitempty"`
	Students     []StudentSimple `json:"students,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// GroupSimple represents a simplified group
type GroupSimple struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	Capacity  int       `json:"capacity"`
}

// CreateTimetableRequest represents a request to create a timetable
type CreateTimetableRequest struct {
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
	Days      string `json:"days" binding:"required"`
	Classroom string `json:"classroom" binding:"required,min=1,max=50"`
}

// UpdateTimetableRequest represents a request to update a timetable
type UpdateTimetableRequest struct {
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
	Days      string `json:"days" binding:"required"`
	Classroom string `json:"classroom" binding:"required,min=1,max=50"`
}

// TimetableResponse represents a timetable response
type TimetableResponse struct {
	ID        uuid.UUID     `json:"id"`
	StartTime string        `json:"start_time"`
	EndTime   string        `json:"end_time"`
	Days      string        `json:"days"`
	Classroom string        `json:"classroom"`
	Groups    []GroupSimple `json:"groups,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// TimetableSimple represents a simplified timetable
type TimetableSimple struct {
	ID        uuid.UUID `json:"id"`
	StartTime string    `json:"start_time"`
	EndTime   string    `json:"end_time"`
	Days      string    `json:"days"`
	Classroom string    `json:"classroom"`
}
