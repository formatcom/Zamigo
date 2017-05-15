package server

import (
  "../config"
  "../player"
  "../zsnes"
  "fmt"
  "net"
  "os"
  "io"
  "bytes"
  "encoding/hex"
)


// Almacena los clientes conectados
var id int
var connections []*player.Player


func Listen() {

  conf := config.Get()

  // Se pone a la escucha en el puerto 7845
  ln, err := net.Listen("tcp", conf.HOST+":"+conf.PORT)

  if err != nil {
    fmt.Println("Error al iniciar el servidor", err)
    os.Exit(1)
  }

  // Termina el listen 7845 al cerrar la aplicacion
  defer ln.Close()
  fmt.Printf("El servidor esta a la escucha en el puerto %v\n", conf.PORT)


  for {

    // Escucha cuando se conecta un cliente
    conn, err := ln.Accept()

    if err != nil {
      fmt.Println("Error aceptando")
      os.Exit(1)
    }

    id += 1
    p := player.Player{id, conn, player.STATE_INIT_0}

    // Grabar la conexion
    connections = append(connections, &p)

    go handleRequest(&p)

  }

}

func handleRequest(p *player.Player) {

  conf := config.Get()

  for {
    buffer := make([]byte, config.MEN)
    _, err := p.Conn.Read(buffer)

    if err != nil {
      if err == io.EOF {
        removeConn(p)
        p.Conn.Close()
        return
      }
      return
    }

    /*if conf.DEBUG && p.State != player.STATE_PLAY {
      fmt.Println(hex.Dump(buffer[:10]))
      fmt.Println(string(buffer))
    }*/

    fmt.Printf("==>  p[%v] ", p.Id)
    fmt.Print(hex.Dump(buffer[:10]))
    fmt.Println(p)


    // Activar el emulador
    if p.State == player.STATE_INIT_0 {
      if bytes.Equal(zsnes.ID142[:], buffer) {
        p.Conn.Write(zsnes.ID142[:])
        p.Conn.Write(zsnes.SAVEDATA[:])
        p.State = player.STATE_INIT_1

        fmt.Printf("<==  p[%v] ", p.Id)
        fmt.Print(hex.Dump(zsnes.ID142[:10]))
        fmt.Println(p)
        fmt.Printf("<==  p[%v] ", p.Id)
        fmt.Print(hex.Dump(zsnes.SAVEDATA[:10]))
        fmt.Println(p)

        continue
      }else{
        p.Conn.Close()
        return
      }
    }


    // Se habilita el jugador 1
    if bytes.Equal(zsnes.ZSET0[:], buffer) {
      p.Conn.Write(zsnes.PLAYER1[:])
      p.State = player.STATE_GUI

      fmt.Printf("<==  p[%v] ", p.Id)
      fmt.Print(hex.Dump(zsnes.PLAYER1[:10]))
      fmt.Println(p)
      continue
    }


    if zsnes.CHAT.MatchString(string(buffer)) {
      broadcast(p, buffer, false)

      fmt.Printf("<==  p[all] sin mi [%v]", p.Id)
      fmt.Print(hex.Dump(buffer[:10]))
      fmt.Println(p)
      continue
    }

    // Se valida si se ha enviado un juego
    if p.State == player.STATE_GUI && zsnes.GAME.MatchString(string(buffer)) {
      broadcast(p, buffer, false) // Se envia el juegos a todos excluyendome
      p.Conn.Write(zsnes.GAME0[:]) // Se indica que se cargara un juego

      bootGame(player.STATE_BOOT)

      if conf.DEBUG {
        fmt.Println("Se esta iniciando el juego")
      }

      fmt.Printf("<==  p[all] sin mi [%v]", p.Id)
      fmt.Print(hex.Dump(buffer[:10]))
      fmt.Println(p)

      fmt.Printf("<==  p[%v]", p.Id)
      fmt.Print(hex.Dump(zsnes.GAME0[:10]))
      fmt.Println(p)
      continue
    }


    // Se lanza el juego
    if p.State == player.STATE_BOOT  || p.State == player.STATE_PLAY {
      broadcast(p, buffer, true)
      fmt.Printf("<==  p[all] [%v]", p.Id)
      fmt.Print(hex.Dump(buffer[:10]))
      fmt.Println(p)
    }

    // Se cambia el estado a jugando
    if p.State == player.STATE_BOOT && bytes.Equal(zsnes.FRAME[:], buffer) {
      p.State = player.STATE_PLAY
      if conf.DEBUG {
        fmt.Println("El juego se ha ejecutado")
      }
      continue
    }

    // Imprime debug cuando lo que recibe no es un frame
    if conf.DEBUG && p.State == player.STATE_PLAY {
      //fmt.Println(hex.Dump(buffer))
    }

  }

}


func bootGame(state int) {
  for i := range connections {
    connections[i].State = state
  }
}


func removeConn(p *player.Player) {
  var i int
  for i = range connections {
    if connections[i].Id == p.Id {
      fmt.Printf("player[%v] desconectado\n", connections[i].Id)
      break
    }
  }
  connections = append(connections[:i], connections[i+1:]...)
}


func broadcast(p *player.Player, data []byte, me bool) {
  for i := range connections {
    if connections[i].Id != p.Id || me {
      connections[i].Conn.Write(data)
    }
  }
}
