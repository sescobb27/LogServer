package main

import (
        "net"
        "net/http"
        "net/rpc"
)

var (
        logger *Logger
)

func init() {
        logger = newLogger()
}

func startServer() {
        listener, err := net.Listen("tcp", ":9000")
        assertNoError(err)
        go http.Serve(listener, nil)
}

func main() {
        rpc.Register(logger)
        rpc.HandleHTTP()
        startServer()
}
