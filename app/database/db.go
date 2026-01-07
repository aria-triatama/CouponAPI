package database

import (
	"CouponAPI/helpers"
	"context"
	"errors"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBstore interface {
	GetDB() *mongo.Database
	GetClient() (*mongo.Client, error)
	GetCtx() *context.Context
	Disconnect() error
}

type DBStore struct {
	client *mongo.Client
	db     *mongo.Database
	ctx    *context.Context
}

func MongoDB() (DBstore, error) {
	client, db, err := Connect()
	if err != nil {
		return nil, err
	}
	cont := context.Background()
	return &DBStore{client: client, db: db, ctx: &cont}, nil
}

func (mdb *DBStore) GetDB() *mongo.Database {
	return mdb.db
}

func (mdb *DBStore) GetClient() (*mongo.Client, error) {
	if mdb.client != nil {
		return mdb.client, nil
	}
	return nil, errors.New("missing client")
}

func (mdb *DBStore) GetCtx() *context.Context {
	return mdb.ctx
}

func (mdb *DBStore) Disconnect() error {
	err := mdb.client.Disconnect(*mdb.ctx)
	if err != nil {
		return err
	}
	return nil
}

func Connect() (*mongo.Client, *mongo.Database, error) {
	var connectOnce sync.Once
	var db *mongo.Database
	var client *mongo.Client
	var err error
	connectOnce.Do(func() {
		client, db, err = MongoConnect()
	})
	return client, db, err
}

func MongoConnect() (*mongo.Client, *mongo.Database, error) {
	clientOptions := options.Client()
	clientOptions.ApplyURI(helpers.GetEnv("MONGO_URI"))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientCon, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, nil, err
	}

	clientDB := clientCon.Database(helpers.GetEnv("MONGO_DB"))
	return clientCon, clientDB, nil
}
