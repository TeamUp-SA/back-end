package router

import (
	"github.com/gin-gonic/gin"
)

func (r Router) AddMemberRouter(rg *gin.RouterGroup) {

	memberCont := r.deps.MemberController

	memberRouter := rg.Group("member")

	memberRouter.POST("/", memberCont.CreateMember)
	memberRouter.GET("/", memberCont.GetMembers)
	memberRouter.GET("/:member_id", memberCont.GetMemberByID)
	memberRouter.PUT("/:member_id", memberCont.UpdateMemberData)

}
