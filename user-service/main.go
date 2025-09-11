package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"encoding/gob"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	oauthConfig *oauth2.Config
)

const (
	sessionName     = "appsession"
	sessionKeyUser  = "user"
	sessionKeyState = "oauthState"
)

type GoogleUser struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Locale        string `json:"locale"`
}

func loadEnv() {
	// Load .env if present
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env vars")
	}
}

func main() {
	loadEnv()
	gob.Register(GoogleUser{})

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL := os.Getenv("OAUTH_REDIRECT_URL")
	sessionSecret := os.Getenv("SESSION_SECRET")

	if clientID == "" || clientSecret == "" || redirectURL == "" {
		log.Fatal("Missing GOOGLE_CLIENT_ID / GOOGLE_CLIENT_SECRET / OAUTH_REDIRECT_URL")
	}

	oauthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"openid", "email", "profile",
		},
		Endpoint: google.Endpoint,
	}

	r := gin.Default()

	// Cookie session store (in prod, use strong key & HTTPS)
	store := cookie.NewStore([]byte(sessionSecret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60 * 8, // 8h
		HttpOnly: true,
		Secure:   false, // set true behind HTTPS
		SameSite: http.SameSiteLaxMode,
	})
	r.Use(sessions.Sessions(sessionName, store))

	// Public routes
	r.GET("/", func(c *gin.Context) {
		user := getUser(c)
		if user != nil {
			c.HTML(http.StatusOK, "home.html", gin.H{
				"loggedIn": true,
				"user":     user,
			})
			return
		}
		c.String(200, "Hello! Go to /login to sign in with Google.")
	})

	r.GET("/login", loginHandler)
	r.GET("/auth/callback", callbackHandler)
	r.GET("/logout", logoutHandler)

	// Protected routes
	auth := r.Group("/app")
	auth.Use(requireAuth())
	auth.GET("/me", func(c *gin.Context) {
		user := getUser(c)
		c.JSON(200, gin.H{
			"email":   user.Email,
			"name":    user.Name,
			"picture": user.Picture,
			"sub":     user.Sub,
		})
	})

	// (Optional) templates or static files if you want nicer pages
	// r.LoadHTMLGlob("templates/*.html")

	r.Run(":8080")
}

// ----- Handlers -----

func loginHandler(c *gin.Context) {
	state := randomKey(20)
	sess := sessions.Default(c)
	sess.Set(sessionKeyState, state)
	_ = sess.Save()

	url := oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func callbackHandler(c *gin.Context) {
	sess := sessions.Default(c)

	state := c.Query("state")
	code := c.Query("code")
	if state == "" || code == "" {
		c.String(http.StatusBadRequest, "missing state or code")
		return
	}

	if storedState, _ := sess.Get(sessionKeyState).(string); storedState == "" || storedState != state {
		c.String(http.StatusBadRequest, "invalid oauth state")
		return
	}
	// clean up state
	sess.Delete(sessionKeyState)
	_ = sess.Save()

	// Exchange code for token
	token, err := oauthConfig.Exchange(context.Background(), code)
	// log.Println("token: ", token)
	if err != nil {
		log.Println("token exchange error:", err)
		c.String(http.StatusInternalServerError, "token exchange failed")
		return
	}

	// Fetch userinfo via Google OIDC endpoint
	client := oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://openidconnect.googleapis.com/v1/userinfo")
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("userinfo error:", err, "status:", resp.StatusCode)
		c.String(http.StatusInternalServerError, "failed to get userinfo")
		return
	}
	defer resp.Body.Close()

	var gu GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&gu); err != nil {
		log.Println("decode userinfo:", err)
		c.String(http.StatusInternalServerError, "parse userinfo failed")
		return
	}
	// log.Println(gu.Email)
	// Store minimal user info in session
	sess.Set(sessionKeyUser, gu)
	_ = sess.Save()

	c.Redirect(http.StatusTemporaryRedirect, "/app/me")
}

func logoutHandler(c *gin.Context) {
	sess := sessions.Default(c)
	sess.Clear()
	_ = sess.Save()
	c.String(200, "Logged out.")
}

// ----- Middleware & helpers -----

func requireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if getUser(c) == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
			return
		}
		c.Next()
	}
}

func getUser(c *gin.Context) *GoogleUser {
	sess := sessions.Default(c)
	if v := sess.Get(sessionKeyUser); v != nil {
		if u, ok := v.(GoogleUser); ok {
			return &u
		}
	}
	return nil
}

func randomKey(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(b)
}
