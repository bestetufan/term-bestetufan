package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	"github.com/bestetufan/beste-store/config"
	_ "github.com/bestetufan/beste-store/docs"
	"github.com/bestetufan/beste-store/internal/api/controller"
	"github.com/bestetufan/beste-store/internal/api/middleware"
	"github.com/bestetufan/beste-store/internal/domain/repo"
	"github.com/bestetufan/beste-store/internal/service"
	"github.com/bestetufan/beste-store/pkg/logger"
)

// NewRouter -.
// Swagger spec:
// @title       Beste Tufan Store API
// @description An API for store management
// @version     1.0
// @host        localhost:8080
// @BasePath    /api/v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func NewRouter(handler *gin.Engine, l *logger.Logger, c *config.Config, db *gorm.DB) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// Health probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Repo
	userRepo := repo.NewUserRepository(db)
	categoryRepo := repo.NewCategoryRepository(db)
	productRepo := repo.NewProductRepository(db)
	basketRepo := repo.NewBasketRepository(db)
	orderRepo := repo.NewOrderRepository(db)

	// Service
	authService := service.NewJWTAuthService(*c)
	userService := service.NewUserService(*userRepo)
	storeService := service.NewStoreService(*categoryRepo, *productRepo, *basketRepo, *orderRepo)

	// Controller
	auth := controller.NewAuth(*userService, *authService, *l)
	category := controller.NewCategory(*storeService, *l)
	product := controller.NewProduct(*storeService, *l)
	basket := controller.NewBasket(*storeService, *l)
	order := controller.NewOrder(*storeService, *l)

	// Middleware
	authMw := middleware.NewJWTAuthMiddleware(*authService, *userService, *l)

	// Routers
	h := handler.Group("/api/v1")
	{
		a := h.Group("/auth")
		{
			a.POST("/login", auth.Login)
			a.POST("/register", auth.Register)
		}
		c := h.Group("/category", authMw.ValidateToken())
		{
			c.GET("", category.GetAllCategories)
			c.GET(":id", category.GetCategory)
			c.POST("", authMw.CheckRole("admin"), category.CreateCategory)
			c.DELETE("", authMw.CheckRole("admin"), category.DeleteCategory)
			c.POST("/bulk", authMw.CheckRole("admin"), category.CreateBulkCategory)
		}
		p := h.Group("/product", authMw.ValidateToken())
		{
			p.GET("", product.GetAllProducts)
			p.GET(":id", product.GetProduct)
			p.GET("/search/:query", product.SearchProducts)
			p.POST("", authMw.CheckRole("admin"), product.CreateProduct)
			p.PUT(":id", authMw.CheckRole("admin"), product.UpdateProduct)
			p.DELETE(":id", authMw.CheckRole("admin"), product.DeleteProduct)
		}
		b := h.Group("/basket", authMw.ValidateToken())
		{
			b.GET("", basket.GetBasket)
			b.POST("", basket.AddBasketItem)
			b.PUT("", basket.UpdateBasketItem)
			b.DELETE("", basket.RemoveBasketItem)
		}
		o := h.Group("/order", authMw.ValidateToken())
		{
			o.GET("", order.GetAllOrders)
			o.POST("", order.CreateOrder)
			o.PATCH(":id/cancel", order.CancelOrder)
		}
	}
}
