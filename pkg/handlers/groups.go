package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/softclub-go-0-0/tamim-crm-service/pkg/models"
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) GetAllCourses(c *gin.Context) {
	var courses []models.Course
	result := h.DB.Find(&courses)
	if result.Error != nil {
		log.Println("DB error - cannot find courses:", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, courses)
}
func (h *Handler) CreateCourse(c *gin.Context) {
	var course models.Course
	err := c.ShouldBindJSON(&course)
	if err != nil {
		log.Fatal("Creating teacher:", err)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": "validation error",
			"err":     err.Error(),
		})
		return
	}
	h.DB.Create(&course)
	c.JSON(http.StatusCreated, course)
}
func (h *Handler) GetOneCourse(c *gin.Context) {
	courseID, err := strconv.Atoi(c.Param("courseID"))
	if err != nil {
		log.Println("client error - bad courseID param:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	var course models.Course
	result := h.DB.First(&course, courseID)
	if result.Error != nil {
		log.Println("client error - cannot find teacher:", result.Error)
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not Found",
		})
		return
	}

	c.JSON(http.StatusOK, course)
}
func (h *Handler) UpdateCourse(c *gin.Context) {
	var course models.Course
	err := h.DB.Where("id = ?", c.Param("courseID")).First(&course).Error
	if err != nil {
		log.Println("getting a course:", err)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "There is no such course",
		})
		return
	}
	err = c.ShouldBindJSON(&course)
	if err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}
	h.DB.Save(&course)
}
func (h *Handler) DeleteCourse(c *gin.Context) {
	courseID, err := strconv.Atoi(c.Param("courseID"))
	if err != nil {
		log.Println("client error - bad teacherID param:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}
	var course models.Course
	result := h.DB.First(&course, courseID)
	if result.Error != nil {
		log.Println("client error - cannot find teacher:", result.Error)
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not Found",
		})
		return
	}
	h.DB.Delete(&course, courseID)
	c.JSON(http.StatusOK, course)
}
