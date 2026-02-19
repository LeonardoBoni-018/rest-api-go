package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-api/model"
	"go-api/usecase"
)

type productController struct {
	//usecase
	productUsecase *usecase.ProductUseCase
}

func NewProductController(usecase usecase.ProductUseCase) *productController {
	return &productController{
		productUsecase: &usecase,
	}
}

func (p *productController) GetProducts(ctx *gin.Context) {
	//Dados mocados
	// products := []model.Product{
	// 	{
	// 		ID:    1,
	// 		Name:  "Produto 1",
	// 		Price: 10.0,
	// 	},
	// }
	products, err := p.productUsecase.GetProducts()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, products)
}

func (p *productController) CreateProduct(ctx *gin.Context) {
	var product model.Product

	err := ctx.BindJSON(&product)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	insertedProduct, err := p.productUsecase.CreateProduct(product)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, insertedProduct)
}

func (p *productController) GetProductById(ctx *gin.Context) {
	// Pegar o ID do produto a partir dos parâmetros da URL
	id := ctx.Param("productId")

	if id == "" {
		response := model.Response{
			Message: "Product ID is required",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	// Converter o ID para um inteiro
	productId, err := strconv.Atoi(id)
	if err != nil {
		response := model.Response{
			Message: "Invalid product ID",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	product, err := p.productUsecase.GetProductById(productId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	if product == nil {
		response := model.Response{
			Message: "Product not found",
		}
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	ctx.JSON(http.StatusOK, product)
}
