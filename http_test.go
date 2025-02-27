// SPDX-License-Identifier: MIT

package jsonrpc

import (
	"net/http/httptest"
	"testing"

	"github.com/issue9/assert/v3"
)

var _ Transport = &httpTransport{}

func TestHTTPConn_ServeHTTP(t *testing.T) {
	a := assert.New(t, false)
	s := initServer(a)
	a.NotNil(s)

	conn := s.NewHTTPConn("", nil)

	srv := httptest.NewServer(conn)
	defer srv.Close()
	conn.url = srv.URL

	a.NotError(conn.Notify("f1", &inType{Age: 18, First: "f", Last: "l"}))
	a.NotError(conn.Notify("not-found", &inType{})) // notify 不返回错误，即使找不到服务

	a.NotError(conn.Send("f1", &inType{Age: 18, First: "f", Last: "l"}, func(out *outType) error {
		a.Equal(out.Age, 18).Equal(out.Name, "fl")
		return nil
	}))

	// 检测抛出错误是否正确
	err := conn.Send("f2", &inType{Age: 18}, func(out *outType) error {
		return nil
	})
	err1, ok := err.(*Error)
	a.True(ok).Equal(err1.Code, CodeInvalidParams) // 由函数 f2 抛出的错误 Error

	// 检测抛出错误是否正确
	err = conn.Send("f3", &inType{Age: 18}, func(out *outType) error {
		return nil
	})
	err1, ok = err.(*Error)
	a.True(ok).Equal(err1.Code, CodeInternalError) // 由函数 f3 抛出的普通错误

	a.Error(conn.Send("not-found", &inType{Age: 18}, func(out *outType) error {
		a.Equal(out.Age, 0)
		return nil
	})) // 不存在的服务名称
}

func TestValidContentType(t *testing.T) {
	a := assert.New(t, false)

	a.NotError(validContentType("application/json"))
	a.NotError(validContentType(""))
	a.NotError(validContentType("application/json;charset=utf-8"))
	a.NotError(validContentType("application/json;;charset=utf-8"))
	a.NotError(validContentType("application/json-rpc;;charset=utf-8"))
	a.NotError(validContentType("application/json-rpc;;charset=UTF-8"))
	a.NotError(validContentType("application/json;charset=utf-8"))
	a.NotError(validContentType("application/jsonrequest;charset=utf-8"))
	a.NotError(validContentType("application/json;"))

	a.Error(validContentType("text/json;"))
	a.Error(validContentType("application/json;charset="))
	a.Error(validContentType("application/json;charset=utf8"))
}
