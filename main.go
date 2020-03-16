package main

import (
	"context"

	"github.com/sleekservices/ServiceRenderer/common/badactor"
	"github.com/sleekservices/ServiceRenderer/common/config"
	"github.com/sleekservices/ServiceRenderer/common/log"
	"github.com/sleekservices/ServiceRenderer/common/redis"
	"github.com/sleekservices/ServiceRenderer/database"
	"github.com/sleekservices/ServiceRenderer/features/accounts"
	"github.com/sleekservices/ServiceRenderer/features/accounts/datastore"
)

func main() {

	err := config.LoadAndWatch()
	if err != nil {
		log.Panic("Load and watch config err=%v", err)
	}

	initialContext := context.Background()

	logLevel := log.ParseLevel(config.MustString("log_level"))
	log.SetLevel(logLevel)

	dbCred := &database.MongoCredentials{
		Username: config.MustString("mongodb.username"),
		Password: config.MustString("mongodb.password"),
		HostURL:  config.MustString("mongodb.uri"),
		AuthName: config.MustString("mongodb.authname"),
		Name:     config.MustString("mongodb.db"),
	}

	redisHost := config.MustString("redis.host")
	redisClient := redis.New(redisHost)

	sessionStore := redis.NewSessionStore(redisClient)

	d, _ := database.New(initialContext, dbCred)

	serviceProvider := datastore.NewDatastore(initialContext, d)
	badactorService := badactor.NewBadactorService(initialContext, d)
	badactorStudio := badactor.NewStudio(initialContext, true, badactorService)
	accounts.NewHandler(serviceProvider, badactorStudio, sessionStore)

}
