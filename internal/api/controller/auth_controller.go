package controller

import (
	"net/http"

	"github.com/bestetufan/beste-store/internal/domain/entity"
	"github.com/bestetufan/beste-store/internal/service"
	"github.com/bestetufan/beste-store/pkg/logger"
	"github.com/gin-gonic/gin"
)

type (
	Auth struct {
		userService service.UserService
		authService service.JWTAuthService
		logger      logger.Logger
	}

	registerRequest struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	loginRequest struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	authResponse struct {
		Token string `json:"token"`
	}
)

func NewAuth(us service.UserService, as service.JWTAuthService, l logger.Logger) *Auth {
	return &Auth{us, as, l}
}

// register godoc
// @Description  Creates a new user account and a JWT token.
// @Tags         User
// @Accept       json
// @Produce      json
// @Param data body registerRequest true "Register Model"
// @Success 200 {object} authResponse
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /auth/register [post]
func (c *Auth) Register(g *gin.Context) {
	var req registerRequest
	if err := g.ShouldBind(&req); err != nil {
		c.logger.Error(err, "http - v1 - register")
		errorResponse(g, http.StatusBadRequest, "invalid request body")
		return
	}

	user := entity.NewUser(req.Email, req.Password, []*entity.Role{{Name: "customer"}})
	err := c.userService.CreateUser(user)
	if err != nil {
		c.logger.Error(err, "http - v1 - register")
		errorResponse(g, http.StatusInternalServerError, err.Error())
		return
	}

	token, err := c.authService.CreateToken(*user)
	if err != nil {
		errorResponse(g, http.StatusNotFound, err.Error())
		return
	}

	g.JSON(http.StatusOK, authResponse{Token: *token})
}

// Register godoc
// @Description  Returns a jwt token related to existing user account.
// @Tags         User
// @Accept       json
// @Produce      json
// @Param data body loginRequest true "Login Model"
// @Success 200 {object} authResponse
// @Failure 400 {object} response
// @Failure 404 {object} response
// @Router /auth/login [post]
func (c *Auth) Login(g *gin.Context) {
	var req loginRequest
	if err := g.ShouldBind(&req); err != nil {
		c.logger.Error(err, "http - v1 - login")
		errorResponse(g, http.StatusBadRequest, "invalid request body")
		return
	}

	user := c.userService.GetUser(req.Email, req.Password)
	if user == nil {
		errorResponse(g, http.StatusNotFound, "invalid email or password")
		return
	}

	token, err := c.authService.CreateToken(*user)
	if err != nil {
		errorResponse(g, http.StatusNotFound, err.Error())
		return
	}

	g.JSON(http.StatusOK, authResponse{Token: *token})
}
