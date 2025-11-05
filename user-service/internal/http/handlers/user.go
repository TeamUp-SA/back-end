package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"user-service/internal/cache"
	"user-service/internal/db"
	"user-service/internal/user"
)

type updateProfileInput struct {
    Username    *string `json:"username"`
    Name        *string `json:"name"`
    Lastname    *string `json:"lastname"`
    PhoneNumber *string `json:"phone_number"`
}

func GetMe(c *gin.Context) {
    email := c.GetString("email")
    if email == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    cached, err := cache.Rdb.Get(cache.Ctx, email).Result()
    if err == nil && cached != "" {
        var u user.User
        if err := json.Unmarshal([]byte(cached), &u); err == nil {
            c.JSON(http.StatusOK, gin.H{
                "id":           u.ID,
                "username":     u.Username,
                "name":         u.Name,
                "lastname":     u.Lastname,
                "phone_number": u.PhoneNumber,
                "email":        u.Email,
                "created_at":   u.CreatedAt,
                "updated_at":   u.UpdatedAt,
                "source":       "cache",
            })
            return
        }
    }

    // 2️⃣ Fallback to DB
    var u user.User
    if err := db.Gorm().Where("email = ?", email).First(&u).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
        return
    }

    // 3️⃣ Cache result for faster next time
    jsonData, _ := json.Marshal(u)
    _ = cache.Rdb.Set(cache.Ctx, email, jsonData, 0).Err()

    c.JSON(http.StatusOK, gin.H{
        "id":           u.ID,
        "username":     u.Username,
        "name":         u.Name,
        "lastname":     u.Lastname,
        "phone_number": u.PhoneNumber,
        "email":        u.Email,
        "created_at":   u.CreatedAt,
        "updated_at":   u.UpdatedAt,
        "source":       "database",
    })
}

// ✅ PUT /api/profile — update user info and refresh cache
func UpdateProfile(c *gin.Context) {
    email := c.GetString("email")
    if email == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    var in updateProfileInput
    if err := c.ShouldBindJSON(&in); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
        return
    }

    updates := map[string]interface{}{}
    if in.Username != nil && *in.Username != "" {
        updates["username"] = *in.Username
    }
    if in.Name != nil && *in.Name != "" {
        updates["name"] = *in.Name
    }
    if in.Lastname != nil && *in.Lastname != "" {
        updates["lastname"] = *in.Lastname
    }
    if in.PhoneNumber != nil { // allow empty string to clear
        updates["phone_number"] = *in.PhoneNumber
    }

    if len(updates) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
        return
    }

    // 1️⃣ Update user record
    if err := db.Gorm().Model(&user.User{}).Where("email = ?", email).Updates(updates).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
        return
    }

    // 2️⃣ Get updated record
    var u user.User
    if err := db.Gorm().Where("email = ?", email).First(&u).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "reload failed"})
        return
    }

    // 3️⃣ Refresh Redis cache
    jsonData, _ := json.Marshal(u)
    _ = cache.Rdb.Set(cache.Ctx, email, jsonData, 0).Err()

    c.JSON(http.StatusOK, gin.H{
        "id":           u.ID,
        "username":     u.Username,
        "name":         u.Name,
        "lastname":     u.Lastname,
        "phone_number": u.PhoneNumber,
        "email":        u.Email,
        "created_at":   u.CreatedAt,
        "updated_at":   u.UpdatedAt,
    })
}
