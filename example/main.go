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

package main

import (
	"context"
	"net/http"
	"time"

	"github.com/cloudwego/hertz-contrib/timeout"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	h := server.Default()
	defaultTimeout := func(ctx context.Context, c *app.RequestContext) error {
		sleepTime, _ := time.ParseDuration(c.Param("sleepTime") + "ms")
		time.Sleep(sleepTime)
		req, _ := http.NewRequest(http.MethodGet, "https://www.baidu.com", nil)
		req = req.WithContext(ctx)
		client := &http.Client{}
		_, err := client.Do(req)
		if err != nil {
			return err
		}
		return nil
	}
	newTimeout := func(ctx context.Context, c *app.RequestContext) error {
		sleepTime, _ := time.ParseDuration(c.Param("sleepTime") + "ms")
		done := make(chan error, 1)
		go func() {
			time.Sleep(sleepTime)
			req, _ := http.NewRequest(http.MethodGet, "https://www.baidu.com", nil)
			req = req.WithContext(ctx)
			client := &http.Client{}
			_, err := client.Do(req)
			if err != nil {
				done <- err
				return
			}
			done <- nil
		}()
		select {
		case <-ctx.Done():
			return context.DeadlineExceeded
		case err := <-done:
			return err
		}
	}
	h.GET("/default/:sleepTime", timeout.Default(defaultTimeout, 2*time.Second))
	h.GET("/new/:sleepTime", timeout.New(newTimeout, 2*time.Second))
	h.Spin()
}
