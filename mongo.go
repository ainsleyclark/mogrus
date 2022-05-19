// Copyright 2022 Ainsley Clark. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mogrus

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoOptions struct {
	Collection *mongo.Collection
	UseAll     bool
	Forget     bool
	Expiry     time.Duration
	// TODO Expiry for each level, with how many seconds
}

// NewWithMongoClient creates a new standard logger and sets logging levels
// dependent on environment variables. Upon a log fire, logs will be sent
// to the mongo database that is passed.
func NewWithMongoClient(ctx context.Context, opts Options, config MongoOptions) error {
	opts.setDefaults()
	initialise(opts)

	if config.Forget {
		addIndexes(ctx, config.Collection)
		addHooks(config.Collection, &config)
	}

	addHooks(config.Collection, &config)

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
