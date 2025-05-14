package repository

import (
	"context"
	"tesodev-product-api/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductRepository struct {
	DB *mongo.Database
}

func (r *ProductRepository) GetAllProducts() ([]models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.DB.Collection("products").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) CreateProduct(product models.Product) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := r.DB.Collection("products").InsertOne(ctx, product)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func (r *ProductRepository) GetProductByID(id primitive.ObjectID) (models.Product, error) {
	var product models.Product
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.DB.Collection("products").FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	return product, err
}

func (r *ProductRepository) UpdateProduct(id primitive.ObjectID, updated models.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$set": bson.M{
		"name":        updated.Name,
		"description": updated.Description,
		"price":       updated.Price,
	}}

	res, err := r.DB.Collection("products").UpdateByID(ctx, id, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *ProductRepository) PatchProduct(id primitive.ObjectID, updateData map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$set": updateData}

	res, err := r.DB.Collection("products").UpdateByID(ctx, id, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *ProductRepository) DeleteProduct(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := r.DB.Collection("products").DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *ProductRepository) SearchProducts(name string, minPrice, maxPrice float64, sortOrder string) ([]models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{}
	if name != "" {
		filter["name"] = bson.M{"$regex": name, "$options": "i"}
	}
	priceRange := bson.M{}
	if minPrice > 0 {
		priceRange["$gte"] = minPrice
	}
	if maxPrice > 0 {
		priceRange["$lte"] = maxPrice
	}
	if len(priceRange) > 0 {
		filter["price"] = priceRange
	}

	sort := bson.D{}
	if sortOrder == "asc" {
		sort = bson.D{{Key: "price", Value: 1}}
	} else if sortOrder == "desc" {
		sort = bson.D{{Key: "price", Value: -1}}
	}

	opts := options.Find().SetSort(sort)
	cursor, err := r.DB.Collection("products").Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
}
