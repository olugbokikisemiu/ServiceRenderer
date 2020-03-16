package datastore

import (
	"context"
	"os"
	"testing"

	"github.com/sleekservices/ServiceRenderer/common/authority"
	"github.com/sleekservices/ServiceRenderer/common/config"
	"github.com/sleekservices/ServiceRenderer/database"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var dStore *Datastore

func init() {
	err := config.LoadFromPath("../../../config.yaml")
	if err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	dbCred := &database.MongoCredentials{
		Username: config.MustString("mongodb.username"),
		Password: config.MustString("mongodb.password"),
		HostURL:  config.MustString("mongodb.uri"),
		AuthName: config.MustString("mongodb.authname"),
		Name:     config.MustString("mongodb.db"),
	}

	d, _ := database.New(ctx, dbCred)

	dStore = NewDatastore(ctx, d)

	exit := m.Run()
	os.Exit(exit)
}

func TestCreateServiceProvider__Should_create_service_provider_details_successfully(t *testing.T) {
	r := &authority.ServiceRenderer{
		ID:        primitive.NewObjectID(),
		FirstName: "Femi",
		LastName:  "Teste",
		Address:   "Niger",
		Email:     "test@rendize.com",
	}

	err := dStore.CreateServiceProvider(context.Background(), r)

	assert.Nil(t, err)
}

func TestFindServicroviderByID__Should_return_service_provider_details(t *testing.T) {
	r := &authority.ServiceRenderer{}
	r.ID = primitive.NewObjectID()
	r.FirstName = "Ton"
	r.LastName = "Joy"
	r.Email = "testing@rendize.com"

	createErr := dStore.CreateServiceProvider(context.Background(), r)

	assert.Nil(t, createErr)

	serviceProvider, err := dStore.FindServiceProviderByID(context.Background(), r.ID)

	assert.Nil(t, err)

	assert.Equal(t, serviceProvider.FirstName, r.FirstName)
	assert.Equal(t, serviceProvider.LastName, r.LastName)

}
