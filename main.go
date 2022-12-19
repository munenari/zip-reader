package main

import (
	"embed"
	"log"

	"github.com/munenari/read-zip/handlers"
)

var (
	readBaseDir = "/mnt/n1"

	//go:embed static/*
	htmls embed.FS
)

func main() {
	e, err := handlers.Routes(htmls, readBaseDir)
	if err != nil {
		log.Fatalln("failed to init handlers,", err)
	}
	log.Fatalln(e.Start(":80"))
}
