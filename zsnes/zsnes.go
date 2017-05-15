package zsnes

import "regexp"

// ZSNES V1.42
var ID142 = [9]byte{0x49, 0x44, 0xDE, 0x31, 0x34, 0x32, 0x20, 0x01, 0x01}
var ZSET0 = [1]byte{0x01}
var SAVEDATA = [1]byte{0x32} // FLAG para colocar el savedata a none

// FLAG para cada jugador
var PLAYER1 = [1]byte{0x03}
var PLAYER2 = [1]byte{0x04}
var PLAYER3 = [1]byte{0x05}
var PLAYER4 = [1]byte{0x06}
var PLAYER5 = [1]byte{0x07}

// FLAG para lanzar y enviar datos del juego
var BOOT = [1]byte{0x0B}
var FRAME = [2]byte{0x00, 0x04}
var BOOT_1 = [2]byte{0x00, 0x01}
var BOOT_2 = [2]byte{0x00, 0x02}
var BOOT_3 = [1]byte{0xE5}
var OK = [2]byte{0x00, 0x00}

// FILTROS EXTRAS
var GAME = regexp.MustCompile(`\.sfc`) // Se valida la extension
var CHAT = regexp.MustCompile(`\>`)
