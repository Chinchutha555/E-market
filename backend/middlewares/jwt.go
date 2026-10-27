package middlewares

import (
	"net/http"
	"strings"

	"emarket-api/helpers"
	"emarket-api/responses"
	"github.com/gin-gonic/gin"
)

const (
	ContextUserID   = "userID"
	ContextUsername = "username"
	ContextRole     = "userRole"
)

func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, responses.FailRespone(responses.NewAppErr(http.StatusUnauthorized, "ไม่พบ Authorization header")))
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, responses.FailRespone(responses.NewAppErr(http.StatusUnauthorized, "รูปแบบ Token ไม่ถูกต้อง (ต้องเป็น Bearer Token)")))
			return
		}

		tokenString := parts[1]
		claims, err := helpers.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, responses.FailRespone(responses.NewAppErr(http.StatusUnauthorized, "Token ไม่ถูกต้องหรือหมดอายุ")))
			return
		}

		// Identity and role always come from the signed token. Client headers must
		// never be allowed to impersonate another account or elevate privileges.
		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextUsername, claims.Username)
		c.Set(ContextRole, claims.Role)
		c.Next()
	}
}

func RequireRoles(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get(ContextRole)
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, responses.FailRespone(responses.NewAppErr(http.StatusForbidden, "ไม่มีสิทธิ์เข้าถึงฟังก์ชันนี้")))
			return
		}

		role, ok := roleVal.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, responses.FailRespone(responses.NewAppErr(http.StatusForbidden, "ไม่มีสิทธิ์เข้าถึงฟังก์ชันนี้")))
			return
		}

		hasPermission := false
		for _, r := range allowedRoles {
			if strings.EqualFold(r, role) || r == role {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.AbortWithStatusJSON(http.StatusForbidden, responses.FailRespone(responses.NewAppErr(http.StatusForbidden, "สิทธิ์ของคุณไม่เพียงพอสำหรับการดำเนินการนี้")))
			return
		}

		c.Next()
	}
}

func GetUserID(c *gin.Context) int {
	val, exists := c.Get(ContextUserID)
	if !exists {
		return 0
	}
	id, ok := val.(int)
	if !ok {
		return 0
	}
	return id
}

func GetUserRole(c *gin.Context) string {
	val, exists := c.Get(ContextRole)
	if !exists {
		return ""
	}
	role, ok := val.(string)
	if !ok {
		return ""
	}
	return role
}
