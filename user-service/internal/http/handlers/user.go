package handlers

import (
	"encoding/json"
	"net/http"

	"user-service/internal/cache"
	"user-service/internal/db"
	"user-service/internal/user"
	"user-service/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type updateProfileInput struct {
	Username    *string `json:"username"`
	Name        *string `json:"name"`
	Lastname    *string `json:"lastname"`
	PhoneNumber *string `json:"phone_number"`
	ImageUrl    *string `json:"image_url"` 
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
				"image_url":    u.ImageUrl,
				"created_at":   u.CreatedAt,
				"updated_at":   u.UpdatedAt,
				"source":       "cache",
			})
			return
		}
	}

	var u user.User
	if err := db.Gorm().Where("email = ?", email).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	jsonData, _ := json.Marshal(u)
	_ = cache.Rdb.Set(cache.Ctx, email, jsonData, 0).Err()

	c.JSON(http.StatusOK, gin.H{
		"id":           u.ID,
		"username":     u.Username,
		"name":         u.Name,
		"lastname":     u.Lastname,
		"phone_number": u.PhoneNumber,
		"email":        u.Email,
		"image_url":    u.ImageUrl,
		"created_at":   u.CreatedAt,
		"updated_at":   u.UpdatedAt,
		"source":       "database",
	})
}

func UpdateProfile(c *gin.Context) {
	email := c.GetString("email")
	if email == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var in updateProfileInput
	_ = c.ShouldBindJSON(&in)

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
	if in.PhoneNumber != nil {
		updates["phone_number"] = *in.PhoneNumber
	}

	fileHeader, err := c.FormFile("image")
	if err == nil && fileHeader != nil {
		// Open the file
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read uploaded file"})
			return
		}
		defer file.Close()
	
		s3URL, err := utils.UploadToS3(file, fileHeader)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload image"})
			return
		}
		updates["image_url"] = s3URL
	}
	

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}

	if err := db.Gorm().Model(&user.User{}).Where("email = ?", email).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}

	// Reload updated user
	var u user.User
	if err := db.Gorm().Where("email = ?", email).First(&u).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "reload failed"})
		return
	}

	jsonData, _ := json.Marshal(u)
	_ = cache.Rdb.Set(cache.Ctx, email, jsonData, 0).Err()

	c.JSON(http.StatusOK, gin.H{
		"id":           u.ID,
		"username":     u.Username,
		"name":         u.Name,
		"lastname":     u.Lastname,
		"phone_number": u.PhoneNumber,
		"email":        u.Email,
		"image_url":    u.ImageUrl,
		"created_at":   u.CreatedAt,
		"updated_at":   u.UpdatedAt,
	})
}

func startsWithHTTP(s string) bool {
	return len(s) > 4 && (s[:7] == "http://" || s[:8] == "https://")
}
