package controller

import (
	"errors"
	"net/http"
	"strings"

	"app-service/internal/dto"
	"app-service/internal/model"
	"app-service/internal/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IGroupController interface {
	CreateGroup(c *gin.Context)
	GetGroups(c *gin.Context)
	GetGroupByID(c *gin.Context)
	GetGroupsByOwnerID(c *gin.Context)
	UpdateGroup(c *gin.Context)
	DeleteGroup(c *gin.Context)
}

type GroupController struct {
	groupService service.IGroupService
}

func NewGroupController(s service.IGroupService) IGroupController {
	return GroupController{
		groupService: s,
	}
}

// CreateGroup godoc
//
//	@Summary		Create a new group
//	@Description	Creates a new group in the database
//	@Tags			group
//	@Accept			json
//	@Produce		json
//	@Param			group	body		model.Group	true	"Group to create"
//	@Success		201		{object}	dto.SuccessResponse{data=dto.Group}
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/group/ [post]
func (s GroupController) CreateGroup(c *gin.Context) {
	var newGroup model.Group

	if err := c.BindJSON(&newGroup); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Error:   "Invalid request body, failed to bind JSON",
			Message: err.Error(),
		})
		return
	}
	res, err := s.groupService.CreateGroup(&newGroup)

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
		Message: "Group created",
		Data:    res,
	})
}

// GetGroupByID godoc
//
//	@Summary		Get a group by ID
//	@Description	Retrieves a group's data by their ID
//	@Tags			group
//	@Accept			json
//	@Produce		json
//	@Param			group_id	path		string	true	"Group ID"
//	@Success		200			{object}	dto.SuccessResponse{data=dto.Group}
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/group/{group_id} [get]
func (s GroupController) GetGroupByID(c *gin.Context) {
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
	res, err := s.groupService.GetGroupByID(groupID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Error:   "No group with this groupID",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Get group success",
		Data:    res,
	})
}

// GetGroupsByOwnerID godoc
//
//	@Summary		Get groups by ownerID
//	@Description	Retrieves each owner's groups-on-display by owner ID
//	@Tags			group
//	@Accept			json
//	@Produce		json
//	@Param			owner_id	path		string	true	"Owner ID"
//	@Success		200			{object}	dto.SuccessResponse{data=[]dto.Group}
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/group/owner/{owner_id} [get]
func (s GroupController) GetGroupsByOwnerID(c *gin.Context) {
	ownerIDstr := c.Param("owner_id")
	ownerID, err := primitive.ObjectIDFromHex(ownerIDstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Error:   "Invalid ownerID format",
			Message: err.Error(),
		})
		return
	}
	res, err := s.groupService.GetGroupsByOwnerID(ownerID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Error:   "No group with this ownerID",
			Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Get group success",
		Data:    res,
	})
}

// GetGroups godoc
//
//	@Summary		Get all groups
//	@Description	Retrieves all groups
//	@Tags			group
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.SuccessResponse{data=[]dto.Group}
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/group/ [get]
func (s GroupController) GetGroups(c *gin.Context) {
	res, err := s.groupService.GetGroups()

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Error:   "No groups found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Get groups success",
		Data:    res,
	})
}

// UpdateGroup godoc
//
//	@Summary		Update a group by ID
//	@Description	Updates an existing group's data by their ID
//	@Tags			group
//	@Accept			json
//	@Produce		json
//	@Param			group_id	path		string			true	"Group ID"
//	@Param			updatedGroup		body		dto.GroupUpdateRequest	true	"Group data to update"
//	@Success		200			{object}	dto.SuccessResponse{data=dto.Group}
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/group/{group_id} [put]
func (s GroupController) UpdateGroup(c *gin.Context) {
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
	var updateReq dto.GroupUpdateRequest
	if err := c.BindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Error:   "Invalid request body, failed to bind JSON",
			Message: err.Error(),
		})
		return
	}

	res, err := s.groupService.UpdateGroup(groupID, &updateReq)

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Error:   "Failed to update group data",
			Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Update group success",
		Data:    res,
	})
}

// DeleteGroup godoc
//
//	@Summary		Delete a group by ID
//	@Description	Delete a group's data by its ID
//	@Tags			group
//	@Accept			json
//	@Produce		json
//	@Param			group_id	path		string	true	"Group ID"
//	@Success		200			{object}	dto.SuccessResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/group/{group_id} [delete]
func (s GroupController) DeleteGroup(c *gin.Context) {
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

	memberIDHeader := strings.TrimSpace(c.GetHeader("X-Member-ID"))
	if memberIDHeader == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Status:  http.StatusUnauthorized,
			Error:   "Unauthorized",
			Message: "missing member identification header",
		})
		return
	}

	requesterID, err := primitive.ObjectIDFromHex(memberIDHeader)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Error:   "Invalid requester ID",
			Message: err.Error(),
		})
		return
	}

	err = s.groupService.DeleteGroup(groupID, requesterID)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrGroupNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Success: false,
				Status:  http.StatusNotFound,
				Error:   "Group not found",
				Message: err.Error(),
			})
			return
		case errors.Is(err, service.ErrGroupForbidden):
			c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Success: false,
				Status:  http.StatusForbidden,
				Error:   "Forbidden",
				Message: "You do not have permission to delete this group",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Error:   "No group with this groupID",
			Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Delete group success",
	})
}

