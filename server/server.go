package main

import (
  "net"
  "fmt"
  "syscall"
  "errors"
  "bytes"
  "encoding/binary"
)

const MAXSIZEMSG = 4096

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

    tcpconn, ok := conn.(*net.TCPConn)
    if !ok {
      fmt.Println("Error converting to TCP connection")
    }
    
    go handleConnection(tcpconn)
  }
}

func handleConnection(conn *net.TCPConn) error {
  defer conn.Close()
  fmt.Println("Handling Connection")

  file, _ := conn.File()
  fd := int(file.Fd())
  for {
    err := oneRequest(fd)
    if err != nil {
      break
    }
  }
  return nil
}
  

func oneRequest(connection_fd int) error {
  // extract the 4 byte header and allow for null terminator (hence the +1)
  read_buffer := make([]byte, 4 + MAXSIZEMSG + 1)
  err := readFull(connection_fd , read_buffer[:4])
  if err != nil {
    return err
  }

  // just checking if the message is too long or not - since the first 4 bytes 
  // indicate the length of the message
  var buf_len uint32
  fmt.Printf("Creating new buffer\n")
  new_buffer := bytes.NewBuffer(read_buffer[:4])
  fmt.Printf("value of new buffer is %s\n", new_buffer)
  err = binary.Read(new_buffer, binary.LittleEndian, &buf_len)
  if err != nil {
    fmt.Printf("Error happened: %s\n", err)
    return err
  }

  if buf_len > MAXSIZEMSG {
    fmt.Printf("Error happened - buf len is too long")
    return nil
  }

  // request body
  fmt.Printf("read buffer before second reading is %s\n", read_buffer[4:4+buf_len])
  err = readFull(connection_fd, read_buffer[4: 4 + buf_len])
  if err != nil {
    return err
  }

  read_buffer[4+buf_len] = 0
  fmt.Printf("client is saying %s\n", read_buffer[4:4 + buf_len])

  // reply from server
  reply := []byte("World\n")
  write_buffer := make([]byte, 4+len(reply))

  buf_len = uint32(len(reply))
  binary.LittleEndian.PutUint32(write_buffer, buf_len)
  copy(write_buffer[4:], reply)
  _, err = syscall.Write(connection_fd, write_buffer)
  return err
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

