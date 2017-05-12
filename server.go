package main

import (
  "fmt"
  "net"
  "os"
  "io"
  "bytes"
  "encoding/hex"
  "regexp"
)

const (
  DEBUG = true
  HOST = "0.0.0.0"
  PORT = "7845"
  MEN = 0xFF

  STATE_GUI = 0
  STATE_INIT_0 = 1
  STATE_INIT_1 = 2
  STATE_BOOT = 3
  STATE_PLAY = 4
)

// ZSNES V1.42
var ID142 = [MEN]byte{0x49, 0x44, 0xDE, 0x31, 0x34, 0x32, 0x20, 0x01, 0x01}
var ZSET0 = [MEN]byte{0x01}
var SAVEDATA = [MEN]byte{0x32} // Se asigna a None
var PLAYER1 = [MEN]byte{0x03}
var PLAYER2 = [MEN]byte{0x04}
var PLAYER3 = [MEN]byte{0x05}
var PLAYER4 = [MEN]byte{0x06}
var PLAYER5 = [MEN]byte{0x07}
var FRAME = [MEN]byte{0x00, 0x04}

// DATOS DEL JUEGO
var id int
var GAME = regexp.MustCompile(`.*\.sfc`) // Se valida la extension
var GAME0 = [MEN]byte{0x0B} // Espera que todos acepten la solicitud


type Player struct {
  id int
  conn net.Conn
  state int
}

var connections []Player

func main() {


  // Se pone a la escucha en el puerto 7845
  ln, err := net.Listen("tcp", HOST+":"+PORT)

  if err != nil {
    fmt.Println("Error al iniciar el servidor")
    os.Exit(1)
  }

  // Termina el listen 7845 al cerrar la aplicacion
  defer ln.Close()
  fmt.Println("El servidor esta a la escucha en el puerto "+ PORT)


  for {

    // Escucha cuando se conecta un cliente
    conn, err := ln.Accept()

    if err != nil {
      fmt.Println("Error aceptando")
      os.Exit(1)
    }

    id += 1
    player := Player{id, conn, STATE_INIT_0}

    // Grabar la conexion
    connections = append(connections, player)

    go handleRequest(player)

  }

}

func handleRequest(player Player) {

  for {
    buffer := make([]byte, MEN)
    _, err := player.conn.Read(buffer)

    if err != nil {
      if err == io.EOF {
        removeConn(player)
        player.conn.Close()
        return
      }
      return
    }

    if DEBUG && player.state != STATE_PLAY {
      fmt.Println(hex.Dump(buffer[:10]))
      fmt.Println(string(buffer))
    }

    // Activar el emulador
    if player.state == STATE_INIT_0 {
      if bytes.Equal(ID142[:], buffer) {
        player.conn.Write(ID142[:])
        player.conn.Write(SAVEDATA[:])
        player.state = STATE_INIT_1
        continue
      }else{
        player.conn.Close()
        return
      }
    }


    // Se habilita el jugador 1
    if bytes.Equal(ZSET0[:], buffer) {
      player.conn.Write(PLAYER1[:])
      player.state = STATE_GUI
      continue
    }

    // Se valida si se ha enviado un juego
    if GAME.MatchString(string(buffer)) {
      player.conn.Write(GAME0[:]) // Se indica que se cargara un juego
      player.state = STATE_BOOT
      if DEBUG {
        fmt.Println("Se esta iniciando el juego")
      }
      continue
    }

    // Se lanza el juego
    if player.state == STATE_BOOT  || player.state == STATE_PLAY {
      player.conn.Write(buffer)
    }

    // Se cambia el estado a jugando
    if player.state == STATE_BOOT && bytes.Equal(FRAME[:], buffer) {
      player.state = STATE_PLAY
      if DEBUG {
        fmt.Println("El juego se ha ejecutado")
      }
      continue
    }

    // Imprime debug cuando lo que recibe no es un frame
    if DEBUG && player.state == STATE_PLAY && bytes.Equal(FRAME[:], buffer) == false {
      fmt.Println(hex.Dump(buffer[:10]))
    }

  }

}


func removeConn(player Player) {
  var i int
  for i = range connections {
    if connections[i].id == player.id {
      fmt.Println("conn["+string(connections[i].id)+"] desconectado")
      break
    }
  }
  connections = append(connections[:i], connections[i+1:]...)
}
