<div align="center">
<img height="250" src="res/logo.svg?size=new" alt="Errors Logo" style="margin-bottom: 1rem" />
</div>

[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
[![Go Report Card](https://goreportcard.com/badge/github.com/ainsleyclark/mogrus)](https://goreportcard.com/report/github.com/ainsleyclark/mogrus)
[![Maintainability](https://api.codeclimate.com/v1/badges/b3afd7bf115341995077/maintainability)](https://codeclimate.com/github/ainsleyclark/mogrus/maintainability)
[![Test](https://github.com/ainsleyclark/mogrus/actions/workflows/test.yml/badge.svg?branch=master)](https://github.com/ainsleyclark/mogrus/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/ainsleyclark/mogrus/branch/master/graph/badge.svg?token=K27L8LS7DA)](https://codecov.io/gh/ainsleyclark/mogrus)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/ainsleyclark/mogrus)

# Mogrus

A Go wrapper for Logrus, Errors and Mongo giving you extremely detailed log reports. This package is designed to be used
with [github.com/ainsleyclark/errors](https://github.com/ainsleyclark/errors) for error reporting with codes, messages and more.

## Overview

- ✅ Add hooks to a Mongo collection.
- ✅ Logs with custom errors featuring codes, messages and lifelines.
- ✅ Customisable expiry times for different Logrus levels.
- ✅ Specify a callback function for when an entry is fired to Mongo.

## Why?



## Installation

```bash
go get -u github.com/ainsleyclark/mogrus
```

## How to use

Below is an example of how to use Mogrus, instantiate a Logrus instance and connect to a Mongo DB and pass the hook to
the `AddHook()` function.

### Add hook

```go
func ExampleMogrus() {
	l := logrus.New()

	clientOptions := options.Client().
	ApplyURI(os.Getenv("MONGO_CONNECTION")).
	SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1))

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalln(err)
	}

	opts := mogrus.Options{
		Collection: client.Database("logs").Collection("col"),
		// FireHook is a hook function called just before an
		// entry is logged to Mongo.
		FireHook: func(e mogrus.Entry) {
			fmt.Printf("%+v\n", e)
		},
		// ExpirationLevels allows for the customisation of expiry
		// time for each Logrus level by default entries do not expire.
		ExpirationLevels: mogrus.ExpirationLevels{
			logrus.DebugLevel: time.Second * 5,
			logrus.InfoLevel:  time.Second * 15,
			logrus.ErrorLevel: time.Second * 30,
		},
	}

	// Create the new Mogrus hook, returns an error if the
	// collection is nil.
	hook, err := mogrus.New(context.Background(), opts)
	if err != nil {
		log.Fatalln(err)
	}

	// Add the hook to the Logrus instance.
	l.AddHook(hook)

	l.Debug("Debug level")
	l.WithField("key", "value").Info("Info level")
	l.WithError(errors.NewInternal(errors.New("error"), "message", "op")).Error("Error level")
}
```

### Entry

The definition of an entry below is the object that is fired to Mongo.

```go
// Entry defines a singular entry sent to Mongo
// when a Logrus event is fired.
Entry struct {
	Level   string               `json:"level" bson:"level"`
	Time    time.Time            `json:"time" bson:"time"`
	Message string               `json:"message" bson:"message"`
	Data    map[string]any       `json:"data" bson:"data"`
	Error   *Error               `json:"error" bson:"error"`
	Expiry  map[string]time.Time `json:"expiry" bson:"expiry"`
}
```

### Error

The definition of an error below is the object that is stored in an Entry when `log.WithError` is used.

```go
// Error defines a custom Error for log entries, detailing
// the file line, errors are returned as strings instead
// of the stdlib error.
Error struct {
	Code      string `json:"code" bson:"code"`
	Message   string `json:"message" bson:"message"`
	Operation string `json:"operation" bson:"op"`
	Err       string `json:"error" bson:"err"`
	FileLine  string `json:"file_line" bson:"file_line"`
}
```

## TODO

- Add global expiration time for all log levels.
- Add constructors for:
  - `WithDB()`
  - `WithAuth()`
  - `New()`

## Contributing

Please feel free to make a pull request if you think something should be added to this package!

## Credits

Shout out to the incredible [Maria Letta](https://github.com/MariaLetta) for her excellent Gopher illustrations.

## Licence

Code Copyright 2022 Mogrus. Code released under the [MIT Licence](LICENSE).
