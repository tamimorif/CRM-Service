package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/paraparadox/datetime"
	"github.com/softclub-go-0-0/crm-service/pkg/helpers"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func (h *handler) GetAllGroups(c *gin.Context) {
	var groups []models.Group
	result := h.DB.Find(&groups)
	if result.Error != nil {
		log.Println("DB error - cannot find groups:", result.Error)
		helpers.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, groups)
}

func (h *handler) CreateGroup(c *gin.Context) {
	var group models.Group

	if err := c.ShouldBindJSON(&group); err != nil {
		log.Println("binding group data:", err)
		helpers.UnprocessableEntity(c, err)
		return
	}

	if err := h.DB.Create(&group).Error; err != nil {
		log.Println("inserting group data to DB:", err)
		helpers.InternalServerError(c)
		return
	}

	c.JSON(http.StatusCreated, group)
}

func (h *handler) GetOneGroup(c *gin.Context) {
	var group models.Group

	if err := h.DB.First(&group, "id = ?", c.Param("groupID")).Error; err != nil {
		log.Println("getting a group:", err)
		helpers.NotFound(c, "group")
		return
	}

	c.JSON(http.StatusOK, group)
}

type groupDataForUpdate struct {
	CourseID    uuid.UUID     `json:"course_id"`
	TeacherID   uuid.UUID     `json:"teacher_id"`
	TimetableID uuid.UUID     `json:"timetable_id"`
	Title       string        `json:"title" binding:"required"`
	StartDate   datetime.Date `json:"start_date" binding:"required"`
}

func (h *handler) UpdateGroup(c *gin.Context) {
	var group models.Group

	if err := h.DB.Where("id = ?", c.Param("groupID")).First(&group).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("getting a group:", err)
			helpers.NotFound(c, "group")
			return
		}
		log.Println("getting a group:", err)
		helpers.InternalServerError(c)
		return
	}

	var groupData groupDataForUpdate

	if err := c.ShouldBindJSON(&groupData); err != nil {
		log.Println("binding group data:", err)
		helpers.UnprocessableEntity(c, err)
		return
	}

	if err := h.DB.Model(&group).Updates(groupData).Error; err != nil {
		log.Println("updating group data in DB:", err)
		helpers.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, group)
}

func (h *handler) DeleteGroup(c *gin.Context) {
	var group models.Group

	if err := h.DB.Where("id = ?", c.Param("groupID")).First(&group).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("getting a group:", err)
			helpers.NotFound(c, "group")
			return
		}
		log.Println("getting a group:", err)
		helpers.InternalServerError(c)
		return
	}

	if err := h.DB.Delete(&group).Error; err != nil {
		log.Println("deleting group data from DB:", err)
		helpers.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, group)
}
