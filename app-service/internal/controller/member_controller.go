package controller

import (
	"net/http"

	"github.com/Ntchah/TeamUp-application-service/internal/dto"
	"github.com/Ntchah/TeamUp-application-service/internal/model"
	"github.com/Ntchah/TeamUp-application-service/internal/service"
	"github.com/Ntchah/TeamUp-application-service/pkg/utils/converter"
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
	var updatedMember dto.MemberUpdateRequest
	if err := c.BindJSON(&updatedMember); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Error:   "Invalid request body, failed to bind JSON",
			Message: err.Error(),
		})
		return
	}
	var groupMemberships []primitive.ObjectID
	for _, m := range updatedMember.GroupMemberships {
		id, err := primitive.ObjectIDFromHex(m)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Success: false,
				Status:  http.StatusInternalServerError,
				Error:   "Invalid groupID format",
				Message: err.Error(),
			})
			return
		}
		groupMemberships = append(groupMemberships, id)
	}
	var experienceModels []model.Experience
	for _, e := range updatedMember.Experience {
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
	for _, e := range updatedMember.Education {
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
		MemberID:         memberID,
		Username:         updatedMember.Username,
		Password:         updatedMember.Password,
		FirstName:        updatedMember.FirstName,
		LastName:         updatedMember.LastName,
		Email:            updatedMember.Email,
		PhoneNumber:      updatedMember.PhoneNumber,
		Bio:              updatedMember.Bio,
		Skills:           updatedMember.Skills,
		LinkedIn:         updatedMember.LinkedIn,
		GitHub:           updatedMember.GitHub,
		Website:          updatedMember.Website,
		ProfileImage:     updatedMember.ProfileImage,
		Experience:       experienceModels,
		Education:        educationModels,
		GroupMemberships: groupMemberships,
	})

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
