// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmlock_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/g/container/garray"
	"github.com/gogf/gf/g/os/gmlock"
	"github.com/gogf/gf/g/test/gtest"
)

func Test_Locker_RLock(t *testing.T) {
	//RLock before Lock
	gtest.Case(t, func() {
		key := "testRLockBeforeLock"
		array := garray.New()
		go func() {
			gmlock.RLock(key)
			array.Append(1)
			time.Sleep(50 * time.Millisecond)
			array.Append(1)
			gmlock.RUnlock(key)
		}()
		go func() {
			time.Sleep(10 * time.Millisecond)
			gmlock.Lock(key)
			array.Append(1)
			gmlock.Unlock(key)
		}()
		time.Sleep(20 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(80 * time.Millisecond)
		gtest.Assert(array.Len(), 3)
	})

	//Lock before RLock
	gtest.Case(t, func() {
		key := "testLockBeforeRLock"
		array := garray.New()
		go func() {
			gmlock.Lock(key)
			array.Append(1)
			time.Sleep(50 * time.Millisecond)
			gmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(10 * time.Millisecond)
			gmlock.RLock(key)
			array.Append(1)
			gmlock.RUnlock(key)
		}()
		time.Sleep(20 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(50 * time.Millisecond)
		gtest.Assert(array.Len(), 2)
	})

	//Lock before RLocks
	gtest.Case(t, func() {
		key := "testLockBeforeRLocks"
		array := garray.New()
		go func() {
			gmlock.Lock(key)
			array.Append(1)
			time.Sleep(50 * time.Millisecond)
			gmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(10 * time.Millisecond)
			gmlock.RLock(key)
			array.Append(1)
			time.Sleep(70 * time.Millisecond)
			gmlock.RUnlock(key)
		}()
		go func() {
			time.Sleep(10 * time.Millisecond)
			gmlock.RLock(key)
			array.Append(1)
			time.Sleep(70 * time.Millisecond)
			gmlock.RUnlock(key)
		}()
		time.Sleep(20 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(120 * time.Millisecond)
		gtest.Assert(array.Len(), 3)
	})
}

func Test_Locker_TryRLock(t *testing.T) {
	//Lock before TryRLock
	gtest.Case(t, func() {
		key := "testLockBeforeTryRLock"
		array := garray.New()
		go func() {
			gmlock.Lock(key)
			array.Append(1)
			time.Sleep(50 * time.Millisecond)
			gmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(10 * time.Millisecond)
			if gmlock.TryRLock(key) {
				array.Append(1)
				gmlock.RUnlock(key)
			}
		}()
		time.Sleep(20 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(50 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
	})

	//Lock before TryRLocks
	gtest.Case(t, func() {
		key := "testLockBeforeTryRLocks"
		array := garray.New()
		go func() {
			gmlock.Lock(key)
			array.Append(1)
			time.Sleep(50 * time.Millisecond)
			gmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(10 * time.Millisecond)
			if gmlock.TryRLock(key) {
				array.Append(1)
				gmlock.RUnlock(key)
			}
		}()
		go func() {
			time.Sleep(70 * time.Millisecond)
			if gmlock.TryRLock(key) {
				array.Append(1)
				gmlock.RUnlock(key)
			}
		}()
		time.Sleep(20 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(80 * time.Millisecond)
		gtest.Assert(array.Len(), 2)
	})
}

func Test_Locker_RLockFunc1(t *testing.T) {
	//RLockFunc before Lock
	gtest.Case(t, func() {
		key := "testRLockFuncBeforeLock"
		array := garray.New()
		go func() {
			gmlock.RLockFunc(key, func() {
				array.Append(1)
				time.Sleep(500 * time.Millisecond)
				array.Append(1)
			})
		}()
		go func() {
			time.Sleep(10 * time.Millisecond)
			gmlock.Lock(key)
			array.Append(1)
			gmlock.Unlock(key)
		}()
		time.Sleep(200 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(800 * time.Millisecond)
		gtest.Assert(array.Len(), 3)
	})

	//Lock before RLockFunc
	gtest.Case(t, func() {
		key := "testLockBeforeRLockFunc"
		array := garray.New()
		go func() {
			gmlock.Lock(key)
			array.Append(1)
			time.Sleep(50 * time.Millisecond)
			gmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(10 * time.Millisecond)
			gmlock.RLockFunc(key, func() {
				array.Append(1)
			})
		}()
		time.Sleep(20 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(50 * time.Millisecond)
		gtest.Assert(array.Len(), 2)
	})

}

func Test_Locker_RLockFunc2(t *testing.T) {
	//Lock before RLockFuncs
	gtest.Case(t, func() {
		key := "testLockBeforeRLockFuncs"
		array := garray.New()
		go func() {
			gmlock.Lock(key)
			array.Append(1)
			//glog.Println("add1")
			time.Sleep(500 * time.Millisecond)
			gmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			gmlock.RLockFunc(key, func() {
				array.Append(1)
				//glog.Println("add2")
				time.Sleep(700 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			gmlock.RLockFunc(key, func() {
				array.Append(1)
				//glog.Println("add3")
				time.Sleep(700 * time.Millisecond)
			})
		}()
		time.Sleep(200 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(700 * time.Millisecond)
		gtest.Assert(array.Len(), 3)
	})
}

func Test_Locker_TryRLockFunc(t *testing.T) {
	//Lock before TryRLockFunc
	gtest.Case(t, func() {
		key := "testLockBeforeTryRLockFunc"
		array := garray.New()
		go func() {
			gmlock.Lock(key)
			array.Append(1)
			time.Sleep(50 * time.Millisecond)
			gmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(10 * time.Millisecond)
			gmlock.TryRLockFunc(key, func() {
				array.Append(1)
			})
		}()
		time.Sleep(20 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(50 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
	})

	//Lock before TryRLockFuncs
	gtest.Case(t, func() {
		key := "testLockBeforeTryRLockFuncs"
		array := garray.New()
		go func() {
			gmlock.Lock(key)
			array.Append(1)
			time.Sleep(50 * time.Millisecond)
			gmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(10 * time.Millisecond)
			gmlock.TryRLockFunc(key, func() {
				array.Append(1)
			})
		}()
		go func() {
			time.Sleep(70 * time.Millisecond)
			gmlock.TryRLockFunc(key, func() {
				array.Append(1)
			})
		}()
		time.Sleep(20 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(70 * time.Millisecond)
		gtest.Assert(array.Len(), 2)
	})
}
