package client

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	fmt.Println("Creating Client connection to server")
	conn, err := net.Dial("tcp", address)
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
		fmt.Println(err)
		return 
	}
	
	fmt.Println(line)
}
