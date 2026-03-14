package handler

import (
	"net/http"
	"product-listing/internal/delivery/dto"
	"product-listing/internal/domain"
	"product-listing/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductImageHandler struct {
	usecase usecase.ProductImageUsecase
}

func NewProductImageHandler(u usecase.ProductImageUsecase) *ProductImageHandler {
	return &ProductImageHandler{usecase: u}
}

func (h *ProductImageHandler) AddImage(c *gin.Context) {
	var req dto.ProductImageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResp{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	puid, err := uuid.Parse(req.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResp{
			Status:  http.StatusBadRequest,
			Message: "invalid product_id",
		})
		return
	}

	input := domain.ProductImageInput{
		ProductID: puid,
		Url:       req.Url,
		IsPrimary: req.IsPrimary,
	}

	img, err := h.usecase.AddImage(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResp{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.Response{
		Status:  http.StatusCreated,
		Message: "Image added",
		Data:    dto.ToProductImageDTO(img),
	})
}

func (h *ProductImageHandler) GetProductImages(c *gin.Context) {
	productID := c.Param("product_id")
	images, err := h.usecase.GetProductImages(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResp{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	resp := make([]dto.ProductImageResp, 0, len(images))
	for _, img := range images {
		resp = append(resp, dto.ToProductImageDTO(&img))
	}

	c.JSON(http.StatusOK, dto.Response{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    resp,
	})
}

func (h *ProductImageHandler) DeleteImage(c *gin.Context) {
	id := c.Param("id")
	if err := h.usecase.DeleteImage(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResp{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResp{
		Status:  http.StatusOK,
		Message: "Image deleted",
	})
}

func (h *ProductImageHandler) SetPrimary(c *gin.Context) {
	productID := c.Param("product_id")
	imageID := c.Param("image_id")

	if err := h.usecase.SetPrimary(c.Request.Context(), productID, imageID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResp{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResp{
		Status:  http.StatusOK,
		Message: "Primary image set",
	})
}
