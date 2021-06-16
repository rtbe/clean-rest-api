// Package database provides
// database specific functions.
package database

import (
	"context"
	"net/url"
	"sync"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoOnce sync.Once
	mongoDB   *mongo.Database
)

type MongoConfig struct {
	User     string
	Password string
	Host     string
	Name     string
}

// NewMongo returns connection to mongoDB.
func NewMongo(cfg MongoConfig) (*mongo.Database, error) {
	var connErr error

	mongoOnce.Do(func() {

		u := url.URL{
			Scheme: "mongodb",
			User:   url.UserPassword(cfg.User, cfg.Password),
			Host:   cfg.Host,
		}

		client, err := mongo.NewClient(options.Client().ApplyURI(u.String()))
		if err != nil {
			connErr = errors.Wrap(err, "auth db")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		err = client.Connect(ctx)
		if err != nil {
			connErr = errors.Wrap(err, "auth db")
		}

		if err := client.Ping(ctx, nil); err != nil {
			connErr = errors.Wrap(err, "auth db")
		}

		mongoDB = client.Database(cfg.Name)
	})
	return mongoDB, connErr
}
