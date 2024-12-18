package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tonymj76/mytheresa-test/handlers"
)

func SetRouter(h *handlers.Handler) *gin.Engine {
	router := gin.Default()
	apiGroupRoute := router.Group("/api")
	apiGroupRoute.GET("/products", h.FetchProducts)
	apiGroupRoute.GET("/", h.Test)
	return router
}
