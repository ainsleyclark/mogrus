// Copyright 2022 Ainsley Clark. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package logger

import (
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"io/ioutil"
	"os"
)

var (
	// logger is an alias for the standard logger.
	logger = logrus.New()
)

// New creates a new standard logger and sets logging levels
// dependent on environment variables.
func New() {
	initialise()
	addHooks(nil, nil)
}

// Trace logs a trace message with args.
func Trace(args ...any) {
	logger.Trace(args...)
}

// Debug logs a debug message with args.
func Debug(args ...any) {
	logger.Debug(args...)
}

// Info logs ab info message with args.
func Info(args ...any) {
	logger.Info(args...)
}

// Warn logs a warn message with args.
func Warn(args ...any) {
	logger.Warn(args...)
}

// Error logs an error message with args.
func Error(args ...any) {
	logger.Error(args...)
}

// Fatal logs a fatal message with args.
func Fatal(args ...any) {
	logger.Fatal(args...)
}

// Panic logs a panic message with args.
func Panic(args ...any) {
	logger.Panic(args...)
}

// WithField logs with field, sets a new map containing
// "fields".
func WithField(key string, value any) *logrus.Entry {
	return logger.WithFields(logrus.Fields{"fields": logrus.Fields{
		key: value,
	}})
}

// WithFields logs with fields, sets a new map containing
// "fields".
func WithFields(fields logrus.Fields) *logrus.Entry {
	return logger.WithFields(logrus.Fields{"fields": fields})
}

// WithError - Logs with a custom error.
func WithError(err any) *logrus.Entry {
	return logger.WithField(logrus.ErrorKey, err)
}

// SetOutput sets the output of the logger to an io.Writer,
// useful for testing.
func SetOutput(writer io.Writer) {
	logger.SetOutput(writer)
}

// SetLevel sets the level of the logger.
func SetLevel(level logrus.Level) {
	logger.SetLevel(level)
}

// SetLogger sets the application logger.
func SetLogger(l *logrus.Logger) {
	logger = l
}

// initialise sets the standard log level, sets the
// log formatter and discards the stdout.
func initialise() {
	logger.SetLevel(logrus.TraceLevel)

	logger.SetFormatter(&Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		Colours:         true,
	})

	// Send all logs to nowhere by default.
	logger.SetOutput(ioutil.Discard)
}

// addHoos sends the various log levels to os.Stderr and
// os.Stdout.
func addHooks(collection *mongo.Collection, config *MongoOptions) { //nolint
	// Send logs with level higher than warning to stderr.
	logger.AddHook(&WriterHook{
		Writer: os.Stderr,
		LogLevels: []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		},
		Collection: collection,
	})

	// Send info and debug logs to stdout
	logger.AddHook(&WriterHook{
		Writer: os.Stdout,
		LogLevels: []logrus.Level{
			logrus.TraceLevel,
			logrus.InfoLevel,
			logrus.DebugLevel,
		},
		Collection: collection,
	})
}
