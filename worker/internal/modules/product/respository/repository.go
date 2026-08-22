package product_respository

import (
	"context"

	"github.com/celio001/product-cqrs/worker/internal/modules/product"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type productRepository struct {
	client *mongo.Client
}

type ProductRepositoryInterface interface {
	CreateProductRepository(ctx context.Context, product product.Product) error
}

func NewProductRepository(client *mongo.Client) ProductRepositoryInterface {
	return &productRepository{
		client: client,
	}
}

func (c *productRepository) CreateProductRepository(ctx context.Context, product product.Product) error {
	collection := c.client.Database("products").Collection("products")

	_, err := collection.UpdateOne(ctx, bson.M{"id": product.ID}, bson.M{"$set": product}, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}

	return nil
}
