package auth

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rtbe/clean-rest-api/domain/entity"
	"github.com/rtbe/clean-rest-api/internal/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Mongo is an abstraction layer that manages auth entities inside MongoDB.
type Mongo struct {
	db *mongo.Collection
	logger.Logger
}

// NewMongoRepo creates a new MongoDB repository for auth entity
func NewMongoRepo(db *mongo.Database, l logger.Logger) *Mongo {
	coll := "authentication"

	return &Mongo{
		db.Collection(coll),
		l,
	}
}

// Create a new refresh token for particular user inside mongoDB.
func (r *Mongo) Create(ctx context.Context, rt entity.RefreshToken) error {
	_, err := r.db.InsertOne(
		ctx,
		bson.D{
			{"_id", rt.UserID},
			{"token", rt.Token},
			{"expires", rt.ExpiresAt},
		})
	if err != nil {
		return errors.Wrap(err, "error inserting tokens into mongoDB")
	}

	return nil
}

// Refresh replaces refresh token for particular user inside mongoDB.
func (r *Mongo) Refresh(ctx context.Context, userID string, rt entity.RefreshToken) error {
	err := r.Delete(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "error refreshing tokens inside mongoDB, delete token")
	}

	err = r.Create(ctx, rt)
	if err != nil {
		return errors.Wrap(err, "error refreshing tokens inside mongoDB, create token")
	}

	return nil
}

// Delete deletes refresh token from mongoDB by given user id.
func (r *Mongo) Delete(ctx context.Context, userID string) error {
	_, err := r.db.DeleteOne(
		ctx,
		bson.D{
			{"_id", userID},
		})
	if err != nil {
		return err
	}

	return nil
}
