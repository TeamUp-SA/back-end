package middleware

import (
    "net/http"
    "strings"
    "time"

    "github.com/gin-gonic/gin"

    "user-service/internal/auth"
)

const cookieName = "access_token"

func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        var token string

        if ah := c.GetHeader("Authorization"); strings.HasPrefix(strings.ToLower(ah), "bearer ") {
            token = strings.TrimSpace(ah[7:])
        }
        if token == "" {
            if ck, err := c.Cookie(cookieName); err == nil {
                token = ck
            }
        }

        if token == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
            return
        }

        claims, err := auth.ParseAndValidate(token)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            return
        }
        if claims.ExpiresAt != nil && time.Until(claims.ExpiresAt.Time) <= 0 {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
            return
        }

        c.Set("userID", claims.UserID)
        c.Set("email", claims.Email)
        c.Set("role", claims.Role)
        c.Next()
    }
}

