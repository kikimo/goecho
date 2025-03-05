// server.go
package main

import (
    "log"
    "net"
    "net/rpc"
    "net/rpc/jsonrpc"
)

// Args holds arguments passed to the RPC method
type Args struct {
    A, B int
}

// Arith provides an RPC service
type Arith int

// Add is an RPC method that adds two integers
func (t *Arith) Add(args *Args, reply *int) error {
    *reply = args.A + args.B
    return nil
}

func main() {
    arith := new(Arith)
    rpc.Register(arith)

    listener, err := net.Listen("tcp", ":1234")
    if err != nil {
        log.Fatal("Listen error:", err)
    }
    defer listener.Close()

    log.Println("Serving RPC handler")
    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Fatal("Accept error:", err)
        }
        go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
    }
}

