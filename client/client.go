package main

import (
	"bufio"
	"fmt"
	"net"
  "io"
)

func main() {
	fmt.Println("Creating Client connection to server")
  conn, err := net.Dial("tcp", "0.0.0.0:8080")
	if err != nil {
		fmt.Println(err)
		return 
	}
	fmt.Println("Connected to server")
	
	defer conn.Close()
	
	// write data to server
	fmt.Println("Sending data to server...")

	_, err = conn.Write([]byte("Hello from the client"))
	if err != nil {
		fmt.Println(err)
		return
	}
	
	reader := bufio.NewReader(conn)
	line, err := reader.ReadBytes('\n')


  if err != nil {
    if err == io.EOF {
      fmt.Println("Connection closed by remote end")
    } else {
      fmt.Println("Error reading from connection:", err)
    }
    return
  }

  response := string(line)
	
	fmt.Println(response)
}
