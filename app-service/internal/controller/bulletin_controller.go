package controller

import (
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"strings"

	"app-service/internal/dto"
	"app-service/internal/model"
	"app-service/internal/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IBulletinController interface {
	GetBulletins(c *gin.Context)
	GetBulletinByID(c *gin.Context)
	GetBulletinsByAuthorID(c *gin.Context)
	GetBulletinsByGroupID(c *gin.Context)
	CreateBulletin(c *gin.Context)
	UpdateBulletin(c *gin.Context)
	DeleteBulletin(c *gin.Context)
}

type BulletinController struct {
	bulletinService service.IBulletinService
}

func NewBulletinController(s service.IBulletinService) IBulletinController {
	return BulletinController{
		bulletinService: s,
	}
}

// GetBulletins godoc
//
//	@Summary		Get all bulletins
//	@Description	Retrieves all bulletins
//	@Tags			bulletin
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.SuccessResponse{data=[]dto.Bulletin}
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/bulletin/ [get]
func (s BulletinController) GetBulletins(c *gin.Context) {
	res, err := s.bulletinService.GetBulletins()

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Error:   "No bulletins",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Get Bulletins success",
		Data:    res,
	})
}

// GetBulletinsByID godoc
//
//	@Summary		Get a bulletin by ID
//	@Description	Retrieves a bulletin's data by its ID
//	@Tags			bulletin
//	@Accept			json
//	@Produce		json
//	@Param			bulletin_id	path		string	true	"Bulletin ID"
//	@Success		200			{object}	dto.SuccessResponse{data=dto.Bulletin}
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/bulletin/{bulletin_id} [get]
func (s BulletinController) GetBulletinByID(c *gin.Context) {
	bulletinIDstr := c.Param("bulletin_id")
	bulletinID, err := primitive.ObjectIDFromHex(bulletinIDstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Error:   "Invalid bulletinID format",
			Message: err.Error(),
		})
		return
	}
	res, err := s.bulletinService.GetBulletinByID(bulletinID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Error:   "No bulletin with this bulletinID",
			Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Get bulletin success",
		Data:    res,
	})
}

// GetBulletinsByAuthorID godoc
//
//	@Summary		Get bulletins by authorID
//	@Description	Retrieves each author's bulletins by author ID
//	@Tags			bulletin
//	@Accept			json
//	@Produce		json
//	@Param			author_id	path		string	true	"Author ID"
//	@Success		200			{object}	dto.SuccessResponse{data=[]dto.Bulletin}
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/bulletin/author/{author_id} [get]
func (s BulletinController) GetBulletinsByAuthorID(c *gin.Context) {
	authorIDstr := c.Param("author_id")
	authorID, err := primitive.ObjectIDFromHex(authorIDstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Error:   "Invalid authorID format",
			Message: err.Error(),
		})
		return
	}
	res, err := s.bulletinService.GetBulletinsByAuthorID(authorID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Error:   "No bulletin with this authorID",
			Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Get bulletin success",
		Data:    res,
	})
}

// GetBulletinsByGroupID godoc
//
//	@Summary		Get bulletins by groupID
//	@Description	Retrieves each group's bulletins by group ID
//	@Tags			bulletin
//	@Accept			json
//	@Produce		json
//	@Param			group_id	path		string	true	"Group ID"
//	@Success		200			{object}	dto.SuccessResponse{data=[]dto.Bulletin}
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/bulletin/group/{group_id} [get]
func (s BulletinController) GetBulletinsByGroupID(c *gin.Context) {
	groupIDstr := c.Param("group_id")
	groupID, err := primitive.ObjectIDFromHex(groupIDstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Error:   "Invalid groupID format",
			Message: err.Error(),
		})
		return
	}
	res, err := s.bulletinService.GetBulletinsByGroupID(groupID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Error:   "No bulletin with this groupID",
			Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Get bulletin success",
		Data:    res,
	})
}

// CreateBulletin godoc

// CreateBulletin godoc
//
//	@Summary		Create a new bulletin
//	@Description	Creates a new bulletin in the database
//	@Tags			bulletin
//	@Accept			json
//	@Produce		json
//	@Param			bulletin	body		dto.BulletinCreateRequest	true	"Bulletin to create"
//	@Success		201		{object}	dto.SuccessResponse{data=dto.Bulletin}
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/bulletin/ [post]
func (s BulletinController) CreateBulletin(c *gin.Context) {
	var newBulletin model.Bulletin
	var imageFile multipart.File
	var fileHeader *multipart.FileHeader

	contentType := c.GetHeader("Content-Type")

	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Success: false,
				Status:  http.StatusBadRequest,
				Error:   "Invalid multipart form data",
				Message: err.Error(),
			})
			return
		}

		payload := c.PostForm("bulletin")
		if payload == "" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Success: false,
				Status:  http.StatusBadRequest,
				Error:   "Missing bulletin payload",
				Message: "expected bulletin payload in multipart form field \"bulletin\"",
			})
			return
		}

		if err := json.Unmarshal([]byte(payload), &newBulletin); err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Success: false,
				Status:  http.StatusBadRequest,
				Error:   "Invalid bulletin payload",
				Message: err.Error(),
			})
			return
		}

		formHeader, err := c.FormFile("image")
		if err != nil {
			if !errors.Is(err, http.ErrMissingFile) {
				c.JSON(http.StatusBadRequest, dto.ErrorResponse{
					Success: false,
					Status:  http.StatusBadRequest,
					Error:   "Invalid image upload",
					Message: err.Error(),
				})
				return
			}
		} else {
			file, err := formHeader.Open()
			if err != nil {
				c.JSON(http.StatusBadRequest, dto.ErrorResponse{
					Success: false,
					Status:  http.StatusBadRequest,
					Error:   "Unable to process image upload",
					Message: err.Error(),
				})
				return
			}
			defer file.Close()
			imageFile = file
			fileHeader = formHeader
		}
	} else {
		if err := c.ShouldBindJSON(&newBulletin); err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Success: false,
				Status:  http.StatusBadRequest,
				Error:   "Invalid request body, failed to bind JSON",
				Message: err.Error(),
			})
			return
		}
	}

	res, err := s.bulletinService.CreateBulletin(&newBulletin, imageFile, fileHeader)

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Error:   "Failed to insert to database",
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, dto.SuccessResponse{
		Success: true,
		Status:  http.StatusCreated,
		Message: "Bulletin created",
		Data:    res,
	})
}

// UpdateBulletin godoc
//
//	@Summary		Update a bulletin by ID
//	@Description	Updates an existing bulletin's data by its ID
//	@Tags			bulletin
//	@Accept			json
//	@Produce		json
//	@Param			bulletin_id	path		string					true	"Bulletin ID"
//	@Param			bulletin		body		dto.BulletinUpdateRequest	true	"Bulletin data to update"
//	@Success		200			{object}	dto.SuccessResponse{data=dto.Bulletin}
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/bulletin/{bulletin_id} [put]
func (s BulletinController) UpdateBulletin(c *gin.Context) {
	bulletinIDstr := c.Param("bulletin_id")
	bulletinID, err := primitive.ObjectIDFromHex(bulletinIDstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Error:   "Invalid bulletinID format",
			Message: err.Error(),
		})
		return
	}
	var updateReq dto.BulletinUpdateRequest
	var imageFile multipart.File
	var fileHeader *multipart.FileHeader

	contentType := c.GetHeader("Content-Type")

	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Success: false,
				Status:  http.StatusBadRequest,
				Error:   "Invalid multipart form data",
				Message: err.Error(),
			})
			return
		}

		payload := c.PostForm("bulletin")
		if payload == "" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Success: false,
				Status:  http.StatusBadRequest,
				Error:   "Missing bulletin payload",
				Message: "expected bulletin payload in multipart form field \"bulletin\"",
			})
			return
		}

		if err := json.Unmarshal([]byte(payload), &updateReq); err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Success: false,
				Status:  http.StatusBadRequest,
				Error:   "Invalid bulletin payload",
				Message: err.Error(),
			})
			return
		}

		formHeader, err := c.FormFile("image")
		if err != nil {
			if !errors.Is(err, http.ErrMissingFile) {
				c.JSON(http.StatusBadRequest, dto.ErrorResponse{
					Success: false,
					Status:  http.StatusBadRequest,
					Error:   "Invalid image upload",
					Message: err.Error(),
				})
				return
			}
		} else {
			file, err := formHeader.Open()
			if err != nil {
				c.JSON(http.StatusBadRequest, dto.ErrorResponse{
					Success: false,
					Status:  http.StatusBadRequest,
					Error:   "Unable to process image upload",
					Message: err.Error(),
				})
				return
			}
			defer file.Close()
			imageFile = file
			fileHeader = formHeader
		}
	} else {
		if err := c.ShouldBindJSON(&updateReq); err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Success: false,
				Status:  http.StatusBadRequest,
				Error:   "Invalid request body, failed to bind JSON",
				Message: err.Error(),
			})
			return
		}
	}

	res, err := s.bulletinService.UpdateBulletin(bulletinID, &updateReq, imageFile, fileHeader)

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Error:   "Failed to update bulletin data",
			Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Update bulletin success",
		Data:    res,
	})
}

// DeleteBulletin godoc
//
//	@Summary		Delete a bulletin by ID
//	@Description	Delete a bulletin's data by its ID
//	@Tags			bulletin
//	@Accept			json
//	@Produce		json
//	@Param			bulletin_id	path		string	true	"Bulletin ID"
//	@Success		200			{object}	dto.SuccessResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/bulletin/{bulletin_id} [delete]
func (s BulletinController) DeleteBulletin(c *gin.Context) {
	bulletinIDstr := c.Param("bulletin_id")
	bulletinID, err := primitive.ObjectIDFromHex(bulletinIDstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Error:   "Invalid bulletinID format",
			Message: err.Error(),
		})
		return
	}
	err = s.bulletinService.DeleteBulletin(bulletinID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Error:   "No bulletin with this bulletinID",
			Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Delete bulletin success",
	})
}
