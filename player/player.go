package player

import "net"

type Player struct {
  Id int
  Conn net.Conn
  State int
}
