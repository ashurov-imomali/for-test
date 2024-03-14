package router

import (
	"github.com/gin-gonic/gin"
	"main/handler"
)

func GetRouter(h *handler.Handler) *gin.Engine {
	router := gin.New()
	router.GET("/test", h.Test)
	router.POST("/commission-profile", h.CreateCommissionProfile)
	router.PUT("/commission-profile", h.UpdateCommissionProfile)
	router.PUT("/commission-rules", h.UpdateCommissionRules)
	return router
}
