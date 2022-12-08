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
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/test/assert"
)

func TestDefaultOption(t *testing.T) {
	opts := NewOptions()
	assert.DeepEqual(t, Handler, opts.TimeoutHandler)
	assert.DeepEqual(t, Timing, opts.Timing)
	assert.DeepEqual(t, TErr, opts.TErr)
}

func TestNewOption(t *testing.T) {
	TimeoutHandler := func(ctx context.Context, c *app.RequestContext) {
		c.String(http.StatusOK, "request timeout")
	}
	testErr := errors.New("test")
	opts := NewOptions(
		WithTimeoutHandler(TimeoutHandler),
		WithTiming(2*time.Second),
		WithTErr(testErr),
	)
	assert.DeepEqual(t, fmt.Sprintf("%p", TimeoutHandler), fmt.Sprintf("%p", opts.TimeoutHandler))
	assert.DeepEqual(t, 2*time.Second, opts.Timing)
	assert.DeepEqual(t, testErr, opts.TErr)
}
