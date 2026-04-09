package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/auth"
	"github.com/softclub-go-0-0/crm-service/pkg/config"
	"github.com/softclub-go-0-0/crm-service/pkg/errors"
	"github.com/softclub-go-0-0/crm-service/pkg/logger"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var authClient auth.AuthClient

// InitAuthClient initializes the gRPC auth client
func InitAuthClient(cfg *config.Config) error {
	conn, err := grpc.Dial(cfg.Auth.ServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	authClient = auth.NewAuthClient(conn)
	return nil
}

// AuthMiddleware handles authentication via gRPC auth service
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip auth for development if configured
		if cfg.Auth.SkipAuth {
			c.Next()
			return
		}

		tokenString := c.GetHeader("X-Auth-Token")
		if tokenString == "" {
			tokenString = c.GetHeader("Authorization")
			if tokenString != "" {
				// Remove "Bearer " prefix if present
				if len(tokenString) > 7 && strings.HasPrefix(tokenString, "Bearer ") {
					tokenString = tokenString[7:]
				}
			}
		}

		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Unauthorized - token required",
				"code":    errors.ErrCodeUnauthorized,
			})
			return
		}

		if authClient == nil {
			logger.Error("Auth service client not initialized", nil)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Authentication service unavailable",
				"code":    errors.ErrCodeAuthServiceUnavailable,
			})
			return
		}

		// Make gRPC call to auth service with timeout
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Auth.Timeout)
		defer cancel()

		authReq := &auth.AuthenticateRequest{
			Token: tokenString,
		}

		authResp, err := authClient.Authenticate(ctx, authReq)
		if err != nil {
			logger.Error("Failed to authenticate token", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Authentication failed",
				"code":    errors.ErrCodeInvalidToken,
			})
			return
		}

		if !authResp.Authenticated {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid token",
				"code":    errors.ErrCodeInvalidToken,
			})
			return
		}

		// Store user info in context for use in handlers and RBAC middleware
		if authResp.User != nil {
			// Create user object for RBAC middleware compatibility
			userRole := models.RoleStaff // Default role
			if len(authResp.User.Roles) > 0 {
				switch authResp.User.Roles[0] {
				case "admin":
					userRole = models.RoleAdmin
				case "teacher":
					userRole = models.RoleTeacher
				case "student":
					userRole = models.RoleStudent
				case "staff":
					userRole = models.RoleStaff
				}
			}

			user := &models.User{
				Email:     authResp.User.Email,
				FirstName: authResp.User.FirstName,
				LastName:  authResp.User.LastName,
				Role:      userRole,
			}
			// Parse UUID if valid
			if uid, err := uuid.Parse(authResp.User.Id); err == nil {
				user.ID = uid
			}

			c.Set("user", user)
			c.Set("user_id", authResp.User.Id)
			c.Set("user_name", authResp.User.FirstName+" "+authResp.User.LastName)
			c.Set("user_email", authResp.User.Email)
			c.Set("user_roles", authResp.User.Roles)
		}

		c.Next()
	}
}
