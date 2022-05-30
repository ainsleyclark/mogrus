// Copyright 2022 Ainsley Clark. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mogrus

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// addHoos sends the various log levels to os.Stderr and
// os.Stdout.
func addHooks(collection *mongo.Collection, config *MongoOptions) { //nolint
	//// Send logs with level higher than warning to stderr.
	//L.AddHook(&WriterHook{
	//	Writer: os.Stderr,
	//	LogLevels: []logrus.Level{
	//		logrus.PanicLevel,
	//		logrus.FatalLevel,
	//		logrus.ErrorLevel,
	//		logrus.WarnLevel,
	//	},
	//	Collection:   collection,
	//	MongoOptions: config,
	//})
	//
	//// Send info and debug logs to stdout
	//L.AddHook(&WriterHook{
	//	Writer: os.Stdout,
	//	LogLevels: []logrus.Level{
	//		logrus.TraceLevel,
	//		logrus.InfoLevel,
	//		logrus.DebugLevel,
	//	},
	//	Collection:   collection,
	//	MongoOptions: config,
	//})
}
