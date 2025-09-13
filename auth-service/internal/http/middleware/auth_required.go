package middleware

import (
	"net/http"

	"auth-service/internal/auth"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AuthRequired ensures a valid authenticated session exists.
// If a user is found in the session, it passes the request through and
// stores the GoogleUser in the Gin context under the key "googleUser".
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := sessions.Default(c)
		v := sess.Get(auth.SessionKeyUser)
		if v == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		if u, ok := v.(auth.GoogleUser); ok {
			c.Set("googleUser", u)
		}
		c.Next()
	}
}
