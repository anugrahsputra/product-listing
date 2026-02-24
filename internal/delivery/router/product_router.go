package router

import (
	"product-listing/internal/delivery/handler"

	"github.com/gin-gonic/gin"
)

func ProductRoutes(r *gin.RouterGroup, h *handler.ProductHandler) {
	route := r.Group("/products")
	{
		route.GET("/", h.GetProducts)
		route.GET("/:id", h.GetProductById)
		route.GET("/category/:category_id", h.GetProductByCategory)
		route.POST("/", h.CreateProduct)
		route.PUT("/:id", h.UpdateProduct)
		route.DELETE("/:id", h.DeleteProduct)
	}

}
