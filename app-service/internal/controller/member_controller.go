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
	"app-service/pkg/utils/converter"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IMemberController interface {
	CreateMember(c *gin.Context)
	GetMembers(c *gin.Context)
	GetMemberByID(c *gin.Context)
	UpdateMemberData(c *gin.Context)
}

type MemberController struct {
	memberService service.IMemberService
}

func NewMemberController(s service.IMemberService) IMemberController {
	return MemberController{
		memberService: s,
	}
}

// CreateMember godoc
//
//	@Summary		Create a new member
//	@Description	Creates a new member in the database
//	@Tags			member
//	@Accept			json
//	@Produce		json
//	@Param			member	body		dto.MemberRegisterRequest	true	"Member to create"
//	@Success		201		{object}	dto.SuccessResponse{data=dto.Member}
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/member/ [post]
func (s MemberController) CreateMember(c *gin.Context) {
	var newMember model.Member

	if err := c.BindJSON(&newMember); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Error:   "Invalid request body, failed to bind JSON",
			Message: err.Error(),
		})
		return
	}
	res, err := s.memberService.CreateMemberData(&newMember)

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
		Message: "Member created",
		Data:    res,
	})
}

// GetMemberByID godoc
//
//	@Summary		Get a member by ID
//	@Description	Retrieves a member's data by their ID
//	@Tags			member
//	@Accept			json
//	@Produce		json
//	@Param			member_id	path		string	true	"Member ID"
//	@Success		200			{object}	dto.SuccessResponse{data=dto.Member}
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/member/{member_id} [get]
func (s MemberController) GetMemberByID(c *gin.Context) {
	memberIDstr := c.Param("member_id")
	memberID, err := primitive.ObjectIDFromHex(memberIDstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Error:   "Invalid memberID format",
			Message: err.Error(),
		})
		return
	}
	res, err := s.memberService.GetMemberByID(memberID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Error:   "No member with this memberID",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Get member success",
		Data:    res,
	})
}

// GetMember godoc
//
//	@Summary		Get all members
//	@Description	Retrieves all members
//	@Tags			member
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.SuccessResponse{data=[]dto.Member}
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/member/ [get]
func (s MemberController) GetMembers(c *gin.Context) {
	res, err := s.memberService.GetMember()

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Error:   "No members",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Get members success",
		Data:    res,
	})
}

// UpdateMember godoc
//
//	@Summary		Update a member by ID
//	@Description	Updates an existing member's data by their ID
//	@Tags			member
//	@Accept			json
//	@Produce		json
//	@Param			member_id	path		string			true	"Member ID"
//	@Param			updatedMember		body		dto.MemberUpdateRequest	true	"Member data to update"
//	@Success		200			{object}	dto.SuccessResponse{data=dto.Member}
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/member/{member_id} [put]
func (s MemberController) UpdateMemberData(c *gin.Context) {
	memberIDstr := c.Param("member_id")
	memberID, err := primitive.ObjectIDFromHex(memberIDstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Error:   "Invalid memberID format",
			Message: err.Error(),
		})
		return
	}
	var updateReq dto.MemberUpdateRequest
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

		payload := c.PostForm("member")
		if payload == "" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Success: false,
				Status:  http.StatusBadRequest,
				Error:   "Missing member payload",
				Message: "expected member payload in multipart form field \"member\"",
			})
			return
		}

		if err := json.Unmarshal([]byte(payload), &updateReq); err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Success: false,
				Status:  http.StatusBadRequest,
				Error:   "Invalid member payload",
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

	var experienceModels []model.Experience
	for _, e := range updateReq.Experience {
		expModel, err := converter.ExperienceDTOToModel(&e)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Success: false,
				Status:  http.StatusBadRequest,
				Error:   "Failed to convert experience",
				Message: err.Error(),
			})
			return
		}
		experienceModels = append(experienceModels, *expModel)
	}

	var educationModels []model.Education
	for _, e := range updateReq.Education {
		eduModel, err := converter.EducationDTOToModel(&e)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Success: false,
				Status:  http.StatusBadRequest,
				Error:   "Failed to convert education",
				Message: err.Error(),
			})
			return
		}
		educationModels = append(educationModels, *eduModel)
	}

	res, err := s.memberService.UpdateMemberData(memberID, &model.Member{
		MemberID:     memberID,
		Username:     updateReq.Username,
		Password:     updateReq.Password,
		FirstName:    updateReq.FirstName,
		LastName:     updateReq.LastName,
		Email:        updateReq.Email,
		PhoneNumber:  updateReq.PhoneNumber,
		Bio:          updateReq.Bio,
		Skills:       updateReq.Skills,
		LinkedIn:     updateReq.LinkedIn,
		GitHub:       updateReq.GitHub,
		Website:      updateReq.Website,
		ProfileImage: updateReq.ProfileImage,
		Experience:   experienceModels,
		Education:    educationModels,
	}, imageFile, fileHeader)

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Error:   "Failed to update member data",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Update member success",
		Data:    res,
	})
}
