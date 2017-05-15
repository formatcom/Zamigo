package zsnes

import (
  "../config"
  "regexp"
)

// ZSNES V1.42
var ID142 = [config.MEN]byte{0x49, 0x44, 0xDE, 0x31, 0x34, 0x32, 0x20, 0x01, 0x01}
var ZSET0 = [config.MEN]byte{0x01}
var SAVEDATA = [config.MEN]byte{0x32} // Se asigna a None
var PLAYER1 = [config.MEN]byte{0x03}
var PLAYER2 = [config.MEN]byte{0x04}
var PLAYER3 = [config.MEN]byte{0x05}
var PLAYER4 = [config.MEN]byte{0x06}
var PLAYER5 = [config.MEN]byte{0x07}
var FRAME = [config.MEN]byte{0x00, 0x04}

// DATOS DEL JUEGO
var GAME = regexp.MustCompile(`\.sfc`) // Se valida la extension
var CHAT = regexp.MustCompile(`\>`) // Se valida la extension
var GAME0 = [config.MEN]byte{0x0B} // Espera que todos acepten la solicitud
