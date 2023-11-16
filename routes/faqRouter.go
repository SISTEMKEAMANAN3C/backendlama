package routes

import (
	controller "golangsidang/controllers"

	"github.com/gin-gonic/gin"
)

func FaqRoutes(faqRoutes *gin.Engine) { // membuat routes auth
	faqRoutes.POST("/faq", controller.CreateFaq()) // membuat routes signin untuk mengani sigin
	faqRoutes.GET("/faq/:id", controller.GetFaq())
	faqRoutes.PUT("faq/:id", controller.EditFaq())
	faqRoutes.DELETE("faq/:id", controller.DeleteFaq())
	faqRoutes.GET("/faqs", controller.GetAllFaq())
}
