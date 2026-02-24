package router

import (
	"product-listing/config"
	"product-listing/internal/delivery/handler"
	"product-listing/internal/repository"
	"product-listing/internal/usecase"

	"github.com/gin-gonic/gin"
)

func SetupRouter(db *config.Database) *gin.Engine {
	route := gin.Default()

	api := route.Group("/api")

	categoryRepo := repository.NewCategoryRepository(db)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryUsecase)
	CategoriesRoute(api, categoryHandler)

	productRepo := repository.NewProductRepository(db)
	productUsecase := usecase.NewProductUsecase(productRepo)
	productHandler := handler.NewProductHandler(productUsecase)
	ProductRoutes(api, productHandler)

	return route
}
