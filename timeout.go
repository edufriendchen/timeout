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
	"errors"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
)

var (
	// ErrDefaultTimeout Error result returned by default timeout
	ErrDefaultTimeout = errors.New("request timeout")
	// DefaultNormalExitResult Default normal return result
	DefaultNormalExitResult = map[string]interface{}{"message": "normal exit"}
)

// SetTimeoutErr Set the error returned by the timeout
func SetTimeoutErr(err error) {
	ErrDefaultTimeout = err
}

// SetNormalExitResult Set the normal return result
func SetNormalExitResult(m map[string]interface{}) {
	DefaultNormalExitResult = m
}

// Handler defines a function to serve HTTP requests, Note that it is different from the hertz default function.
type Handler = func(ctx context.Context, c *app.RequestContext) error

// New implementation of timeout middleware.
// To use this method, you need to add the following code to the Handler you defined to listen for the context timeout.
// select {
//	case <-ctx.Done():
//		return context.DeadlineExceeded
//	}
func New(h Handler, t time.Duration) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.Next(ctx)
		timeoutContext, cancel := context.WithTimeout(context.Background(), t)
		defer cancel()
		if err := h(timeoutContext, c); err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				// context timeout exits
				c.JSON(http.StatusRequestTimeout, ErrDefaultTimeout.Error())
				return
			} else {
				// Throw out actively
				c.JSON(http.StatusOK, err.Error())
				return
			}
		} else {
			// Normal exit
			c.JSON(http.StatusOK, DefaultNormalExitResult)
			return
		}
	}
}

// Default timeout middleware.
// This method helps the user implement code to listen for context timeouts
// It is important to note, however, it should be noted that this method can only feed back the timeout results in advance when the context timeout occurs.
// The user-defined handler continues to execute because it is on another goroutine and goroutine cannot be forced to end from outside.
func Default(h Handler, t time.Duration) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.Next(ctx)
		done := make(chan error, 1)
		timeoutContext, cancel := context.WithTimeout(context.Background(), t)
		defer cancel()
		go func() {
			done <- h(timeoutContext, c)
		}()
		select {
		case <-timeoutContext.Done():
			// context timeout exits
			c.JSON(http.StatusRequestTimeout, ErrDefaultTimeout.Error())
			return
		case err := <-done:
			if err != nil {
				// Throw out actively
				c.JSON(http.StatusOK, err.Error())
				return
			} else {
				// Normal exit
				c.JSON(http.StatusOK, DefaultNormalExitResult)
				return
			}
		}
	}
}
