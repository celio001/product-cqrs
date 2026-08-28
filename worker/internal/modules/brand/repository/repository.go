package brand_repository

import (
	"context"

	"github.com/celio001/product-cqrs/worker/internal/modules/brand"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type brandRepository struct {
	client *mongo.Client
}

type BrandRepositoryInterface interface {
	CreateBrandRepository(ctx context.Context, brand brand.Brand) error
}

func NewProductRepository(client *mongo.Client) BrandRepositoryInterface {
	return &brandRepository{
		client: client,
	}
}

func (c *brandRepository) CreateBrandRepository(ctx context.Context, brand brand.Brand) error {
	collection := c.client.Database("brands").Collection("brands")

	_, err := collection.UpdateOne(ctx, bson.M{"id": brand.ID}, bson.M{"$set": brand}, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}

	return nil
}
