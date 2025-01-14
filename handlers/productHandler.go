package handlers

import (
	"encoding/csv"
	"go-ecommerce/db"
	"go-ecommerce/models"
	"go-ecommerce/services"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetAllProducts(c echo.Context) error {
	products, err := services.GetAllProducts()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, products)
}

func GetCachedProducts(c echo.Context) error {
	products, err := services.GetcachedProducts()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, products)
}

func AddProductHandler(c echo.Context) error {
	var product models.Product

	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"Message": "Invalid input"})
	}
	result := db.DB.Create(&product)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Error": "Failed to add product"})
	}
	return c.JSON(http.StatusOK, product)
}

func GetProductHandler(c echo.Context) error {
	id := c.Param("id")
	product, err := services.GetProductByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, product)
}

func UpdateProductHandler(c echo.Context) error {
	id := c.Param("id")
	var product models.Product
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	err := services.UpdateProduct(id, &product)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Error": err.Error()})
	}
	return c.JSON(http.StatusOK, product)
}

func DeleteProductHandler(c echo.Context) error {
	id := c.Param("id")
	err := services.DeleteProduct(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"Error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"Message": "Product deleted successfully"})
}

func SearchProductHandler(c echo.Context) error {
	name := c.QueryParam("name")
	products, err := services.SearchProductsByName(name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Error": err.Error()})
	}
	return c.JSON(http.StatusOK, products)
}

func DeleteAllProductsHandler(c echo.Context) error {
	err := services.DeleteAllProducts()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"Message": "All products deleted successfully"})
}

func UploadProducts(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"Error": "Invalid file upload"})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Error": "could not open the file"})
	}
	defer src.Close()
	products := []models.Product{}
	reader := csv.NewReader(src)
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		price, _ := strconv.ParseFloat(record[2], 64)
		product := models.Product{
			Name:        record[0],
			Description: record[1],
			Price:       price,
		}
		products = append(products, product)

	}
	err = services.AddProducts(products)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Error": "Products couold not be added"})
	}
	return c.JSON(http.StatusOK, map[string]string{"Message": "Products added successfully"})
}

type DeleteRequest struct {
	ProductIDs []string `json:"product_ids"`
}

func DeleteProductsByIDs(c echo.Context) error {
	var req DeleteRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	err := services.DeleteProductsByID(req.ProductIDs)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete products"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Products deleted successfully"})
}
