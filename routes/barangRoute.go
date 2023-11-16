package routes

import (
	"golangsidang/controllers"

	"github.com/gin-gonic/gin"
)

func BarangGetAll(router *gin.Engine) {
	router.POST("/barang", controllers.CreateBarang())
	router.GET("/barangs/:id", controllers.GetBarang())
	router.PUT("/barang/:id", controllers.UpdateBarang())
	router.DELETE("/barang/:id", controllers.DeleteBarang())
	router.GET("/barangs", controllers.GetAllBarang())
}
