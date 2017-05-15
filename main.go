package main

import (
  "./server"
  "log"
  "os"
)


func main() {
  file, err := os.OpenFile("Zamigo.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
  if err != nil {
    log.Fatal("Error al crear el archivo Zamigo.log\n")
  }

  defer file.Close()

  log.SetOutput(file)

  server.Listen()
}
