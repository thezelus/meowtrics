graceful [![GoDoc](https://godoc.org/github.com/stretchr/graceful?status.png)](http://godoc.org/github.com/stretchr/graceful) [![wercker status](https://app.wercker.com/status/2729ba763abf87695a17547e0f7af4a4/s "wercker status")](https://app.wercker.com/project/bykey/2729ba763abf87695a17547e0f7af4a4)
========

Graceful is a Go package enabling graceful shutdown of http.Handler servers.

## Usage

Usage of Graceful is simple. Create your http.Handler and pass it to the `Run` function:

```go

import (
  "github.com/stretchr/graceful"
  "net/http"
  "fmt"
)

func main() {
  mux := http.NewServeMux()
  mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, "Welcome to the home page!")
  })

  graceful.Run(":3001",10*time.Second,mux)
}
```

 Another example, using [Negroni](https://github.com/codegangsta/negroni), functions in much the same manner:

```go
package main

import (
  "github.com/codegangsta/negroni"
  "github.com/stretchr/graceful"
  "net/http"
  "fmt"
)

func main() {
  mux := http.NewServeMux()
  mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, "Welcome to the home page!")
  })

  n := negroni.Classic()
  n.UseHandler(mux)
  //n.Run(":3000")
  graceful.Run(":3001",10*time.Second,n)
}
```



When Graceful is sent a SIGINT (ctrl+c), it:

1. Disables keepalive connections.
2. Closes the listening socket, allowing another process to listen on that port immediately.
3. Starts a timer of `timeout` duration to give active requests a chance to finish.
4. When timeout expires, closes all active connections.
5. Returns from the function, allowing the server to terminate.

## Notes

If the `timeout` argument to `Run` is 0, the server never times out, allowing all active requests to complete.

Graceful relies on functionality in [Go 1.3](http://tip.golang.org/doc/go1.3) which has not yet been released. If you wish to use it, you
must [install the beta](https://code.google.com/p/go/wiki/Downloads) of Go 1.3. Once 1.3 is released, this note will be removed.

