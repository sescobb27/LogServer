package main

import (
        "encoding/json"
        "fmt"
        "net"
)

var (
        logger *Logger
)

func init() {
        logger = newLogger()
}

func request_handler(conn net.Conn) {
        buffer := make([]byte, 0, 4096)
        var request map[string]string
        for {
                _, err := conn.Read(buffer)
                if err != nil {
                        conn.Close()
                        fmt.Println(err)
                }
                conn.Close()
                json.Unmarshal(buffer, &request)
                logger.AddLogFile(request["logfile"])
                logger.AsyncWrite(request["logfile"], []byte(request["msg"]))

        }
}

func main() {
        listener, err := net.Listen("tcp", ":9000")
        if err != nil {
                fmt.Println(err)
        }
        defer listener.Close()
        for {
                conn, err := listener.Accept()
                if err != nil {
                        fmt.Println(err)
                }
                go request_handler(conn)
        }
}
