package authority

import (
	"context"
	"github.com/sleekservices/ServiceRenderer/common/password"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var Anonymous = ServiceRenderer{}

type ServiceRenderer struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	FirstName    string             `bson:"first_name,omitempty"`
	LastName     string             `bosn:"last_name,omitempty"`
	Email        string             `bson:"email,omitempty"`
	PhoneNumber  string             `bson:"phone_number,omitempty"`
	Address      string             `bson:"address,omitempty"`
	LGA          string             `bson:"lga,omitempty"`
	State        string             `bson:"state,omitempty"`
	Location     gEOLocation        `bson:"location,omitempty"`
	Gender       string             `bson:"gender,omitempty"`
	DOB          time.Time          `bson:"dob,omitempty"`
	BVN          string             `bson:"bvn,omitempty"` // Do we really need this, Will user b willing to provide
	Guarantors   []guarantor        `bson:"guarantors,omitempty"`
	Photo        string             `bson:"photo,omitempty"`
	Verification verification       `bson:"verification,omitempty"`
	UserName     string             `bson:"user_name,omitempty"`
	Pin          *password.Hash     `bson:"pin"`
}

type guarantor struct {
	FirstName    string `bson:"first_name,omitempty"`
	LastName     string `bson:"last_name,omitempty"`
	Address      string `bson:"address,omitempty"`
	LGA          string `bson:"lga,omitempty"`
	State        string `bson:"state,omitempty"`
	PhoneNo      string `bson:"phone_number,omitempty"`
	Relationship string `bson:"relationship,omitempty"`
}

type verification struct {
	Type       string    `bson:"type,omitempty"`
	Number     string    `bson:"number,omitempty"`
	IssuedDate time.Time `bson:"issued_date,omitempty"`
	ExpiryDate time.Time `bson:"expiry_date,omitempty"`
}

type coordinate struct {
	Longitude float64
	Latitude  float64
}

type gEOLocation struct {
	Type        string    `bson:"type"`
	Coordinates []float64 `bson:"coordinates"`
}

type ServiceProvider interface {
	CreateServiceProvider(ctx context.Context, r *ServiceRenderer) error
	FindServiceProviderByID(ctx context.Context, Id primitive.ObjectID) (*ServiceRenderer, error)
	FindServiceProviderByEmail(ctx context.Context, email string) (*ServiceRenderer, error)
	GetAllServiceProvider(ctx context.Context) ([]ServiceRenderer, error)
	UpdateServiceProviderByID(ctx context.Context, Id primitive.ObjectID) error
}

func (g gEOLocation) Longitude() float64 {
	if len(g.Coordinates) < 1 {
		return 0
	}
	return g.Coordinates[0]
}

func (g gEOLocation) Latitude() float64 {
	if len(g.Coordinates) < 2 {
		return 0
	}
	return g.Coordinates[1]
}
