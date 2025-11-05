package router

import (
	"github.com/gin-gonic/gin"
)

func (r Router) AddBulletinRouter(rg *gin.RouterGroup) {

	bulletinCont := r.deps.BulletinController
	bulletinRouter := rg.Group("bulletin")

	bulletinRouter.POST("/", bulletinCont.CreateBulletin)
	bulletinRouter.GET("/", bulletinCont.GetBulletins)
	bulletinRouter.GET("/:bulletin_id", bulletinCont.GetBulletinByID)
	bulletinRouter.GET("/author/:author_id", bulletinCont.GetBulletinsByAuthorID)
	bulletinRouter.GET("/group/:group_id", bulletinCont.GetBulletinsByGroupID)
	bulletinRouter.PUT("/:bulletin_id", bulletinCont.UpdateBulletin)
	bulletinRouter.DELETE("/:bulletin_id", bulletinCont.DeleteBulletin)

}
