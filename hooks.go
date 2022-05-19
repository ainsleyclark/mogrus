// Copyright 2022 Ainsley Clark. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mogrus

import (
	"context"
	"github.com/ainsleyclark/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"time"
)

// WriterHook is a hook that writes logs of specified
// LogLevels to specified Writer.
type WriterHook struct {
	// The io.Writer, this can be stdout or stderr.
	Writer io.Writer
	// The slice of log levels the writer can too.
	LogLevels []logrus.Level
	// The collection to send logs too.
	Collection *mongo.Collection

	Options MongoOptions
}

// Fire will be called when some logging function is
// called with current hook. It will format log
// entry to string and write it to
// appropriate writer
func (hook *WriterHook) Fire(entry *logrus.Entry) error {
	const op = "Logger.Hook.Fire"

	line, err := entry.String()
	if err != nil {
		return &errors.Error{Code: errors.INTERNAL, Message: "Error obtaining the entry string", Operation: op, Err: err}
	}

	if hook.Collection != nil {
		err = hook.FireMongo(entry)
		if err != nil {
			return err
		}
	}

	_, err = hook.Writer.Write([]byte(line))
	if err != nil {
		return &errors.Error{Code: errors.INTERNAL, Message: "Error writing entry to io.Writer", Operation: op, Err: err}
	}

	return nil
}

// Levels Define on which log levels this hook would
// trigger.
func (hook *WriterHook) Levels() []logrus.Level {
	return hook.LogLevels
}

// TODO add firehook!

// FireMongo sends the entry time, level and message and any
// other entry fields to the database.
func (hook *WriterHook) FireMongo(entry *logrus.Entry) error {
	const op = "Logger.Hook.FireMongo"

	data := make(logrus.Fields)
	data["level"] = entry.Level.String()
	data["time"] = entry.Time
	data["message"] = entry.Message

	for k, v := range entry.Data {
		if logrus.ErrorKey == k && v != nil {
			err := errors.ToError(v)
			data["code"] = err.Code
			data["message"] = err.Message
			data["operation"] = err.Operation
			data["error"] = err.Err.Error()
			continue
		}
		data[k] = v
	}

	data["expiry"] = map[string]time.Time{}

	if data["level"] == "panic" {
		data["expiry"] = map[string]time.Time{
			"ttl60s": time.Now(),
		}
	} else {
		data["expiry"] = map[string]time.Time{
			"ttl5s": time.Now(),
		}
	}

	_, err := hook.Collection.InsertOne(context.Background(), data)
	if err != nil {
		return &errors.Error{Code: errors.INTERNAL, Message: "Error writing entry to Mongo Collection", Operation: op, Err: err}
	}

	return nil
}
