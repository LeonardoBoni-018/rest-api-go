package repository

import (
	"database/sql"
	"fmt"

	"go-api/model"
)

type ProductRepository struct {
	connection *sql.DB
}

func NewProductRepository(connection *sql.DB) *ProductRepository {
	return &ProductRepository{
		connection: connection,
	}
}

func (pr *ProductRepository) GetProducts() ([]model.Product, error) {
	query := "SELECT id, product_name, price FROM product"

	rows, err := pr.connection.Query(query)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return []model.Product{}, err
	}
	var productList []model.Product
	var productObject model.Product

	for rows.Next() {
		err = rows.Scan(
			&productObject.ID,
			&productObject.Name,
			&productObject.Price)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return []model.Product{}, err
		}
		productList = append(productList, productObject)
	}

	rows.Close()

	return productList, nil

}

func (pr *ProductRepository) CreateProduct(product model.Product) (int, error) {
	var id int
	query, err := pr.connection.Prepare("INSERT INTO product" +
		"(product_name, price)" +
		"VALUES ($1, $2) RETURNING id")

	if err != nil {
		fmt.Println("Error preparing query:", err)
		return 0, err
	}

	err = query.QueryRow(product.Name, product.Price).Scan(&id)
	if err != nil {
		fmt.Println("Error preparing query:", err)
		return 0, err
	}
	query.Close()

	return id, nil
}

func (pr *ProductRepository) GetProductById(id_product int) (*model.Product, error) {
	query, err := pr.connection.Prepare("SELECT * FROM product WHERE id = $1")

	if err != nil {
		fmt.Println("Error preparing query:", err)
		return nil, err
	}

	var produto model.Product
	err = query.QueryRow(id_product).Scan(
		&produto.ID,
		&produto.Name,
		&produto.Price,
	)

	if err != nil {
		// Banco de dados não encontrou o produto, retornando um produto vazio
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	query.Close()

	return &produto, nil
}

func (pr *ProductRepository) UpdateProduct(id_product int, product model.Product) (*model.Product, error) {
	query, err := pr.connection.Prepare("UPDATE product SET product_name = $1, price = $2 WHERE id = $3 RETURNING id, product_name, price")

	if err != nil {
		fmt.Println("Error preparing query:", err)
		return nil, err
	}

	var updatedProduct model.Product
	err = query.QueryRow(product.Name, product.Price, id_product).Scan(
		&updatedProduct.ID,
		&updatedProduct.Name,
		&updatedProduct.Price,
	)

	if err != nil {
		fmt.Println("Error executing query:", err)
		return nil, err
	}

	query.Close()

	return &updatedProduct, nil
}

func (pr *ProductRepository) DeleteProduct(id_product int) error {
	query, err := pr.connection.Prepare("DELETE FROM product WHERE id = $1")

	if err != nil {
		fmt.Println("Error preparing query:", err)
		return err
	}

	_, err = query.Exec(id_product)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return err
	}
	query.Close()
	return nil
}
