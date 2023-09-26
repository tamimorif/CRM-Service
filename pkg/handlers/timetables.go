package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/paraparadox/datetime"
	"github.com/softclub-go-0-0/crm-service/pkg/helpers"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func (h *handler) GetAllTimetables(c *gin.Context) {
	var timetables []models.Timetable
	result := h.DB.Find(&timetables)
	if result.Error != nil {
		log.Println("DB error - cannot find timetables:", result.Error)
		helpers.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, timetables)
}

func (h *handler) CreateTimetable(c *gin.Context) {
	var timetable models.Timetable

	if err := c.ShouldBindJSON(&timetable); err != nil {
		log.Println("binding timetable data:", err)
		helpers.UnprocessableEntity(c, err)
		return
	}

	if err := h.DB.Create(&timetable).Error; err != nil {
		log.Println("inserting timetable data to DB:", err)
		helpers.InternalServerError(c)
		return
	}

	c.JSON(http.StatusCreated, timetable)
}

func (h *handler) GetOneTimetable(c *gin.Context) {
	var timetable models.Timetable

	if err := h.DB.First(&timetable, "id = ?", c.Param("timetableID")).Error; err != nil {
		log.Println("getting a timetable:", err)
		helpers.NotFound(c, "timetable")
		return
	}

	c.JSON(http.StatusOK, timetable)
}

type timetableDataForUpdate struct {
	Classroom string        `json:"classroom" binding:"required"`
	Start     datetime.Time `json:"start" binding:"required"`
	Finish    datetime.Time `json:"finish" binding:"required"`
}

func (h *handler) UpdateTimetable(c *gin.Context) {
	var timetable models.Timetable

	if err := h.DB.Where("id = ?", c.Param("timetableID")).First(&timetable).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("getting a timetable:", err)
			helpers.NotFound(c, "timetable")
			return
		}
		log.Println("getting a timetable:", err)
		helpers.InternalServerError(c)
		return
	}

	var timetableData timetableDataForUpdate

	if err := c.ShouldBindJSON(&timetableData); err != nil {
		log.Println("binding timetable data:", err)
		helpers.UnprocessableEntity(c, err)
		return
	}

	if err := h.DB.Model(&timetable).Updates(timetableData).Error; err != nil {
		log.Println("updating timetable data in DB:", err)
		helpers.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, timetable)
}

func (h *handler) DeleteTimetable(c *gin.Context) {
	var timetable models.Timetable

	if err := h.DB.Where("id = ?", c.Param("timetableID")).First(&timetable).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("getting a timetable:", err)
			helpers.NotFound(c, "timetable")
			return
		}
		log.Println("getting a timetable:", err)
		helpers.InternalServerError(c)
		return
	}

	if err := h.DB.Delete(&timetable).Error; err != nil {
		log.Println("deleting timetable data from DB:", err)
		helpers.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, timetable)
}
