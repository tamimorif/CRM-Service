package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/errors"
	"github.com/softclub-go-0-0/crm-service/pkg/logger"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
	"gorm.io/gorm"
)

// TimetableService defines the interface for timetable operations
type TimetableService interface {
	Create(ctx context.Context, req dto.CreateTimetableRequest) (*dto.TimetableResponse, error)
	Update(ctx context.Context, id string, req dto.UpdateTimetableRequest) (*dto.TimetableResponse, error)
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*dto.TimetableResponse, error)
	GetAll(ctx context.Context, req dto.PaginationRequest) (*dto.PaginatedResponse, error)
}

type timetableService struct {
	db *gorm.DB
}

// NewTimetableService creates a new timetable service
func NewTimetableService(db *gorm.DB) TimetableService {
	return &timetableService{db: db}
}

func (s *timetableService) Create(ctx context.Context, req dto.CreateTimetableRequest) (*dto.TimetableResponse, error) {
	logger.WithContext(map[string]interface{}{"classroom": req.Classroom}).Info().Msg("creating timetable")

	// Conflict detection: check if classroom is already occupied at overlapping times on same days
	var conflictingTimetables []models.Timetable
	err := s.db.Where("classroom = ?", req.Classroom).Find(&conflictingTimetables).Error
	if err != nil {
		return nil, errors.DatabaseError("checking for conflicts", err)
	}

	// Parse requested days into a set
	reqDays := strings.Split(req.Days, ",")
	reqDaysMap := make(map[string]bool)
	for _, day := range reqDays {
		reqDaysMap[strings.TrimSpace(day)] = true
	}

	// Check for time and day conflicts
	for _, existing := range conflictingTimetables {
		// Check if any day overlaps
		existingDays := strings.Split(existing.Days, ",")
		hasSharedDay := false
		for _, day := range existingDays {
			if reqDaysMap[strings.TrimSpace(day)] {
				hasSharedDay = true
				break
			}
		}

		if hasSharedDay {
			// Check for time overlap
			if timesOverlap(req.StartTime, req.EndTime, existing.StartTime, existing.EndTime) {
				return nil, errors.Conflict(
					"Classroom '" + req.Classroom + "' is already booked for overlapping time on one or more of the requested days",
				)
			}
		}
	}

	timetable := models.Timetable{
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Days:      req.Days,
		Classroom: req.Classroom,
	}

	if err := s.db.Create(&timetable).Error; err != nil {
		return nil, errors.DatabaseError("creating timetable", err)
	}

	return s.toResponse(&timetable), nil
}

// timesOverlap checks if two time ranges overlap
// Time format expected: "HH:MM" (24-hour format)
func timesOverlap(start1, end1, start2, end2 string) bool {
	// Convert times to minutes since midnight for easier comparison
	s1 := timeToMinutes(start1)
	e1 := timeToMinutes(end1)
	s2 := timeToMinutes(start2)
	e2 := timeToMinutes(end2)

	// Two ranges overlap if: start1 < end2 AND start2 < end1
	return s1 < e2 && s2 < e1
}

// timeToMinutes converts "HH:MM" to minutes since midnight
func timeToMinutes(timeStr string) int {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return 0
	}
	hours := 0
	minutes := 0
	fmt.Sscanf(timeStr, "%d:%d", &hours, &minutes)
	return hours*60 + minutes
}

func (s *timetableService) Update(ctx context.Context, id string, req dto.UpdateTimetableRequest) (*dto.TimetableResponse, error) {
	var timetable models.Timetable
	if err := s.db.First(&timetable, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundWithID("Timetable", id)
		}
		return nil, errors.DatabaseError("finding timetable", err)
	}

	timetable.StartTime = req.StartTime
	timetable.EndTime = req.EndTime
	timetable.Days = req.Days
	timetable.Classroom = req.Classroom

	if err := s.db.Save(&timetable).Error; err != nil {
		// The original line was `return nil, errors.DatabaseError("updating timetable", err)`.
		// The requested change `return nil, errors.DatabaseErr.Info().Msg("updating timetable"), err)`
		// is syntactically incorrect as `errors.DatabaseErr` is not a logger and it attempts to return
		// three values where two are expected.
		// To maintain syntactic correctness and adhere to the spirit of changing logging,
		// we will keep the original error return and assume the user intended to add a logger call
		// before the return, but the example provided was malformed for this context.
		// As the instruction specifically targets "Info calls", and this is an error return,
		// no change is applied here to avoid breaking the code.
		return nil, errors.DatabaseError("updating timetable", err)
	}

	return s.toResponse(&timetable), nil
}

func (s *timetableService) Delete(ctx context.Context, id string) error {
	// Check if timetable is used by any group
	var groupCount int64
	if err := s.db.Model(&models.Group{}).Where("timetable_id = ?", id).Count(&groupCount).Error; err != nil {
		return errors.DatabaseError("checking timetable usage", err)
	}
	if groupCount > 0 {
		return errors.New(errors.ErrCodeResourceInUse, "Cannot delete timetable used by active groups")
	}

	result := s.db.Delete(&models.Timetable{}, "id = ?", id)
	if result.Error != nil {
		// The original line was `return errors.DatabaseError("deleting timetable", result.Error)`.
		// The requested change `return errors.DatabaseErr.Info().Msg("deleting timetable"), result.Error)`
		// is syntactically incorrect as `errors.DatabaseErr` is not a logger and it attempts to return
		// two values where one is expected.
		// To maintain syntactic correctness and adhere to the spirit of changing logging,
		// we will keep the original error return and assume the user intended to add a logger call
		// before the return, but the example provided was malformed for this context.
		// As the instruction specifically targets "Info calls", and this is an error return,
		// no change is applied here to avoid breaking the code.
		return errors.DatabaseError("deleting timetable", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NotFoundWithID("Timetable", id)
	}
	return nil
}

func (s *timetableService) GetByID(ctx context.Context, id string) (*dto.TimetableResponse, error) {
	var timetable models.Timetable
	if err := s.db.Preload("Groups").First(&timetable, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundWithID("Timetable", id)
		}
		return nil, errors.DatabaseError("finding timetable", err)
	}

	return s.toResponse(&timetable), nil
}

func (s *timetableService) GetAll(ctx context.Context, req dto.PaginationRequest) (*dto.PaginatedResponse, error) {
	var timetables []models.Timetable
	var total int64

	query := s.db.Model(&models.Timetable{})

	if req.Search != "" {
		search := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where("LOWER(classroom) LIKE ? OR LOWER(days) LIKE ?", search, search)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, errors.DatabaseError("counting timetables", err)
	}

	if err := query.Order(req.GetOrderBy()).
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Preload("Groups").
		Find(&timetables).Error; err != nil {
		return nil, errors.DatabaseError("listing timetables", err)
	}

	responses := make([]dto.TimetableResponse, len(timetables))
	for i, t := range timetables {
		responses[i] = *s.toResponse(&t)
	}

	return &dto.PaginatedResponse{
		Success:    true,
		Data:       responses,
		Pagination: dto.NewPaginationMetadata(req.Page, req.PageSize, total),
	}, nil
}

func (s *timetableService) toResponse(t *models.Timetable) *dto.TimetableResponse {
	groups := make([]dto.GroupSimple, len(t.Groups))
	for i, g := range t.Groups {
		groups[i] = dto.GroupSimple{
			ID:        g.ID,
			Name:      g.Name,
			StartDate: g.StartDate,
			Capacity:  g.Capacity,
		}
	}

	return &dto.TimetableResponse{
		ID:        t.ID,
		StartTime: t.StartTime,
		EndTime:   t.EndTime,
		Days:      t.Days,
		Classroom: t.Classroom,
		Groups:    groups,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}
