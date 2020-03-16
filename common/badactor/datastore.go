package badactor

import (
	"context"
	"time"

	"github.com/sleekservices/ServiceRenderer/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BadactorDB struct {
	database *database.Database
}

func NewBadactorService(ctx context.Context, db *database.Database) *BadactorDB {
	badactor := &BadactorDB{db}
	expireAfter := int32(1440)
	background := true

	database.EnsureOrUpgradeIndex(ctx, badactor.infractions(), database.IndexModel{
		Keys: bson.M{
			"actor":     1,
			"rule.name": 1,
			"expire_by": 1,
		},
		Options: database.NewIndexOptions().SetUnique(true),
	})
	database.EnsureOrUpgradeIndex(ctx, badactor.infractions(), database.IndexModel{

		Keys: bson.M{"expire_by": 1},
		Options: &options.IndexOptions{
			ExpireAfterSeconds: &expireAfter,
			Background:         &background,
		},
	})
	database.EnsureOrUpgradeIndex(ctx, badactor.jails(), database.IndexModel{
		Keys: bson.M{
			"actor":      1,
			"rule.name":  1,
			"release_by": 1,
		},
	})
	return badactor
}

func (self *BadactorDB) infractions() *mongo.Collection {
	return self.database.Collection("badactor_infractions")
}

func (self *BadactorDB) jails() *mongo.Collection {
	return self.database.Collection("badactor_jails")
}

func (self *BadactorDB) CreateInfraction(ctx context.Context, infraction *Infraction) (*Infraction, error) {
	_, err := self.infractions().InsertOne(ctx, infraction)
	if err != nil {
		return nil, err
	}
	return infraction, nil
}

func (self *BadactorDB) CountInfraction(ctx context.Context, actorName, ruleName string, expireTerm time.Time) (int64, error) {
	query := bson.M{
		"actor":     actorName,
		"rule.name": ruleName,
		"expire_by": bson.M{"$gt": expireTerm},
	}

	return self.infractions().CountDocuments(ctx, query)
}

func (self *BadactorDB) CreateJail(ctx context.Context, jail *Jail) error {
	jail.ID = primitive.NewObjectID()
	_, err := self.jails().InsertOne(ctx, jail)
	return err
}

func (self *BadactorDB) FindJail(ctx context.Context, actorName, ruleName string, releaseTerm time.Time) (*Jail, error) {
	jail := &Jail{}
	query := bson.M{
		"actor":      actorName,
		"rule.name":  ruleName,
		"release_by": bson.M{"$gt": releaseTerm},
	}

	err := self.jails().FindOne(ctx, query).Decode(jail)
	return jail, err
}

func (self *BadactorDB) UpdateJail(ctx context.Context, jail *Jail) error {
	_, err := self.jails().UpdateOne(ctx,
		bson.M{"_id": jail.ID},
		bson.M{"$set": bson.M{
			"release_by": jail.ReleaseBy,
		},
		})
	return err
}
