package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"go-ecommerce/models"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

const (
	lruCacheKey   = "recently_viewed_products"
	maxCacheSize  = 15
	cacheDuration = 24 * time.Hour
)

func GetRecentlyViewedProducts(rdb *redis.Client) ([]models.Product, error) {
	productIDs, err := rdb.ZRevRange(ctx, lruCacheKey, 0, -1).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	var products []models.Product
	for _, id := range productIDs {
		val, err := rdb.Get(ctx, id).Result()
		if err == redis.Nil {
			continue
		} else if err != nil {
			return nil, err
		}

		var product models.Product
		err = json.Unmarshal([]byte(val), &product)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func AddProductsToRecentlyViewed(rdb *redis.Client, product models.Product) error {
	productJSON, err := json.Marshal(product)
	if err != nil {
		return err
	}

	productID := strconv.Itoa(int(product.ID))
	productKey := "product:" + productID

	err = rdb.Set(ctx, productKey, productJSON, cacheDuration).Err()
	if err != nil {
		return err
	}

	err = rdb.ZAdd(ctx, lruCacheKey, &redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: productKey,
	}).Err()
	if err != nil {
		return err
	}

	err = rdb.ZRemRangeByRank(ctx, lruCacheKey, 0, int64(-maxCacheSize-1)).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetProductByID(rdb *redis.Client, id string) (*models.Product, error) {
	productKey := "product:" + id
	val, err := rdb.Get(ctx, productKey).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("cache miss")
	} else if err != nil {
		return nil, err
	}

	var product models.Product
	err = json.Unmarshal([]byte(val), &product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func DeleteProductFromCache(rdb *redis.Client, id string) error {
	productKey := "product:" + id
	err := rdb.Del(ctx, productKey).Err()
	if err != nil && err != redis.Nil {
		return err
	}

	err = rdb.ZRem(ctx, lruCacheKey, productKey).Err()
	if err != nil {
		return err
	}
	return nil
}
