// SPDX-License-Identifier: MIT

package jsonrpc

import (
	"log"
	"testing"

	"github.com/issue9/assert"
)

// 用于测试的数据类型
type (
	inType struct {
		Last  string `json:"last"`
		First string `json:"first"`
		Age   int
	}

	outType struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
)

var (
	f1 = func(notify bool, params *inType, result *outType) error {
		if notify {
			return nil
		}

		result.Name = params.First + params.Last
		result.Age = params.Age
		return nil
	}
)

func initConn(a *assert.Assertion, errlog *log.Logger) *Conn {
	conn := NewConn(errlog)
	a.NotNil(conn)

	a.True(conn.Register("f1", f1))
	a.False(conn.Register("f1", f1))
	a.False(conn.Register("f1", f1))

	return conn
}

func TestNewHandler(t *testing.T) {
	a := assert.New(t)

	type f func()

	a.Panic(func() {
		newHandler(5)
	})

	// 参数数量不正确
	a.Panic(func() {
		newHandler(func(*int) error { return nil })
	})

	// 参数类型不正确
	a.Panic(func() {
		newHandler(func(bool, *f, *int) error { return nil })
	})

	// 参数类型不正确
	a.Panic(func() {
		newHandler(func(bool, *int, *f) error { return nil })
	})

	// 参数类型不正确
	a.Panic(func() {
		newHandler(func(*bool, *int, *int) error { return nil })
	})

	// 返回值不正确
	a.Panic(func() {
		newHandler(func(bool, *int, *int) int { return 0 })
	})

	// 返回值不正确
	a.NotPanic(func() {
		newHandler(func(bool, *int, *int) *Error { return nil })
	})

	// 正常签名
	a.NotPanic(func() {
		newHandler(func(bool, *int, *int) error { return nil })
	})
}
