package server

import (
  "../config"
  "../player"
  "../zsnes"
  "fmt"
  "log"
  "net"
  "io"
  "bytes"
  "encoding/hex"
)


// Almacena los clientes conectados
var id int
var connections []*player.Player


func Listen() {

  conf := config.Get()

  // Se pone a la escucha en el puerto asignado en la configuracion
  ln, err := net.Listen("tcp", conf.HOST+":"+conf.PORT)

  if err != nil {
    fmt.Printf("Error al iniciar el servidor [%v]\n", conf.PORT)
    log.Fatalf("Error al iniciar el servidor [%v]\n", conf.PORT)
  }

  // Termina el listen en el puerto configurado, al cerrar la aplicacion
  defer ln.Close()

  fmt.Printf("El servidor esta a la escucha en el puerto %v\n", conf.PORT)
  log.Printf("El servidor esta a la escucha en el puerto %v\n\n", conf.PORT)


  for {

    // Escucha cuando se conecta un cliente
    conn, err := ln.Accept()

    if err != nil {
      fmt.Println("Error al incorporar un nuevo cliente")
      log.Fatal("Error al incorporar un nuevo cliente\n")
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
    buffer := make([]byte, 256)
    readLen, err := p.Conn.Read(buffer)

    if err != nil {
      if err == io.EOF {
        removeConn(p)
        p.Conn.Close()
        return
      }
      return
    }


    log.Printf("[%v][RECIBO] cliente STATE[%v]\n%v\n\n", p.Id, p.State, hex.Dump(buffer[:readLen]))


    // Se valida la version del cliente. Solo soportamos la version 1.42
    if p.State == player.STATE_INIT_0 {
      if bytes.Equal(zsnes.ID142[:], buffer[:readLen]) {

        p.Conn.Write(zsnes.ID142[:])
        p.State = player.STATE_INIT_1

        log.Printf("[%v][ENVIO] cliente STATE[%v]\n%v\n\n", p.Id, p.State, hex.Dump(zsnes.ID142[:]))

        continue
      }else{
        p.Conn.Close()
        return
      }
    }

    // Se valida si recibimos el byte 0x01
    if p.State == player.STATE_INIT_1 {
      if bytes.Equal(zsnes.ZSET0[:], buffer[:readLen]) {
        p.Conn.Write(zsnes.ZSET0[:]) // Se confirma al cliente que se ha recibido
        p.Conn.Write(zsnes.SAVEDATA[:]) // Se establece por defecto en NONE el savedata
        p.Conn.Write(zsnes.PLAYER1[:]) // Se habilita el jugador 1

        p.State = player.STATE_GUI

        log.Printf("[%v][ENVIO] cliente STATE[%v]\n%v\n\n", p.Id, p.State, hex.Dump(zsnes.ZSET0[:]))
        log.Printf("[%v][ENVIO] cliente STATE[%v]\n%v\n\n", p.Id, p.State, hex.Dump(zsnes.SAVEDATA[:]))
        log.Printf("[%v][ENVIO] cliente STATE[%v]\n%v\n\n", p.Id, p.State, hex.Dump(zsnes.PLAYER1[:]))

        continue
      }
    }


    // Se valida si el request es un mensaje de chat
    if zsnes.CHAT.MatchString(string(buffer[:readLen])) {
      // Se reenvia el mensaje a todos los clientes excluyendome
      broadcast(p, buffer[:readLen], false)

      log.Printf("[%v][ENVIO][TODOS SIN MI] cliente STATE[%v]\n%v\n\n", p.Id, p.State, hex.Dump(buffer[:readLen]))

      continue
    }


    // Se valida si un cliente intenta lanzar un juego
    if p.State == player.STATE_GUI && zsnes.GAME.MatchString(string(buffer[:readLen])) {
      broadcast(p, buffer[:readLen], false) // Se envia el juegos a todos excluyendome
      p.Conn.Write(zsnes.BOOT[:]) // Se indica que se cargara un juego

      bootGame(player.STATE_BOOT) // Se asigna a todos el estado de STATE_BOOT

      log.Printf("[%v][ENVIO][TODOS SIN MI] cliente STATE[%v]\n%v\n\n", p.Id, p.State, hex.Dump(buffer[:readLen]))
      log.Printf("[%v][ENVIO] cliente STATE[%v]\n%v\n\n", p.Id, p.State, hex.Dump(zsnes.BOOT[:]))

      continue
    }


    // Se lanza el juego
    if p.State == player.STATE_BOOT {

      if bytes.Equal(zsnes.BOOT_3[:], buffer[:readLen]) {
        p.State = player.STATE_PLAY

        if conf.DEBUG {
          fmt.Println("El juego se ha ejecutado")
        }

      }

      broadcast(p, buffer[:readLen], true)

      log.Printf("[%v][ENVIO][TODOS] cliente STATE[%v]\n%v\n\n", p.Id, p.State, hex.Dump(buffer[:readLen]))

      continue
    }

    if p.State == player.STATE_PLAY {

      p.Conn.Write(zsnes.OK[:])

      if bytes.Equal(zsnes.FRAME[:], buffer[:readLen]) == false {
        broadcast(p, buffer[:readLen], false)

        log.Printf("[%v][ENVIO][TODOS SIN MI] cliente STATE[%v]\n%v\n\n", p.Id, p.State, hex.Dump(buffer[:readLen]))
      }

      log.Printf("[%v][ENVIO] cliente STATE[%v]\n%v\n\n", p.Id, p.State, hex.Dump(zsnes.OK[:]))

      continue
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
