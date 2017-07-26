go-servers
==========
A collection of HTTP servers in Go with a 
specific functionality.

## Server Types
* **ConnectionLimitServer**: Listens and accepts up to 
specific number of connections. A connection is not
accepted if the server is full.

```go
    srv := &servers.ConnectionLimitServer{
      ListenLimit: 5,
      Server: &http.Server{
        Addr: ":1234",
        Handler: mux,
      },
    }
    
    srv.ListenAndServe() // Accepts up to 5 connections
```

## Licence
License

MIT © [Theo Despoudis](https://theodespoudis.com/)

