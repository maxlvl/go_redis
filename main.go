package main
import (
  "fmt"
  "time"
)


func main() {
  fmt.Println("Starting server and client...")
  StartServer()
  time.Sleep(3 * time.Second)
  StartClient("0.0.0.0:8080")
}
