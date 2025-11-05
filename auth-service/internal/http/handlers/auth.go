package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"auth-service/internal/auth"
)

type Handler struct {
	Redis *redis.Client
}

func NewHandler(rdb *redis.Client) *Handler {
	return &Handler{Redis: rdb}
}

var ctx = context.Background()


func (h *Handler) LoginHandler(c *gin.Context) {
	if tok, err := c.Cookie("access_token"); err == nil && tok != "" {
		if _, err := auth.ParseAndValidate(tok); err == nil {
			exists, _ := h.Redis.Exists(ctx, "session:"+tok).Result()
			if exists == 1 {
				c.String(http.StatusOK, "Already logged in.")
				return
			}
		}
	}

	state := auth.RandomKey(20)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   300,
	})
	c.Redirect(http.StatusTemporaryRedirect, auth.GetAuthURL(state))
}


func (h *Handler) CallbackHandler(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	if state == "" || code == "" {
		c.String(http.StatusBadRequest, "Missing state or code")
		return
	}

	// Validate OAuth state
	if ck, err := c.Cookie("oauth_state"); err != nil || ck != state {
		c.String(http.StatusBadRequest, "Invalid OAuth state")
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:   "oauth_state",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	gu, err := auth.ExchangeCodeForUser(code)
	if err != nil {
		log.Println("Token exchange error:", err)
		c.String(http.StatusInternalServerError, "Token exchange failed")
		return
	}

	// Save or update user info in your DB
	if err := auth.UpsertUserInfo(gu); err != nil {
		log.Println("DB upsert error:", err)
		c.String(http.StatusInternalServerError, "Failed to update user info")
		return
	}

	token, err := auth.GenerateAccessToken(gu.Sub, gu.Email, "user", 15*time.Minute)
	if err != nil {
		log.Println("JWT generation error:", err)
		c.String(http.StatusInternalServerError, "Failed to sign JWT")
		return
	}

	err = h.Redis.Set(ctx, "session:"+token, gu.Email, 15*time.Minute).Err()
	if err != nil {
		log.Println("Redis set session error:", err)
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // change to true in production (HTTPS)
		SameSite: http.SameSiteLaxMode,
		MaxAge:   15 * 60,
	})

	c.String(http.StatusOK, "✅ Logged in successfully.")
}


func (h *Handler) LogoutHandler(c *gin.Context) {
	if tok, err := c.Cookie("access_token"); err == nil && tok != "" {
		h.Redis.Del(ctx, "session:"+tok)
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	c.String(http.StatusOK, "Logged out successfully.")
}
