package routes

import (
	"golangsidang/controllers"

	"github.com/gin-gonic/gin"
)

func GetAllRuangan(router *gin.Engine) {
	router.POST("/ruangan", controllers.CreateRuang())
	router.GET("/ruangans/:id", controllers.GetRuangan())
	router.PUT("/ruangan/:id", controllers.EditRuangan())
	router.DELETE("/ruangan/:id", controllers.DeleteRuangan())
	router.GET("/ruangans", controllers.GetAllRuangan())
}
