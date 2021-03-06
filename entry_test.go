// Copyright 2022 Ainsley Clark. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mogrus

import (
	"github.com/ainsleyclark/errors"
	"reflect"
	"testing"
)

func TestEntry_HasError(t *testing.T) {
	tt := map[string]struct {
		input Entry
		want  bool
	}{
		"False": {
			Entry{Error: nil},
			false,
		},
		"True": {
			Entry{Error: &Error{Err: "error"}},
			true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			got := test.input.HasError()
			if !reflect.DeepEqual(test.want, got) {
				t.Fatalf("expecting %t, got %t", test.want, got)
			}
		})
	}
}

func TestError_Error(t *testing.T) {
	tt := map[string]struct {
		input Error
		want  error
	}{
		"Nil": {
			Error{},
			nil,
		},
		"True": {
			Error{Err: "error"},
			errors.New("error"),
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			got := test.input.Error()
			if !reflect.DeepEqual(test.want, got) {
				t.Fatalf("expecting %t, got %t", test.want, got)
			}
		})
	}
}
