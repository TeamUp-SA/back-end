package handlers

import (
	"net/http"

	"user-service/internal/auth"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetMe(c *gin.Context) {
	session := sessions.Default(c)
	userData := session.Get(auth.SessionKeyUser)
	if userData == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, ok := userData.(auth.GoogleUser)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"email":   user.Email,
		"name":    user.Name,
		"picture": user.Picture,
		"locale":  user.Locale,
	})
}
