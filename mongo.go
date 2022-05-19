// Copyright 2022 Ainsley Clark. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package logger

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoOptions struct {
	UseAll bool
	Forget bool
	Expiry time.Duration
}

// NewWithMongoClient creates a new standard logger and sets logging levels
// dependent on environment variables. Upon a log fire, logs will be sent
// to the mongo database that is passed.
//
// Info messages will be sent to the CollectionStdOut collection.
// Errors will be sent to the CollectionStdErr collection.
func NewWithMongoClient(ctx context.Context, collection *mongo.Collection, config MongoOptions) error {
	initialise()

	if config.Forget {
		addIndexes(ctx, collection)
		addHooks(collection, &config)
	}

	addHooks(collection, &config)

	return nil
}

func addIndexes(ctx context.Context, collection *mongo.Collection) {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"expiry.ttl60s": 1},
			Options: options.Index().SetExpireAfterSeconds(60).SetSparse(true),
		},
		{
			Keys:    bson.M{"expiry.ttl5s": 1},
			Options: options.Index().SetExpireAfterSeconds(5).SetSparse(true),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		logger.Debug("Error creating index" + err.Error())
	}
}
