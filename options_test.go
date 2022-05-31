// Copyright 2022 Ainsley Clark. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mogrus

import (
	"github.com/ainsleyclark/errors"
	"reflect"
	"testing"
)

func TestOptions_Validate(t *testing.T) {
	tt := map[string]struct {
		input Options
		want  error
	}{
		"Nil Collection": {
			Options{},
			errors.New("nil mongo collection"),
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			got := test.input.Validate()
			if !reflect.DeepEqual(test.want, got) {
				t.Fatalf("expecting %s, got %s", test.want, got)
			}
		})
	}
}
