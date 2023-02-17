package main

import (
  "net"
  "fmt"
)

func main() {
  addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:8080")
  if err != nil {
    fmt.Println(err)
    return 
  }

  listener, err := net.ListenTCP("tcp", addr)
  if err != nil {
    fmt.Println(err)
    return 
  }
  
  fmt.Println("Server Listening on port 8080...")
  
  for {
    conn, err := listener.Accept()
    if err != nil {
      fmt.Println(err)
      continue 
    }
    
    go handleConnection(conn)
  }
}

func handleConnection(conn net.Conn) error {
  defer conn.Close()
  fmt.Println("Handling Connection")
  buf := make([]byte, 1024)
  request, err := conn.Read(buf)
  if err != nil {
    fmt.Println(err)
    return err
  }
  
  fmt.Println(string(buf[:request]))
  
  _, err = conn.Write([]byte("Message ackd\n"))
  if err != nil {
    fmt.Println(err)
    return err
  }
  
  return nil
}
