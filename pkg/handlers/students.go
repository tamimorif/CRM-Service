package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/helpers"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func (h *handler) GetAllStudents(c *gin.Context) {
	var students []models.Student
	result := h.DB.Find(&students)
	if result.Error != nil {
		log.Println("DB error - cannot find students:", result.Error)
		helpers.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, students)
}

func (h *handler) CreateStudent(c *gin.Context) {
	var group models.Group

	if err := h.DB.First(&group, "id = ?", c.Param("groupID")).Error; err != nil {
		log.Println("getting a group:", err)
		helpers.NotFound(c, "group")
		return
	}

	var student models.Student

	if err := c.ShouldBindJSON(&student); err != nil {
		log.Println("binding student data:", err)
		helpers.UnprocessableEntity(c, err)
		return
	}

	student.GroupID = group.ID

	if err := h.DB.Create(&student).Error; err != nil {
		log.Println("inserting student data to DB:", err)
		helpers.InternalServerError(c)
		return
	}

	c.JSON(http.StatusCreated, student)
}

func (h *handler) GetOneStudent(c *gin.Context) {
	var student models.Student

	if err := h.DB.First(&student, "id = ?", c.Param("studentID")).Error; err != nil {
		log.Println("getting a student:", err)
		helpers.NotFound(c, "student")
		return
	}

	c.JSON(http.StatusOK, student)
}

type studentDataForUpdate struct {
	GroupID uuid.UUID `json:"group_id"`
	Name    string    `json:"name" binding:"required,alphaunicode"`
	Surname string    `json:"surname" binding:"required,alphaunicode"`
	Phone   string    `json:"phone" binding:"required,len=12,numeric"`
	Email   string    `json:"email" binding:"omitempty,email"`
}

func (h *handler) UpdateStudent(c *gin.Context) {
	var student models.Student

	if err := h.DB.Where("id = ?", c.Param("studentID")).First(&student).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("getting a student:", err)
			helpers.NotFound(c, "student")
			return
		}
		log.Println("getting a student:", err)
		helpers.InternalServerError(c)
		return
	}

	var studentData studentDataForUpdate

	if err := c.ShouldBindJSON(&studentData); err != nil {
		log.Println("binding student data:", err)
		helpers.UnprocessableEntity(c, err)
		return
	}

	if err := h.DB.Model(&student).Updates(studentData).Error; err != nil {
		log.Println("updating student data in DB:", err)
		helpers.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, student)
}

func (h *handler) DeleteStudent(c *gin.Context) {
	var student models.Student

	if err := h.DB.Where("id = ?", c.Param("studentID")).First(&student).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("getting a student:", err)
			helpers.NotFound(c, "student")
			return
		}
		log.Println("getting a student:", err)
		helpers.InternalServerError(c)
		return
	}

	if err := h.DB.Delete(&student).Error; err != nil {
		log.Println("deleting student data from DB:", err)
		helpers.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, student)
}
