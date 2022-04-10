package controller

import (
	"net/http"

	"github.com/bestetufan/beste-store/internal/domain/entity"
	"github.com/bestetufan/beste-store/internal/service"
	"github.com/bestetufan/beste-store/pkg/logger"
	"github.com/gin-gonic/gin"
)

type (
	Basket struct {
		storeService service.StoreService
		logger       logger.Logger
	}

	newBasketItemRequest struct {
		ProductId uint32 `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}

	removeBasketItemRequest struct {
		ProductId uint32 `json:"product_id"`
	}

	basketResponse struct {
		ID    string               `json:"id"`
		Items []*entity.BasketItem `json:"items"`
	}
)

func NewBasket(cs service.StoreService, l logger.Logger) *Basket {
	return &Basket{cs, l}
}

// getBasket godoc
// @Description  Returns user's basket.
// @Tags         Basket
// @Accept       json
// @Produce      json
// @Success 200 {object} basketResponse
// @Failure 400 {object} response
// @Failure 404 {object} response
// @Router /basket [get]
// @Security Bearer
func (c *Basket) GetBasket(g *gin.Context) {
	userName := g.GetString("Email")
	if len(userName) == 0 {
		errorResponse(g, http.StatusBadRequest, "unable to get parameters")
		return
	}

	basket := c.storeService.GetBasket(userName)
	if basket == nil {
		errorResponse(g, http.StatusNotFound, "no record found")
		return
	}

	g.JSON(http.StatusOK, basketResponse{ID: basket.ID, Items: basket.Items})
}

// addBasketItem godoc
// @Description  Adds an item to basket.
// @Tags         Basket
// @Accept       json
// @Produce      json
// @Param data body newBasketItemRequest true "New Basket Item Model"
// @Success 200 {object} response
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /basket [post]
// @Security Bearer
func (c *Basket) AddBasketItem(g *gin.Context) {
	userName := g.GetString("Email")
	if len(userName) == 0 {
		errorResponse(g, http.StatusBadRequest, "unable to get parameters")
		return
	}

	var req newBasketItemRequest
	if err := g.ShouldBind(&req); err != nil {
		c.logger.Error(err, "http - v1 - addBasketItem")
		errorResponse(g, http.StatusBadRequest, "invalid request body")
		return
	}

	err := c.storeService.AddItemToBasket(userName, req.ProductId, req.Quantity)
	if err != nil {
		errorResponse(g, http.StatusNotFound, err.Error())
		return
	}

	successResponse(g, http.StatusOK, "Operation completed successfully.")
}

// updateBasketItem godoc
// @Description  Updates an item in basket.
// @Tags         Basket
// @Accept       json
// @Produce      json
// @Param data body newBasketItemRequest true "New Basket Item Model"
// @Success 200 {object} response
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /basket [put]
// @Security Bearer
func (c *Basket) UpdateBasketItem(g *gin.Context) {
	userName := g.GetString("Email")
	if len(userName) == 0 {
		errorResponse(g, http.StatusBadRequest, "unable to get parameters")
		return
	}

	var req newBasketItemRequest
	if err := g.ShouldBind(&req); err != nil {
		c.logger.Error(err, "http - v1 - updateBasketItem")
		errorResponse(g, http.StatusBadRequest, "invalid request body")
		return
	}

	err := c.storeService.UpdateItemInBasket(userName, req.ProductId, req.Quantity)
	if err != nil {
		errorResponse(g, http.StatusNotFound, err.Error())
		return
	}

	successResponse(g, http.StatusOK, "Operation completed successfully.")
}

// removeBasketItem godoc
// @Description  Removes an item from basket.
// @Tags         Basket
// @Accept       json
// @Produce      json
// @Param data body removeBasketItemRequest true "Remove Basket Item Model"
// @Success 200 {object} response
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /basket [delete]
// @Security Bearer
func (c *Basket) RemoveBasketItem(g *gin.Context) {
	userName := g.GetString("Email")
	if len(userName) == 0 {
		errorResponse(g, http.StatusBadRequest, "unable to get parameters")
		return
	}

	var req removeBasketItemRequest
	if err := g.ShouldBind(&req); err != nil {
		c.logger.Error(err, "http - v1 - removeBasketItem")
		errorResponse(g, http.StatusBadRequest, "invalid request body")
		return
	}

	err := c.storeService.RemoveItemFromBasket(userName, req.ProductId)
	if err != nil {
		errorResponse(g, http.StatusNotFound, err.Error())
		return
	}

	successResponse(g, http.StatusOK, "Operation completed successfully.")
}
