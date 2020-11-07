// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gutil_test

import (
	"github.com/gogf/gf/frame/g"
	"testing"

	"github.com/gogf/gf/test/gtest"
	"github.com/gogf/gf/util/gutil"
)

func Test_Dump(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		gutil.Dump(map[int]int{
			100: 100,
		})
	})
<<<<<<< HEAD
=======
}

func Test_Try(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := `gutil Try test`
		t.Assert(gutil.Try(func() {
			panic(s)
		}), s)
	})
>>>>>>> 4ae89dc9f62ced2aaf3c7eeb2eaf438c65c1521c
}

func Test_TryCatch(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		gutil.TryCatch(func() {
			panic("gutil TryCatch test")
		})
	})

	gtest.C(t, func(t *gtest.T) {
		gutil.TryCatch(func() {
			panic("gutil TryCatch test")

<<<<<<< HEAD
		}, func(err interface{}) {
=======
		}, func(err error) {
>>>>>>> 4ae89dc9f62ced2aaf3c7eeb2eaf438c65c1521c
			t.Assert(err, "gutil TryCatch test")
		})
	})
}

func Test_IsEmpty(t *testing.T) {
<<<<<<< HEAD

=======
>>>>>>> 4ae89dc9f62ced2aaf3c7eeb2eaf438c65c1521c
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.IsEmpty(1), false)
	})
}

func Test_Throw(t *testing.T) {
<<<<<<< HEAD

=======
>>>>>>> 4ae89dc9f62ced2aaf3c7eeb2eaf438c65c1521c
	gtest.C(t, func(t *gtest.T) {
		defer func() {
			t.Assert(recover(), "gutil Throw test")
		}()

		gutil.Throw("gutil Throw test")
	})
}

func Test_Keys(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		keys := gutil.Keys(map[int]int{
			1: 10,
			2: 20,
		})
		t.AssertIN("1", keys)
		t.AssertIN("2", keys)
	})

	gtest.C(t, func(t *gtest.T) {
		type T struct {
			A string
			B int
		}
		keys := gutil.Keys(new(T))
		t.Assert(keys, g.SliceStr{"A", "B"})
	})

	gtest.C(t, func(t *gtest.T) {
		type T struct {
			A string
			B int
		}
		var pointer *T
		keys := gutil.Keys(pointer)
		t.Assert(keys, g.SliceStr{"A", "B"})
	})
}

func Test_Values(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		values := gutil.Keys(map[int]int{
			1: 10,
			2: 20,
		})
		t.AssertIN("1", values)
		t.AssertIN("2", values)
	})

	gtest.C(t, func(t *gtest.T) {
		type T struct {
			A string
			B int
		}
		keys := gutil.Values(T{
			A: "1",
			B: 2,
		})
		t.Assert(keys, g.Slice{"1", 2})
	})
}
