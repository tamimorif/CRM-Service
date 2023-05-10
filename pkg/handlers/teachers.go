package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/softclub-go-0-0/tamim-crm-service/pkg/models"
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) GetAllTeachers(c *gin.Context) {
	var teachers []models.Teacher
	result := h.DB.Find(&teachers)
	if result.Error != nil {
		log.Println("DB error - cannot find teachers:", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, teachers)
}
func (h *Handler) CreateTeacher(c *gin.Context) {
	var teacher models.Teacher
	err := c.ShouldBindJSON(&teacher)
	if err != nil {
		log.Fatal("Creating teacher:", err)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": "validation error",
			"err":     err.Error(),
		})
		return
	}
	h.DB.Create(&teacher)
	c.JSON(http.StatusCreated, teacher)
}
func (h *Handler) GetOneTeacher(c *gin.Context) {
	teacherID, err := strconv.Atoi(c.Param("teacherID"))
	if err != nil {
		log.Println("client error - bad teacherID param:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	var teacher models.Teacher
	result := h.DB.First(&teacher, teacherID)
	if result.Error != nil {
		log.Println("client error - cannot find teacher:", result.Error)
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not Found",
		})
		return
	}

	c.JSON(http.StatusOK, teacher)
}
func (h *Handler) UpdateTeacher(c *gin.Context) {
	var teacher models.Teacher
	err := h.DB.Where("id = ?", c.Param("teacherID")).First(&teacher).Error
	if err != nil {
		log.Println("getting a teacher:", err)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "There is no such teacher",
		})
		return
	}
	err = c.ShouldBindJSON(&teacher)
	if err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}
	h.DB.Save(&teacher)
}
func (h *Handler) DeleteTeacher(c *gin.Context) {
	teacherID, err := strconv.Atoi(c.Param("teacherID"))
	if err != nil {
		log.Println("client error - bad teacherID param:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}
	var teacher models.Teacher
	result := h.DB.First(&teacher, teacherID)
	if result.Error != nil {
		log.Println("client error - cannot find teacher:", result.Error)
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not Found",
		})
		return
	}
	h.DB.Delete(&teacher, teacherID)
	c.JSON(http.StatusOK, teacher)
}
