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
	"net/http"
	"testing"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/test/assert"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

// go test -run TestDefaultOptions
func TestDefaultOptions(t *testing.T) {
	h := server.Default()
	ping := func(ctx context.Context, c *app.RequestContext) {
		sleepTime, _ := time.ParseDuration(c.Param("sleepTime") + "ms")
		if err := sleepWithContext(ctx, sleepTime, context.DeadlineExceeded); err != nil {
			_ = c.Error(err)
			return
		} else {
			c.JSON(http.StatusOK, utils.H{"message": "test"})
			return
		}
	}
	h.Use(New())
	h.GET("/ping/:sleepTime", ping)
	testTimeoutFunc := func(timeoutStr string) {
		w := ut.PerformRequest(h.Engine, "GET", "/ping/"+timeoutStr, nil)
		resp := w.Result()
		assert.DeepEqual(t, http.StatusOK, resp.StatusCode())
		assert.DeepEqual(t, "", string(resp.Body()))
	}
	testNormalFunc := func(timeoutStr string) {
		w := ut.PerformRequest(h.Engine, "GET", "/ping/"+timeoutStr, nil)
		resp := w.Result()
		data := map[string]interface{}{}
		err := json.Unmarshal(resp.Body(), &data)
		assert.DeepEqual(t, http.StatusOK, resp.StatusCode())
		assert.Nil(t, err)
		assert.DeepEqual(t, map[string]interface{}{"message": "test"}, data)
	}
	testTimeoutFunc("4000")
	testNormalFunc("2000")
}

// go test -run TestWithTimeoutHandler
func TestWithTimeoutHandler(t *testing.T) {
	h := server.Default()
	ping := func(ctx context.Context, c *app.RequestContext) {
		sleepTime, _ := time.ParseDuration(c.Param("sleepTime") + "ms")
		if err := sleepWithContext(ctx, sleepTime, context.DeadlineExceeded); err != nil {
			_ = c.Error(err)
			return
		} else {
			c.JSON(http.StatusOK, utils.H{"message": "test"})
			return
		}
	}
	h.Use(New(WithTimeoutHandler(func(c context.Context, ctx *app.RequestContext) {
		ctx.String(http.StatusOK, "request timeout")
	})))
	h.GET("/ping/:sleepTime", ping)
	testTimeoutFunc := func(timeoutStr string) {
		w := ut.PerformRequest(h.Engine, "GET", "/ping/"+timeoutStr, nil)
		resp := w.Result()
		assert.DeepEqual(t, http.StatusOK, resp.StatusCode())
		assert.DeepEqual(t, "request timeout", string(resp.Body()))
	}
	testNormalFunc := func(timeoutStr string) {
		w := ut.PerformRequest(h.Engine, "GET", "/ping/"+timeoutStr, nil)
		resp := w.Result()
		data := map[string]interface{}{}
		err := json.Unmarshal(resp.Body(), &data)
		assert.DeepEqual(t, http.StatusOK, resp.StatusCode())
		assert.Nil(t, err)
		assert.DeepEqual(t, map[string]interface{}{"message": "test"}, data)
	}
	testTimeoutFunc("4000")
	testNormalFunc("2000")
}

// go test -run TestWithTiming
func TestWithTiming(t *testing.T) {
	h := server.Default()
	ping := func(ctx context.Context, c *app.RequestContext) {
		sleepTime, _ := time.ParseDuration(c.Param("sleepTime") + "ms")
		if err := sleepWithContext(ctx, sleepTime, context.DeadlineExceeded); err != nil {
			_ = c.Error(err)
			return
		} else {
			c.JSON(http.StatusOK, utils.H{"message": "test"})
			return
		}
	}
	h.Use(New(WithTiming(2 * time.Second)))
	h.GET("/ping/:sleepTime", ping)
	testTimeoutFunc := func(timeoutStr string) {
		w := ut.PerformRequest(h.Engine, "GET", "/ping/"+timeoutStr, nil)
		resp := w.Result()
		assert.DeepEqual(t, http.StatusOK, resp.StatusCode())
		assert.DeepEqual(t, "", string(resp.Body()))
	}
	testNormalFunc := func(timeoutStr string) {
		w := ut.PerformRequest(h.Engine, "GET", "/ping/"+timeoutStr, nil)
		resp := w.Result()
		data := map[string]interface{}{}
		err := json.Unmarshal(resp.Body(), &data)
		assert.DeepEqual(t, http.StatusOK, resp.StatusCode())
		assert.Nil(t, err)
		assert.DeepEqual(t, map[string]interface{}{"message": "test"}, data)
	}
	testTimeoutFunc("3000")
	testNormalFunc("1000")
}

// go test -run TestWithTErr
func TestWithTErr(t *testing.T) {
	h := server.Default()
	ping := func(ctx context.Context, c *app.RequestContext) {
		sleepTime, _ := time.ParseDuration(c.Param("sleepTime") + "ms")
		if err := sleepWithContext(ctx, sleepTime, errors.New("test")); err != nil {
			_ = c.Error(err)
			return
		} else {
			c.JSON(http.StatusOK, utils.H{"message": "test"})
			return
		}
	}
	h.Use(New(WithTErr(errors.New("test"))))
	h.GET("/ping/:sleepTime", ping)
	testTimeoutFunc := func(timeoutStr string) {
		w := ut.PerformRequest(h.Engine, "GET", "/ping/"+timeoutStr, nil)
		resp := w.Result()
		assert.DeepEqual(t, http.StatusOK, resp.StatusCode())
		assert.DeepEqual(t, "", string(resp.Body()))
	}
	testNormalFunc := func(timeoutStr string) {
		w := ut.PerformRequest(h.Engine, "GET", "/ping/"+timeoutStr, nil)
		resp := w.Result()
		data := map[string]interface{}{}
		err := json.Unmarshal(resp.Body(), &data)
		assert.DeepEqual(t, http.StatusOK, resp.StatusCode())
		assert.Nil(t, err)
		assert.DeepEqual(t, map[string]interface{}{"message": "test"}, data)
	}
	testTimeoutFunc("4000")
	testNormalFunc("2000")
}

// go test -run TestWithAll
func TestWithAll(t *testing.T) {
	h := server.Default()
	testErr := errors.New("test")
	ping := func(ctx context.Context, c *app.RequestContext) {
		sleepTime, _ := time.ParseDuration(c.Param("sleepTime") + "ms")
		if err := sleepWithContext(ctx, sleepTime, testErr); err != nil {
			_ = c.Error(err)
			return
		} else {
			c.JSON(http.StatusOK, utils.H{"message": "test"})
			return
		}
	}
	h.Use(New(WithTErr(testErr), WithTiming(2*time.Second), WithTimeoutHandler(func(c context.Context, ctx *app.RequestContext) {
		ctx.String(http.StatusOK, "request timeout")
	})))
	h.GET("/ping/:sleepTime", ping)
	testTimeoutFunc := func(timeoutStr string) {
		w := ut.PerformRequest(h.Engine, "GET", "/ping/"+timeoutStr, nil)
		resp := w.Result()
		assert.DeepEqual(t, http.StatusOK, resp.StatusCode())
		assert.DeepEqual(t, "request timeout", string(resp.Body()))
	}
	testNormalFunc := func(timeoutStr string) {
		w := ut.PerformRequest(h.Engine, "GET", "/ping/"+timeoutStr, nil)
		resp := w.Result()
		data := map[string]interface{}{}
		err := json.Unmarshal(resp.Body(), &data)
		assert.DeepEqual(t, http.StatusOK, resp.StatusCode())
		assert.Nil(t, err)
		assert.DeepEqual(t, map[string]interface{}{"message": "test"}, data)
	}
	testTimeoutFunc("3000")
	testNormalFunc("1000")
}

func sleepWithContext(ctx context.Context, d time.Duration, tErr error) error {
	timer := time.NewTimer(d)
	select {
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
		return tErr
	case <-timer.C:
	}
	return nil
}
