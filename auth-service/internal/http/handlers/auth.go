package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"auth-service/internal/auth"
	"auth-service/internal/config"
	"auth-service/internal/db"
)

type Handler struct {
	Redis      *redis.Client
	Config     *config.Config
	HTTPClient *http.Client
}

func NewHandler(rdb *redis.Client, cfg *config.Config) *Handler {
	client := &http.Client{
		Timeout: memberAPITimeout,
	}
	return &Handler{
		Redis:      rdb,
		Config:     cfg,
		HTTPClient: client,
	}
}

var (
	ctx               = context.Background()
	usernameSanitizer = regexp.MustCompile(`[^a-z0-9_]+`)
)

const (
	sessionKeyPrefix  = "session:"
	minPasswordLength = 8

	defaultSessionTTL = 24 * time.Hour
	oauthSessionTTL   = 15 * time.Minute
	oauthStateTTL     = 5 * time.Minute
	memberAPITimeout  = 5 * time.Second
)

var (
	errEmailAlreadyRegistered = errors.New("email already registered")
)

type registerRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Bio      string `json:"bio"`
	Github   string `json:"github"`
	LinkedIn string `json:"linkedin"`
	Website  string `json:"website"`

	Educations  []educationInput  `json:"educations"`
	Experiences []experienceInput `json:"experiences"`
	Skills      []string          `json:"skills"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type educationInput struct {
	School    string `json:"school"`
	Degree    string `json:"degree"`
	Field     string `json:"field"`
	StartYear *int   `json:"startYear"`
	EndYear   *int   `json:"endYear"`
}

type experienceInput struct {
	Title       string `json:"title"`
	Company     string `json:"company"`
	Description string `json:"description"`
	StartYear   *int   `json:"startYear"`
	EndYear     *int   `json:"endYear"`
}

type memberProfilePayload struct {
	Username    string                `json:"username"`
	FirstName   string                `json:"firstName"`
	LastName    string                `json:"lastName"`
	Email       string                `json:"email"`
	PhoneNumber string                `json:"phoneNumber,omitempty"`
	Bio         string                `json:"bio,omitempty"`
	LinkedIn    string                `json:"linkedIn,omitempty"`
	GitHub      string                `json:"github,omitempty"`
	Website     string                `json:"website,omitempty"`
	Skills      []string              `json:"skills,omitempty"`
	Experience  []memberExperienceDTO `json:"experience,omitempty"`
	Education   []memberEducationDTO  `json:"education,omitempty"`
}

type memberExperienceDTO struct {
	Title       string `json:"title"`
	Company     string `json:"company"`
	Description string `json:"description,omitempty"`
	StartYear   int    `json:"startYear,omitempty"`
	EndYear     int    `json:"endYear,omitempty"`
}

type memberEducationDTO struct {
	School    string `json:"school"`
	Degree    string `json:"degree,omitempty"`
	Field     string `json:"field,omitempty"`
	StartYear int    `json:"startYear,omitempty"`
	EndYear   int    `json:"endYear,omitempty"`
}

func sessionKey(token string) string {
	return sessionKeyPrefix + token
}

func splitName(fullName string) (string, string) {
	fullName = strings.TrimSpace(fullName)
	if fullName == "" {
		return "", ""
	}
	parts := strings.Fields(fullName)
	if len(parts) == 1 {
		return parts[0], parts[0]
	}
	return parts[0], strings.Join(parts[1:], " ")
}

func generateUniqueUsername(gormDB *gorm.DB, email, name string) (string, error) {
	base := strings.TrimSpace(strings.ToLower(email))
	if base != "" {
		if at := strings.Index(base, "@"); at > 0 {
			base = base[:at]
		}
	}
	if base == "" {
		base = strings.ToLower(strings.ReplaceAll(name, " ", ""))
	}
	base = usernameSanitizer.ReplaceAllString(base, "")
	if base == "" {
		base = "user"
	}
	if len(base) > 32 {
		base = base[:32]
	}

	username := base
	for {
		var count int64
		if err := gormDB.Model(&auth.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
			return "", err
		}
		if count == 0 {
			return username, nil
		}

		suffix := strings.ReplaceAll(uuid.NewString(), "-", "")
		if len(suffix) > 6 {
			suffix = suffix[:6]
		}
		maxBaseLen := 50 - len(suffix) - 1
		if maxBaseLen < 1 {
			maxBaseLen = 1
		}
		trimmed := base
		if len(trimmed) > maxBaseLen {
			trimmed = trimmed[:maxBaseLen]
		}
		username = fmt.Sprintf("%s_%s", trimmed, suffix)
	}
}

func (h *Handler) RegisterHandler(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid registration payload"})
		return
	}
	if len(req.Password) < minPasswordLength {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Password must be at least %d characters long", minPasswordLength)})
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	name := strings.TrimSpace(req.Name)
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Name is required"})
		return
	}

	gormDB := db.Gorm()

	first, last := splitName(name)
	if first == "" {
		first = name
	}
	if last == "" {
		last = first
	}

	username, err := generateUniqueUsername(gormDB, email, name)
	if err != nil {
		log.Println("register username error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to create account"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("register hash error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to create account"})
		return
	}

	user := auth.User{
		Username: username,
		Name:     first,
		Lastname: last,
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := h.ensureUserRecord(gormDB, &user, email); err != nil {
		if errors.Is(err, errEmailAlreadyRegistered) {
			c.JSON(http.StatusConflict, gin.H{"message": "Email already registered"})
			return
		}
		log.Println("register create error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to create account"})
		return
	}

	if err := h.syncMemberProfile(c.Request.Context(), user, req, first, last); err != nil {
		log.Println("register sync member error:", err)
		if delErr := gormDB.Delete(&auth.User{}, "id = ?", user.ID).Error; delErr != nil {
			log.Println("register cleanup error:", delErr)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to create account"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Account created successfully"})
}

func (h *Handler) LoginHandler(c *gin.Context) {
	if tok, err := c.Cookie("access_token"); err == nil && tok != "" {
		if _, err := auth.ParseAndValidate(tok); err == nil {
			exists, _ := h.Redis.Exists(ctx, sessionKey(tok)).Result()
			if exists == 1 {
				c.JSON(http.StatusOK, gin.H{"message": "Already logged in."})
				return
			}
		}
	}

	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid login payload"})
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	gormDB := db.Gorm()

	var user auth.User
	if err := gormDB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password"})
			return
		}
		log.Println("login lookup error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to sign in"})
		return
	}

	if user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Account uses Google sign-in. Continue with Google instead."})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password"})
		return
	}

	token, err := auth.GenerateAccessToken(user.ID.String(), user.Email, "user", defaultSessionTTL)
	if err != nil {
		log.Println("login jwt error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to sign in"})
		return
	}

	if err := h.Redis.Set(ctx, sessionKey(token), user.Email, defaultSessionTTL).Err(); err != nil {
		log.Println("login redis error:", err)
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(defaultSessionTTL.Seconds()),
	})

	c.JSON(http.StatusOK, gin.H{"message": "Signed in successfully"})
}

func (h *Handler) LoginWithGoogleHandler(c *gin.Context) {
	state := auth.RandomKey(20)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(oauthStateTTL.Seconds()),
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

	err = h.Redis.Set(ctx, sessionKey(token), gu.Email, oauthSessionTTL).Err()
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
		MaxAge:   int(oauthSessionTTL.Seconds()),
	})

	c.String(http.StatusOK, "✅ Logged in successfully.")
}


func (h *Handler) LogoutHandler(c *gin.Context) {
	if tok, err := c.Cookie("access_token"); err == nil && tok != "" {
		h.Redis.Del(ctx, sessionKey(tok))
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

func (h *Handler) ensureUserRecord(gormDB *gorm.DB, user *auth.User, email string) error {
	var existing auth.User
	if err := gormDB.Where("email = ?", email).First(&existing).Error; err == nil {
		return errEmailAlreadyRegistered
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if err := gormDB.Create(user).Error; err != nil {
		if isUniqueViolation(err) {
			return errEmailAlreadyRegistered
		}
		return err
	}
	return nil
}

func (h *Handler) syncMemberProfile(ctx context.Context, user auth.User, req registerRequest, first, last string) error {
	endpoint := h.memberEndpoint()
	if endpoint == "" {
		return nil
	}

	payload := memberProfilePayload{
		Username:  user.Username,
		FirstName: first,
		LastName:  last,
		Email:     user.Email,
		Bio:       strings.TrimSpace(req.Bio),
		LinkedIn:  strings.TrimSpace(req.LinkedIn),
		GitHub:    strings.TrimSpace(req.Github),
		Website:   strings.TrimSpace(req.Website),
	}

	payload.Skills = sanitizeSkills(req.Skills)
	if len(payload.Skills) == 0 {
		payload.Skills = nil
	}
	payload.Experience = mapExperience(req.Experiences)
	payload.Education = mapEducation(req.Educations)

	if payload.Bio == "" {
		payload.Bio = ""
	}
	if payload.LinkedIn == "" {
		payload.LinkedIn = ""
	}
	if payload.GitHub == "" {
		payload.GitHub = ""
	}
	if payload.Website == "" {
		payload.Website = ""
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	reqCtx := ctx
	var cancel context.CancelFunc
	reqCtx, cancel = context.WithTimeout(reqCtx, memberAPITimeout)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(reqCtx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := h.HTTPClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("app-service member create failed: status %d: %s", resp.StatusCode, strings.TrimSpace(string(b)))
	}
	return nil
}

func (h *Handler) memberEndpoint() string {
	base := "http://app-service:3001"
	if h.Config != nil && h.Config.AppServiceBaseURL != "" {
		base = h.Config.AppServiceBaseURL
	}
	base = strings.TrimRight(base, "/")
	if base == "" {
		return ""
	}
	return base + "/member/"
}

func sanitizeSkills(in []string) []string {
	if len(in) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, skill := range in {
		s := strings.TrimSpace(skill)
		if s == "" {
			continue
		}
		key := strings.ToLower(s)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, s)
	}
	return out
}

func mapExperience(in []experienceInput) []memberExperienceDTO {
	if len(in) == 0 {
		return nil
	}
	out := make([]memberExperienceDTO, 0, len(in))
	for _, exp := range in {
		title := strings.TrimSpace(exp.Title)
		company := strings.TrimSpace(exp.Company)
		description := strings.TrimSpace(exp.Description)
		if title == "" && company == "" && description == "" {
			continue
		}
		dto := memberExperienceDTO{
			Title:   title,
			Company: company,
		}
		if description != "" {
			dto.Description = description
		}
		if exp.StartYear != nil && *exp.StartYear > 0 {
			dto.StartYear = *exp.StartYear
		}
		if exp.EndYear != nil && *exp.EndYear > 0 {
			dto.EndYear = *exp.EndYear
		}
		out = append(out, dto)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func mapEducation(in []educationInput) []memberEducationDTO {
	if len(in) == 0 {
		return nil
	}
	out := make([]memberEducationDTO, 0, len(in))
	for _, edu := range in {
		school := strings.TrimSpace(edu.School)
		if school == "" {
			continue
		}
		dto := memberEducationDTO{
			School: school,
		}
	if trimmed := strings.TrimSpace(edu.Degree); trimmed != "" {
		dto.Degree = trimmed
	}
	if trimmed := strings.TrimSpace(edu.Field); trimmed != "" {
		dto.Field = trimmed
		}
		if edu.StartYear != nil && *edu.StartYear > 0 {
			dto.StartYear = *edu.StartYear
		}
		if edu.EndYear != nil && *edu.EndYear > 0 {
			dto.EndYear = *edu.EndYear
		}
		out = append(out, dto)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func isUniqueViolation(err error) bool {
	return strings.Contains(err.Error(), "duplicate key")
}
