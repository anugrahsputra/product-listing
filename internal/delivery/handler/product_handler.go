package handler

import (
	"net/http"
	"product-listing/internal/delivery/dto"
	"product-listing/internal/domain"
	"product-listing/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductHandler struct {
	usecase usecase.ProductUsecase
}

func NewProductHandler(u usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{usecase: u}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.ProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResp{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
		})
		return
	}

	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResp{
			Status:  http.StatusBadRequest,
			Message: "invalid category_id",
		})
		return
	}

	input := domain.ProductInput{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		CategoryID:  categoryID,
		Price:       req.Price,
	}

	if err := h.usecase.CreateProduct(ctx, input); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResp{
			Status:  http.StatusInternalServerError,
			Message: "failed to create product",
		})
		return
	}

	c.JSON(http.StatusCreated, dto.Response{
		Status:  http.StatusCreated,
		Message: "Product created",
	})
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	ctx := c.Request.Context()
	total, err := h.usecase.GetProductCount(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResp{
			Status:  http.StatusInternalServerError,
			Message: "Failed to get total counts",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	products, err := h.usecase.GetProducts(ctx, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResp{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	var productResp []dto.ProductResp
	for _, p := range products {
		productResp = append(productResp, toProductDTO(&p))
	}
	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Status:     http.StatusOK,
		Message:    "Success get categories",
		Data:       productResp,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: total / page,
	})
}

func (h *ProductHandler) GetProductById(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	product, err := h.usecase.GetProductsById(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResp{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	result := toProductDTO(product)
	c.JSON(http.StatusOK, dto.Response{
		Status:  http.StatusOK,
		Message: "Success get product",
		Data:    result,
	})
}

func (h *ProductHandler) GetProductByCategory(c *gin.Context) {
	ctx := c.Request.Context()
	categoryID := c.Param("category_id")

	products, err := h.usecase.GetProductsByCategory(ctx, categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResp{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	var productResp []dto.ProductResp
	for _, p := range products {
		productResp = append(productResp, toProductDTO(&p))
	}

	c.JSON(http.StatusOK, dto.Response{
		Status:  http.StatusOK,
		Message: "Success get products by category",
		Data:    productResp,
	})
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	var req dto.ProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResp{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
		})
		return
	}

	CategoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResp{
			Status:  http.StatusBadRequest,
			Message: "invalid category_id",
		})
		return
	}

	input := domain.ProductInput{
		Name:        req.Name,
		Description: req.Description,
		CategoryID:  CategoryID,
	}

	if err := h.usecase.UpdateProduct(ctx, id, input); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResp{
			Status:  http.StatusInternalServerError,
			Message: "failed to update product",
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResp{
		Status:  http.StatusOK,
		Message: "Success update product",
	})
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")
	if err := h.usecase.DeleteProduct(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResp{
			Status:  http.StatusInternalServerError,
			Message: "Failed delete product",
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResp{
		Status:  http.StatusOK,
		Message: "Success delete product",
	})
}

func toProductDTO(p *domain.Product) dto.ProductResp {
	return dto.ProductResp{
		ID:          p.ID.String(),
		Name:        p.Name,
		Slug:        p.Slug,
		Description: p.Description,
		CategoryID:  p.CategoryID.String(),
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
