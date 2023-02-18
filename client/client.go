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

  send_message("4 this", conn)
  // send_message("Hello2", conn)
  // send_message("Hello3", conn)

}

func send_message(msg string, conn net.Conn) error {

  _, err := conn.Write([]byte(msg))
	if err != nil {
		fmt.Println(err)
		return err
	}
	
	reader := bufio.NewReader(conn)
	line, err := reader.ReadBytes('\n')


  if err != nil {
    if err == io.EOF {
      fmt.Println("Connection closed by remote end")
      return err
    } else {
      fmt.Println("Error reading from connection:", err)
      return err
    }
    return err
  }

  response := string(line)
	
	fmt.Println(response)
  return nil
}
