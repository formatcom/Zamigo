package main

import (
  "encoding/hex"
  "net"
  "os"
  "fmt"
  "bytes"
)


func main() {

  // Me conecto con el socket
  conn, _ := net.Dial("tcp", "127.0.0.1:7845")

  next := 0
  var ZERO  = make([]byte, 2)
  var FRAME = [2]byte{0x00, 0x04}
  var nFrame int

  for {
    buffer := make([]byte, 256)

    readLen, err := conn.Read(buffer)

    if err != nil {
      os.Exit(1)
    }

    // Saltar los frames
    if bytes.Equal(FRAME[:], buffer[:readLen]) == false {
      fmt.Println(readLen)
      fmt.Println(hex.Dump(buffer[:readLen]))
      fmt.Printf("Frames: %v\n", nFrame)
      nFrame = 0
    }else{
      nFrame++
    }

    if next < 10 {
      conn.Write(buffer[:readLen])
      next++
    }else{
      conn.Write(ZERO)
    }

  }
}
