package routes

import (
	"golangsidang/controllers"

	"github.com/gin-gonic/gin"
)

func GetAlluser(router *gin.Engine) {

	router.POST("/users", controllers.CreateUser())
	router.GET("/users/:usersGetId", controllers.GetAuser())
	router.PUT("/users/:usersID", controllers.EditAbsensi())
	router.DELETE("/users/:usersID", controllers.DeleteAabsensi())
	router.GET("/userss", controllers.GetAllAbssenis())
}
