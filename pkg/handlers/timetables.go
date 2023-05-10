package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/softclub-go-0-0/tamim-crm-service/pkg/models"
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) GetAllTimetables(c *gin.Context) {
	var timetable []models.TimeTable
	result := h.DB.Find(&timetable)
	if result.Error != nil {
		log.Println("DB error - cannot find teachers:", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, timetable)
}
func (h *Handler) CreateTimetable(c *gin.Context) {
	var timetable models.TimeTable
	err := c.ShouldBindJSON(&timetable)
	if err != nil {
		log.Fatal("Creating teacher:", err)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": "validation error",
			"err":     err.Error(),
		})
		return
	}
	h.DB.Create(&timetable)
	c.JSON(http.StatusCreated, timetable)
}
func (h *Handler) GetOneTimetable(c *gin.Context) {
	timetableID, err := strconv.Atoi(c.Param("timetableID"))
	if err != nil {
		log.Println("client error - bad teacherID param:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	var timetable models.TimeTable
	result := h.DB.First(&timetable, timetableID)
	if result.Error != nil {
		log.Println("client error - cannot find teacher:", result.Error)
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not Found",
		})
		return
	}

	c.JSON(http.StatusOK, timetable)
}
func (h *Handler) UpdateTimetable(c *gin.Context) {
	var timetable models.TimeTable
	err := h.DB.Where("id = ?", c.Param("timetableID")).First(&timetable).Error
	if err != nil {
		log.Println("getting a timetable:", err)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "There is no such timetable",
		})
		return
	}
	err = c.ShouldBindJSON(&timetable)
	if err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}
	h.DB.Save(&timetable)
}
func (h *Handler) DeleteTimetable(c *gin.Context) {
	timetableID, err := strconv.Atoi(c.Param("timetableID"))
	if err != nil {
		log.Println("client error - bad teacherID param:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}
	var timetable models.TimeTable
	result := h.DB.First(&timetable, timetableID)
	if result.Error != nil {
		log.Println("client error - cannot find teacher:", result.Error)
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not Found",
		})
		return
	}
	h.DB.Delete(&timetable, timetableID)
	c.JSON(http.StatusOK, timetable)
}
