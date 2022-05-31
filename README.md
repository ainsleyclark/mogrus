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
A Go wrapper for Logrus, Errors and Mongo giving you extremely detailed log reports.

## Contributing

Please feel free to make a pull request if you think something should be added to this package!

## Installation

```bash
go get -u github.com/ainsleyclark/mogrus
```

## How to use

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

## Credits

Shout out to the incredible [Maria Letta](https://github.com/MariaLetta) for her excellent Gopher illustrations.

## Licence

Code Copyright 2022 Mogrus. Code released under the [MIT Licence](LICENSE).
