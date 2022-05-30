// Copyright 2022 Ainsley Clark. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mogrus

import (
	"context"
	"github.com/ainsleyclark/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type (
	// hooker is a hook that writes logs of specified
	// LogLevels to specified Writer.
	hooker struct {
		// The collection to send logs too.
		Collection *mongo.Collection
		// MongoOptions - TODO
		MongoOptions MongoOptions
	}
	// MongoOptions - TODO
	MongoOptions struct {
		Collection *mongo.Collection
		UseAll     bool
		Forget     bool
		Expiry     time.Duration
		// TODO Expiry for each level, with how many seconds
	}
	// Entry - TODD
	Entry struct {
		Level   string         `json:"level" bson:"level"`
		Time    time.Time      `json:"time" bson:"time"`
		Message string         `json:"string" bson:"string"`
		Data    map[string]any `json:"data" bson:"data"`
		Error   *errors.Error  `json:"error" bson:"error"`
		Expiry  time.Time      `json:"expiry" bson:"expiry"`
	}
)

// New creates a new Mogrus hooker.
// Returns errors.INVALID if the collection is nil.
// Returns errors.INTERNAL if the indexes could not be added.
func New(ctx context.Context, collection *mongo.Collection, opts MongoOptions) (*hooker, error) {
	const op = "Mogrus.New"

	if collection == nil {
		return nil, errors.NewInvalid(errors.New("mongo collection nil"), "Mongo collection cannot be nil", op)
	}

	err := addIndexes(ctx, collection)
	if err != nil {
		return nil, errors.NewInternal(err, "Error creating indexes", op)
	}

	return &hooker{
		Collection:   collection,
		MongoOptions: opts,
	}, nil
}

// TODO add firehook!

// Fire sends the entry time, level and message and any
// other entry fields to the database.
// Returns errors.INTERNAL if the entry could not be written.
func (hook *hooker) Fire(entry *logrus.Entry) error {
	const op = "Mogrus.Fire"

	formatted := Entry{
		Level:   entry.Level.String(),
		Time:    entry.Time,
		Message: entry.Message,
		Data:    make(map[string]any),
		Error:   nil,
		Expiry:  time.Time{},
	}

	for k, v := range entry.Data {
		if logrus.ErrorKey == k && v != nil {
			formatted.Error = errors.ToError(v)
			continue
		}
		entry.Data[k] = v
	}

	//data["expiry"] = map[string]time.Time{}
	//
	//if data["level"] == "panic" {
	//	data["expiry"] = map[string]time.Time{
	//		"ttl60s": time.Now(),
	//	}
	//} else {
	//	data["expiry"] = map[string]time.Time{
	//		"ttl5s": time.Now(),
	//	}
	//}

	_, err := hook.Collection.InsertOne(context.Background(), formatted)
	if err != nil {
		return errors.NewInternal(err, "Error writing entry to Mongo Collection", op)
	}

	return nil

}

// Levels define on which log levels this hook would
// trigger.
func (hook *hooker) Levels() []logrus.Level {
	return logrus.AllLevels
}

func addIndexes(ctx context.Context, collection *mongo.Collection) error {
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
		return err
	}

	return nil
}
