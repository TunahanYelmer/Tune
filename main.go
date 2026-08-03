/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"log"

	"github.com/tunahanyelmer/Tune/cmd"
	"github.com/tunahanyelmer/Tune/internal/providers/dummy"
	"github.com/tunahanyelmer/Tune/internal/service"
)

func main() {
	
	cmd.Music = service.NewMusicService(dummy.NewProvider())

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}

}
