package auth

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// SetUserToSession stores the authenticated GoogleUser into the current session.
func SetUserToSession(c *gin.Context, user GoogleUser) {
	sess := sessions.Default(c)
	sess.Set(SessionKeyUser, user)
	_ = sess.Save()
}

// GetUserFromSession retrieves the GoogleUser from the current session.
// Returns (user, true) when found; otherwise (zeroUser, false).
func GetUserFromSession(c *gin.Context) (GoogleUser, bool) {
	sess := sessions.Default(c)
	v := sess.Get(SessionKeyUser)
	if v == nil {
		return GoogleUser{}, false
	}
	u, ok := v.(GoogleUser)
	return u, ok
}

// ClearSession removes all data from the current session and saves it.
func ClearSession(c *gin.Context) {
	sess := sessions.Default(c)
	sess.Clear()
	_ = sess.Save()
}

// SetState saves the OAuth state in the session (used for CSRF protection).
func SetState(c *gin.Context, state string) {
	sess := sessions.Default(c)
	sess.Set(SessionKeyState, state)
	_ = sess.Save()
}

// PopState retrieves and deletes the OAuth state from the session.
// Returns (state, true) if present; otherwise ("", false).
func PopState(c *gin.Context) (string, bool) {
	sess := sessions.Default(c)
	v := sess.Get(SessionKeyState)
	s, ok := v.(string)
	if !ok || s == "" {
		return "", false
	}
	sess.Delete(SessionKeyState)
	_ = sess.Save()
	return s, true
}
