package router

import (
	"github.com/gin-gonic/gin"
)

func (r Router) AddGroupRouter(rg *gin.RouterGroup) {

	groupCont := r.deps.GroupController
	groupRouter := rg.Group("group")

	groupRouter.POST("/", groupCont.CreateGroup)
	groupRouter.GET("/", groupCont.GetGroups)
	groupRouter.GET("/:group_id", groupCont.GetGroupByID)
	groupRouter.GET("/owner/:owner_id", groupCont.GetGroupsByOwnerID)
	groupRouter.PUT("/:group_id", groupCont.UpdateGroup)
	groupRouter.DELETE("/:group_id", groupCont.DeleteGroup)
}
