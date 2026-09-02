package product_repository

import (
	"context"

	"github.com/celio001/product-cqrs/product-query/internal/modules/product"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type productRepo struct {
	c *mongo.Client
}

type ProductRepositoryInterface interface {
	GetProductByID(ctx context.Context, id uuid.UUID) (product.Product, error)
}

func NewProductRepository(c *mongo.Client) ProductRepositoryInterface {
	return &productRepo{
		c: c,
	}
}

func (r *productRepo) GetProductByID(ctx context.Context, id uuid.UUID) (product.Product, error) {

	var p product.Product

	collection := r.c.Database("products").Collection("products")
	filter := bson.D{primitive.E{Key: "id", Value: id.String()}}

	err := collection.FindOne(ctx, filter).Decode(&p)
	if err!= nil{
		return product.Product{}, err
	}

	return p, nil
}
