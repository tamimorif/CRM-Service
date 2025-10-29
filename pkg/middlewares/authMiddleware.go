package middlewares

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/softclub-go-0-0/crm-service/pkg/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"os"
	"time"
)

var authClient auth.AuthClient

func init() {
	authServiceAddr := os.Getenv("AUTH_SERVICE_ADDR")
	if authServiceAddr == "" {
		authServiceAddr = "localhost:50051" // default
	}

	conn, err := grpc.Dial(authServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to auth service: %v", err)
		return
	}

	authClient = auth.NewAuthClient(conn)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("X-Auth-Token")
		if tokenString == "" {
			tokenString = c.GetHeader("Authorization")
			if tokenString != "" {
				// Remove "Bearer " prefix if present
				if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
					tokenString = tokenString[7:]
				}
			}
		}

		if tokenString == "" {
			log.Print("empty token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized - token required",
			})
			return
		}

		// Skip auth for development if SKIP_AUTH is set
		if os.Getenv("SKIP_AUTH") == "true" {
			log.Println("Skipping auth for development")
			c.Next()
			return
		}

		if authClient == nil {
			log.Print("auth service client not initialized")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Authentication service unavailable",
			})
			return
		}

		// Make gRPC call to auth service
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		authReq := &auth.AuthenticateRequest{
			Token: tokenString,
		}

		authResp, err := authClient.Authenticate(ctx, authReq)
		if err != nil {
			log.Printf("Failed to authenticate token: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Authentication failed",
			})
			return
		}

		if !authResp.Authenticated {
			log.Print("token authentication failed")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token",
			})
			return
		}

		// Store user info in context for use in handlers
		if authResp.User != nil {
			c.Set("user_id", authResp.User.Id)
			c.Set("user_name", authResp.User.FirstName+" "+authResp.User.LastName)
			c.Set("user_email", authResp.User.Email)
			c.Set("user_roles", authResp.User.Roles)
		}

		c.Next()
	}
}
