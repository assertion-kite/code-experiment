// +build linux
package main

/*
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

// 调用 C 语言的系统函数，获取当前进程的 ID
int getProcessID() {
    return getpid();
}

// 调用 C 语言的系统函数，执行一个简单的命令
void executeCommand(const char* command) {
    system(command);
}
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L. -lcallee
#include "callee.h"
*/
import "C"

import (
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
)

type A interface {
	Is() bool
}

type B struct{}

func (b *B) Is() bool {
	return true
}

type C struct{}

func (b *C) Is() bool {
	return false
}

type D struct{}

func (b *D) Is() bool {
	return true
}

// 迭代器接口
type Iterator interface {
	HasNext() bool
	Next() A
}

// 自定义迭代器类型
type MyIterator struct {
	data  []A
	index int
}

// 实现迭代器接口的方法
func (it *MyIterator) HasNext() bool {
	return it.index < len(it.data)
}

func (it *MyIterator) Next() A {
	if !it.HasNext() {
		return nil
	}
	d := it.data[it.index]
	it.index++
	return d
}

// mode 0 change to size                  0x0
// FALLOC_FL_KEEP_SIZE                  = 0x1
// FALLOC_FL_PUNCH_HOLE                 = 0x2
//func main() {

	//iterator := &MyIterator{data: []A{&B{}, &C{}, &D{}}}
	//
	//for iterator.HasNext() {
	//	if iterator.Next().Is() {
	//		fmt.Println(iterator.data[iterator.index-1])
	//	}
	//	fmt.Println(iterator.index)
	//}
	//for i := 0; i < 5; i++ {
	//	i := i
	//	go func() {
	//		sl := GetSpinLock("11")
	//		sl.Lock()
	//		defer sl.Unlock()
	//		fmt.Println(i)
	//		time.Sleep(time.Second)
	//	}()
	//}
	//time.Sleep(time.Second * 10)

	C.SayHello()
	// 调用 C 语言的系统函数，获取当前进程的 ID
	pid := C.getProcessID()
	fmt.Printf("Process ID: %d\n", pid)

	// 调用 C 语言的系统函数，执行一个简单的命令
	command := C.CString("ls -l")
	defer C.free(unsafe.Pointer(command))
	C.executeCommand(command)
	fmt.Println("Success!")
	//spew.Sdump(uintptr(1))
	//r := clock.RealClock{}
	//now := r.Now()
	//t := r.NewTimer(time.Second * 5)
	//t.Reset(time.Second * 2)
	//go func() {
	//	select {
	//	case <-t.C():
	//		fmt.Println(r.Since(now))
	//		fmt.Println(3333)
	//		t.Stop()
	//		return
	//	}
	//}()
	//r.Sleep(time.Second * 10)

	//cache := cache2.NewExpiring()
	//
	//if result, ok := cache.Get("foo"); ok || result != nil {
	//	fmt.Printf("expected null, false, got %#v, %v", result, ok)
	//}
	//
	//record1 := "bob"
	//record2 := struct {
	//	Id int64
	//}{
	//	Id: 12223,
	//}
	//
	//// when empty, record is stored
	//cache.Set("foo", record1, time.Hour)
	//if result, ok := cache.Get("foo"); !ok || result != record1 {
	//	fmt.Printf("Expected %#v, true, got %#v, %v", record1, result, ok)
	//}
	//
	//// newer record overrides
	//cache.Set("foo", record2, time.Hour)
	//if result, ok := cache.Get("foo"); !ok || result != record2 {
	//	fmt.Printf("Expected %#v, true, got %#v, %v", record2, result, ok)
	//}
	//
	//time.Sleep(time.Second * 2)
	//
	//// delete the current value
	////cache.Delete("foo")
	//if result, ok := cache.Get("foo"); ok || result != nil {
	//	fmt.Printf("Expected null, false, got %#v, %v", result, ok)
	//}
//}

func MapKeys[Map map[*Key]Value, Key any, Value any](m Map) []Key {
	keys := make([]Key, 0, len(m))
	for key := range m {
		keys = append(keys, *key)
	}
	return keys
}

type SpinLock struct {
	lock     sync.Mutex
	spinLock uint32
}

var lockMap = make(map[string]*SpinLock)
var lockMapMutex sync.Mutex

func GetSpinLock(key string) *SpinLock {
	lockMapMutex.Lock()
	defer lockMapMutex.Unlock()

	if sl, ok := lockMap[key]; ok {
		return sl
	}

	sl := &SpinLock{}
	lockMap[key] = sl
	return sl
}

func (sl *SpinLock) Lock() {
	for {
		if atomic.CompareAndSwapUint32(&sl.spinLock, 0, 1) {
			return
		}
		runtime.Gosched()
	}
}

func (sl *SpinLock) Unlock() {
	if atomic.CompareAndSwapUint32(&sl.spinLock, 1, 0) {
		return
	}

	sl.lock.Unlock()
}
