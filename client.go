package main

import "net"
import "fmt"

var PLAYER1 = [0x400]byte{0x0B}

func main() {

  // connect to this socket
  conn, _ := net.Dial("tcp", "127.0.0.1:7845")
  for {
    buffer := make([]byte, 1024)

    fmt.Print("Recibo: ")
    numRead, _ := conn.Read(buffer)
    fmt.Println(buffer[:numRead])
    conn.Write(buffer[:10])
  }
}
