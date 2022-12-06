// The MIT License (MIT)
//
// Copyright (c) 2022 Friend Chen
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//
// This file may have been modified by CloudWeGo authors. All CloudWeGo
// Modifications are Copyright 2022 CloudWeGo Authors.

package timeout

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/test/assert"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"net/http"
	"testing"
	"time"
)

// go test -run TestTimeoutNew
func TestTimeoutNew(t *testing.T) {
	h := server.Default()
	ping := New(func(ctx context.Context, c *app.RequestContext) error {
		sleepTime, _ := time.ParseDuration(c.Param("sleepTime") + "ms")
		if err := sleepWithContext(ctx, sleepTime); err != nil {
			return err
		}
		return nil
	}, 100*time.Millisecond)
	throw := New(func(ctx context.Context, c *app.RequestContext) error {
		return errors.New("throw")
	}, 100*time.Millisecond)
	h.GET("/ping/:sleepTime", ping)
	h.GET("/throw", throw)
	testTimeoutFunc := func(timeoutStr string) {
		w := ut.PerformRequest(h.Engine, "GET", "/ping/"+timeoutStr, nil)
		resp := w.Result()
		assert.DeepEqual(t, http.StatusRequestTimeout, resp.StatusCode())
		assert.DeepEqual(t, "\""+ErrDefaultTimeout.Error()+"\"", string(resp.Body()))
	}
	testNormalFunc := func(timeoutStr string) {
		w := ut.PerformRequest(h.Engine, "GET", "/ping/"+timeoutStr, nil)
		resp := w.Result()
		data := map[string]interface{}{}
		err := json.Unmarshal(resp.Body(), &data)
		assert.DeepEqual(t, http.StatusOK, resp.StatusCode())
		assert.Nil(t, err)
		assert.DeepEqual(t, DefaultNormalExitResult, data)
	}
	testThrowFunc := func(timeoutStr string) {
		w := ut.PerformRequest(h.Engine, "GET", "/throw", nil)
		resp := w.Result()
		assert.DeepEqual(t, http.StatusOK, resp.StatusCode())
		assert.NotNil(t, resp.Body())
	}
	testTimeoutFunc("800")
	testTimeoutFunc("500")
	testNormalFunc("50")
	testNormalFunc("30")
	testThrowFunc("50")
	testThrowFunc("30")
}

// go test -run TestTimeoutNew_Set
func TestTimeoutNew_Set(t *testing.T) {
	h := server.Default()
	// set
	SetTimeoutErr(errors.New("test"))
	SetNormalExitResult(map[string]interface{}{"message": "test"})
	ping := New(func(ctx context.Context, c *app.RequestContext) error {
		sleepTime, _ := time.ParseDuration(c.Param("sleepTime") + "ms")
		if err := sleepWithContext(ctx, sleepTime); err != nil {
			return err
		}
		return nil
	}, 100*time.Millisecond)
	throw := New(func(ctx context.Context, c *app.RequestContext) error {
		return errors.New("throw")
	}, 100*time.Millisecond)
	h.GET("/ping/:sleepTime", ping)
	h.GET("/throw", throw)
	testTimeoutFunc := func(timeoutStr string) {
		w := ut.PerformRequest(h.Engine, "GET", "/ping/"+timeoutStr, nil)
		resp := w.Result()
		assert.DeepEqual(t, http.StatusRequestTimeout, resp.StatusCode())
		assert.DeepEqual(t, "\""+ErrDefaultTimeout.Error()+"\"", string(resp.Body()))
	}
	testNormalFunc := func(timeoutStr string) {
		w := ut.PerformRequest(h.Engine, "GET", "/ping/"+timeoutStr, nil)
		resp := w.Result()
		data := map[string]interface{}{}
		err := json.Unmarshal(resp.Body(), &data)
		assert.DeepEqual(t, http.StatusOK, resp.StatusCode())
		assert.Nil(t, err)
		assert.DeepEqual(t, DefaultNormalExitResult, data)
	}
	testThrowFunc := func(timeoutStr string) {
		w := ut.PerformRequest(h.Engine, "GET", "/throw", nil)
		resp := w.Result()
		assert.DeepEqual(t, http.StatusOK, resp.StatusCode())
		assert.NotNil(t, resp.Body())
	}
	testTimeoutFunc("800")
	testTimeoutFunc("500")
	testNormalFunc("50")
	testNormalFunc("30")
	testThrowFunc("50")
	testThrowFunc("30")
}

// go test -run TestTimeoutDefault
func TestTimeoutDefault(t *testing.T) {
	h := server.Default()
	ping := Default(func(ctx context.Context, c *app.RequestContext) error {
		sleepTime, _ := time.ParseDuration(c.Param("sleepTime") + "ms")
		if err := sleepWithContext(ctx, sleepTime); err != nil {
			return err
		}
		return nil
	}, 100*time.Millisecond)
	throw := Default(func(ctx context.Context, c *app.RequestContext) error {
		return errors.New("throw")
	}, 100*time.Millisecond)
	h.GET("/ping/:sleepTime", ping)
	h.GET("/throw", throw)
	testTimeoutFunc := func(timeoutStr string) {
		w := ut.PerformRequest(h.Engine, "GET", "/ping/"+timeoutStr, nil)
		resp := w.Result()
		assert.DeepEqual(t, http.StatusRequestTimeout, resp.StatusCode())
		assert.DeepEqual(t, "\""+ErrDefaultTimeout.Error()+"\"", string(resp.Body()))
	}
	testNormalFunc := func(timeoutStr string) {
		w := ut.PerformRequest(h.Engine, "GET", "/ping/"+timeoutStr, nil)
		resp := w.Result()
		data := map[string]interface{}{}
		err := json.Unmarshal(resp.Body(), &data)
		assert.DeepEqual(t, http.StatusOK, resp.StatusCode())
		assert.Nil(t, err)
		assert.DeepEqual(t, DefaultNormalExitResult, data)
	}
	testThrowFunc := func(timeoutStr string) {
		w := ut.PerformRequest(h.Engine, "GET", "/throw", nil)
		resp := w.Result()
		assert.DeepEqual(t, http.StatusOK, resp.StatusCode())
		assert.NotNil(t, resp.Body())
	}
	testTimeoutFunc("800")
	testTimeoutFunc("500")
	testNormalFunc("50")
	testNormalFunc("30")
	testThrowFunc("50")
	testThrowFunc("30")
}

// go test -run TestTimeoutDefault_Set
func TestTimeoutDefault_Set(t *testing.T) {
	h := server.Default()
	// set
	SetTimeoutErr(errors.New("test"))
	SetNormalExitResult(map[string]interface{}{"message": "test"})
	ping := Default(func(ctx context.Context, c *app.RequestContext) error {
		sleepTime, _ := time.ParseDuration(c.Param("sleepTime") + "ms")
		if err := sleepWithContext(ctx, sleepTime); err != nil {
			return err
		}
		return nil
	}, 100*time.Millisecond)
	throw := Default(func(ctx context.Context, c *app.RequestContext) error {
		return errors.New("throw")
	}, 100*time.Millisecond)
	h.GET("/ping/:sleepTime", ping)
	h.GET("/throw", throw)
	testTimeoutFunc := func(timeoutStr string) {
		w := ut.PerformRequest(h.Engine, "GET", "/ping/"+timeoutStr, nil)
		resp := w.Result()
		assert.DeepEqual(t, http.StatusRequestTimeout, resp.StatusCode())
		assert.DeepEqual(t, "\""+ErrDefaultTimeout.Error()+"\"", string(resp.Body()))
	}
	testNormalFunc := func(timeoutStr string) {
		w := ut.PerformRequest(h.Engine, "GET", "/ping/"+timeoutStr, nil)
		resp := w.Result()
		data := map[string]interface{}{}
		err := json.Unmarshal(resp.Body(), &data)
		assert.DeepEqual(t, http.StatusOK, resp.StatusCode())
		assert.Nil(t, err)
		assert.DeepEqual(t, DefaultNormalExitResult, data)
	}
	testThrowFunc := func(timeoutStr string) {
		w := ut.PerformRequest(h.Engine, "GET", "/throw", nil)
		resp := w.Result()
		assert.DeepEqual(t, http.StatusOK, resp.StatusCode())
		assert.NotNil(t, resp.Body())
	}
	testTimeoutFunc("800")
	testTimeoutFunc("500")
	testNormalFunc("50")
	testNormalFunc("30")
	testThrowFunc("50")
	testThrowFunc("30")
}

func sleepWithContext(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	select {
	case <-ctx.Done():
		return context.DeadlineExceeded
	case <-timer.C:
	}
	return nil
}
