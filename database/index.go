package database

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IndexModel struct {
	Keys    interface{}
	Options *options.IndexOptions
}

func NewIndexOptions() *options.IndexOptions {
	return options.Index()
}
func (im IndexModel) GetModel() mongo.IndexModel {
	return mongo.IndexModel{
		Keys:    im.Keys,
		Options: im.Options,
	}
}
