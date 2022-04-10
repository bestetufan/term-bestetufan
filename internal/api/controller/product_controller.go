package controller

import (
	"net/http"
	"strconv"

	"github.com/bestetufan/beste-store/internal/domain/entity"
	"github.com/bestetufan/beste-store/internal/service"
	"github.com/bestetufan/beste-store/pkg/logger"
	"github.com/bestetufan/beste-store/pkg/pagination"
	"github.com/gin-gonic/gin"
)

type (
	Product struct {
		storeService service.StoreService
		logger       logger.Logger
	}

	createProductRequest struct {
		Name       string  `json:"name" binding:"required"`
		Sku        string  `json:"sku" binding:"required"`
		UnitPrice  float64 `json:"unit_price" binding:"required"`
		Quantity   int     `json:"quantity" binding:"required"`
		CategoryID uint32  `json:"category_id" binding:"required"`
	}

	updateProductRequest struct {
		Name      string  `json:"name" binding:"required"`
		UnitPrice float64 `json:"unit_price" binding:"required"`
		Quantity  int     `json:"quantity" binding:"required"`
	}

	productResponse struct {
		ID           uint32  `json:"id"`
		Name         string  `json:"name"`
		Sku          string  `json:"sku"`
		UnitPrice    float64 `json:"unit_price"`
		Quantity     int     `json:"quantity"`
		CategoryID   uint32  `json:"category_id"`
		CategoryName string  `json:"category_name"`
	}

	searchResponse struct {
		Items []entity.Product `json:"items"`
		Count int              `json:"count"`
	}
)

func NewProduct(cs service.StoreService, l logger.Logger) *Product {
	return &Product{cs, l}
}

// getAllProducts godoc
// @Description  Returns all products with pagination.
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param page query int false "Page Index"
// @Param pageSize query int false "Page Size"
// @Success 200 {object} pagination.Pages
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /product [get]
// @Security Bearer
func (c *Product) GetAllProducts(g *gin.Context) {
	pageIndex, pageSize := pagination.GetPaginationParametersFromRequest(g)
	items, count := c.storeService.GetAllProducts(pageIndex, pageSize)
	paginatedResult := pagination.NewFromGinRequest(g, count)
	paginatedResult.Items = items

	g.JSON(http.StatusOK, paginatedResult)
}

// getProduct godoc
// @Description  Returns one product by id.
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        id path int true "Product ID"
// @Success 200 {object} productResponse
// @Failure 400 {object} response
// @Failure 404 {object} response
// @Router /product/{id} [get]
// @Security Bearer
func (c *Product) GetProduct(g *gin.Context) {
	id, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		errorResponse(g, http.StatusBadRequest, "unable to get parameters")
		return
	}

	product := c.storeService.GetProduct(uint32(id))
	if product == nil {
		errorResponse(g, http.StatusNotFound, "no record found")
		return
	}

	g.JSON(http.StatusOK, productResponse{
		ID: product.ID, Name: product.Name, Sku: product.Sku,
		UnitPrice: product.UnitPrice, Quantity: product.Quantity,
		CategoryName: product.Category.Name, CategoryID: product.Category.ID})
}

// searchProducts godoc
// @Description  Returns searched products.
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        query path string true "Search Query"
// @Success 200 {object} searchResponse
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /product/search/{query} [get]
// @Security Bearer
func (c *Product) SearchProducts(g *gin.Context) {
	searchQuery := g.Param("query")
	if len(searchQuery) == 0 {
		errorResponse(g, http.StatusBadRequest, "unable to get query")
		return
	}

	products := c.storeService.SearchProducts(searchQuery)
	g.JSON(http.StatusOK, searchResponse{Items: products, Count: len(products)})
}

// createProduct godoc
// @Description  Creates a new product.
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param 		 data body createProductRequest true "Create Product Model"
// @Success 200 {object} productResponse
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /product [post]
// @Security Bearer
func (c *Product) CreateProduct(g *gin.Context) {
	var req createProductRequest
	if err := g.ShouldBind(&req); err != nil {
		c.logger.Error(err, "http - v1 - createProduct")
		errorResponse(g, http.StatusBadRequest, "invalid request body")
		return
	}

	product := entity.NewProduct(req.Name, req.Sku, req.UnitPrice, req.Quantity, req.CategoryID)
	err := c.storeService.CreateProduct(product)
	if err != nil {
		c.logger.Error(err, "http - v1 - createProduct")
		errorResponse(g, http.StatusInternalServerError, err.Error())
		return
	}

	g.JSON(http.StatusOK, productResponse{
		ID: product.ID, Name: product.Name, Sku: product.Sku,
		UnitPrice: product.UnitPrice, Quantity: product.Quantity,
		CategoryName: product.Category.Name, CategoryID: product.Category.ID})
}

// updateProduct godoc
// @Description  Updates a product.
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        id path int true "Product ID"
// @Param 		 data body updateProductRequest true "Update Product Model"
// @Success 200 {object} productResponse
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /product/{id} [put]
// @Security Bearer
func (c *Product) UpdateProduct(g *gin.Context) {
	id, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		errorResponse(g, http.StatusBadRequest, "unable to get id")
		return
	}

	var req updateProductRequest
	if err := g.ShouldBind(&req); err != nil {
		c.logger.Error(err, "http - v1 - updateProduct")
		errorResponse(g, http.StatusBadRequest, "invalid request body")
		return
	}

	product := c.storeService.GetProduct(uint32(id))
	product.Name = req.Name
	product.UnitPrice = req.UnitPrice
	product.Quantity = req.Quantity
	err = c.storeService.UpdateProduct(product)
	if err != nil {
		c.logger.Error(err, "http - v1 - updateProduct")
		errorResponse(g, http.StatusInternalServerError, err.Error())
		return
	}

	g.JSON(http.StatusOK, productResponse{
		ID: product.ID, Name: product.Name, Sku: product.Sku,
		UnitPrice: product.UnitPrice, Quantity: product.Quantity,
		CategoryName: product.Category.Name, CategoryID: product.Category.ID})
}

// deleteProduct godoc
// @Description  Deletes a product.
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        id path int true "Product ID"
// @Success 200 {object} response
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /product/{id} [delete]
// @Security Bearer
func (c *Product) DeleteProduct(g *gin.Context) {
	id, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		errorResponse(g, http.StatusBadRequest, "unable to get id")
		return
	}

	err = c.storeService.DeleteProduct(uint32(id))
	if err != nil {
		c.logger.Error(err, "http - v1 - deleteProduct")
		errorResponse(g, http.StatusInternalServerError, err.Error())
		return
	}

	successResponse(g, http.StatusOK, "Operation completed successfully.")
}
