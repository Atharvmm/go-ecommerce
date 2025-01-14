package services

import (
	"errors"
	"fmt"
	"go-ecommerce/cache"
	"go-ecommerce/db"
	"go-ecommerce/models"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

func GetAllProducts() ([]models.Product, error) {
	var products []models.Product
	result := db.DB.Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}

func GetcachedProducts() ([]models.Product, error) {
	products, err := cache.GetRecentlyViewedProducts(rdb)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func AddProduct(product *models.Product) error {
	result := db.DB.Create(product)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetProductByID(id string) (*models.Product, error) {
	cacheProduct, err := cache.GetProductByID(rdb, id)
	if err == nil {
		return cacheProduct, nil
	} else if err.Error() != "cache miss" {
		return nil, err
	}

	var product models.Product
	result := db.DB.First(&product, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("product not found")
		}
		return nil, err
	}

	err = cache.AddProductsToRecentlyViewed(rdb, product)
	if err != nil {
		return nil, fmt.Errorf("failed to update the cache with the new product")
	}

	return &product, nil
}

func UpdateProduct(id string, updatedProduct *models.Product) error {
	var product models.Product

	result := db.DB.First(&product, "id=?", id)
	if result.Error != nil {
		return result.Error
	}
	product.Name = updatedProduct.Name
	product.Description = updatedProduct.Description
	product.Price = updatedProduct.Price
	db.DB.Save(&product)

	err := cache.AddProductsToRecentlyViewed(rdb, product)
	if err != nil {
		return fmt.Errorf("failed to update the cache with the new product")
	}

	return nil
}

func DeleteProduct(id string) error {
	result := db.DB.Delete(&models.Product{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func SearchProductsByName(name string) ([]models.Product, error) {
	var products []models.Product
	result := db.DB.Where("name LIKE ?", "%"+name+"%").Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}

	return products, nil
}

func DeleteAllProducts() error {
	result := db.DB.Exec("DELETE FROM products")
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func AddProducts(products []models.Product) error {
	result := db.DB.Create(&products)
	return result.Error
}

func DeleteProductsByID(ids []string) error {
	tx := db.DB.Begin()

	for _, id := range ids {
		result := tx.Delete(&models.Product{}, id)
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}

		err := cache.DeleteProductFromCache(rdb, id)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}
