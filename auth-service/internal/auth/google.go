package auth

import (
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> d99d6e8 (chore: add todo)
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"os"
<<<<<<< HEAD

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"auth-service/internal/config"
	"auth-service/internal/userclient"
	userv1 "user-service/pb/userv1"
=======
    "context"
    "crypto/rand"
    "encoding/base64"
    "encoding/json"
    "log"
    "os"
=======
>>>>>>> d99d6e8 (chore: add todo)

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

<<<<<<< HEAD
    "auth-service/internal/config"
    "auth-service/internal/userclient"
    userv1 "user-service/pb/userv1"
>>>>>>> d809285 (feat: connect user-service and auth-service using gprc)
=======
	"auth-service/internal/config"
	"auth-service/internal/userclient"
	userv1 "user-service/pb/userv1"
>>>>>>> d99d6e8 (chore: add todo)
)

const (
	SessionKeyUser  = "user"
	SessionKeyState = "state"
	SessionName     = "appsession"
)

var (
	oauthConfig *oauth2.Config
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

func InitGoogleOAuth() {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL := os.Getenv("OAUTH_REDIRECT_URL")

	if clientID == "" || clientSecret == "" || redirectURL == "" {
		log.Fatal("Google OAuth environment variables not set")
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
}

func GetAuthURL(state string) string {
	return oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func ExchangeCodeForUser(code string) (*GoogleUser, error) {
	ctx := context.Background()
	token, err := oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	client := oauthConfig.Client(ctx, token)
	resp, err := client.Get("https://openidconnect.googleapis.com/v1/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func RandomKey(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalf("failed to generate random key: %v", err)
	}
	return base64.URLEncoding.EncodeToString(b)
}

func UpsertUserInfo(gu *GoogleUser) error {
	// Prefer calling user-service via gRPC instead of DB write here
	cfg := config.Load()
	cli, err := userclient.New(cfg)
	if err != nil {
		log.Println("userclient dial:", err)
		return err
	}
	defer cli.Close()

	//TODO: check if user already have info in db to avoid overwrite updated info
	_, err = cli.UpsertUser(context.Background(), &userv1.UpsertUserRequest{
		Email:           gu.Email,
		Name:            gu.Name,
		Lastname:        gu.FamilyName,
		Username:        gu.Email,
		OauthProvider:   "google",
		OauthProviderId: gu.Sub,
		PhoneNumber:     "",
	})
	if err != nil {
		log.Println("grpc upsert user:", err)
		return err
	}
	return nil
}
