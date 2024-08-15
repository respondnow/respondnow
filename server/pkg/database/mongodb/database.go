package mongodb

import (
	"context"

	"github.com/respondnow/respond/server/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	lockKey = "__lock:iro:database_setup"
)

func InitMongoClient() error {
	server := config.EnvConfig.MongoConfig.MongoURL
	db := config.EnvConfig.MongoConfig.MongoDbName
	dbUser := config.EnvConfig.MongoConfig.MongoDbUsername
	dbPassword := config.EnvConfig.MongoConfig.MongoDbPassword

	opts := options.Client().ApplyURI(server)
	if dbUser != "" && dbPassword != "" {
		credential := options.Credential{
			Username: dbUser,
			Password: dbPassword,
		}
		opts = opts.SetAuth(credential)
	}

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return err
	}

	MgoClient = client
	MongoSession = MgoClient
	MClient, err = MClient.Initialize(MgoClient, db)
	if err != nil {
		return err
	}

	return nil
}
