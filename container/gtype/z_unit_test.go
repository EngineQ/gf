// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtype_test

import (
	"sync"
	"testing"

	"github.com/gogf/gf/container/gtype"
	"github.com/gogf/gf/test/gtest"
)

type Temp struct {
	Name string
	Age  int
}

func Test_Bool(t *testing.T) {
	gtest.Case(t, func() {
		i := gtype.NewBool(true)
		iClone := i.Clone()
		gtest.AssertEQ(iClone.Set(false), true)
		gtest.AssertEQ(iClone.Val(), false)

		i1 := gtype.NewBool(false)
		iClone1 := i1.Clone()
		gtest.AssertEQ(iClone1.Set(true), false)
		gtest.AssertEQ(iClone1.Val(), true)

		//空参测试
		i2 := gtype.NewBool()
		gtest.AssertEQ(i2.Val(), false)
		gtest.AssertEQ(i2.Cas(false, true), true)
		gtest.AssertEQ(i2.Cas(false, true), false)

		// json测试
		i3 := gtype.NewBool(true)
		d, err := i3.MarshalJSON()
		if err != nil {
			gtest.Fatal(err)
		}
		err = i3.UnmarshalJSON(d)
		if err != nil {
			gtest.Fatal(err)
		}
		gtest.AssertEQ(i3.Val(), true)

	})
}

func Test_Byte(t *testing.T) {
	gtest.Case(t, func() {
		var wg sync.WaitGroup
		addTimes := 127
		i := gtype.NewByte(byte(0))
		iClone := i.Clone()
		gtest.AssertEQ(iClone.Set(byte(1)), byte(0))
		gtest.AssertEQ(iClone.Val(), byte(1))
		for index := 0; index < addTimes; index++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				i.Add(1)
			}()
		}
		wg.Wait()
		gtest.AssertEQ(byte(addTimes), i.Val())

		//空参测试
		i1 := gtype.NewByte()
		gtest.AssertEQ(i1.Val(), byte(0))
		gtest.AssertEQ(i1.Cas(byte(0), byte(1)), true)
		gtest.AssertEQ(i1.Cas(byte(0), byte(1)), false)

		// json测试
		i3 := gtype.NewByte(byte(123))
		d, err := i3.MarshalJSON()
		if err != nil {
			gtest.Fatal(err)
		}
		err = i3.UnmarshalJSON(d)
		if err != nil {
			gtest.Fatal(err)
		}
		gtest.AssertEQ(i3.Val(), byte(123))
	})
}

func Test_Bytes(t *testing.T) {
	gtest.Case(t, func() {
		i := gtype.NewBytes([]byte("abc"))
		iClone := i.Clone()
		gtest.AssertEQ(iClone.Set([]byte("123")), []byte("abc"))
		gtest.AssertEQ(iClone.Val(), []byte("123"))

		//空参测试
		i1 := gtype.NewBytes()
		gtest.AssertEQ(i1.Val(), nil)

		// json测试
		i2 := gtype.NewBytes([]byte("123"))
		d, err := i2.MarshalJSON()
		if err != nil {
			gtest.Fatal(err)
		}
		err = i2.UnmarshalJSON(d)
		if err != nil {
			gtest.Fatal(err)
		}
		gtest.AssertEQ(i2.Val(), []byte("123"))
	})
}

func Test_String(t *testing.T) {
	gtest.Case(t, func() {
		i := gtype.NewString("abc")
		iClone := i.Clone()
		gtest.AssertEQ(iClone.Set("123"), "abc")
		gtest.AssertEQ(iClone.Val(), "123")

		//空参测试
		i1 := gtype.NewString()
		gtest.AssertEQ(i1.Val(), "")

		// json测试
		i2 := gtype.NewString("123")
		d, err := i2.MarshalJSON()
		if err != nil {
			gtest.Fatal(err)
		}
		err = i2.UnmarshalJSON(d)
		if err != nil {
			gtest.Fatal(err)
		}
		gtest.AssertEQ(i2.Val(), "123")
	})
}

func Test_Interface(t *testing.T) {
	gtest.Case(t, func() {
		t := Temp{Name: "gf", Age: 18}
		t1 := Temp{Name: "gf", Age: 19}
		i := gtype.New(t)
		iClone := i.Clone()
		gtest.AssertEQ(iClone.Set(t1), t)
		gtest.AssertEQ(iClone.Val().(Temp), t1)

		//空参测试
		i1 := gtype.New()
		gtest.AssertEQ(i1.Val(), nil)
	})
}

func Test_Float32(t *testing.T) {
	gtest.Case(t, func() {
		//var wg sync.WaitGroup
		//addTimes := 100
		i := gtype.NewFloat32(0)
		iClone := i.Clone()
		gtest.AssertEQ(iClone.Set(0.1), float32(0))
		gtest.AssertEQ(iClone.Val(), float32(0.1))
		// for index := 0; index < addTimes; index++ {
		// 	wg.Add(1)
		// 	go func() {
		// 	defer wg.Done()
		// 	i.Add(0.2)
		// 	fmt.Println(i.Val())
		// 	}()
		// }
		// wg.Wait()
		// gtest.AssertEQ(100.0, i.Val())

		//空参测试
		i1 := gtype.NewFloat32()
		gtest.AssertEQ(i1.Val(), float32(0))
		gtest.AssertEQ(i1.Cas(0, 0.1), true)
		gtest.AssertEQ(i1.Val(), float32(0.1))
		gtest.AssertEQ(i1.Cas(0, 0.1), false)
		gtest.AssertEQ(i1.Add(0.1), float32(0.2))

		// json测试
		i2 := gtype.NewFloat32(0.123)
		d, err := i2.MarshalJSON()
		if err != nil {
			gtest.Fatal(err)
		}
		err = i2.UnmarshalJSON(d)
		if err != nil {
			gtest.Fatal(err)
		}
		gtest.AssertEQ(i2.Val(), float32(0.123))
	})
}

func Test_Float64(t *testing.T) {
	gtest.Case(t, func() {
		//var wg sync.WaitGroup
		//addTimes := 100
		i := gtype.NewFloat64(0)
		iClone := i.Clone()
		gtest.AssertEQ(iClone.Set(0.1), float64(0))
		gtest.AssertEQ(iClone.Val(), float64(0.1))
		// for index := 0; index < addTimes; index++ {
		// 	wg.Add(1)
		// 	go func() {
		// 	defer wg.Done()
		// 	i.Add(0.1)
		// 	fmt.Println(i.Val())
		// 	}()
		// }
		// wg.Wait()
		// gtest.AssertEQ(100.0, i.Val())

		//空参测试
		i1 := gtype.NewFloat64()
		gtest.AssertEQ(i1.Val(), float64(0))
		gtest.AssertEQ(i1.Cas(0, 0.1), true)
		gtest.AssertEQ(i1.Val(), float64(0.1))
		gtest.AssertEQ(i1.Cas(0, 0.1), false)
		gtest.AssertEQ(i1.Add(0.1), float64(0.2))

		// json测试
		i2 := gtype.NewFloat64(0.123)
		d, err := i2.MarshalJSON()
		if err != nil {
			gtest.Fatal(err)
		}
		err = i2.UnmarshalJSON(d)
		if err != nil {
			gtest.Fatal(err)
		}
		gtest.AssertEQ(i2.Val(), 0.123)
	})
}

func Test_Int(t *testing.T) {
	gtest.Case(t, func() {
		var wg sync.WaitGroup
		addTimes := 1000
		i := gtype.NewInt(0)
		iClone := i.Clone()
		gtest.AssertEQ(iClone.Set(1), 0)
		gtest.AssertEQ(iClone.Val(), 1)
		for index := 0; index < addTimes; index++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				i.Add(1)
			}()
		}
		wg.Wait()
		gtest.AssertEQ(addTimes, i.Val())

		//空参测试
		i1 := gtype.NewInt()
		gtest.AssertEQ(i1.Val(), 0)
		gtest.AssertEQ(i1.Cas(0, 1), true)
		gtest.AssertEQ(i1.Val(), 1)
		gtest.AssertEQ(i1.Cas(0, 1), false)

		// json测试
		i2 := gtype.NewInt(123)
		d, err := i2.MarshalJSON()
		if err != nil {
			gtest.Fatal(err)
		}
		err = i2.UnmarshalJSON(d)
		if err != nil {
			gtest.Fatal(err)
		}
		gtest.AssertEQ(i2.Val(), 123)
	})
}

func Test_Int32(t *testing.T) {
	gtest.Case(t, func() {
		var wg sync.WaitGroup
		addTimes := 1000
		i := gtype.NewInt32(0)
		iClone := i.Clone()
		gtest.AssertEQ(iClone.Set(1), int32(0))
		gtest.AssertEQ(iClone.Val(), int32(1))
		for index := 0; index < addTimes; index++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				i.Add(1)
			}()
		}
		wg.Wait()
		gtest.AssertEQ(int32(addTimes), i.Val())

		//空参测试
		i1 := gtype.NewInt32()
		gtest.AssertEQ(i1.Val(), int32(0))

		// json测试
		i2 := gtype.NewInt32(123)
		d, err := i2.MarshalJSON()
		if err != nil {
			gtest.Fatal(err)
		}
		err = i2.UnmarshalJSON(d)
		if err != nil {
			gtest.Fatal(err)
		}
		gtest.AssertEQ(i2.Val(), int32(123))
	})
}

func Test_Int64(t *testing.T) {
	gtest.Case(t, func() {
		var wg sync.WaitGroup
		addTimes := 1000
		i := gtype.NewInt64(0)
		iClone := i.Clone()
		gtest.AssertEQ(iClone.Set(1), int64(0))
		gtest.AssertEQ(iClone.Val(), int64(1))
		for index := 0; index < addTimes; index++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				i.Add(1)
			}()
		}
		wg.Wait()
		gtest.AssertEQ(int64(addTimes), i.Val())

		//空参测试
		i1 := gtype.NewInt64()
		gtest.AssertEQ(i1.Val(), int64(0))

		// json测试
		i2 := gtype.NewInt64(123)
		d, err := i2.MarshalJSON()
		if err != nil {
			gtest.Fatal(err)
		}
		err = i2.UnmarshalJSON(d)
		if err != nil {
			gtest.Fatal(err)
		}
		gtest.AssertEQ(i2.Val(), int64(123))
	})
}

func Test_Uint(t *testing.T) {
	gtest.Case(t, func() {
		var wg sync.WaitGroup
		addTimes := 1000
		i := gtype.NewUint(0)
		iClone := i.Clone()
		gtest.AssertEQ(iClone.Set(1), uint(0))
		gtest.AssertEQ(iClone.Val(), uint(1))
		for index := 0; index < addTimes; index++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				i.Add(1)
			}()
		}
		wg.Wait()
		gtest.AssertEQ(uint(addTimes), i.Val())

		//空参测试
		i1 := gtype.NewUint()
		gtest.AssertEQ(i1.Val(), uint(0))

		// json测试
		i2 := gtype.NewUint(123)
		d, err := i2.MarshalJSON()
		if err != nil {
			gtest.Fatal(err)
		}
		err = i2.UnmarshalJSON(d)
		if err != nil {
			gtest.Fatal(err)
		}
		gtest.AssertEQ(i2.Val(), uint(123))
	})
}

func Test_Uint32(t *testing.T) {
	gtest.Case(t, func() {
		var wg sync.WaitGroup
		addTimes := 1000
		i := gtype.NewUint32(0)
		iClone := i.Clone()
		gtest.AssertEQ(iClone.Set(1), uint32(0))
		gtest.AssertEQ(iClone.Val(), uint32(1))
		for index := 0; index < addTimes; index++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				i.Add(1)
			}()
		}
		wg.Wait()
		gtest.AssertEQ(uint32(addTimes), i.Val())

		//空参测试
		i1 := gtype.NewUint32()
		gtest.AssertEQ(i1.Val(), uint32(0))

		// json测试
		i2 := gtype.NewUint32(123)
		d, err := i2.MarshalJSON()
		if err != nil {
			gtest.Fatal(err)
		}
		err = i2.UnmarshalJSON(d)
		if err != nil {
			gtest.Fatal(err)
		}
		gtest.AssertEQ(i2.Val(), uint32(123))
	})
}

func Test_Uint64(t *testing.T) {
	gtest.Case(t, func() {
		var wg sync.WaitGroup
		addTimes := 1000
		i := gtype.NewUint64(0)
		iClone := i.Clone()
		gtest.AssertEQ(iClone.Set(1), uint64(0))
		gtest.AssertEQ(iClone.Val(), uint64(1))
		for index := 0; index < addTimes; index++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				i.Add(1)
			}()
		}
		wg.Wait()
		gtest.AssertEQ(uint64(addTimes), i.Val())

		//空参测试
		i1 := gtype.NewUint64()
		gtest.AssertEQ(i1.Val(), uint64(0))

		// json测试
		i2 := gtype.NewUint64(123)
		d, err := i2.MarshalJSON()
		if err != nil {
			gtest.Fatal(err)
		}
		err = i2.UnmarshalJSON(d)
		if err != nil {
			gtest.Fatal(err)
		}
		gtest.AssertEQ(i2.Val(), uint64(123))
	})
}
