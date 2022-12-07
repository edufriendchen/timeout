# Timeout
Timeout middleware for Hertz. A wrapper that acts as the Hertz handler (note: This handler is somewhat different from Hertz's default handler),It creates a context with a timeout and passes it down.
If the context passed executions (eg. DB ops, Http calls) takes longer than the given duration to return, If the set time is exceeded, the default error display will be returned, and this error display can be set.



### Table of Contents
- [Concise](#Concise)

- [Examples](#examples)

  

### Concise

1、Default timeout middleware.This method helps the user implement code to listen for context timeouts It is important to note, however, it should be noted that this method can only feed back the timeout results in advance when the context timeout occurs.The user-defined handler continues to execute because it is on another goroutine and goroutine cannot be forced to end from outside.

```go
func Default(h Handler, t time.Duration) app.HandlerFunc
```

2、New implementation of timeout middleware. To use this method, you need to add the following code to the Handler you defined to listen for the context timeout.
// select {
//	case <-ctx.Done():
//		return context.DeadlineExceeded
//	}

```go
func New(h Handler, t time.Duration) app.HandlerFunc
```

3、SetTimeoutErr Set the error returned by the timeout.

```go
func SetTimeoutErr(err error)
```

4、SetNormalExitResult Set the normal return result.

```go
func SetNormalExitResult(m map[string]interface{})
```



### Examples

**Import the middleware package that is part of the Hertz web framework**

```go
import "github.com/hertz-contrib/timeout"
```

**Default Examples:**

The default method of timeout middleware uses samples.

```go
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
	h.GET("/default/:sleepTime", timeout.Default(defaultTimeout, 2*time.Second))
	h.Spin()
}
```

Test http 200 with curl:
```bash
curl --location -I --request GET 'http://localhost:8888/default/1000' 
```

Test http 408 with curl :
(notes-Hertz retries twice by default by returning a timeout, so you will observe that the total time of this request is 6s.)

```bash
curl --location -I --request GET 'http://localhost:8888/default/4000' 
```



**New Examples:**

New implementation of timeout middleware.To use this method, you need to add the following code to the Handler you defined to listen for the context timeout.

```
select {
	case <-ctx.Done():
		return context.DeadlineExceeded
}
```

The new method of timeout middleware uses samples.

```go
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
	h.GET("/new/:sleepTime", timeout.New(newTimeout, 2*time.Second))
	h.Spin()
}
```



**Set Examples:**

You just need to set it up before using timeout middleware.

```go
package main

import (
	"context"
	"errors"
	"net/http"
	"time"
	
	"github.com/cloudwego/hertz-contrib/timeout"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	h := server.Default()
	// set
	timeout.SetTimeoutErr(errors.New("example"))
	timeout.SetNormalExitResult(map[string]interface{}{"message": "example"})
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
	h.GET("/new/:sleepTime", timeout.New(newTimeout, 2*time.Second))
	h.Spin()
}
```

