package utils

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DataStore struct {
	C *mongo.Client
}

var (
	dataStore DataStore
)

type IMongoDataStore interface {
}

func BuildDataStore() *DataStore {
	dataStore = DataStore{getMongoDBClient(config.DB)}
	return &dataStore
}

func GetDataStore() *DataStore {
	return &dataStore
}

func getMongoDBClient(mongoConfig Database) *mongo.Client {
	ctx, cancel := GetCtxInSec(10)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().
		ApplyURI("mongodb://"+mongoConfig.Host+":"+mongoConfig.Port).
		SetAuth(options.Credential{Username: mongoConfig.Username, Password: mongoConfig.Password}))
	if err != nil {
		panic(err)
	}
	ctx, cancel = GetCtxInSec(60)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(errors.New("Failed to connect to mongodb : " + err.Error()))
	}
	return client
}

func GetCtxInSec(t int) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(t)*time.Second)
	return ctx, cancel
}
