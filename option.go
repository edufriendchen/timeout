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
	"time"

	"github.com/cloudwego/hertz/pkg/app"
)

// Option is the only struct that can be used to set Options.
type Option struct {
	F func(o *Options)
}

var (
	Handler app.HandlerFunc = nil
	TErr                    = context.DeadlineExceeded
)

const Timing = 3 * time.Second

// Options defines the config for timeout middleware.
type Options struct {
	// TimeoutHandler is the handler after the timeout
	//
	// Optional. Default: nil
	TimeoutHandler app.HandlerFunc

	// Timing is used to set the timeout period
	//
	// Optional. Default:  3 * time.Second
	Timing time.Duration

	// TErr is used to customize the error indicating timeout
	//
	// Optional. Default: context.DeadlineExceeded
	TErr error
}

func (o *Options) Apply(opts []Option) {
	for _, op := range opts {
		op.F(o)
	}
}

// OptionsDefault is the default options.
var OptionsDefault = Options{
	Handler,
	Timing,
	TErr,
}

func NewOptions(opts ...Option) *Options {
	options := &Options{
		OptionsDefault.TimeoutHandler,
		OptionsDefault.Timing,
		OptionsDefault.TErr,
	}
	options.Apply(opts)
	return options
}

// WithTimeoutHandler sets TimeoutHandler.
func WithTimeoutHandler(h app.HandlerFunc) Option {
	return Option{
		F: func(o *Options) {
			o.TimeoutHandler = h
		},
	}
}

// WithTiming sets Timing.
func WithTiming(t time.Duration) Option {
	return Option{
		F: func(o *Options) {
			o.Timing = t
		},
	}
}

// WithTErr sets TErr.
func WithTErr(err error) Option {
	return Option{
		F: func(o *Options) {
			o.TErr = err
		},
	}
}
