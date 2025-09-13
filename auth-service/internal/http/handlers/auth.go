package handlers

import (
	"log"
	"net/http"

	"auth-service/internal/auth"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// GET /login
func LoginHandler(c *gin.Context) {
	state := auth.RandomKey(20)
	sess := sessions.Default(c)
	sess.Set(auth.SessionKeyState, state)
	_ = sess.Save()

	url := auth.GetAuthURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GET /auth/callback
func CallbackHandler(c *gin.Context) {
	sess := sessions.Default(c)

	state := c.Query("state")
	code := c.Query("code")
	if state == "" || code == "" {
		c.String(http.StatusBadRequest, "missing state or code")
		return
	}

	if storedState, _ := sess.Get(auth.SessionKeyState).(string); storedState == "" || storedState != state {
		c.String(http.StatusBadRequest, "invalid oauth state")
		return
	}
	// clean up state
	sess.Delete(auth.SessionKeyState)
	_ = sess.Save()

	gu, err := auth.ExchangeCodeForUser(code)
	if err != nil {
		log.Println("token exchange error:", err)
		c.String(http.StatusInternalServerError, "token exchange failed")
		return
	}

	err = auth.UpsertUserInfo(gu)

	if err != nil {
		log.Println("Upsert user info:", err)
		c.String(http.StatusInternalServerError, "upsert user info failed")
		return
	}

	// Store minimal user info in session
	sess.Set(auth.SessionKeyUser, gu)
	_ = sess.Save()

	// c.Redirect(http.StatusTemporaryRedirect, "/app/me")
}

// GET /logout
func LogoutHandler(c *gin.Context) {
	sess := sessions.Default(c)
	sess.Clear()
	_ = sess.Save()
	c.String(http.StatusOK, "Logged out.")
}
