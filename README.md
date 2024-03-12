# SGS - Simple Game Server

This is a simple game server written in Golang. You can adopt it to your needs for small multiplayer games.

# How it works

You should provide only two funcs for SGS.

The first is handler function - it is called in separated goroutines on every incoming data packet. You will have netip.AddrPort, io.Reader and io.Writer as arguments to provide your own logic. Here netip.AddrPort is client's address, Reader holds incoming from client request data and Writer for your response to this client if needed.
```go
type Handler interface {  
    ServeUDP(netip.AddrPort, io.Reader, io.Writer)  
}
```

The second is sender function. This function also runs in a separate goroutine and used to send the same data to many clients. SGS sends data to clients in separated goroutines .
```go
type Sender interface {  
    Send() ([]netip.AddrPort, io.Reader)  
}
```

## Sender timeout

At this time you can't call sender function when you want but to provide better game logic you can set timeout for periodic calls.

## Protocols

SGS works only with UDP but we have future plans to provide TCP for non-game data.


```sh
$ go get github.com/oktavarium/sgs
```