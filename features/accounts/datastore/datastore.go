package datastore

import (
	"context"

	"github.com/sleekservices/ServiceRenderer/common/authority"
	"github.com/sleekservices/ServiceRenderer/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Datastore struct {
	database *database.Database
}

func NewDatastore(ctx context.Context, d *database.Database) *Datastore {
	dStore := &Datastore{database: d}
	database.EnsureOrUpgradeIndex(ctx, dStore.serviceRenderCollection(), database.IndexModel{
		Keys: bson.M{
			"email": 1,
		},
		Options: database.NewIndexOptions().SetUnique(true),
	})
	return dStore
}

func (d *Datastore) serviceRenderCollection() *mongo.Collection {
	return d.database.Collection("ServiceRender")
}

func (d *Datastore) CreateServiceProvider(ctx context.Context, render *authority.ServiceRenderer) error {
	_, err := d.serviceRenderCollection().InsertOne(ctx, render)
	return err
}

func (d *Datastore) FindServiceProviderByID(ctx context.Context, ID primitive.ObjectID) (*authority.ServiceRenderer, error) {
	var render authority.ServiceRenderer
	if err := d.serviceRenderCollection().FindOne(ctx, bson.M{"_id": ID}).Decode(&render); err != nil {
		return nil, err
	}
	return &render, nil
}

func (d *Datastore) FindServiceProviderByEmail(ctx context.Context, email string) (*authority.ServiceRenderer, error) {
	var user authority.ServiceRenderer
	if err := d.serviceRenderCollection().FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *Datastore) GetAllServiceProvider(ctx context.Context) ([]authority.ServiceRenderer, error) {
	return nil, nil
}

func (d *Datastore) UpdateServiceProviderByID(ctx context.Context, Id primitive.ObjectID) error {
	return nil
}
