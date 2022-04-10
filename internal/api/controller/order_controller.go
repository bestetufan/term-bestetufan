package controller

import (
	"net/http"

	"github.com/bestetufan/beste-store/internal/domain/entity"
	"github.com/bestetufan/beste-store/internal/service"
	"github.com/bestetufan/beste-store/pkg/logger"
	"github.com/gin-gonic/gin"
)

type (
	Order struct {
		storeService service.StoreService
		logger       logger.Logger
	}

	newOrderRequest struct {
		Name        string `json:"name" binding:"required"`
		Address     string `json:"address" binding:"required"`
		PhoneNumber string `json:"phone_number" binding:"required"`
		CardNumber  string `json:"card_number" binding:"required"`
		CardExp     string `json:"card_exp" binding:"required"`
		CardCVV     int    `json:"card_cvv" binding:"required"`
	}

	orderResponse struct {
		ID    string              `json:"id"`
		Items []*entity.OrderItem `json:"items"`
	}
)

func NewOrder(cs service.StoreService, l logger.Logger) *Order {
	return &Order{cs, l}
}

// getAllOrders godoc
// @Description  Returns user's orders.
// @Tags         Order
// @Accept       json
// @Produce      json
// @Success 200 {object} orderResponse
// @Failure 400 {object} response
// @Failure 404 {object} response
// @Router /order [get]
// @Security Bearer
func (c *Order) GetAllOrders(g *gin.Context) {
	userName := g.GetString("Email")
	if len(userName) == 0 {
		errorResponse(g, http.StatusBadRequest, "unable to get parameters")
		return
	}

	orders := c.storeService.GetAllOrders(userName)
	g.JSON(http.StatusOK, orders)
}

// createOrder godoc
// @Description  Creates an order with basket items.
// @Tags         Order
// @Accept       json
// @Produce      json
// @Param data body newOrderRequest true "New Order Model"
// @Success 200 {object} response
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /order [post]
// @Security Bearer
func (c *Order) CreateOrder(g *gin.Context) {
	userName := g.GetString("Email")
	if len(userName) == 0 {
		errorResponse(g, http.StatusBadRequest, "unable to get parameters")
		return
	}

	var req newOrderRequest
	if err := g.ShouldBind(&req); err != nil {
		c.logger.Error(err, "http - v1 - createOrder")
		errorResponse(g, http.StatusBadRequest, "invalid request body")
		return
	}

	err := c.storeService.CreateOrder(userName, req.Name, req.Address, req.PhoneNumber,
		req.CardNumber, req.CardExp, req.CardCVV)
	if err != nil {
		errorResponse(g, http.StatusNotFound, err.Error())
		return
	}

	successResponse(g, http.StatusOK, "Operation completed successfully.")
}

// cancelOrder godoc
// @Description  Cancels an order.
// @Tags         Order
// @Accept       json
// @Produce      json
// @Param        id path string true "Order ID"
// @Success 200 {object} response
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /order/{id}/cancel [patch]
// @Security Bearer
func (c *Order) CancelOrder(g *gin.Context) {
	userName := g.GetString("Email")
	if len(userName) == 0 {
		errorResponse(g, http.StatusBadRequest, "unable to get parameters")
		return
	}

	id := g.Param("id")
	if len(id) == 0 {
		errorResponse(g, http.StatusBadRequest, "unable to get parameters")
		return
	}

	err := c.storeService.CancelOrder(userName, id)
	if err != nil {
		errorResponse(g, http.StatusNotFound, err.Error())
		return
	}

	successResponse(g, http.StatusOK, "Operation completed successfully.")
}
