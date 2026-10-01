package product_respository

import (
	"context"
	"errors"

	"github.com/celio001/product-cqrs/worker/internal/modules/product"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrProductNotFound = errors.New("product not found")
)

type productRepository struct {
	client *mongo.Client
}

type ProductRepositoryInterface interface {
	CreateProductRepository(ctx context.Context, product product.Product) error
	SoftDeleteProductRepository(ctx context.Context, id string) error
	UpdateProductRepository(ctx context.Context, product product.Product) error
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

func (r *productRepository) SoftDeleteProductRepository(ctx context.Context, id string) error {
	collection := r.client.Database("products").Collection("products")

	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"status": "DISABLED"}}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *productRepository) UpdateProductRepository(ctx context.Context, product product.Product) error {
	collection := r.client.Database("products").Collection("products")

	filter := bson.M{"id": product.ID}
	update := bson.M{"$set": product}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrProductNotFound
		}
		return err
	}

	return nil
}
