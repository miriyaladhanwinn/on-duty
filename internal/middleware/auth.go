package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/miriyaladhanwinn/on-duty/internal/auth"
	"github.com/miriyaladhanwinn/on-duty/internal/db"
)

const ClaimsKey = "currentUserClaims"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format. Expected Bearer <token>"})
			c.Abort()
			return
		}

		claims, err := auth.ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired session token"})
			c.Abort()
			return
		}

		c.Set(ClaimsKey, claims)
		c.Next()
	}
}

func RequireRoles(allowedRoles ...db.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get(ClaimsKey)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthenticated context"})
			c.Abort()
			return
		}

		claims, ok := val.(*auth.Claims)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid context claims"})
			c.Abort()
			return
		}

		for _, role := range allowedRoles {
			if claims.Role == role || claims.Role == db.RoleDeptAdmin {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error": "Forbidden: insufficient permissions for this workflow action",
			"user_role": claims.Role,
		})
		c.Abort()
	}
}
