package database

import (
	"context"
	"fmt"

	"github.com/sleekservices/ServiceRenderer/common/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func New(ctx context.Context, params *MongoCredentials) (*Database, error) {
	db := &Database{}
	if err := db.Client(ctx, params); err != nil {
		return nil, err
	}

	if err := db.Connect(ctx); err != nil {
		return nil, err
	}

	if err := db.Database(params.Name); err != nil {
		return nil, err
	}

	return db, nil
}

func (d *Database) Client(ctx context.Context, params *MongoCredentials) error {
	if d.MongoClient != nil {
		log.Debug("DB client already exist")
	}
	dbcredentials := fmt.Sprintf(`mongodb://%s:%s@%s/%s?retryWrites=false`,
		params.Username,
		params.Password,
		params.HostURL,
		params.AuthName,
	)

	client, err := mongo.NewClient(options.Client().ApplyURI(dbcredentials))
	if err != nil {
		return err
	}

	d.MongoClient = client

	return nil
}

func (d *Database) Connect(ctx context.Context) error {
	if d.MongoClient == nil {
		return fmt.Errorf("No mongo client running")
	}
	return d.MongoClient.Connect(ctx)
}

func (d *Database) Database(name string) error {
	if d.MongoClient == nil {
		return fmt.Errorf("No mongo client running")
	}

	d.MongoDatabase = d.MongoClient.Database(name)
	return nil
}

func (d *Database) Collection(name string) *mongo.Collection {
	return d.MongoDatabase.Collection(name)
}

func EnsureOrUpgradeIndex(ctx context.Context, coll *mongo.Collection, model IndexModel) {
	name, err := coll.Indexes().CreateOne(ctx, model.GetModel())
	if err != nil {
		log.Debug("issue creating index err=%v", err)
	}
	log.Debug("index created: %s", name)
}
