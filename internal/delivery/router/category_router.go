package router

import (
	"product-listing/internal/delivery/handler"

	"github.com/gin-gonic/gin"
)

func CategoriesRoute(r *gin.RouterGroup, h *handler.CategoryHandler) {
	router := r.Group("/category")
	{
		router.GET("", h.GetCategories)
		router.GET("/:id", h.GetCategoryByID)
		router.GET("/slug/:slug", h.GetCategoryBySlug)
		router.POST("", h.CreateCategory)
		router.PUT("/:id", h.UpdateCategory)
		router.DELETE("/:id", h.DeleteCategory)
	}
}
