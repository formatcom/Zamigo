package config

import (
  "fmt"
  "os"
  "encoding/json"
)

const MEN = 0xFF

type Config struct {
  DEBUG bool
  HOST string
  PORT string
  MEN int
}

func Get() Config {

  // Se lee la configuracion
  file, _ := os.Open("config.json")
  decoder := json.NewDecoder(file)
  conf := Config{}
  err := decoder.Decode(&conf)
  if err != nil {
    fmt.Println("Error al cargar la configuracion")
    os.Exit(1)
  }
  conf.MEN = MEN

  return conf

}
