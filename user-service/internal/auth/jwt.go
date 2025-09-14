package auth

import (
    "errors"
    "os"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID string `json:"uid"`
    Email  string `json:"email"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

func jwtSecret() []byte {
    s := os.Getenv("JWT_SECRET")
    if s == "" {
        s = os.Getenv("SESSION_SECRET")
    }
    return []byte(s)
}

func GenerateAccessToken(userID, email, role string, ttl time.Duration) (string, error) {
    now := time.Now()
    claims := Claims{
        UserID: userID,
        Email:  email,
        Role:   role,
        RegisteredClaims: jwt.RegisteredClaims{
            IssuedAt:  jwt.NewNumericDate(now),
            ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
            NotBefore: jwt.NewNumericDate(now),
            Issuer:    "user-service",
            Audience:  []string{"topup-app"},
        },
    }
    tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return tok.SignedString(jwtSecret())
}

func ParseAndValidate(tokenStr string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return jwtSecret(), nil
    })
    if err != nil {
        return nil, err
    }
    if c, ok := token.Claims.(*Claims); ok && token.Valid {
        return c, nil
    }
    return nil, errors.New("invalid token")
}

