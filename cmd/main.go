package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"go-api/controller"
	"go-api/db"
	"go-api/middleware"
	"go-api/repository"
	"go-api/usecase"
)

func main() {
	server := gin.Default()

	dbConnection, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}
	ProductRepository := repository.NewProductRepository(dbConnection)

	ProductUseCase := usecase.NewProductUseCase(*ProductRepository)

	ProductController := controller.NewProductController(*ProductUseCase)

	UserRepo := repository.NewUserRepository(dbConnection)
	AuthUseCase := usecase.NewAuthUserCase(UserRepo)
	AuthController := controller.NewAuthController(AuthUseCase)

	server.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// auth routes
	server.POST("/login", AuthController.Login)
	server.POST("/logout", AuthController.Logout)

	// public product routes
	server.GET("/products", ProductController.GetProducts)
	server.POST("/product", ProductController.CreateProduct)
	server.GET("/product/:productId", ProductController.GetProductById)

	// protected routes example (update/delete require auth)
	protected := server.Group("/")
	protected.Use(middleware.AuthMiddlewareFromEnv())
	protected.PUT("/product/:productId", ProductController.UpdateProduct)
	protected.DELETE("/product/:productId", ProductController.DeleteProduct)

	// debug: listar rotas
	for _, r := range server.Routes() {
		fmt.Println(r.Method, r.Path)
	}

	server.Run(":8000")
}
