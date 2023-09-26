package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/softclub-go-0-0/crm-service/pkg/helpers"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func (h *handler) GetAllCourses(c *gin.Context) {
	var courses []models.Course
	result := h.DB.Find(&courses)
	if result.Error != nil {
		log.Println("DB error - cannot find courses:", result.Error)
		helpers.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, courses)
}

func (h *handler) CreateCourse(c *gin.Context) {
	var course models.Course

	if err := c.ShouldBindJSON(&course); err != nil {
		log.Println("binding course data:", err)
		helpers.UnprocessableEntity(c, err)
		return
	}

	if err := h.DB.Create(&course).Error; err != nil {
		log.Println("inserting course data to DB:", err)
		helpers.InternalServerError(c)
		return
	}

	c.JSON(http.StatusCreated, course)
}

func (h *handler) GetOneCourse(c *gin.Context) {
	var course models.Course

	if err := h.DB.First(&course, "id = ?", c.Param("courseID")).Error; err != nil {
		log.Println("getting a course:", err)
		helpers.NotFound(c, "course")
		return
	}

	c.JSON(http.StatusOK, course)
}

type courseDataForUpdate struct {
	Title      string `json:"title" binding:"required"`
	MonthlyFee uint   `json:"monthly_fee" binding:"omitempty,number"`
	Duration   uint   `json:"duration" binding:"omitempty,number"`
}

func (h *handler) UpdateCourse(c *gin.Context) {
	var course models.Course

	if err := h.DB.Where("id = ?", c.Param("courseID")).First(&course).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("getting a course:", err)
			helpers.NotFound(c, "course")
			return
		}
		log.Println("getting a course:", err)
		helpers.InternalServerError(c)
		return
	}

	var courseData courseDataForUpdate

	if err := c.ShouldBindJSON(&courseData); err != nil {
		log.Println("binding course data:", err)
		helpers.UnprocessableEntity(c, err)
		return
	}

	if err := h.DB.Model(&course).Updates(courseData).Error; err != nil {
		log.Println("updating course data in DB:", err)
		helpers.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, course)
}

func (h *handler) DeleteCourse(c *gin.Context) {
	var course models.Course

	if err := h.DB.Where("id = ?", c.Param("courseID")).First(&course).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("getting a course:", err)
			helpers.NotFound(c, "course")
			return
		}
		log.Println("getting a course:", err)
		helpers.InternalServerError(c)
		return
	}

	if err := h.DB.Delete(&course).Error; err != nil {
		log.Println("deleting course data from DB:", err)
		helpers.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, course)
}
