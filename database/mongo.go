package database

import "context"

type MongoCredentials struct {
	HostURL  string
	Password string
	Username string
	Name     string
	AuthName string
}

type DBConnection interface {
	Client(ctx context.Context, params *MongoCredentials) error
	Connect(ctx context.Context) error
	Database(name string) error
}
