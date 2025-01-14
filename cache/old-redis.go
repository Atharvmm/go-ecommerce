package cache

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"go-ecommerce/models"

// 	"github.com/go-redis/redis/v8"
// )

// var ctx = context.Background()

// // GetProducts retrieves cached products from Redis
// func GetProducts(rdb *redis.Client) ([]models.Product, error) {
// 	val, err := rdb.Get(ctx, "products").Result()
// 	if err == redis.Nil {
// 		// Cache miss: products are not found in Redis
// 		return nil, fmt.Errorf("cache miss")
// 	} else if err != nil {
// 		return nil, err
// 	}

// 	var products []models.Product
// 	err = json.Unmarshal([]byte(val), &products)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return products, nil
// }

// // SetProducts caches products in Redis
// func SetProducts(rdb *redis.Client, products []models.Product) error {
// 	productsJSON, err := json.Marshal(products)
// 	if err != nil {
// 		return err
// 	}
// 	err = rdb.Set(ctx, "products", productsJSON, 0).Err()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
