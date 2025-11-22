package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
	"gorm.io/gorm"
)

// RequireRole creates a middleware that checks if the user has the required role
func RequireRole(allowedRoles ...models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context (set by AuthMiddleware)
		userInterface, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "User not authenticated",
			})
			c.Abort()
			return
		}

		user, ok := userInterface.(*models.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Invalid user data",
			})
			c.Abort()
			return
		}

		// Check if user's role is in the allowed roles
		for _, role := range allowedRoles {
			if user.Role == role {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Insufficient permissions",
		})
		c.Abort()
	}
}

// RequirePermission creates a middleware that checks if the user has the required permission
// Permission format: "resource:action" e.g., "students:create", "courses:update"
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "User not authenticated",
			})
			c.Abort()
			return
		}

		user, ok := userInterface.(*models.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Invalid user data",
			})
			c.Abort()
			return
		}

		// Admin has all permissions
		if user.Role == models.RoleAdmin {
			c.Next()
			return
		}

		// Check permission in database
		// This would normally use the UserService, but for middleware we'll access DB directly
		db, exists := c.Get("db")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Database not available",
			})
			c.Abort()
			return
		}

		// Parse permission
		parts := strings.Split(permission, ":")
		if len(parts) != 2 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Invalid permission format",
			})
			c.Abort()
			return
		}

		resource, action := parts[0], parts[1]

		// Check if role has permission
		var count int64
		err := db.(*gorm.DB).Table("role_permissions").
			Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
			Where("role_permissions.role = ? AND permissions.resource = ? AND permissions.action = ?", user.Role, resource, action).
			Count(&count).Error

		if err != nil || count == 0 {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireOwnership checks if the user owns the resource (for students/teachers accessing their own data)
func RequireOwnership(resourceType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "User not authenticated",
			})
			c.Abort()
			return
		}

		user, ok := userInterface.(*models.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Invalid user data",
			})
			c.Abort()
			return
		}

		// Admin can access everything
		if user.Role == models.RoleAdmin {
			c.Next()
			return
		}

		// Get resource ID from URL
		var resourceID string
		switch resourceType {
		case "student":
			resourceID = c.Param("studentID")
			if user.StudentID != nil && user.StudentID.String() == resourceID {
				c.Next()
				return
			}
		case "teacher":
			resourceID = c.Param("teacherID")
			if user.TeacherID != nil && user.TeacherID.String() == resourceID {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "You can only access your own data",
		})
		c.Abort()
	}
}
