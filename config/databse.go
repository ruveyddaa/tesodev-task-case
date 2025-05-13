package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Mongo URI
	clientOptions := options.Client().ApplyURI("mongodb://root:example@localhost:27017")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("MongoDB bağlantı hatası: %v", err)
	}

	// Ping ile bağlantıyı test et
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB erişilemiyor: %v", err)
	}

	fmt.Println("✅ MongoDB bağlantısı başarılı!")

	// tesodevdb adlı veritabanını seç
	DB = client.Database("tesodevdb")
}
