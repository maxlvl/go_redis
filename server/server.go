package main

import (
  "net"
  "fmt"
  "syscall"
  "errors"
  "bytes"
  "encoding/binary"
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
  for {
    err := oneRequest(request)
    if err != nil {
      break
    }
  }
  return nil
}
  

func readFull(fd int, buf []byte) error {
  for buffer_length := len(buf); buffer_length > 0; {
    bytes_read, err := syscall.Read(fd, buf[:buffer_length])
    if err != nil {
      return err
    }

    if bytes_read <= 0 {
      return errors.New("Unexpected EOF")
    }
    buffer_length -= bytes_read
    buf = buf[bytes_read:]
  }
  return nil
}

func writeAll(fd int, buf []byte) error {
  for buffer_length := len(buf); buffer_length > 0; {
    bytes_written, err := syscall.Write(fd, buf[:buffer_length])
    if err != nil {
      return err
    }

    if bytes_written <= 0 {
      return errors.New("Unexpected Write Error")
    }

    buffer_length -= bytes_written
    buf = buf[bytes_written:]
  }
  return nil
}

const MAXSIZEMSG = 4096

func oneRequest(connection_fd int) error {
  fmt.Println("Handling OneRequest")
  // extract the 4 byte header and allow for null terminator (hence the +1)
  read_buffer := make([]byte, 4 + MAXSIZEMSG + 1)
  _, err := syscall.Read(connection_fd, read_buffer[:4])
  if err != nil {
    return err
  }

  // create a new buffer that reads from the 4 byte header onward
  var buf_len uint32
  new_buffer := bytes.NewBuffer(read_buffer[4:])
  err = binary.Read(new_buffer, binary.LittleEndian, &buf_len)
  if err != nil {
    return err
  }

  if buf_len > MAXSIZEMSG {
    return errors.New("Message too long")
  }

  // request body
  _, err = syscall.Read(connection_fd, read_buffer[4:4+buf_len])
  if err != nil {
    return err
  }

  read_buffer[4 + buf_len] = 0
  fmt.Println("client is saying %s\n", read_buffer[4: 4 + buf_len])

  // reply from server
  reply := []byte("World")
  write_buffer := make([]byte, 4+len(reply))

  buf_len = uint32(len(reply))
  binary.LittleEndian.PutUint32(write_buffer, buf_len)
  copy(write_buffer[4:], reply)
  _, err = syscall.Write(connection_fd, write_buffer)
  return err
}
