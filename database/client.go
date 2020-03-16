package database

import "go.mongodb.org/mongo-driver/mongo"

type Database struct {
	MongoClient   *mongo.Client
	MongoDatabase *mongo.Database
}
