package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"tesodev-product-api/config"
	"tesodev-product-api/models"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetAllProducts veritabanından tüm ürünleri çeker
func GetAllProducts(c echo.Context) error {
	collection := config.DB.Collection("products")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Veritabanı hatası"})
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err := cursor.All(ctx, &products); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Veriler çözümlenemedi"})
	}

	return c.JSON(http.StatusOK, products)
}

// CreateProduct yeni bir ürün oluşturur
func CreateProduct(c echo.Context) error {
	var product models.Product

	// JSON'u parse et
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Geçersiz JSON formatı"})
	}

	// Geçerli alanlar mı?
	if product.Name == "" || product.Price <= 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Name ve Price zorunludur"})
	}

	// ID ve zaman ayarla
	product.ID = primitive.NewObjectID()
	product.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	// MongoDB'ye ekle
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := config.DB.Collection("products")

	res, err := collection.InsertOne(ctx, product)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "MongoDB ekleme hatası"})
	}

	product.ID = res.InsertedID.(primitive.ObjectID)
	return c.JSON(http.StatusCreated, product)
}

// GetProductByID ID'ye göre bir ürünü getirir
func GetProductByID(c echo.Context) error {
	idParam := c.Param("id")

	// ObjectID formatına çevir
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Geçersiz ürün ID formatı",
		})
	}

	var product models.Product
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := config.DB.Collection("products")
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&product)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Ürün bulunamadı",
		})
	}

	return c.JSON(http.StatusOK, product)
}

func UpdateProduct(c echo.Context) error {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Geçersiz ürün ID"})
	}

	var updatedProduct models.Product
	if err := c.Bind(&updatedProduct); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Geçersiz JSON formatı"})
	}

	if updatedProduct.Name == "" || updatedProduct.Price <= 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Name ve Price zorunludur"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := config.DB.Collection("products")

	update := bson.M{
		"$set": bson.M{
			"name":        updatedProduct.Name,
			"description": updatedProduct.Description,
			"price":       updatedProduct.Price,
		},
	}

	result, err := collection.UpdateByID(ctx, objID, update)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Güncelleme hatası"})
	}
	if result.MatchedCount == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Ürün bulunamadı"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Ürün başarıyla güncellendi"})
}

// PatchProduct bir ürünün belirli alanlarını günceller
func PatchProduct(c echo.Context) error {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Geçersiz ürün ID formatı",
		})
	}

	// JSON'dan sadece gelen alanları al (map olarak)
	var updateData map[string]interface{}
	if err := c.Bind(&updateData); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Geçersiz JSON formatı",
		})
	}

	// Geçersiz price varsa çıkar
	if price, ok := updateData["price"]; ok {
		if _, ok := price.(float64); !ok {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error": "price değeri sayı olmalı",
			})
		}
	}

	// Update yapısı oluştur
	update := bson.M{"$set": updateData}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := config.DB.Collection("products")

	res, err := collection.UpdateByID(ctx, objID, update)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Güncelleme sırasında hata oluştu",
		})
	}
	if res.MatchedCount == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Ürün bulunamadı",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Ürün başarıyla güncellendi"})

}

func DeleteProduct(c echo.Context) error {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Geçersiz ürün ID"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := config.DB.Collection("products")

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Silme hatası"})
	}
	if result.DeletedCount == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Silinecek ürün bulunamadı"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Ürün silindi"})
}

// SearchProducts ürünleri isme ve fiyata göre filtreler, sıralar
func SearchProducts(c echo.Context) error {
	// Query parametrelerini al
	name := c.QueryParam("name")
	minPrice := c.QueryParam("minPrice")
	maxPrice := c.QueryParam("maxPrice")
	sortOrder := c.QueryParam("sort")

	filter := bson.M{}

	// İsme göre arama (kısmi eşleşme)
	if name != "" {
		filter["name"] = bson.M{"$regex": name, "$options": "i"}
	}

	// Fiyat aralığı
	priceRange := bson.M{}
	if minPrice != "" {
		if min, err := strconv.ParseFloat(minPrice, 64); err == nil {
			priceRange["$gte"] = min
		}
	}
	if maxPrice != "" {
		if max, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			priceRange["$lte"] = max
		}
	}
	if len(priceRange) > 0 {
		filter["price"] = priceRange
	}

	// Sıralama
	sort := bson.D{}
	if sortOrder == "asc" {
		sort = bson.D{{Key: "price", Value: 1}}
	} else if sortOrder == "desc" {
		sort = bson.D{{Key: "price", Value: -1}}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := config.DB.Collection("products")

	opts := options.Find().SetSort(sort)

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Veri aranırken hata oluştu"})
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err := cursor.All(ctx, &products); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Veriler çözümlenemedi"})
	}

	return c.JSON(http.StatusOK, products)
}
