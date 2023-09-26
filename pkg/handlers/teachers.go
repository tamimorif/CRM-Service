package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/softclub-go-0-0/crm-service/pkg/helpers"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func (h *handler) GetAllTeachers(c *gin.Context) {
	var teachers []models.Teacher
	result := h.DB.Find(&teachers)
	if result.Error != nil {
		log.Println("DB error - cannot find teachers:", result.Error)
		helpers.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, teachers)
}

func (h *handler) CreateTeacher(c *gin.Context) {
	var teacher models.Teacher

	if err := c.ShouldBindJSON(&teacher); err != nil {
		log.Println("binding teacher data:", err)
		helpers.UnprocessableEntity(c, err)
		return
	}

	if err := h.DB.Create(&teacher).Error; err != nil {
		log.Println("inserting teacher data to DB:", err)
		helpers.InternalServerError(c)
		return
	}

	c.JSON(http.StatusCreated, teacher)
}

func (h *handler) GetOneTeacher(c *gin.Context) {
	var teacher models.Teacher

	if err := h.DB.First(&teacher, "id = ?", c.Param("teacherID")).Error; err != nil {
		log.Println("getting a teacher:", err)
		helpers.NotFound(c, "teacher")
		return
	}

	c.JSON(http.StatusOK, teacher)
}

type teacherDataForUpdate struct {
	Name    string `json:"name" binding:"required,alphaunicode"`
	Surname string `json:"surname" binding:"required,alphaunicode"`
	Phone   string `json:"phone" binding:"required,len=12,numeric"`
	Email   string `json:"email" binding:"omitempty,email"`
}

func (h *handler) UpdateTeacher(c *gin.Context) {
	var teacher models.Teacher

	if err := h.DB.Where("id = ?", c.Param("teacherID")).First(&teacher).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("getting a teacher:", err)
			helpers.NotFound(c, "teacher")
			return
		}
		log.Println("getting a teacher:", err)
		helpers.InternalServerError(c)
		return
	}

	var teacherData teacherDataForUpdate

	if err := c.ShouldBindJSON(&teacherData); err != nil {
		log.Println("binding teacher data:", err)
		helpers.UnprocessableEntity(c, err)
		return
	}

	if err := h.DB.Model(&teacher).Updates(teacherData).Error; err != nil {
		log.Println("updating teacher data in DB:", err)
		helpers.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, teacher)
}

func (h *handler) DeleteTeacher(c *gin.Context) {
	var teacher models.Teacher

	if err := h.DB.Where("id = ?", c.Param("teacherID")).First(&teacher).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("getting a teacher:", err)
			helpers.NotFound(c, "teacher")
			return
		}
		log.Println("getting a teacher:", err)
		helpers.InternalServerError(c)
		return
	}

	if err := h.DB.Delete(&teacher).Error; err != nil {
		log.Println("deleting teacher data from DB:", err)
		helpers.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, teacher)
}
