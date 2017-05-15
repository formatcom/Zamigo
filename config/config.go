package config

import (
  "os"
  "encoding/json"
)

type Config struct {
  DEBUG bool
  HOST string
  PORT string
}

func Get() Config {

  // Se lee la configuracion
  file, _ := os.Open("config.json")
  decoder := json.NewDecoder(file)
  conf := Config{}
  err := decoder.Decode(&conf)
  if err != nil {
    conf = Config{true, "0.0.0.0", "7845"}
  }

  return conf

}
