package main

import (
	"github.com/gin-gonic/gin"

	"go-api/controller"
	"go-api/db"
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

	server.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Rota para obter todos os produtos
	server.GET("/products", ProductController.GetProducts)

	server.POST("/product", ProductController.CreateProduct)

	// Rota com parâmetro de rota para obter um produto por ID
	server.GET("/product/:productId", ProductController.GetProductById)

	// Rota para atualizar um produto por ID
	server.PUT("/product/:productId", ProductController.UpdateProduct)

	// Rota para deletar um produto por ID
	server.DELETE("/product/:productId", ProductController.DeleteProduct)

	server.Run(":8000")
}
