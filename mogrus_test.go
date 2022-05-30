// Copyright 2022 Ainsley Clark. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mogrus

import (
	"context"
	"github.com/ainsleyclark/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"reflect"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	m := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock).CreateClient(true))

	col := m.CreateCollection(mtest.Collection{
		Name:       "test",
		DB:         "db",
		Client:     m.Client,
		Opts:       nil,
		CreateOpts: nil,
	}, false)

	tt := map[string]struct {
		input Options
		want  any
	}{
		"Success": {
			Options{Collection: &mongo.Collection{}},
			nil,
		},
		"Validation Failed": {
			Options{},
			"Error validating Options",
		},
		"Index Error": {
			Options{
				Collection: col,
				ExpirationLevels: ExpirationLevels{
					logrus.PanicLevel: LevelIndex{
						Expire:   false,
						Duration: time.Second * 10,
					},
				},
			},
			"Error creating indexes",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			_, err := New(context.TODO(), test.input)
			if err != nil {
				msg := errors.Message(err)
				if !reflect.DeepEqual(test.want, msg) {
					t.Fatalf("expecting %s, got %s", test.want, msg)
				}
			}
		})
	}
}
