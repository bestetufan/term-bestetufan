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
	Category struct {
		storeService service.StoreService
		logger       logger.Logger
	}

	createCategoryRequest struct {
		Name string `json:"name" binding:"required"`
	}

	categoryResponse struct {
		ID   uint32 `json:"id"`
		Name string `json:"name"`
	}

	createBulkCategoryResponse struct {
		Added    int `json:"added_count"`
		Existing int `json:"existing_count"`
	}
)

func NewCategory(cs service.StoreService, l logger.Logger) *Category {
	return &Category{cs, l}
}

// getAllCategories godoc
// @Description  Returns all categories with pagination.
// @Tags         Category
// @Accept       json
// @Produce      json
// @Param page query int false "Page Index"
// @Param pageSize query int false "Page Size"
// @Success 200 {object} pagination.Pages
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /category [get]
// @Security Bearer
func (c *Category) GetAllCategories(g *gin.Context) {
	pageIndex, pageSize := pagination.GetPaginationParametersFromRequest(g)
	items, count := c.storeService.GetAllCategories(pageIndex, pageSize, true)
	paginatedResult := pagination.NewFromGinRequest(g, count)
	paginatedResult.Items = items

	g.JSON(http.StatusOK, paginatedResult)
}

// getCategory godoc
// @Description  Returns one category by id.
// @Tags         Category
// @Accept       json
// @Produce      json
// @Param        id path int true "Category ID"
// @Success 200 {object} categoryResponse
// @Failure 400 {object} response
// @Failure 404 {object} response
// @Router /category/{id} [get]
// @Security Bearer
func (c *Category) GetCategory(g *gin.Context) {
	id, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		errorResponse(g, http.StatusBadRequest, "unable to get parameters")
		return
	}

	category := c.storeService.GetCategory(uint32(id))
	if category == nil || !category.IsActive {
		errorResponse(g, http.StatusNotFound, "no record found")
		return
	}

	g.JSON(http.StatusOK, categoryResponse{ID: category.ID, Name: category.Name})
}

// createCategory godoc
// @Description  Creates a new category.
// @Tags         Category
// @Accept       json
// @Produce      json
// @Param data body createCategoryRequest true "Create Category Model"
// @Success 200 {object} categoryResponse
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /category [post]
// @Security Bearer
func (c *Category) CreateCategory(g *gin.Context) {
	var req createCategoryRequest
	if err := g.ShouldBind(&req); err != nil {
		c.logger.Error(err, "http - v1 - createCategory")
		errorResponse(g, http.StatusBadRequest, "invalid request body")
		return
	}

	category := entity.NewCategory(req.Name, true)
	err := c.storeService.CreateCategory(category)
	if err != nil {
		c.logger.Error(err, "http - v1 - createCategory")
		errorResponse(g, http.StatusInternalServerError, err.Error())
		return
	}

	g.JSON(http.StatusOK, categoryResponse{ID: category.ID, Name: category.Name})
}

// createBulkCategory godoc
// @Description  Creates new categories in bulk.
// @Tags         Category
// @Accept       json
// @Produce      json
// @Param   	 file formData file true "Category CSV"
// @Success 200 {object} nil
// @Failure 400 {object} nil
// @Failure 500 {object} nil
// @Router /category/bulk [post]
// @Security Bearer
func (c *Category) CreateBulkCategory(g *gin.Context) {

	file, fileHead, err := g.Request.FormFile("file")
	if err != nil {
		c.logger.Error(err, "http - v1 - createBulkCategory")
		errorResponse(g, http.StatusBadRequest, err.Error())

		return
	}

	contentType := fileHead.Header.Values("Content-Type")
	if contentType[0] != "text/csv" {
		errorResponse(g, http.StatusBadRequest, "invalid file type")
	}

	added, existing, err := c.storeService.CreateBulkCategory(file)
	if err != nil {
		errorResponse(g, http.StatusInternalServerError, err.Error())
	}

	g.JSON(http.StatusOK, createBulkCategoryResponse{Added: added, Existing: existing})
}

// deleteCategory godoc
// @Description  Deletes a category.
// @Tags         Category
// @Accept       json
// @Produce      json
// @Param        id path int true "Category ID"
// @Success 200 {object} nil
// @Failure 400 {object} nil
// @Failure 500 {object} nil
// @Router /category [delete]
// @Security Bearer
func (c *Category) DeleteCategory(g *gin.Context) {
	c.logger.Fatal("unimplemented exception")
}
