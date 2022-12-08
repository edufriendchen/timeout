# Timeout (This is a community driven project)
Timeout middleware for Hertz. A wrapper that acts as the Hertz handler (note: This handler is somewhat different from Hertz's default handler),It creates a context with a timeout and passes it down.
If the context passed executions (eg. DB ops, Http calls) takes longer than the given duration to return, If the set time is exceeded, the default error display will be returned, and this error display can be set.

## Table of Contents
- [Usage](#Usage)
- [Examples](#examples)
- [Options](#Options)

## Usage

**Install**

```go
go get github.com/hertz-contrib/timeout
```

**Import**

```go
import "github.com/hertz-contrib/timeout"
```

## Examples

New implementation of timeout middleware.To use this method, you need to add the following code to the Handler you defined to listen for the context timeout.

```go
select {
case <-ctx.Done():
	_ = c.Error(context.DeadlineExceeded)
	return
}
```

examples

```go
package main

import (
	"context"
	"net/http"
	"time"

	"github.com/cloudwego/hertz-contrib/timeout"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

func main() {
	h := server.Default()
	example := func(ctx context.Context, c *app.RequestContext) {
		sleepTime, _ := time.ParseDuration(c.Param("sleepTime") + "ms")
		done := make(chan error, 1)
		go func() {
			time.Sleep(sleepTime)
			done <- nil
		}()
		select {
		case <-ctx.Done():
			_ = c.Error(context.DeadlineExceeded)
			return
		case err := <-done:
			if err != nil {
				_ = c.Error(err)
				return
			}
			c.JSON(http.StatusOK, utils.H{"info": "success"})
			return
		}
	}
	h.Use(timeout.New(timeout.WithTiming(2*time.Second), timeout.WithTimeoutHandler(func(c context.Context, ctx *app.RequestContext) {
		ctx.String(http.StatusRequestTimeout, "request timeout")
	})))
	h.GET("/ping/:sleepTime", example)
	h.Spin()
}
```

Test http 200 with curl:
```bash
curl --location -I --request GET 'http://localhost:8888/ping/1000' 
```

Test http 408 with curl :
(notes-Hertz retries twice by default by returning a timeout, so you will observe that the total time of this request is 6s.)

```bash
curl --location -I --request GET 'http://localhost:8888/ping/3000' 
```

## Options

You can configure the timeout middleware using different configuration options.


| Option         | Type            | Default                  | Description                                            |
| -------------- | --------------- | ------------------------ | ------------------------------------------------------ |
| TimeoutHandler | app.HandlerFunc | nil                      | TimeoutHandler is the handler after the timeout        |
| Timing         | time.Duration   | 3 * time.Second          | Timing is used to set the timeout period               |
| TErr           | error           | context.DeadlineExceeded | TErr is used to customize the error indicating timeout |
