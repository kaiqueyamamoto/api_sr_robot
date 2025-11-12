package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client   *mongo.Client
	Database *mongo.Database
)

// Connect estabelece conexão com o MongoDB
func Connect(mongoURI, dbName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURI)
	
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// Verificar a conexão
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	Client = client
	Database = client.Database(dbName)

	log.Println("✅ Conectado ao MongoDB Atlas!")
	return nil
}

// Disconnect desconecta do MongoDB
func Disconnect() error {
	if Client == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return Client.Disconnect(ctx)
}

// GetCollection retorna uma collection do MongoDB
func GetCollection(collectionName string) *mongo.Collection {
	return Database.Collection(collectionName)
}
