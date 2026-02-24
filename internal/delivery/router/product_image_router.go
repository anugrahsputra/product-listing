package router

import (
	"product-listing/internal/delivery/handler"

	"github.com/gin-gonic/gin"
)

func ProductImageRoutes(r *gin.RouterGroup, h *handler.ProductImageHandler) {
	route := r.Group("/product-images")
	{
		route.POST("", h.AddImage)
		route.GET("/product/:product_id", h.GetProductImages)
		route.DELETE("/:id", h.DeleteImage)
		route.PUT("/primary/:product_id/:image_id", h.SetPrimary)
	}
}
