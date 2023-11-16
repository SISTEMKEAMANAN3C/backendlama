package routes

import (
	"golangsidang/controllers"

	"github.com/gin-gonic/gin"
)

func GetAllProduct(router *gin.Engine) {

	router.POST("/products", controllers.CreateProduct())
	router.GET("/products/:productGetId", controllers.GetProduct())
	router.PUT("/products/:productID", controllers.EditProduct())
	router.DELETE("/products/:productID", controllers.DeleteProduct())
	router.GET("/productss", controllers.GetAllProduct())
}
