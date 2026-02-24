package handler

import (
	"net/http"
	"product-listing/internal/delivery/dto"
	"product-listing/internal/domain"
	"product-listing/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	usecase usecase.CategoryUsecase
}

func NewCategoryHandler(u usecase.CategoryUsecase) *CategoryHandler {
	return &CategoryHandler{usecase: u}
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req dto.CategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResp{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	input := domain.CategoryInput{
		Name: req.Name,
		Slug: req.Slug,
	}

	ctx := c.Request.Context()
	if err := h.usecase.CreateCategory(ctx, input); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResp{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResp{
		Status:  http.StatusCreated,
		Message: "Category created",
	})
}

func (h *CategoryHandler) GetCategories(c *gin.Context) {
	total, err := h.usecase.GetCategoryCount(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResp{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	ctx := c.Request.Context()
	categories, err := h.usecase.GetCategories(ctx, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResp{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	var result []dto.CategoryResp
	for _, c := range categories {
		result = append(result, toCategoryDTO(&c))
	}

	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Status:     http.StatusOK,
		Message:    "Success get categories",
		Data:       result,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: total / page,
	})
}

func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	id := c.Param("id")

	ctx := c.Request.Context()
	category, err := h.usecase.GetCategoryById(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResp{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Status:  http.StatusOK,
		Message: "Success get category",
		Data:    toCategoryDTO(category),
	})
}

func (h *CategoryHandler) GetCategoryBySlug(c *gin.Context) {
	slug := c.Param("slug")

	ctx := c.Request.Context()
	category, err := h.usecase.GetCategoryBySlug(ctx, slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResp{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Status:  http.StatusOK,
		Message: "Success get category",
		Data:    toCategoryDTO(category),
	})
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var req dto.CategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResp{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	input := domain.CategoryInput{
		Name: req.Name,
		Slug: req.Slug,
	}

	ctx := c.Request.Context()
	if err := h.usecase.UpdateCategory(ctx, id, input); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResp{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResp{
		Status:  http.StatusOK,
		Message: "Category updated",
	})
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	ctx := c.Request.Context()
	if err := h.usecase.DeleteCategory(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResp{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResp{
		Status:  http.StatusOK,
		Message: "Category deleted",
	})
}

func toCategoryDTO(c *domain.Category) dto.CategoryResp {
	return dto.CategoryResp{
		ID:        c.ID.String(),
		Name:      c.Name,
		Slug:      c.Slug,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.CreatedAt,
	}
}
