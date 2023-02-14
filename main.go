package main

import (
  "net"
  "os"
  "fmt"
)

func main() {
  arguments := os.Args
  if len(arguments) == 1 {
    fmt.Println("Please provide a host:port string to bind to")
  }

  addr, err := net.ResolveTCPAddr("tcp", "localhost:8080")
  if err != nil {
    fmt.Println(err)
    return
  }

  listener, err := net.ListenTCP("tcp", addr)
  if err != nil {
    fmt.Println(err)
    return
  }
  
  for {
    conn, err := listener.Accept()
    if err != nil {
      fmt.Println(err)
      continue
    }
    
    go handleConnection(conn)
  }
}

func handleConnection(conn net.Conn) {
  defer conn.Close()
  buf := make([]byte, 1024)
  request, err := conn.Read(buf)
  if err != nil {
    fmt.Println(err)
    return
  }
  
  fmt.Println(string(buf[:request]))
  
  _, err = conn.Write([]byte("Message ack'd"))
  if err != nil {
    fmt.Println(err)
    return
  }
}
