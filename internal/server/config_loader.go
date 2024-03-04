package server

import (
	"go-starter/config"
	"log"
)

func LoadConfig() *config.Config {
	c, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Can not load configs...", err)
	}
	return c
}
