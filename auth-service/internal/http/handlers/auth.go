package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"auth-service/internal/auth"
)

func LoginHandler(c *gin.Context) {
	// If already has valid JWT cookie, skip Google
	if tok, err := c.Cookie("access_token"); err == nil && tok != "" {
		if _, err := auth.ParseAndValidate(tok); err == nil {
			// c.Redirect(http.StatusFound, "/app/me")
			return
		}
	}
	state := auth.RandomKey(20)
	http.SetCookie(c.Writer, &http.Cookie{Name: "oauth_state", Value: state, Path: "/", HttpOnly: true, MaxAge: 300})
	c.Redirect(http.StatusTemporaryRedirect, auth.GetAuthURL(state))
}

func CallbackHandler(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")
	if state == "" || code == "" {
		c.String(http.StatusBadRequest, "missing state or code")
		return
	}
	if ck, err := c.Cookie("oauth_state"); err != nil || ck != state {
		c.String(http.StatusBadRequest, "invalid oauth state")
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{Name: "oauth_state", Value: "", Path: "/", MaxAge: -1})

	gu, err := auth.ExchangeCodeForUser(code)
	if err != nil {
		log.Println("token exchange error:", err)
		c.String(http.StatusInternalServerError, "token exchange failed")
		return
	}

	// Upsert DB (you already have these functions)
	err = auth.UpsertUserInfo(gu)
	if err != nil {
		log.Println("upsert identity error:", err)
		c.String(http.StatusInternalServerError, "db upsert failed")
		return
	}

	// Create JWT (15 minutes). Claim `uid` can be your internal UUID or provider `sub`.
	token, err := auth.GenerateAccessToken(gu.Sub, gu.Email, "user", 15*time.Minute)
	if err != nil {
		log.Println("jwt sign error:", err)
		c.String(http.StatusInternalServerError, "cannot sign token")
		return
	}

	// Set HttpOnly cookie (or return JSON if SPA/mobile wants bearer)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // true in production behind HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   15 * 60,
	})

	// c.Redirect(http.StatusTemporaryRedirect, "/app/me")
}

func LogoutHandler(c *gin.Context) {
	// Stateless JWT: deleting cookie is enough (or maintain a denylist if needed)
	http.SetCookie(c.Writer, &http.Cookie{Name: "access_token", Value: "", Path: "/", MaxAge: -1, HttpOnly: true})
	c.String(http.StatusOK, "Logged out.")
}
